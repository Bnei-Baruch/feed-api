package common

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	"github.com/Bnei-Baruch/feed-api/instrumentation"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type InitDb func() (*sql.DB, error)
type ShutdownDb func(db *sql.DB) error

type Connection struct {
	DB       *sql.DB
	init     InitDb
	shutdown ShutdownDb

	m sync.Mutex
}

func MakeConnection(i InitDb, s ShutdownDb) *Connection {
	return (&Connection{nil, i, s, sync.Mutex{}}).MustConnect()
}

func (c *Connection) MustConnect() *Connection {
	c.m.Lock()
	defer c.m.Unlock()
	var err error
	c.DB, err = c.init()
	utils.Must(err)
	return c
}

func (c *Connection) Shutdown() error {
	c.m.Lock()
	defer c.m.Unlock()
	err := c.shutdown(c.DB)
	c.DB = nil
	return err
}

type ConnectionWithQuery struct {
	Connection *Connection
	Query      *queries.Query

	m sync.RWMutex
}

func (c *Connection) With(q *queries.Query) *ConnectionWithQuery {
	return &ConnectionWithQuery{c, q, sync.RWMutex{}}
}

func (c *ConnectionWithQuery) reconnect() error {
	c.m.Lock()
	defer c.m.Unlock()
	if err := c.Connection.Shutdown(); err != nil {
		return errors.Wrap(err, "Error shutting down while re-establishing connection.")
	}
	var err error
	if c.Connection.DB, err = c.Connection.init(); err != nil {
		return errors.Wrap(err, "Error initializing while re-establishing connection.")
	}
	return nil
}

func (c *ConnectionWithQuery) handleBadConnection(action func() error) error {
	var err error
	if err = action(); err != nil && strings.Contains(err.Error(), "could not establish connection") {
		if err = c.reconnect(); err != nil {
			return err
		}
		log.Infof("Re-established connection.")
		return action()
	}
	return err
}

func (c *ConnectionWithQuery) Bind(ctx context.Context, obj interface{}) error {
	return c.handleBadConnection(func() error {
		c.m.RLock()
		defer c.m.RUnlock()
		return c.Query.Bind(ctx, c.Connection.DB, obj)
	})
}

func (c *ConnectionWithQuery) Exec() (sql.Result, error) {
	var result sql.Result
	return result, c.handleBadConnection(func() error {
		c.m.RLock()
		defer c.m.RUnlock()
		var err error
		result, err = c.Query.Exec(c.Connection.DB)
		return err
	})
}

var (
	RemoteMdb         *sql.DB     // Readonly MDB
	LocalMdb          *sql.DB     // Local MDB
	LocalChroniclesDb *sql.DB     // Chronicles
	ModelsDb          *Connection // Models
)

func Init() time.Time {
	return InitWithDefault()
}

func GetUserPasswordFromConnectionString(cs string) (string, string, error) {
	re := regexp.MustCompile(`postgres://(.*):(.*)@`)
	match := re.FindStringSubmatch(cs)
	if len(match) != 3 {
		return "", "", errors.New("Unpexpected connection string.")
	}
	return match[1], match[2], nil
}

func InitModelsDb() (db *sql.DB, err error) {
	log.Info("Setting up connection to Models")
	if db, err = sql.Open("postgres", viper.GetString("data_models.url")); err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	var username, password string
	if username, password, err = GetUserPasswordFromConnectionString(viper.GetString("mdb.local_url")); err != nil {
		return
	} else {
		if _, err = queries.Raw(fmt.Sprintf("select dblink_connect('mdb_conn', 'dbname=mdb user=%s password=%s');", username, password)).Exec(db); err != nil {
			return
		}
		log.Infof("mdb_conn dblink_connected")
	}
	if username, password, err = GetUserPasswordFromConnectionString(viper.GetString("chronicles.local_url")); err != nil {
		return
	} else {
		if _, err = queries.Raw(fmt.Sprintf("select dblink_connect('chronicles_conn', 'dbname=chronicles user=%s password=%s');", username, password)).Exec(db); err != nil {
			return
		}
		log.Infof("chronicles_conn dblink_connected")
	}
	return
}

func InitWithDefault() time.Time {
	var err error
	clock := time.Now()

	log.Info("Setting up connection to MDB")

	RemoteMdb, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(RemoteMdb.Ping())
	log.Info("Initializing type registries")
	utils.Must(mdb.InitTypeRegistries(RemoteMdb))

	LocalMdb, err = sql.Open("postgres", viper.GetString("mdb.local_url"))
	utils.Must(err)
	utils.Must(LocalMdb.Ping())
	result, err := queries.Raw("SET session_replication_role = 'replica';").Exec(LocalMdb)
	log.Infof("Set local MDB as replice: %+v", result)
	utils.Must(err)

	log.Info("Setting up connection to Chronicles")
	LocalChroniclesDb, err = sql.Open("postgres", viper.GetString("chronicles.local_url"))
	utils.Must(err)
	utils.Must(LocalChroniclesDb.Ping())

	ModelsDb = MakeConnection(InitModelsDb, func(db *sql.DB) error { return db.Close() })

	boil.SetDB(LocalMdb)
	boil.DebugMode = viper.GetString("server.boiler-mode") == "debug"
	log.Infof("boil.DebugMode: %+v", boil.DebugMode)

	log.Infof("Settin up instrumentation")
	instrumentation.Stats.Init()

	return clock
}

func Shutdown() {
	log.Infof("Shutting down all databases.")
	utils.Must(RemoteMdb.Close())
	utils.Must(LocalMdb.Close())
	utils.Must(LocalChroniclesDb.Close())
	utils.Must(ModelsDb.Shutdown())
}

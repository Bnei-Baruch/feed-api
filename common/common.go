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

type InitDb func() (*sql.DB, map[string]string, error)
type ShutdownDb func(db *sql.DB) error

type Connection struct {
	DB       *sql.DB
	init     InitDb
	shutdown ShutdownDb
	Params   map[string]string

	m sync.Mutex
}

func MakeConnection(i InitDb, s ShutdownDb) *Connection {
	return (&Connection{nil, i, s, make(map[string]string), sync.Mutex{}}).MustConnect()
}

func (c *Connection) FillParams(p map[string]string) {
	for k, v := range c.Params {
		p[k] = v
	}
}

func (c *Connection) MustConnect() *Connection {
	c.m.Lock()
	defer c.m.Unlock()
	var err error
	c.DB, c.Params, err = c.init()
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
	if c.Connection.DB, c.Connection.Params, err = c.Connection.init(); err != nil {
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

func InitModelsDb() (db *sql.DB, params map[string]string, err error) {
	log.Info("Setting up connection to Models")
	params = make(map[string]string)
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
		params["mdb_conn"] = fmt.Sprintf("dbname=mdb user=%s password=%s", username, password)
	}
	if username, password, err = GetUserPasswordFromConnectionString(viper.GetString("chronicles.local_url")); err != nil {
		return
	} else {
		params["chronicles_conn"] = fmt.Sprintf("dbname=chronicles user=%s password=%s", username, password)
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

package common

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"

	"github.com/Bnei-Baruch/feed-api/databases/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

type InitDb func() (*sql.DB, error)
type ShutdownDb func(db *sql.DB) error

type Connection struct {
	db       *sql.DB
	init     InitDb
	shutdown ShutdownDb
}

func MakeConnection(i InitDb, s ShutdownDb) *Connection {
	return (&Connection{nil, i, s}).MustConnect()
}

func (c *Connection) MustConnect() *Connection {
	var err error
	c.db, err = c.init()
	utils.Must(err)
	return c
}

func (c *Connection) Shutdown() error {
	return c.shutdown(c.db)
}

type ConnectionWithQuery struct {
	Connection *Connection
	Query      *queries.Query
}

func (c *Connection) With(q *queries.Query) *ConnectionWithQuery {
	return &ConnectionWithQuery{c, q}
}

func (c *ConnectionWithQuery) Bind(ctx context.Context, obj interface{}) error {
	return c.Query.Bind(ctx, c.Connection.db, obj)
}

func (c *ConnectionWithQuery) Exec() (sql.Result, error) {
	return c.Query.Exec(c.Connection.db)
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

func InitModelsDb() (db *sql.DB, err error) {
	log.Info("Setting up connection to Models")
	if db, err = sql.Open("postgres", viper.GetString("data_models.url")); err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	if _, err = queries.Raw("select dblink_connect('mdb_conn', 'dbname=mdb user=postgres password=YjQ0MD');").Exec(db); err != nil {
		return
	}
	log.Infof("mdb_conn dblink_connect")
	if _, err = queries.Raw("select dblink_connect('chronicles_conn', 'dbname=chronicles user=postgres password=YjQ0MD');").Exec(db); err != nil {
		return
	}
	log.Infof("chronicles_conn dblink_connect")
	return
}

func InitWithDefault() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

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

	return clock
}

func Shutdown() {
	log.Infof("Shutting down all databases.")
	utils.Must(RemoteMdb.Close())
	utils.Must(LocalMdb.Close())
	utils.Must(LocalChroniclesDb.Close())
	utils.Must(ModelsDb.Shutdown())
}

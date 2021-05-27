package common

import (
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

var (
	RemoteMDB         *sql.DB // Readonly MDB
	LocalMDB          *sql.DB // Local MDB
	LocalChroniclesDB *sql.DB // Chronicles
	ModelsDB          *sql.DB // Models
)

func Init() time.Time {
	return InitWithDefault()
}

func InitWithDefault() time.Time {
	var err error
	clock := time.Now()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//log.SetLevel(log.WarnLevel)

	log.Info("Setting up connection to MDB")

	RemoteMDB, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(RemoteMDB.Ping())
	log.Info("Initializing type registries")
	utils.Must(mdb.InitTypeRegistries(RemoteMDB))

	LocalMDB, err = sql.Open("postgres", viper.GetString("mdb.local_url"))
	utils.Must(err)
	utils.Must(LocalMDB.Ping())
	result, err := queries.Raw("SET session_replication_role = 'replica';").Exec(LocalMDB)
	log.Infof("Set local MDB as replice: %+v", result)
	utils.Must(err)

	log.Info("Setting up connection to Chronicles")
	LocalChroniclesDB, err = sql.Open("postgres", viper.GetString("chronicles.local_url"))
	utils.Must(err)
	utils.Must(LocalChroniclesDB.Ping())

	log.Info("Setting up connection to Models")
	ModelsDB, err = sql.Open("postgres", viper.GetString("data_models.url"))
	utils.Must(err)
	utils.Must(ModelsDB.Ping())
	result, err = queries.Raw("select dblink_connect('mdb_conn', 'dbname=mdb user=postgres password=YjQ0MD');").Exec(ModelsDB)
	log.Infof("dblink_connect: %+v", result)
	utils.Must(err)
	result, err = queries.Raw("select dblink_connect('chronicles_conn', 'dbname=chronicles user=postgres password=YjQ0MD');").Exec(ModelsDB)
	log.Infof("dblink_connect: %+v", result)
	utils.Must(err)

	boil.SetDB(LocalMDB)
	boil.DebugMode = viper.GetString("server.boiler-mode") == "debug"
	log.Infof("boil.DebugMode: %+v", boil.DebugMode)

	return clock
}

func Shutdown() {
	utils.Must(RemoteMDB.Close())
	utils.Must(LocalMDB.Close())
	utils.Must(LocalChroniclesDB.Close())
	utils.Must(ModelsDB.Close())
}

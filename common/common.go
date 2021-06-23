package common

import (
	"database/sql"
	"time"

	"github.com/Bnei-Baruch/sqlboiler/boil"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/mdb"
	"github.com/Bnei-Baruch/feed-api/utils"
)

var (
	DB *sql.DB // MDB
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
	DB, err = sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	utils.Must(DB.Ping())

	boil.SetDB(DB)
	boil.DebugMode = viper.GetString("server.boiler-mode") == "debug"
	log.Info("Initializing type registries")
	utils.Must(mdb.InitTypeRegistries(DB))

	return clock
}

func Shutdown() {
	utils.Must(DB.Close())
}

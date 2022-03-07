package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/learn"
	"github.com/Bnei-Baruch/feed-api/utils"
	"github.com/Bnei-Baruch/feed-api/version"
)

var learnCmd = &cobra.Command{
	Use:   "learn",
	Short: "Learn recommendations",
	Run:   learnFn,
}

var prodChronicles bool

func init() {
	learnCmd.PersistentFlags().BoolVar(&prodChronicles, "prod_chronicles", false, "If true will scan prod for entries.")
	viper.BindPFlag("chronicles.prod_chronicles", learnCmd.PersistentFlags().Lookup("prod_chronicles"))
	RootCmd.AddCommand(learnCmd)
}

func learnFn(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	logLevelStr := viper.GetString("server.log-level")
	if logLevelStr == "" {
		logLevelStr = "info"
	}
	logLevel, err := log.ParseLevel(logLevelStr)
	utils.Must(err)
	log.Infof("Setting log level: %+v", logLevel)
	log.SetLevel(logLevel)

	log.Infof("Initializing learn common connections version %s", version.Version)
	common.Init()
	defer common.Shutdown()

	if err = learn.InitFeatures(); err != nil {
		log.Error("Failed initializing features: ", err)
	}

	if err = learn.Learn(prodChronicles, viper.GetString("chronicles.remote_api")); err != nil {
		log.Error("Failed learning: ", err)
	}
}

package cmd

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/api"
	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/events"
	"github.com/Bnei-Baruch/feed-api/utils"
	"github.com/Bnei-Baruch/feed-api/version"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Feed api server",
	Run:   serverFn,
}

var bindAddress string

func init() {
	serverCmd.PersistentFlags().StringVar(&bindAddress, "bind_address", "", "Bind address for server.")
	viper.BindPFlag("server.bind-address", serverCmd.PersistentFlags().Lookup("bind_address"))
	RootCmd.AddCommand(serverCmd)
}

func serverFn(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	logLevelStr := viper.GetString("server.log-level")
	if logLevelStr == "" {
		logLevelStr = "info"
	}
	logLevel, err := log.ParseLevel(logLevelStr)
	utils.Must(err)
	log.Infof("Setting log level: %+v", logLevel)
	log.SetLevel(logLevel)

	log.Infof("Starting feed api server version %s", version.Version)
	common.Init()
	defer common.Shutdown()
	shutDownEvents := events.RunListener()
	defer shutDownEvents()

	dataModels := data_models.MakeDataModels(common.LocalMdb, common.RemoteMdb, common.LocalChroniclesDb, common.ModelsDb, viper.GetString("chronicles.remote_api"))

	// TODO: Setup Rollbar
	// rollbar.Token = viper.GetString("server.rollbar-token")
	// rollbar.Environment = viper.GetString("server.rollbar-environment")
	// rollbar.CodeVersion = version.Version

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router := gin.New()
	router.Use(
		utils.LoggerMiddleware(),
		utils.DataStoresMiddleware(common.RemoteMdb, common.LocalMdb, common.LocalChroniclesDb, common.ModelsDb, dataModels),
		utils.ErrorHandlingMiddleware(),
		cors.New(corsConfig),
		utils.RecoveryMiddleware())

	api.SetupRoutes(router)

	log.Infoln("Running application")
	if cmd != nil {
		router.Run(viper.GetString("server.bind-address"))
	}

	// This would be reasonable once we'll have graceful shutdown implemented
	// if len(rollbar.Token) > 0 {
	// 	rollbar.Wait()
	// }
}

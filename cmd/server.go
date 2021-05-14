package cmd

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Bnei-Baruch/feed-api/api"
	"github.com/Bnei-Baruch/feed-api/common"
	"github.com/Bnei-Baruch/feed-api/data_models"
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

// Set DataModels in context.
func DataModelsMiddleware(dataModels *data_models.DataModels) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DATA_MODELS", dataModels)
		c.Next()
	}
}

// Set Chronicles in context.
func ChroniclesMiddleware(chroniclesDB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("CHRONICLES_DB", chroniclesDB)
		c.Next()
	}
}

func serverFn(cmd *cobra.Command, args []string) {
	log.Infof("Starting feed api server version %s", version.Version)
	common.Init()
	defer common.Shutdown()

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
		utils.DataStoresMiddleware(common.DB),
		ChroniclesMiddleware(common.CDB),
		DataModelsMiddleware(data_models.MakeDataModels(common.DB, common.CDB)),
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

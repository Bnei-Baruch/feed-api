package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hellofresh/health-go/v4"

	"github.com/Bnei-Baruch/feed-api/utils"
)

func HealthCheckHandler(c *gin.Context) {
	h, _ := health.New(health.WithChecks(
		health.Config{
			Name:    "remote_mdb",
			Timeout: time.Second,
			Check:   utils.PostgresNoOpenCheck(c.MustGet("MDB_DB").(*sql.DB)),
		},
		health.Config{
			Name:    "local_mdb",
			Timeout: time.Second,
			Check:   utils.PostgresNoOpenCheck(c.MustGet("LOCAL_MDB").(*sql.DB)),
		},
		health.Config{
			Name:    "local_chronicles",
			Timeout: time.Second,
			Check:   utils.PostgresNoOpenCheck(c.MustGet("LOCAL_CHRONICLES_DB").(*sql.DB)),
		},
		health.Config{
			Name:    "models_db",
			Timeout: time.Second,
			Check:   utils.PostgresNoOpenCheck(c.MustGet("MODELS_DB").(*sql.DB)),
		},
	))

	check := h.Measure(c.Request.Context())
	code := http.StatusOK
	if check.Status == health.StatusUnavailable {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, check)
}

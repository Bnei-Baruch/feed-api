package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/instrumentation"
)

// DataStoresMiddleware set data stores connections in context.
func DataStoresMiddleware(remoteMDB, localMDB, localChroniclesDB *sql.DB, ModelsDb, dataModels interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("MDB_DB", remoteMDB)
		c.Set("LOCAL_MDB", localMDB)
		c.Set("LOCAL_CHRONICLES_DB", localChroniclesDB)
		c.Set("MODELS_DB", ModelsDb)
		c.Set("DATA_MODELS", dataModels)
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.RequestURI() // some evil middleware modify this values

		c.Next()

		latency := time.Now().Sub(start)
		status := c.Writer.Status()
		instrumentation.Stats.RequestDurationHistogram.
			WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).
			Observe(float64(latency.Seconds()))

		log.WithFields(log.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"latency":    latency,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
		}).Info()
	}
}

// Recover with error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()
				err, ok := rval.(error)
				if !ok {
					err = errors.Errorf("panic: %s", rval)
				}
				c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
			}
		}()

		c.Next()
	}
}

func ValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "required"
	case "max":
		return fmt.Sprintf("cannot be longer than %s", e.Param())
	case "min":
		return fmt.Sprintf("must be longer than %s", e.Param())
	case "len":
		return fmt.Sprintf("must be %s characters long", e.Param())
	case "email":
		return "invalid email format"
	case "hexadecimal":
		return "invalid hexadecimal value"
	default:
		return "invalid value"
	}
}

func BindErrorMessage(err error) string {
	switch err.(type) {
	case *json.SyntaxError:
		e := err.(*json.SyntaxError)
		return fmt.Sprintf("json: %s [offset: %d]", e.Error(), e.Offset)
	case *json.UnmarshalTypeError:
		e := err.(*json.UnmarshalTypeError)
		return fmt.Sprintf("json: expecting %s got %s [offset: %d]", e.Type.String(), e.Value, e.Offset)
	default:
		return err.Error()
	}
}

// Handle all errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				switch e.Type {
				case gin.ErrorTypePublic:
					if e.Err != nil {
						log.Warnf("Public error: %s", e.Error())
						c.JSON(c.Writer.Status(), gin.H{"status": "error", "error": e.Error()})
					}

				case gin.ErrorTypeBind:
					// Keep the preset response status
					status := http.StatusBadRequest
					if c.Writer.Status() != http.StatusOK {
						status = c.Writer.Status()
					}

					switch e.Err.(type) {
					case validator.ValidationErrors:
						errs := e.Err.(validator.ValidationErrors)
						errMap := make(map[string]string)
						for field, err := range errs {
							msg := ValidationErrorMessage(err)
							log.WithFields(log.Fields{
								"field": field,
								"error": msg,
							}).Warn("Validation error")
							errMap[err.Field()] = msg
						}
						c.JSON(status, gin.H{"status": "error", "errors": errMap})
					default:
						log.WithFields(log.Fields{
							"error": e.Err.Error(),
						}).Warn("Bind error")
						c.JSON(status, gin.H{
							"status": "error",
							"error":  BindErrorMessage(e.Err),
						})
					}

				default:
					// Log all other errors
					log.Error(e.Err)
					// TODO: Uncomment after Rollbar integration.
					// LogRequestError(c.Request, e.Err)
				}
			}

			// If there was no public or bind error, display default 500 message
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError,
					gin.H{"status": "error", "error": "Internal Server Error"})
			}
		}
	}
}

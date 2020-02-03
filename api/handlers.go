package api

import (
	"database/sql"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/Bnei-Baruch/feed-api/core"
)

// Responds with JSON of given response or aborts the request with the given error.
func concludeRequest(c *gin.Context, resp interface{}, err *HttpError) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

type MoreResponse struct {
	Feed []core.ContentItem `json:"feed"`
}

func MoreHandler(c *gin.Context) {
	r := core.MoreRequest{}
	if c.Bind(&r) != nil {
		return
	}

	resp, err := handleMore(c.MustGet("MDB_DB").(*sql.DB), r)
	concludeRequest(c, resp, err)
}

func handleMore(db *sql.DB, r core.MoreRequest) (*MoreResponse, *HttpError) {
	log.Infof("r: %+v", r)
	feed := core.MakeFeed(db)
	if cis, err := feed.More(r); err != nil {
		return nil, NewInternalError(err)
	} else {
		return &MoreResponse{Feed: cis}, nil
	}
}

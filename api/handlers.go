package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Responds with JSON of given response or aborts the request with the given error.
func concludeRequest(c *gin.Context, resp interface{}, err *HttpError) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

type ItemsRequest struct {
}

type ItemsResponse struct {
}

func ItemsHandler(c *gin.Context) {
	r := ItemsRequest{}
	if c.Bind(&r) != nil {
		return
	}

	resp := ItemsResponse{}
	concludeRequest(c, resp, nil)
}

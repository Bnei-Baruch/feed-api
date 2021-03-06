package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Code int
	Err  error
	Type gin.ErrorType
}

func (e HttpError) Error() string {
	return e.Err.Error()
}

func (e HttpError) Abort(c *gin.Context) {
	c.String(e.Code, "%+v", e.Err)
}

func NewHttpError(code int, err error, t gin.ErrorType) *HttpError {
	return &HttpError{Code: code, Err: err, Type: t}
}

func NewNotFoundError() *HttpError {
	return &HttpError{Code: http.StatusNotFound, Type: gin.ErrorTypePublic}
}

func NewBadRequestError(err error) *HttpError {
	return NewHttpError(http.StatusBadRequest, err, gin.ErrorTypePublic)
}

func NewInternalError(err error) *HttpError {
	return NewHttpError(http.StatusInternalServerError, err, gin.ErrorTypePrivate)
}

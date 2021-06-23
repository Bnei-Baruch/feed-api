package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/more", MoreHandler)
	router.POST("/recommend", RecommendHandler)
	router.POST("/views", ViewsHandler)
}

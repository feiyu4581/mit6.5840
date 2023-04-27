package http_server

import (
	"github.com/gin-gonic/gin"
	"mapreduce/internal/pkg/middle"
)

func InitRouter(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(middle.CORSMiddleware())
	router.Use(middle.LogMiddleware())

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})
}

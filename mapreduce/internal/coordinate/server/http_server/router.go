package http_server

import (
	"github.com/gin-gonic/gin"
	"mapreduce/internal/coordinate/model"
	"mapreduce/internal/pkg"
	"mapreduce/internal/pkg/middle"
)

func InitRouter(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(middle.CORSMiddleware())
	router.Use(middle.LogMiddleware())

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(200, "pong")
	})

	router.POST("/task", func(context *gin.Context) {
		var task model.Task
		if err := context.BindJSON(&task); err != nil {
			context.JSON(400, err.Error())
			return
		}

		var allLines []pkg.Line
		for _, fileName := range task.FileNames {
			lines, err := pkg.ReadFile(fileName)
			if err != nil {
				context.JSON(400, err.Error())
				return
			}

			allLines = append(allLines, lines...)
		}

		if err := pkg.SplitLinesToFile(allLines, task.MNums, task.Name); err != nil {
			context.JSON(400, err.Error())
			return
		}

		context.JSON(200, "success")
	})
}

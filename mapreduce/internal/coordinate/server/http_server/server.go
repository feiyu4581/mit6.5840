package http_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mapreduce/internal/coordinate/model"
	"mapreduce/internal/coordinate/service/coordinator"
	"mapreduce/internal/pkg/utils"
	"net/http"
)

func NewHttpServer(port int) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	InitRouter(router)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
}

func NewTask(context *gin.Context) {
	var request model.TaskRequest
	if err := context.BindJSON(&request); err != nil {
		context.JSON(400, err.Error())
		return
	}

	var allLines []utils.Line
	for _, fileName := range request.FileNames {
		lines, err := utils.ReadFile(fileName)
		if err != nil {
			context.JSON(400, err.Error())
			return
		}

		allLines = append(allLines, lines...)
	}

	fileNames, err := utils.SplitLinesToFile(allLines, request.MNums, request.Name)
	if err != nil {
		context.JSON(400, err.Error())
		return
	}

	err = coordinator.GetService().AddTask(
		model.NewTask(fileNames, request.Name, request.Worker, request.MNums, request.RNums))

	if err != nil {
		context.JSON(400, err.Error())
		return
	}

	context.JSON(200, "success")
}

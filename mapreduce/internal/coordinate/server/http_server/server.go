package http_server

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

}

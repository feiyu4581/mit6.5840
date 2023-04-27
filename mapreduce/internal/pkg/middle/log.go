package middle

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mapreduce/internal/pkg/log"
	"net/http"
)

func getRequestBody(ctx *gin.Context) (string, error) {
	res := fmt.Sprintf("url=%s", ctx.Request.URL.String())
	if ctx.Request.Method != http.MethodGet {
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			return "", nil
		}

		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		res = fmt.Sprintf("%s, body=%s", res, string(bodyBytes))
	}

	return res, nil
}

func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		content, err := getRequestBody(ctx)
		if err != nil {
			log.Warn("handle log record error: %s", err.Error())
			return
		}

		log.Info(content)
	}
}

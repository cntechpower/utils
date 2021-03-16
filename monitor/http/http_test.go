package http

import (
	"net/http"
	"testing"

	"github.com/cntechpower/utils/log"

	"github.com/cntechpower/utils/os"
	"github.com/gin-gonic/gin"
)

func TestHttp(t *testing.T) {
	log.Init(log.WithStd(log.OutputTypeText),
		log.WithEs("main.unit-test.http", "10.0.0.2:9200"))
	s := gin.New()
	s.Use(GinMiddleware(WithLog()))
	s.GET("ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	go s.Run("0.0.0.0:8888")
	select {
	case <-os.ListenKillSignal():
		return
	}
}

package http

import (
	"net/http"
	"testing"

	"github.com/cntechpower/utils/tracing"

	"github.com/cntechpower/utils/log"

	"github.com/cntechpower/utils/os"
	"github.com/gin-gonic/gin"
)

func TestHttp(t *testing.T) {
	tracing.Init("unit-test", "")
	log.Init(
		log.WithStd(log.OutputTypeJson),
		//log.WithEs("main.unit-test.http", "http://10.0.0.2:9200"),
	)
	s := gin.New()
	s.Use(GinMiddleware(
		WithLog(false, true),
		WithTrace()))
	s.GET("ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	go s.Run("0.0.0.0:8888")
	select {
	case <-os.ListenKillSignal():
		return
	}
}

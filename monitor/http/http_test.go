package http

import (
	"net/http"
	"testing"

	"github.com/cntechpower/utils/log"
	"github.com/cntechpower/utils/os"
	"github.com/cntechpower/utils/tracing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHttp(t *testing.T) {
	tracing.Init("unit-test", "10.0.0.2:6831")
	log.Init(
		log.WithStd(log.OutputTypeJson),
		log.WithEs("main.unit-test.http", "http://10.0.0.2:9200"),
	)
	defer log.Close()
	s := gin.New()
	s.Use(GinMiddleware(
		WithLog(false, true),
		WithTrace(),
	))
	s.GET("ping", func(context *gin.Context) {
		log.NewHeader("ping").Infoc(context, "hello")
		context.String(http.StatusOK, "pong")
	})
	s.GET("metrics", gin.WrapH(promhttp.Handler()))

	go s.Run("0.0.0.0:8888")
	select {
	case <-os.ListenKillSignal():
		return
	}
}

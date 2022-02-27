package gc

import (
	"net/http"
	"testing"

	"github.com/cntechpower/utils/log"
	"github.com/cntechpower/utils/os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestHttp(t *testing.T) {
	s := gin.New()
	MetricsHandler()
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

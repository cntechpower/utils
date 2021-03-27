package main

import (
	"flag"
	"fmt"

	"github.com/cntechpower/utils/tracing"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cntechpower/utils/log"
	mGrpc "github.com/cntechpower/utils/monitor/grpc"
	mHttp "github.com/cntechpower/utils/monitor/http"
	unet "github.com/cntechpower/utils/net"
	uos "github.com/cntechpower/utils/os"
	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var grpcPort int
var httpPort int
var version string

const app = "MonitorTestService"

func init() {
	flag.IntVar(&grpcPort, "grpc-port", 2233, "grpc listen port")
	flag.IntVar(&httpPort, "http-port", 2234, "http listen port")
	flag.StringVar(&version, "version", version, "app version")
}

func StartGrpc(addr string) chan error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(mGrpc.GetUnaryServerInterceptor(
			//mGrpc.WithBlackList([]string{"/grpc.health.v1.Health/Check"}),
			mGrpc.WithLog(false, true),
			mGrpc.WithTrace(),
		)),
	)
	grpcExitChan := make(chan error, 1)
	h := log.NewHeader("grpc")
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	h.Infof("starting grpc serve at %v...", addr)
	l, err := unet.ListenTcp(addr)
	if err != nil {
		h.Fatalf("grpc listen tcp error: %v", err)
	}
	go func() {
		grpcExitChan <- server.Serve(l)
	}()

	return grpcExitChan
}

func StartHttp(addr string) chan error {
	h := log.NewHeader("http")
	httpExitChan := make(chan error, 1)
	h.Infof("starting http serve at %v...", addr)
	r := gin.New()
	r.Use(gin.Recovery(), mHttp.GinMiddleware(mHttp.WithBlackList([]string{"/ping", "/metrics"})))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	gin.SetMode(gin.ReleaseMode)

	go func() {
		httpExitChan <- r.Run(addr)
	}()
	return httpExitChan
}

func main() {
	tracing.Init("unit-test-server", "10.0.0.2:6831")
	log.Init(
		log.WithStd(log.OutputTypeJson),
		//log.WithEs("main.unit-test.grpc", "http://10.0.0.2:9200"),
	)
	h := log.NewHeader(app)
	grpcExitChan := StartGrpc(fmt.Sprintf(":%v", grpcPort))
	httpExitChan := StartHttp(fmt.Sprintf(":%v", httpPort))

	go uos.ListenTTINSignalLoop()

	serverExitChan := uos.ListenKillSignal()
	select {
	case <-serverExitChan:
		h.Fatalf("server existing...")
	case err := <-grpcExitChan:
		h.Fatalf("grpc exit with error: %v", err)
	case err := <-httpExitChan:
		h.Fatalf("http exit with error: %v", err)
	}
	return

}

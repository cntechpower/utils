package main

import (
	"flag"
	"fmt"

	"github.com/cntechpower/utils/consul"
	grpcRegister "github.com/cntechpower/utils/consul/discovery/register/grpc"
	uos "github.com/cntechpower/utils/os"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/hashicorp/consul/api"

	"github.com/cntechpower/utils/log"
	unet "github.com/cntechpower/utils/net"
	"google.golang.org/grpc"
)

var port int
var version string

const app = "TestService"

func init() {
	flag.IntVar(&port, "port", 2233, "listen port")
	flag.StringVar(&version, "version", version, "app version")
}

func main() {
	flag.Parse()
	log.InitLogger("")
	h := log.NewHeader(app)
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "10.0.0.2:8500"
	consul.Init(consulConfig)

	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	localIp, err := unet.GetFirstLocalIp()
	if err != nil {
		h.Fatalf("get ip error: %v", err)
	}

	if err := grpcRegister.Register(app, localIp, port); err != nil {
		panic(err)
	}

	go uos.ListenTTINSignalLoop()

	grpcExitChan := make(chan error, 1)
	go func() {
		addr := fmt.Sprintf("%v:%v", localIp, port)
		h.Infof("starting serve at %v...", addr)
		l, err := unet.ListenTcp(addr)
		if err != nil {
			h.Fatalf("listen tcp error: %v", err)
		}
		grpcExitChan <- server.Serve(l)
	}()

	serverExitChan := uos.ListenKillSignal()
	select {
	case <-serverExitChan:
		h.Infof("server existing...")
	case err := <-grpcExitChan:
		h.Fatalf("grpc exit with error: %v", err)
	}
	return

}

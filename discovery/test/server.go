package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

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
	server := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	localIp, err := unet.GetFirstLocalIp()
	if err != nil {
		h.Fatalf("get ip error: %v", err)
	}

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v:%v", app, localIp, strconv.Itoa(port)),
		Name:    app,
		Port:    port,
		Address: localIp,
		Check: &api.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("health-%v-%v:%v", app, localIp, strconv.Itoa(port)),
			Name:                           app,
			Interval:                       (time.Duration(5) * time.Second).String(),
			Timeout:                        (time.Duration(2) * time.Second).String(),
			GRPC:                           fmt.Sprintf("%v:%v", localIp, port),
			DeregisterCriticalServiceAfter: (time.Duration(10) * time.Second).String(),
		},
	}
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "10.0.0.2:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		h.Fatalf("get consul client error: %v", err)
	}
	if err := consulClient.Agent().ServiceRegister(reg); err != nil {
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

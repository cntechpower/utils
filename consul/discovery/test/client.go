package main

import (
	"context"
	"time"

	uos "github.com/cntechpower/utils/os"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc"

	consulResolver "github.com/cntechpower/utils/consul/discovery/resolver"
	"github.com/cntechpower/utils/log"

	"google.golang.org/grpc/resolver"
)

func main() {
	log.Init("")
	h := log.NewHeader("resolver_test_client")
	resolver.Register(consulResolver.NewBuilder("10.0.0.2:8500", time.Second*5))
	resolver.SetDefaultScheme(consulResolver.NAME)
	// Connect to the "userinfo" Consul service.
	_, err := grpc.DialContext(context.Background(), "consul:///TestService", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		panic(err)
	}
	serverExitChan := uos.ListenKillSignal()
	select {
	case <-serverExitChan:
		h.Infof("server existing...")
	}
	return
}

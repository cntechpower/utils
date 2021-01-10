package main

import (
	"context"
	"fmt"
	"time"

	mgrpc "github.com/cntechpower/utils/monitor/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	cc, err := grpc.Dial("127.0.0.1:2233", grpc.WithInsecure(), grpc.WithUnaryInterceptor(mgrpc.GetUnaryClientInterceptor()))
	if err != nil {
		panic(err)
	}
	for {
		_, err := grpc_health_v1.NewHealthClient(cc).Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "EMPTY"})
		fmt.Println(err)
		time.Sleep(time.Millisecond * 500)
	}

}

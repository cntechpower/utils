package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cntechpower/utils/tracing"

	"github.com/cntechpower/utils/log"

	mgrpc "github.com/cntechpower/utils/monitor/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	tracing.Init("unit-test-client", "")
	log.Init(
		log.WithStd(log.OutputTypeJson),
		//log.WithEs("main.unit-test.grpc", "http://10.0.0.2:9200"),
	)
	cc, err := grpc.Dial("127.0.0.1:2233",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(mgrpc.GetUnaryClientInterceptor(
			//mgrpc.WithBlackList([]string{"/grpc.health.v1.Health/Check"}),
			mgrpc.WithLog(false, true),
			mgrpc.WithTrace(),
		)))
	if err != nil {
		panic(err)
	}
	for {
		_, err := grpc_health_v1.NewHealthClient(cc).Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "EMPTY"})
		fmt.Println(err)
		time.Sleep(time.Millisecond * 500)
	}

}

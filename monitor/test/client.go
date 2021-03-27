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
	tracing.Init("unit-test-client", "10.0.0.2:6831")
	log.Init(
		log.WithStd(log.OutputTypeJson),
		//log.WithEs("main.unit-test.grpc", "http://10.0.0.2:9200"),
	)
	defer log.Close()
	cc, err := grpc.Dial("127.0.0.1:2233",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(mgrpc.GetUnaryClientInterceptor(
			//mgrpc.WithBlackList([]string{"/grpc.health.v1.Health/Check"}),
			mgrpc.WithLog(false, true),
			mgrpc.WithTrace(),
		)),
	)
	if err != nil {
		panic(err)
	}
	for {
		_, ctx := tracing.New(context.Background(), "test")
		_, _ = grpc_health_v1.NewHealthClient(cc).Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "/grpc.health.v1.Health/Check"})
		tracing.Do(ctx, "operation-1", func() error {
			time.Sleep(time.Millisecond)
			return nil
		})
		tracing.Do(ctx, "operation-2", func() error {
			time.Sleep(time.Millisecond)
			return nil
		})
		tracing.Do(ctx, "operation-3", func() error {
			time.Sleep(2 * time.Millisecond)
			return fmt.Errorf("fake error")
		})
		time.Sleep(time.Millisecond * 500)
	}

}

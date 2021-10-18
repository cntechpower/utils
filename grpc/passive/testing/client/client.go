package main

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/cntechpower/utils/grpc/passive/testing/pb"

	"github.com/cntechpower/utils/log"
	"github.com/cntechpower/utils/tracing"

	passive "github.com/cntechpower/utils/grpc/passive/client"
	xos "github.com/cntechpower/utils/os"
)

func call(h *log.Header) {
	for range time.Tick(time.Second) {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		span, ctx := tracing.New(ctx, "call")
		gc, err := passive.GetClientConn(ctx, "passive:///test-server", grpc.WithInsecure())
		if err != nil {
			h.Errorc(ctx, "passive.DialContext error: %v", err)
			span.Finish()
			continue
		}
		res, err := pb.NewServiceClient(gc).Ping(ctx, &pb.PingReq{})
		if err != nil {
			h.Errorc(ctx, "call ping error: %v", err)
		} else {
			h.Infoc(ctx, "got res %v", res.HostName)
		}
	}

}

func main() {
	log.Init(log.WithStd(log.OutputTypeJson))
	defer log.Close()

	tracing.Init("test-client", "10.0.0.2:6831")
	defer tracing.Close()

	header := log.NewHeader("test-client")
	err := passive.New(1).Start(2211)
	if err != nil {
		header.Fatalf("passive start error: %v", err)
	}
	go call(header)

	<-xos.ListenKillSignal()
}

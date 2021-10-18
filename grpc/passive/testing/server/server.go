package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cntechpower/utils/log"

	"github.com/cntechpower/utils/tracing"

	"google.golang.org/grpc"

	passive "github.com/cntechpower/utils/grpc/passive/server"
	"github.com/cntechpower/utils/grpc/passive/testing/pb"
	mGrpc "github.com/cntechpower/utils/monitor/grpc"
	xos "github.com/cntechpower/utils/os"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	header *log.Header
}

func (s *Server) Ping(ctx context.Context, _ *pb.PingReq) (resp *pb.PingResp, err error) {
	resp = &pb.PingResp{}
	s.header.Infoc(ctx, "got ping request")
	resp.HostName, err = os.Hostname()
	return
}

func main() {
	log.Init(log.WithStd(log.OutputTypeJson))
	defer log.Close()

	tracing.Init("test-server", "10.0.0.2:6831")
	defer tracing.Close()

	header := log.NewHeader("test-server")

	passive.New("test-server")
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(mGrpc.GetUnaryServerInterceptor(mGrpc.WithTrace())))
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	pb.RegisterServiceServer(server, &Server{header: header})
	grpcExitChan := make(chan error, 1)
	go func() {
		grpcExitChan <- server.Serve(passive.Listener())
	}()
	err := passive.AddClient("127.0.0.1:2211")
	if err != nil {
		header.Fatalf("addClient error: %v", err)
	}
	select {
	case <-xos.ListenKillSignal():
		fmt.Printf("existing...\n")
	case err := <-grpcExitChan:
		fmt.Printf("grpc exist with error: %v", err)
	}

}

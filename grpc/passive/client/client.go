package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/cntechpower/utils/log"
	"github.com/cntechpower/utils/tracing"

	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

const RegisterBytesLen = 128

var cli *clientD
var cliDMu sync.Mutex

var gcPool map[string]*grpc.ClientConn
var gcMu sync.Mutex

type clientD struct {
	listener     net.Listener
	listenerQuit chan struct{}
	size         int
	mu           sync.Mutex
	connPool     map[string]chan net.Conn
	header       *log.Header
	onRegister   func(string)
}

func New(size int) *clientD {
	c := &clientD{
		listenerQuit: make(chan struct{}, 1),
		size:         size,
		mu:           sync.Mutex{},
		connPool:     make(map[string]chan net.Conn),
		header:       log.NewHeader("grpc.passive.client"),
	}
	cli = c
	gcPool = make(map[string]*grpc.ClientConn)
	gcMu = sync.Mutex{}
	return c
}

func (c *clientD) Start(port int) (err error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return
	}
	c.listener = l
	go c.start()
	return
}

func (c *clientD) start() {
	defer func() {
		c.listenerQuit <- struct{}{}
	}()
	for {
		conn, err := c.listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				c.header.Infof("accept temporary error: %v", err)
				continue
			}
			c.header.Errorf("accept error: %v", err)
		}
		keyBs := make([]byte, RegisterBytesLen)

		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if _, err := io.ReadFull(conn, keyBs); nil != err {
			_ = conn.Close()
			c.header.Errorf("accept read key error: %v", err)
			continue
		}
		_ = conn.SetReadDeadline(time.Time{})

		key := string(bytes.Trim(keyBs, "\x00"))
		addr := conn.RemoteAddr().String()
		c.header.Infof("accept new connection from %v %v", key, addr)

		if c.onRegister != nil {
			c.onRegister(key)
		}

		c.mu.Lock()
		if c.connPool[key] == nil {
			c.connPool[key] = make(chan net.Conn, c.size)
		}
		pool := c.connPool[key]
		c.mu.Unlock()

	PUT:
		select {
		case pool <- conn:
		default:
			select {
			case oldConn := <-pool:
				_ = oldConn.Close()
				goto PUT
			default:
				_ = conn.Close()
			}
		}
	}
}

func (c *clientD) Stop() {
	_ = c.listener.Close()
	<-c.listenerQuit
	c.mu.Lock()
OUTER:
	for _, pool := range c.connPool {
		for {
			select {
			case conn := <-pool:
				_ = conn.Close()
			default:
				continue OUTER
			}
		}
	}
	c.mu.Unlock()
}

func (c *clientD) getConn(ctx context.Context, target string) (conn net.Conn, err error) {
	span, _ := tracing.New(ctx, "getConn")
	defer func() {
		if err != nil {
			c.header.Errorf("%+v", err)
			ext.LogError(span, err)
		}
		span.Finish()
	}()
	c.mu.Lock()
	connPool := c.connPool[target]
	c.mu.Unlock()

	if connPool == nil {
		err = fmt.Errorf("pgrpc error: no connection exists")
		c.header.Errorf("connect to %v error: no connection exists", target)
		return
	}

	select {
	case conn = <-connPool:
		return
	case <-ctx.Done():
		err = fmt.Errorf("pgrpc error: wait pgrpc connection timeout")
		c.header.Errorf("connect to %v error: wait pgrpc connection timeout", target)
		return
	}
}

func dialer(ctx context.Context, target string) (net.Conn, error) {
	cliDMu.Lock()
	c := cli
	cliDMu.Unlock()

	if nil == c {
		return nil, fmt.Errorf("pgrpc is stopping")
	}
	return c.getConn(ctx, target)
}

func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (gc *grpc.ClientConn, err error) {
	span, _ := tracing.New(ctx, "grpc.passive.DialContext")
	defer func() {
		if err != nil {
			ext.LogError(span, err)
		}
		span.Finish()
	}()
	kpParam := keepalive.ClientParameters{
		Time:                time.Second,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}
	options := make([]grpc.DialOption, 0, 2+len(opts))
	options = append(options, grpc.WithContextDialer(dialer), grpc.WithResolvers(builder), grpc.WithKeepaliveParams(kpParam))
	options = append(options, opts...)

	gc, err = grpc.DialContext(ctx, target, options...)
	if err == nil {
		_, err = health.NewHealthClient(gc).Check(ctx, &health.HealthCheckRequest{})
		if err != nil {
			_ = gc.Close()
		}
	}

	return
}

func GetClientConn(ctx context.Context, target string, opts ...grpc.DialOption) (gc *grpc.ClientConn, err error) {
	gcMu.Lock()
	gc, ok := gcPool[target]
	gcMu.Unlock()
	if ok {
		_, err = health.NewHealthClient(gc).Check(ctx, &health.HealthCheckRequest{})
		if err != nil {
			gcMu.Lock()
			delete(gcPool, target)
			gcMu.Unlock()
			_ = gc.Close()
		}
		return
	}
	gc, err = DialContext(ctx, target, opts...)
	if err == nil {
		_, err = health.NewHealthClient(gc).Check(ctx, &health.HealthCheckRequest{})
		if err == nil {
			gcMu.Lock()
			gcPool[target] = gc
			gcMu.Unlock()
		}
	}
	return
}

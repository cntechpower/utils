package server

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type connWithCloseSignal struct {
	conn   net.Conn
	target string

	closeChan  chan struct{}
	firstClose sync.Once
	server     *serverD

	firstRead int32 //atomic
}

func newConnWithCloseSignal(conn net.Conn, serverD *serverD, target string) *connWithCloseSignal {
	c := &connWithCloseSignal{
		conn:      conn,
		server:    serverD,
		target:    target,
		closeChan: make(chan struct{}, 0),
	}
	go func() {
		select {
		case <-c.closeChan:
		case <-serverD.closed:
			//when listener is closed, should close pgrpc connections which's waiting for first read
			//otherwise, these connections, are not considered as grpc-live connection, and will block grpc.GracefulStop()
			if 0 == atomic.LoadInt32(&c.firstRead) {
				_ = c.Close()
			}
		}
	}()
	return c
}

func (c *connWithCloseSignal) Read(b []byte) (n int, err error) {
	n, err = c.conn.Read(b)
	atomic.StoreInt32(&c.firstRead, 1)
	return n, err
}

func (c *connWithCloseSignal) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *connWithCloseSignal) Close() error {
	c.firstClose.Do(func() {
		c.server.onConnClose(c.target)
		close(c.closeChan)
	})
	return c.conn.Close()
}

func (c *connWithCloseSignal) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connWithCloseSignal) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connWithCloseSignal) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *connWithCloseSignal) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connWithCloseSignal) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cntechpower/utils/log"
)

const (
	RegisterBytesLen = 128
	MaxSizePerTarget = 2
)

var svr *serverD
var svrMu sync.Mutex

type serverD struct {
	serverID   string
	serverIDBS []byte
	mu         sync.Mutex
	connCounts map[string]int

	newConns map[ /*ip:port*/ string]chan net.Conn
	stopping bool
	dialQuit chan struct{}
	closed   chan struct{}
	header   *log.Header
}

func New(serverID string) *serverD {
	s := &serverD{
		serverID:   serverID,
		mu:         sync.Mutex{},
		connCounts: make(map[string]int),
		newConns:   make(map[string]chan net.Conn),
		dialQuit:   make(chan struct{}, 0),
		closed:     make(chan struct{}, 0),
		header:     log.NewHeader("grpc.passive.server"),
	}
	bs := make([]byte, RegisterBytesLen)
	copy(bs, serverID)
	s.serverIDBS = bs
	svr = s
	go s.dial()
	return s
}

func Listener() *serverD {
	svrMu.Lock()
	defer svrMu.Unlock()
	return svr
}

func (s *serverD) dial() {
	defer func() {
		close(s.dialQuit)
	}()
	for {
		s.mu.Lock()
		if s.stopping {
			s.mu.Unlock()
			return
		}
		for k, c := range s.connCounts {
			if c >= MaxSizePerTarget {
				continue
			}

			conn, err := net.DialTimeout("tcp", k, time.Second)
			if err != nil {
				s.header.Errorf("dial %v error: %v", k, err)
				continue
			}
			{
				tcpConn, _ := conn.(*net.TCPConn)
				_ = tcpConn.SetLinger(1)

				/*
					pgrpc still need tcp-layer keepalive
					tcp.Dial is called by pgrpc.Server, and if pgrpc.Client didn't use this connection in pgrpc-layer,
					the connection has no pgrpc-layer keepalive
				*/
				_ = tcpConn.SetKeepAlive(true)
				_ = tcpConn.SetKeepAlivePeriod(1 * time.Second)
			}

			if n, err := conn.Write(s.serverIDBS); nil != err || n < RegisterBytesLen {
				_ = conn.Close()
				s.header.Errorf("write target error: %v", err)
				continue
			}
			s.header.Infof("new connection to %v", k)

			s.connCounts[k]++
			s.newConns[k] <- newConnWithCloseSignal(conn, s, k)
		}
		s.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

//net.Listener impl
func (s *serverD) Accept() (net.Conn, error) {
	for {
		s.mu.Lock()
		if s.stopping {
			s.mu.Unlock()
			return nil, fmt.Errorf("pgrpc: listener closed")
		}

		for _, ch := range s.newConns {
			select {
			case conn := <-ch:
				s.mu.Unlock()
				return conn, nil
			default:
			}
		}
		s.mu.Unlock()

		time.Sleep(100 * time.Millisecond)
	}
}

// net.Listener impl
func (s *serverD) Close() error {
	s.stop()
	return nil
}

func (s *serverD) stop() {
	s.mu.Lock()
	if s.stopping {
		s.mu.Unlock()
		return
	}
	s.stopping = true
	s.mu.Unlock()

	<-s.dialQuit

	s.mu.Lock()
	for _, chs := range s.newConns {
	LOOP:
		for {
			select {
			case conn := <-chs:
				_ = conn.Close()
			default:
				break LOOP
			}
		}
	}
	s.mu.Unlock()

	close(s.closed)
}

//net.Listener impl
func (s *serverD) Addr() net.Addr {
	a, _ := net.ResolveIPAddr("tcp", "0.0.0.0:0") //mock addr
	return a
}

func (s *serverD) addClient(ipPort string) (err error) {
	defer func() {
		if err != nil {
			s.header.Errorf("add client %v error: %v", ipPort, err)
		} else {
			s.header.Infof("add client %v success", ipPort)
		}
	}()
	s.mu.Lock()
	if nil != s.newConns[ipPort] {
		s.mu.Unlock()
		err = fmt.Errorf("pgrpc error: serverd already has client %v", ipPort)
		return
	}
	s.newConns[ipPort] = make(chan net.Conn, MaxSizePerTarget)
	s.connCounts[ipPort] = 0
	s.mu.Unlock()
	return nil
}

func AddClient(ipPort string) error {
	svrMu.Lock()
	s := svr
	svrMu.Unlock()

	if nil == s {
		return fmt.Errorf("pgrpc error: no serverd")
	}
	if err := s.addClient(ipPort); nil != err {
		return err
	}
	return nil
}

func (s *serverD) removeClient(ipPort string) (err error) {
	defer func() {
		if err != nil {
			s.header.Infof("remove client %v success", ipPort)
		} else {
			s.header.Errorf("remove client %v error: %v", ipPort, err)
		}
	}()
	s.mu.Lock()
	if nil == s.newConns[ipPort] {
		s.mu.Unlock()
		return nil
	}
	delete(s.connCounts, ipPort)
	chs := s.newConns[ipPort]
	delete(s.newConns, ipPort)
	s.mu.Unlock()

	for {
		select {
		case conn := <-chs:
			_ = conn.Close()
		default:
			return nil
		}
	}
}

func RemoveClient(ipPort string) error {
	svrMu.Lock()
	s := svr
	svrMu.Unlock()

	if nil == s {
		return fmt.Errorf("pgrpc error: no serverd")
	}
	return s.removeClient(ipPort)
}

func (s *serverD) onConnClose(k string) {
	fmt.Println("onConnClose")
	s.mu.Lock()
	if s.connCounts[k] > 0 {
		s.connCounts[k]--
	}
	s.mu.Unlock()
}

package output

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type es struct {
	appId     string
	addr      string
	buffer    chan []byte
	esClient  *elasticsearch.Client
	closeChan chan struct{}
}

var HTTPTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   3 * time.Second,   // 连接超时时间
		KeepAlive: 300 * time.Second, // 保持长连接的时间
	}).DialContext, // 设置连接的参数
	MaxIdleConns:          100,               // 最大空闲连接
	IdleConnTimeout:       300 * time.Second, // 空闲连接的超时时间
	ExpectContinueTimeout: 30 * time.Second,  // 等待服务第一个响应的超时时间
	MaxIdleConnsPerHost:   100,               // 每个host保持的空闲连接数
}

func NewES(appId, addr string) (w *es, closer func()) {
	w = &es{
		appId:     appId,
		addr:      addr,
		buffer:    make(chan []byte, 10000),
		closeChan: make(chan struct{}, 0),
	}
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses:  []string{addr},
		MaxRetries: 3,
		Transport:  HTTPTransport,
	})
	if err != nil {
		panic(err)
	}
	w.esClient = c
	closer = w.Close
	go w.do()

	return
}
func (w *es) Write(p []byte) (n int, err error) {
	select {
	case w.buffer <- p:
		return
	case <-time.After(time.Millisecond):
		fmt.Printf("DROP LOG: %v", string(p))
		return
	}

}

func (w *es) do() {
	var err error
	var resp *esapi.Response
	for p := range w.buffer {

		resp, err = w.esClient.Index(w.appId, bytes.NewReader(p),
			w.esClient.Index.WithTimeout(time.Millisecond*100))
		if err == nil && !resp.IsError() {
			break
		}

		if err != nil || resp.IsError() {
			extraMsg := ""
			if resp != nil {
				extraMsg = resp.String()
			}
			fmt.Printf("es do error: %v, extra: %v\n", err, extraMsg)
		}
		//https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
		if err == nil && resp.Body != nil {
			_, _ = io.Copy(ioutil.Discard, resp.Body)
			_ = resp.Body.Close()
		}
	}
	close(w.closeChan)
}

func (w *es) Close() {
	close(w.buffer)

	select {
	case <-w.closeChan:
	case <-time.After(time.Second * 5):
	}
}

package output

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type esOutput struct {
	appId    string
	addr     string
	esClient *elasticsearch.Client
}

var HTTPTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,  // 连接超时时间
		KeepAlive: 300 * time.Second, // 保持长连接的时间
	}).DialContext, // 设置连接的参数
	MaxIdleConns:          100,               // 最大空闲连接
	IdleConnTimeout:       300 * time.Second, // 空闲连接的超时时间
	ExpectContinueTimeout: 30 * time.Second,  // 等待服务第一个响应的超时时间
	MaxIdleConnsPerHost:   100,               // 每个host保持的空闲连接数
}

func NewESOutput(appId, addr string) *esOutput {
	w := &esOutput{
		appId: appId,
		addr:  addr,
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

	return w
}
func (w *esOutput) Write(p []byte) (n int, err error) {
	var resp *esapi.Response
	for i := 0; i < 3; i++ {
		resp, err = w.esClient.Index(w.appId, bytes.NewReader(p),
			w.esClient.Index.WithTimeout(time.Second))
		if err == nil && !resp.IsError() {
			break
		}
	}
	//https://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang
	if err == nil && resp.Body != nil {
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
	}
	return
}
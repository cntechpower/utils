package tracing

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/cntechpower/utils/os"

	"github.com/opentracing/opentracing-go"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func TestSimpleTracing(t *testing.T) {
	ctx := context.Background()
	Init("unit-test", "10.0.0.2:6831")
	defer Close()

	span, ctx := New(ctx, "hello-to")
	span.SetTag("hello-to", "dujinyang")
	defer span.Finish()

	go httpServer()
	time.Sleep(5 * time.Millisecond)

	helloStr := formatString(ctx, "dujinyang")
	printHello(ctx, helloStr)
	println(TraceIdFromContext(ctx))
	select {
	case <-os.ListenKillSignal():
		return
	}

}

func httpServer() {
	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil {
			span = tracer.StartSpan("format", ext.RPCServerOption(spanCtx))
		} else {
			span, _ = New(context.Background(), "format")
		}
		w.Header().Set(TraceID, TraceIdFromSpan(span))
		defer span.Finish()
		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		span.LogFields(
			log.String("event", "string-format"),
			log.String("value", helloStr))
		w.Write([]byte(helloStr))
	})

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var span opentracing.Span
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err == nil {
			span = tracer.StartSpan("publish", ext.RPCServerOption(spanCtx))
		} else {
			span, _ = New(context.Background(), "publish")
		}
		w.Header().Set(TraceID, TraceIdFromSpan(span))
		defer span.Finish()
		helloStr := r.FormValue("helloStr")
		println(helloStr)
	})
	go http.ListenAndServe(":8081", nil)
}

func formatString(ctx context.Context, helloTo string) string {
	span, ctx := New(ctx, "formatString")
	defer span.Finish()
	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://localhost:8081/format?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	resp, err := DoTest(req)
	if err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}

	helloStr := string(resp)
	return helloStr
}

func printHello(ctx context.Context, helloStr string) {
	span, ctx := New(ctx, "printHello")
	defer span.Finish()
	v := url.Values{}
	v.Set("helloStr", helloStr)
	url := "http://localhost:8081/publish?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	if _, err := DoTest(req); err != nil {
		ext.LogError(span, err)
		panic(err.Error())
	}
}

func DoTest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}

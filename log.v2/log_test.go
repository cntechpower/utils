package log_v2

import (
	"context"
	"testing"
	"time"

	"github.com/cntechpower/utils/tracing"
)

func TestLog(t *testing.T) {
	InitWithES("main.anywhere-agent", "http://10.0.0.2:9200")
	a := map[string]interface{}{"a": 123}
	tracing.Init("test-agent", "")
	for i := 0; i < 100; i++ {
		_, ctx := tracing.New(context.Background(), "TestLog")
		InfoC(ctx, a, "hello world, %v", "dujinyang")
	}
	time.Sleep(5 * time.Second)
}

package log_v2

import (
	"context"
	"testing"

	"github.com/cntechpower/utils/tracing"
)

func TestLog(t *testing.T) {
	Init()
	tracing.Init("test-agent","")
	for i := 0; i < 100; i++ {
		_, ctx := tracing.New(context.Background(), "TestLog")
		InfoC(ctx, nil, "hello world, %v", "dujinyang")
	}
}

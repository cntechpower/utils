package log_v2

import (
	"context"
	"fmt"
	"testing"

	"github.com/cntechpower/utils/tracing"
)

func TestLog(t *testing.T) {

	tracing.Init("test-agent", "")
	defer tracing.Close()

	InitWithES("main.anywhere-agent", "http://10.0.0.2:9200")
	defer Close()

	a := map[string]interface{}{"a": 123}

	for i := 0; i < 10; i++ {
		_, ctx := tracing.New(context.Background(), fmt.Sprintf("TestLog-%v", i))
		InfoC(ctx, a, "hello world, %v", "dujinyang")
	}
}

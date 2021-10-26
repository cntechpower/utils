package log_v2

import (
	"context"
	"fmt"
	"testing"

	"github.com/cntechpower/utils/tracing"
)

func TestLog(t *testing.T) {

	tracing.Init("main.unit-test", "10.0.0.2:6831")
	defer tracing.Close()

	InitWithES("main.unit-test", "http://10.0.0.2:9200")
	//Init()
	defer Close()

	a := map[string]interface{}{"a": 123}

	span, ctx := tracing.New(context.Background(), fmt.Sprintf("TestLog"))
	tracing.SetSpanWithFields(span, a)
	for i := 0; i < 10; i++ {
		_ = tracing.DoCtxF(ctx, fmt.Sprintf("TestLog-%v", i), func(c context.Context) error {
			InfoC(c, a, "hello world, %v No.%v", "dujinyang", i)
			return nil
		}, map[string]interface{}{"loop": i})
	}
	span.Finish()
}

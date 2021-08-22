package benchmark

import (
	"context"
	"testing"
	"time"
)

func fakeWorker(ctx context.Context) (err error) {
	time.Sleep(time.Millisecond * 20)
	return
}

func TestBenchmark(t *testing.T) {
	Run(fakeWorker)
}

package benchmark

import (
	"testing"
	"time"
)

func fakeWorker() (err error) {
	time.Sleep(time.Millisecond * 20)
	return
}

func TestBenchmark(t *testing.T) {
	Run(fakeWorker)
}

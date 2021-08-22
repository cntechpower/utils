package benchmark

import (
	"context"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	uos "github.com/cntechpower/utils/os"
)

var report *LatencyReport

func worker(ctx context.Context, wg *sync.WaitGroup, fn func(c context.Context) error) {
	defer func() {
		wg.Done()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		start := time.Now()
		err := fn(ctx)
		if err != nil {
			log.Printf("%v\n", err)
		}
		report.Add(time.Since(start))
	}
}

func Run(fn func(ctx context.Context) error) {
	workerNum := 10
	wg := &sync.WaitGroup{}
	report = NewLatencyReport(10240000)
	if w := os.Getenv("WORKER"); w != "" {
		wi, err := strconv.Atoi(w)
		if err != nil {
			log.Fatalf("ENV WORKER = %v format error: %v", w, err)
		}
		workerNum = wi
	}
	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go worker(ctx, wg, fn)
	}
	log.Printf("Started %v worker\n", workerNum)
	<-uos.ListenKillSignal()
	cancel()
	wg.Wait()
	report.Report()
}

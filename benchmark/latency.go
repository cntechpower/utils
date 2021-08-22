package benchmark

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
)

type LatencyReport struct {
	stopSingle        bool
	totalRequestCount int64
	histogram         metrics.Histogram
}

func NewLatencyReport(size int) *LatencyReport {
	r := &LatencyReport{
		histogram: metrics.NewHistogram(metrics.NewUniformSample(size)),
	}
	go r.singleReport()
	return r
}

func (r *LatencyReport) Add(d time.Duration) {
	r.totalRequestCount++
	r.histogram.Update(d.Nanoseconds())
}
func (r *LatencyReport) singleReport() {
	for range time.NewTicker(time.Second * 2).C {
		if r.stopSingle {
			continue
		}
		fmt.Printf("Total Request Count: %v\n", r.totalRequestCount)
	}
}

func (r *LatencyReport) Report() {
	r.stopSingle = true
	fmt.Println("-------------------------Latency Report-------------------------")
	fmt.Printf("Total Request Count: %v\n", r.histogram.Count())
	fmt.Printf("Max Request Latency: %.3fms\n", float64(r.histogram.Max())/1e6)
	fmt.Printf("Min Request Latency: %.3fms\n", float64(r.histogram.Min())/1e6)
	res := r.histogram.Percentiles([]float64{0.1, 0.5, 0.9, 0.99, 0.999})
	fmt.Printf("Pt 10 Request Latency: %.3fms\n", res[0]/1e6)
	fmt.Printf("Pt 50 Request Latency: %.3fms\n", res[1]/1e6)
	fmt.Printf("Pt 90 Request Latency: %.3fms\n", res[2]/1e6)
	fmt.Printf("Pt 99 Request Latency: %.3fms\n", res[3]/1e6)
	fmt.Printf("Pt 999 Request Latency: %.3fms\n", res[4]/1e6)
	r.stopSingle = false
}

package gc

import (
	"runtime"
	"runtime/debug"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	numGCTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_nums",
			Help:        "gc_nums",
			ConstLabels: nil,
		}, []string{})

	totalGCPauseTimes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_pause_time_total",
			Help:        "gc_pause_time_total",
			ConstLabels: nil,
		}, []string{})
	lastGCPauseTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_pause_time_last",
			Help:        "gc_pause_time_last",
			ConstLabels: nil,
		}, []string{})

	totalObjectMalloc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_objects_malloc_total",
			Help:        "gc_objects_malloc_total",
			ConstLabels: nil,
		}, []string{})

	totalObjectFree = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_objects_free_total",
			Help:        "gc_objects_free_total",
			ConstLabels: nil,
		}, []string{})

	totalMemSys = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_mem_sys",
			Help:        "gc_mem_sys",
			ConstLabels: nil,
		}, []string{})

	totalHeapAlloc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_heap_alloc",
			Help:        "gc_heap_alloc",
			ConstLabels: nil,
		}, []string{})

	totalHeapSys = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_heap_sys",
			Help:        "gc_heap_sys",
			ConstLabels: nil,
		}, []string{})

	totalHeapInuse = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_heap_inuse",
			Help:        "gc_heap_inuse",
			ConstLabels: nil,
		}, []string{})

	totalHeapIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "gc_heap_idle",
			Help:        "gc_heap_idle",
			ConstLabels: nil,
		}, []string{})
)

func init() {
	prometheus.MustRegister(
		numGCTotal,
		totalGCPauseTimes,
		lastGCPauseTime,
		totalObjectMalloc,
		totalObjectFree,
		totalMemSys,
		totalHeapAlloc,
		totalHeapSys,
		totalHeapInuse,
		totalHeapIdle,
	)
}

func MetricsHandler() {
	gcStats := &debug.GCStats{}
	memStats := &runtime.MemStats{}
	go func() {
		for range time.NewTicker(time.Second).C {
			debug.ReadGCStats(gcStats)
			numGCTotal.WithLabelValues().Set(float64(gcStats.NumGC))
			totalGCPauseTimes.WithLabelValues().Set(float64(gcStats.PauseTotal.Milliseconds()))
			if len(gcStats.Pause) > 0 {
				lastGCPauseTime.WithLabelValues().Set(float64(gcStats.Pause[0].Milliseconds()))
			}

			runtime.ReadMemStats(memStats)
			totalObjectMalloc.WithLabelValues().Set(float64(memStats.Mallocs))
			totalObjectFree.WithLabelValues().Set(float64(memStats.Frees))
			totalMemSys.WithLabelValues().Set(float64(memStats.Sys))
			totalHeapAlloc.WithLabelValues().Set(float64(memStats.HeapAlloc))
			totalHeapSys.WithLabelValues().Set(float64(memStats.HeapSys))
			totalHeapInuse.WithLabelValues().Set(float64(memStats.HeapInuse))
			totalHeapIdle.WithLabelValues().Set(float64(memStats.HeapIdle))
		}
	}()

}

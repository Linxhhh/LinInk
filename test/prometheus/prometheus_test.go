package prometheus

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func TestCounter(t *testing.T) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "LinInk",
		Subsystem: "test",
		Name:      "counter_test",
	})
	prometheus.MustRegister(counter) // 注册

	go func() {
		for {
			counter.Inc()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
}

func TestGauge(t *testing.T) {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "LinInk",
		Subsystem: "test",
		Name:      "gauge_test",
	})
	prometheus.MustRegister(gauge) // 注册

	go func() {
		for {
			gauge.Inc()
			time.Sleep(2 * time.Second)
			gauge.Dec()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
}

func TestHistogram(t *testing.T) {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "LinInk",
		Subsystem: "test",
		Name:      "histogram_test",
		// Buckets 用于表示直方图中的区间
		// 该直方图的区间: [ 0~10, 10~20, 20~30, 30~40, 40~50 ]
		Buckets: []float64{10, 20, 30, 40, 50},
	})
	prometheus.MustRegister(histogram) // 注册

	go func() {
		for {
			histogram.Observe(15.0) // 该观测值，会被添加到 10~20 的区间
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
}

func TestSummary(t *testing.T) {
	summary := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "LinInk",
		Subsystem: "test",
		Name:      "summary_test",

		Objectives: map[float64]float64{
			0.50: 0.010,    // 50% 的情况，误差为 1%
			0.75: 0.010,    // 75% 的情况，误差为 1%
			0.90: 0.010,	// 90% 的情况，误差为 1%
			0.95: 0.005,	// 95% 的情况，误差为 0.5%
			0.99: 0.001,    // 99% 的情况，误差为 0.1%
		},
	})
	prometheus.MustRegister(summary) // 注册

	go func() {
		for {
			summary.Observe(float64(50 + rand.Int31n(300)))  // 模拟请求耗时
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
}

func TestSummaryVec(t *testing.T) {
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "LinInk",
		Subsystem: "test",
		Name:      "summary_test",

		Objectives: map[float64]float64{
			0.50: 0.010,    // 50% 的情况，误差为 1%
			0.75: 0.010,    // 75% 的情况，误差为 1%
			0.90: 0.010,	// 90% 的情况，误差为 1%
			0.95: 0.005,	// 95% 的情况，误差为 0.5%
			0.99: 0.001,    // 99% 的情况，误差为 0.1%
		},
	}, []string{"path", "method", "status"})
	prometheus.MustRegister(summary) // 注册

	go func() {
		for {
			summary.WithLabelValues("user/profile", "get", "200").Observe(float64(50 + rand.Int31n(300)))  // 模拟请求耗时
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
}
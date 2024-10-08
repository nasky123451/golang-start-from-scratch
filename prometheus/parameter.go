package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	// Counter metric
	requestCountBase = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method"},
	)

	// Histogram metric
	requestDurationBase = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	// 定義指標，包括 Gauge 和 Summary
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "status", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"}, // Expecting three labels: method, endpoint, and status
	)

	// 追蹤當前連接數的 Gauge
	currentConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_current_connections",
			Help: "Current number of active connections",
		},
	)

	// 追蹤請求的延遲 Summary
	requestLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_latency_seconds",
			Help:       "Summary of HTTP request latencies",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}, // 自定義分位數
		},
		[]string{"method", "path"},
	)
)

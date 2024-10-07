package prometheus

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Counter metric
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method"},
	)

	// Histogram metric
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

// init function registers metrics
func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func PrometheusBase() {
	// Start the Prometheus web UI
	go startPrometheus()

	// Start the application's HTTP server
	http.HandleFunc("/", handler)
	http.Handle("/metrics", promhttp.Handler()) // Route to serve Prometheus metrics
	log.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // Start timing

	// Process the request
	w.Write([]byte("Hello, World!"))

	// Update metrics
	requestCount.WithLabelValues(r.Method).Inc()
	requestDuration.WithLabelValues(r.Method).Observe(time.Since(start).Seconds())
}

// startPrometheus starts the Prometheus process
func startPrometheus() {
	cmd := exec.Command("prometheus", "--config.file=prometheus.yml")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Prometheus: %v", err)
	}
	log.Println("Prometheus started. Access it at http://localhost:9090")
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Prometheus exited with error: %v", err)
	}
}

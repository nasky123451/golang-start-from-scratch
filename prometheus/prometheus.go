package prometheus

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	requestCountBase.WithLabelValues(r.Method).Inc()
	requestDurationBase.WithLabelValues(r.Method).Observe(time.Since(start).Seconds())
}

package prometheus

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusApiApplication() {
	// Use sync.WaitGroup to manage multiple concurrent goroutines
	var wg sync.WaitGroup

	// Create an HTTP server that supports graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// Initialize database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := checkAndCreateTable(db); err != nil {
		log.Fatal(err)
	}

	// Start Prometheus UI in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		startPrometheus()
	}()

	// Handle HTTP request paths
	http.HandleFunc("/api/v1/resource", func(w http.ResponseWriter, r *http.Request) {
		resourceHandler(w, r, db)
	})

	http.HandleFunc("/api/v1/login", loginHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler()) // Provide Prometheus metrics

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	gracefulShutdown(server)

	wg.Wait() // Wait for all goroutines to finish
}

// Handle resource requests
func resourceHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	start := time.Now()

	// Query the database
	rows, err := db.Query("SELECT id, name, type, created_at, updated_at FROM resources")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		requestCount.WithLabelValues(r.Method, "500", r.URL.Path).Inc()
		return
	}
	defer rows.Close()

	// Process query results
	var resources []map[string]interface{}
	for rows.Next() {
		var id int
		var name, resourceType string
		var createdAt, updatedAt string
		if err := rows.Scan(&id, &name, &resourceType, &createdAt, &updatedAt); err != nil {
			http.Error(w, "Data error", http.StatusInternalServerError)
			requestCount.WithLabelValues(r.Method, "500", r.URL.Path).Inc()
			return
		}

		// Add resource information to the slice
		resource := map[string]interface{}{
			"id":         id,
			"name":       name,
			"type":       resourceType,
			"created_at": createdAt,
			"updated_at": updatedAt,
		}
		resources = append(resources, resource)
	}

	// Encode response as JSON
	response, err := json.Marshal(resources)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		requestCount.WithLabelValues(r.Method, "500", r.URL.Path).Inc()
		return
	}

	// Set response headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	// Update metrics
	requestCount.WithLabelValues(r.Method, http.StatusText(http.StatusOK), r.URL.Path).Inc()
	requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	requestLatency.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
}

// Handle login requests
func loginHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Simulate login processing
	w.Write([]byte("Login successful"))

	// Update metrics
	statusCode := http.StatusOK
	requestCount.WithLabelValues(r.Method, http.StatusText(statusCode), r.URL.Path).Inc()
	requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	requestLatency.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Handle graceful shutdown
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // Wait for shutdown signal
	log.Println("Shutting down server...")

	// Set a 5-second shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

package prometheus

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
)

// Initialize Prometheus metrics
func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(currentConnections)
	prometheus.MustRegister(requestLatency)
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

// Initialize database, connect to PostgreSQL
func initDB() (*sql.DB, error) {
	// Read the DATABASE_URL environment variable
	databaseURL := os.Getenv("DATABASE_URL")

	// If DATABASE_URL is not set, the default value is used
	if databaseURL == "" {
		// Use local connection by default
		databaseURL = "localhost"
	}

	// Build the complete connection string
	url := "postgres://postgres:henry@" + databaseURL + ":5432/test?sslmode=disable"

	// Use DATABASE_URL to connect directly to the database
	return sql.Open("postgres", url)
}

// checkAndCreateTable checks if the resources table exists, and creates it if it does not
func checkAndCreateTable(db *sql.DB) error {
	var exists bool
	// Check if the table exists
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'resources');").Scan(&exists)
	if err != nil {
		return err
	}

	// If the table does not exist, create it
	if !exists {
		createTableSQL := `
		CREATE TABLE resources (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			type VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
		_, err = db.Exec(createTableSQL)
		if err != nil {
			return err
		}
		fmt.Println("Table 'resources' created.")

		_, err = db.Exec("INSERT INTO resources (name, type) VALUES ($1, $2), ($3, $4), ($5, $6) ON CONFLICT DO NOTHING",
			"Resource A", "Type 1",
			"Resource B", "Type 2",
			"Resource C", "Type 3")
		if err != nil {
			return err
		}
	}

	return nil
}

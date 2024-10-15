package redis

import (
	"database/sql"
	"fmt"
)

// Initialize database, connect to PostgreSQL
func initDB() (*sql.DB, error) {
	connStr := "host=postgres-container user=postgres password=henry dbname=test sslmode=disable" // Configure as needed
	return sql.Open("postgres", connStr)
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

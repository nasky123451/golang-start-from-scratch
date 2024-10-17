package redis

import (
	"database/sql"
	"fmt"
	"os"
)

func initRedis() (*RedisClient, error) {
	// Read the DATABASE_URL environment variable
	redisURL := os.Getenv("REDIS_URL")

	// If REDIS_URL is not set, the default value is used
	if redisURL == "" {
		// Use local connection by default
		redisURL = "localhost"
	}

	// Build the complete connection string
	url := redisURL + ":6379"

	rdb := NewRedisClient(url, "", 0)

	// Use REDIS_URL to connect directly to the Redis
	return rdb, nil

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
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users');").Scan(&exists)
	if err != nil {
		return err
	}

	// If the table does not exist, create it
	if !exists {
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);`
		_, err = db.Exec(createTableSQL)
		if err != nil {
			return err
		}
		fmt.Println("Table 'users' created.")
	}

	// Check if the table exists
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'access_logs');").Scan(&exists)
	if err != nil {
		return err
	}

	// If the table does not exist, create it
	if !exists {
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS access_logs (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			access_time TIMESTAMP DEFAULT NOW()
		);`
		_, err = db.Exec(createTableSQL)
		if err != nil {
			return err
		}
		fmt.Println("Table 'access_logs' created.")
	}

	return nil
}

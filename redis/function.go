package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initRedis() (*redis.Client, error) {
	// Read the DATABASE_URL environment variable
	redisURL := os.Getenv("REDIS_URL")

	// If REDIS_URL is not set, the default value is used
	if redisURL == "" {
		// Use local connection by default
		redisURL = "localhost"
	}

	// Build the complete connection string
	url := redisURL + ":6379"

	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // 如果有密碼則設置
		DB:       0,  // 使用默認的 Redis 數據庫
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Use REDIS_URL to connect directly to the Redis
	return rdb, nil

}

// Initialize the database and connect to PostgreSQL
func initDB() (*pgxpool.Pool, error) {
	// Read the DATABASE_URL environment variable
	databaseURL := os.Getenv("DATABASE_URL")

	// If DATABASE_URL is not set, the default value is used
	if databaseURL == "" {
		databaseURL = "localhost"
	}

	// Build the complete connection string
	url := "postgres://postgres:henry@" + databaseURL + ":5432/test?sslmode=disable"

	// Use DATABASE_URL to connect directly to the database
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return pool, nil
}

// checkAndCreateTable checks if a table exists and creates it if it does not
func checkAndCreateTable(db *pgxpool.Pool, tableName, createTableSQL string) error {
	var exists bool
	// Check if the table exists
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1);", tableName).Scan(&exists)
	if err != nil {
		return err
	}

	// If the table does not exist, create it
	if !exists {
		_, err = db.Exec(context.Background(), createTableSQL)
		if err != nil {
			return err
		}
		fmt.Printf("Table '%s' created.\n", tableName)
	}

	return nil
}

// checkAndCreateTableBase checks and creates necessary tables for the base functionality
func checkAndCreateTableBase(db *pgxpool.Pool) error {
	// Check and create the users table
	usersTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`
	if err := checkAndCreateTable(db, "users", usersTableSQL); err != nil {
		return err
	}

	// Check and create the access_logs table
	accessLogsTableSQL := `
	CREATE TABLE IF NOT EXISTS access_logs (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		access_time TIMESTAMP DEFAULT NOW()
	);`
	if err := checkAndCreateTable(db, "access_logs", accessLogsTableSQL); err != nil {
		return err
	}

	return nil
}

// handleUserRegistration handles registration logic with proper feedback
func handleUserRegistration(ctx context.Context, db *pgxpool.Pool, name, email string) error {
	err := RegisterUser(ctx, db, name, email)
	if err != nil && err.Error() == "user with email "+email+" already exists" {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	log.Printf("User %s registered successfully", name)
	return nil
}

// RegisterUser registers a new user in PostgreSQL
func RegisterUser(ctx context.Context, db *pgxpool.Pool, name, email string) error {
	exists, err := userExistsBase(ctx, db, email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if exists {
		return fmt.Errorf("user with email %s already exists", email)
	}

	_, err = db.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}
	return nil
}

// userExists checks if a user with the given email already exists
func userExistsBase(ctx context.Context, db *pgxpool.Pool, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)"
	err := db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// checkAndCreateTableMoney checks and creates the money table
func checkAndCreateTableMoney(db *pgxpool.Pool) error {
	// Check and create the money table
	moneyTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		balance DECIMAL(10, 2) DEFAULT 0.00 NOT NULL
	);`
	if err := checkAndCreateTable(db, "users", moneyTableSQL); err != nil {
		return err
	}

	// Insert default users if the table is created for the first time
	users := []struct {
		username string
		balance  float64
	}{
		{"alice", 100.00},
		{"bob", 100.00},
	}

	for _, user := range users {
		exists, err := userExistsMoney(db, user.username)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.Exec(context.Background(), "INSERT INTO users (username, balance) VALUES ($1, $2)", user.username, user.balance)
			if err != nil {
				return fmt.Errorf("failed to insert user %s: %w", user.username, err)
			}
			fmt.Printf("Inserted user %s with balance %.2f.\n", user.username, user.balance)
		} else {
			fmt.Printf("User %s already exists, skipping insertion.\n", user.username)
		}
	}

	return nil
}

// handleUserRegistration handles registration logic with proper feedback
func handleUserRegistrationMoney(ctx context.Context, db *pgxpool.Pool, username string) error {
	err := RegisterUserMoney(ctx, db, username)
	if err != nil && err.Error() == "user "+username+" already exists" {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	log.Printf("User %s registered successfully", username)
	return nil
}

// RegisterUser registers a new user in PostgreSQL
func RegisterUserMoney(ctx context.Context, db *pgxpool.Pool, username string) error {
	exists, err := userExistsMoney(db, username)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if exists {
		return fmt.Errorf("user %s already exists", username)
	}

	_, err = db.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", username)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}
	return nil
}

// userExists checks if a user with the given username already exists
func userExistsMoney(db *pgxpool.Pool, username string) (bool, error) {
	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

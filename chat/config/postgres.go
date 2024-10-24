package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatMessage struct {
	ID      int       `json:"id"`      // Message ID
	Room    string    `json:"room"`    // Room name
	Sender  string    `json:"sender"`  // Sender name
	Content string    `json:"content"` // Message content
	Time    time.Time `json:"time"`    // Message sending time
}

func InitDB() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "localhost"
	}
	url := "postgres://postgres:henry@" + databaseURL + ":5432/test?sslmode=disable"

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

// checkAndCreateTableChat checks and creates the chat table
func CheckAndCreateTableChat(db *pgxpool.Pool) error {
	// Check and create the chat table
	chatTableSQL := `
		CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);`
	if err := checkAndCreateTable(db, "users", chatTableSQL); err != nil {
		return err
	}

	chatTableSQL = `
		CREATE TABLE chat_messages (
		id SERIAL PRIMARY KEY,
		room VARCHAR(255),
		sender VARCHAR(255),
		content TEXT,
		time TIMESTAMPTZ DEFAULT NOW()
	);`
	if err := checkAndCreateTable(db, "chat_messages", chatTableSQL); err != nil {
		return err
	}

	return nil
}

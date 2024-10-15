package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

var ctx = context.Background()

// RedisClient struct to wrap Redis operations
type RedisClient struct {
	Client *redis.Client
}

// PostgreSQLClient struct to wrap PostgreSQL operations
type PostgreSQLClient struct {
	Conn *pgx.Conn
}

// NewRedisClient initializes a Redis client
func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{Client: rdb}
}

// NewPostgreSQLClient initializes a PostgreSQL client
func NewPostgreSQLClient(connStr string) (*PostgreSQLClient, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &PostgreSQLClient{Conn: conn}, nil
}

// ClosePostgreSQLClient closes the PostgreSQL connection
func ClosePostgreSQLClient(pgClient *PostgreSQLClient) {
	pgClient.Conn.Close(context.Background())
}

// RegisterUser registers a new user in PostgreSQL
func RegisterUser(pgClient *PostgreSQLClient, name string, email string) error {
	_, err := pgClient.Conn.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	return err
}

// LogUserAccess logs a user's access in PostgreSQL and caches it in Redis
func LogUserAccess(pgClient *PostgreSQLClient, rdb *RedisClient, userID int) error {
	_, err := pgClient.Conn.Exec(ctx, "INSERT INTO access_logs (user_id, access_time) VALUES ($1, NOW())", userID)
	if err != nil {
		return err
	}

	// Cache the latest access time in Redis
	err = rdb.Client.Set(ctx, fmt.Sprintf("user:%d:last_access", userID), time.Now().Format(time.RFC3339), 0).Err()
	return err
}

// GetUserLastAccess retrieves the last access time of a user from Redis
func GetUserLastAccess(rdb *RedisClient, userID int) (string, error) {
	lastAccess, err := rdb.Client.Get(ctx, fmt.Sprintf("user:%d:last_access", userID)).Result()
	if err != nil {
		return "", err
	}
	return lastAccess, nil
}

func main() {
	// Initialize Redis and PostgreSQL clients
	rdb := NewRedisClient("localhost:6379", "", 0)
	pgClient, err := NewPostgreSQLClient("postgres://postgres:henry@localhost:5432/test")
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	defer ClosePostgreSQLClient(pgClient)

	// 1. Register a new user
	err = RegisterUser(pgClient, "Henry", "Henry@example.com")
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}

	// 2. Log user access
	err = LogUserAccess(pgClient, rdb, 1) // Assuming user ID is 1
	if err != nil {
		log.Fatalf("Failed to log user access: %v", err)
	}

	// 3. Get the user's last access time
	lastAccess, err := GetUserLastAccess(rdb, 1)
	if err != nil {
		log.Fatalf("Failed to get user's last access time: %v", err)
	}

	log.Printf("User 1's last access time: %s", lastAccess)
}

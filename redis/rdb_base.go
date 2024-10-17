package redis

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient wraps Redis operations
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient initializes a Redis client
func NewRedisClient(addr, password string, db int) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

// RegisterUser registers a new user in PostgreSQL
func RegisterUser(ctx context.Context, db *sql.DB, name, email string) error {
	exists, err := userExists(ctx, db, email)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if exists {
		return fmt.Errorf("user with email %s already exists", email)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return fmt.Errorf("failed to insert new user: %w", err)
	}
	return nil
}

// userExists checks if a user with the given email already exists
func userExists(ctx context.Context, db *sql.DB, email string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)"
	err := db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// LogUserAccess logs a user's access in PostgreSQL and caches it in Redis
func LogUserAccess(ctx context.Context, db *sql.DB, rdb *RedisClient, userID int) error {
	_, err := db.ExecContext(ctx, "INSERT INTO access_logs (user_id, access_time) VALUES ($1, NOW())", userID)
	if err != nil {
		return fmt.Errorf("failed to log user access: %w", err)
	}

	// Cache the latest access time in Redis
	err = rdb.Client.Set(ctx, fmt.Sprintf("user:%d:last_access", userID), time.Now().Format(time.RFC3339), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to cache access time in Redis: %w", err)
	}
	return nil
}

// GetUserLastAccess retrieves the last access time of a user from Redis
func GetUserLastAccess(ctx context.Context, rdb *RedisClient, userID int) (string, error) {
	key := fmt.Sprintf("user:%d:last_access", userID)
	lastAccess, err := rdb.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("no access log found for user %d", userID)
	}
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access time from Redis: %w", err)
	}
	return lastAccess, nil
}

// handleUserRegistration handles registration logic with proper feedback
func handleUserRegistration(ctx context.Context, db *sql.DB, name, email string) error {
	err := RegisterUser(ctx, db, name, email)
	if err != nil && err.Error() == "user with email "+email+" already exists" {
		log.Printf("Warning: %v", err)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	log.Printf("User %s registered successfully", name)
	return nil
}

// RedisBase handles the entire Redis and PostgreSQL process flow
func RedisBase() {
	rdb, _ := initRedis()
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := checkAndCreateTable(db); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Register user with better error handling
	if err := handleUserRegistration(ctx, db, "Henry", "Henry@example.com"); err != nil {
		log.Fatal(err)
	}

	// Log user access and cache in Redis
	if err := LogUserAccess(ctx, db, rdb, 1); err != nil {
		log.Fatalf("Failed to log user access: %v", err)
	}

	// Retrieve user's last access time
	if lastAccess, err := GetUserLastAccess(ctx, rdb, 1); err != nil {
		log.Fatalf("Failed to get user's last access time: %v", err)
	} else {
		log.Printf("User 1's last access time: %s", lastAccess)
	}
}

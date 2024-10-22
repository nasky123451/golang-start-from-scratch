package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	redisClient *redis.Client
	pgConn      *pgxpool.Pool
	ctx         = context.Background()
	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients     = make(map[*websocket.Conn]string)
	messages    []Message
	sessionTTL  = 10 * time.Minute
	mu          sync.Mutex
	logger      = logrus.New()
	authKey     = "YOUR_GENERATED_AUTH_KEY"
	secretKey   = "YOUR_GENERATED_SECRET_KEY"

	// Prometheus metrics
	registerUserCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_registration_total",
			Help: "Total number of user registrations",
		},
		[]string{"status"},
	)

	loginCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_login_total",
			Help: "Total number of user logins",
		},
		[]string{"status"},
	)
	messageCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "chat_messages_total",
		Help: "Total number of chat messages sent",
	}, []string{"room"})
)

// SetKey sets a key-value pair and optionally sets an expiration time
func SetKey(r *redis.Client, ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration).Err()
}

// GetKey retrieves the value of a key
func GetKey(r *redis.Client, ctx context.Context, key string) (string, error) {
	return r.Get(ctx, key).Result()
}

// KeyExists checks if a key exists
func KeyExists(r *redis.Client, ctx context.Context, key string) (bool, error) {
	res, err := r.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// DeleteKey deletes a key
func DeleteKey(r *redis.Client, ctx context.Context, key string) error {
	return r.Del(ctx, key).Err()
}

// ExpireKey sets an expiration time for a key
func ExpireKey(r *redis.Client, ctx context.Context, key string, expiration time.Duration) error {
	return r.Expire(ctx, key, expiration).Err()
}

func printRedisKeys(ctx context.Context, r *redis.Client) {
	keys, err := r.Keys(ctx, "*").Result() // 獲取所有鍵
	if err != nil {
		log.Printf("Error fetching keys: %v", err)
		return
	}

	for _, key := range keys {
		value, err := r.Get(ctx, key).Result() // 獲取鍵的值
		if err != nil {
			log.Printf("Error fetching value for key %s: %v", key, err)
			continue
		}
		log.Printf("Key: %s, Value: %s", key, value)
	}
}

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

// checkAndCreateTableChat checks and creates the chat table
func checkAndCreateTableChat(db *pgxpool.Pool) error {
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

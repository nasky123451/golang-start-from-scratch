package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
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

// LogUserAccess logs a user's access in PostgreSQL and caches it in Redis
func LogUserAccess(ctx context.Context, db *pgxpool.Pool, rdb *redis.Client, userID int) error {
	_, err := db.Exec(ctx, "INSERT INTO access_logs (user_id, access_time) VALUES ($1, NOW())", userID)
	if err != nil {
		return fmt.Errorf("failed to log user access: %w", err)
	}

	// 在 Redis 中緩存最新的訪問時間
	err = rdb.Set(ctx, fmt.Sprintf("user:%d:last_access", userID), time.Now().Format(time.RFC3339), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to cache access time in Redis: %w", err)
	}
	return nil
}

// GetUserLastAccess retrieves the last access time of a user from Redis
func GetUserLastAccess(ctx context.Context, rdb *redis.Client, userID int) (string, error) {
	key := fmt.Sprintf("user:%d:last_access", userID)
	lastAccess, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("no access log found for user %d", userID)
	}
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access time from Redis: %w", err)
	}
	return lastAccess, nil
}

// RedisBase handles the entire Redis and PostgreSQL process flow
func RedisBase() {
	rdb, _ := initRedis()
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := checkAndCreateTableBase(db); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 註冊用戶
	if err := handleUserRegistration(ctx, db, "Henry", "Henry@example.com"); err != nil {
		log.Fatal(err)
	}

	// 記錄用戶訪問
	if err := LogUserAccess(ctx, db, rdb, 1); err != nil {
		log.Fatalf("Failed to log user access: %v", err)
	}

	// 獲取用戶的最後訪問時間
	if lastAccess, err := GetUserLastAccess(ctx, rdb, 1); err != nil {
		log.Fatalf("Failed to get user's last access time: %v", err)
	} else {
		log.Printf("User 1's last access time: %s", lastAccess)
	}
}

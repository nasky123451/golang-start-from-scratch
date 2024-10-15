package api

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisClient encapsulates Redis operations
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient initializes the Redis client
func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{Client: rdb}
}

// SetKey sets a key-value pair and optionally sets an expiration time
func (r *RedisClient) SetKey(key string, value string, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// GetKey retrieves the value of a key
func (r *RedisClient) GetKey(key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// KeyExists checks if a key exists
func (r *RedisClient) KeyExists(key string) (bool, error) {
	res, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// DeleteKey deletes a key
func (r *RedisClient) DeleteKey(key string) error {
	return r.Client.Del(ctx, key).Err()
}

// ExpireKey sets an expiration time for a key
func (r *RedisClient) ExpireKey(key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

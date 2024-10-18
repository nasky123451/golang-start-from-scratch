package redis

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Global Redis and PostgreSQL clients
var (
	redisClient *redis.Client // Redis Client
	pgConn      *pgxpool.Pool
	ctx         = context.Background()
	sessionTTL  = 10 * time.Minute
	mu          sync.Mutex
)

// Use Redis distributed locks to handle sensitive operations (e.g., transfers)
func transferFunds(fromUser, toUser string, amount float64) error {
	// Use Redis distributed lock to protect the transfer operation
	lockKey := fmt.Sprintf("lock:%s:%s", fromUser, toUser)
	success, err := redisClient.SetNX(ctx, lockKey, "locked", 30*time.Second).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !success {
		return fmt.Errorf("unable to acquire lock, operation in progress")
	}
	defer redisClient.Del(ctx, lockKey)

	// Perform transfer operation, ensuring PostgreSQL transaction consistency
	tx, err := pgConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx) // Rollback the transaction if the operation fails

	var fromBalance, toBalance float64
	// Query the balances of fromUser and toUser
	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE username=$1", fromUser).Scan(&fromBalance)
	if err != nil {
		return fmt.Errorf("failed to query balance: %w", err)
	}

	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE username=$1", toUser).Scan(&toBalance)
	if err != nil {
		return fmt.Errorf("failed to query balance: %w", err)
	}

	// Check if the balance is sufficient
	if fromBalance < amount {
		return fmt.Errorf("insufficient funds for user %s", fromUser)
	}

	// Update balances
	_, err = tx.Exec(ctx, "UPDATE users SET balance=$1 WHERE username=$2", fromBalance-amount, fromUser)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	_, err = tx.Exec(ctx, "UPDATE users SET balance=$1 WHERE username=$2", toBalance+amount, toUser)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Transferred %.2f from %s to %s successfully.\n", amount, fromUser, toUser)
	return nil
}

// Subscribe to Redis Pub/Sub messages and handle expiration events
func subscribeToRedis() {
	pubsub := redisClient.Subscribe(ctx, "__keyevent@0__:expired")
	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Printf("Received expired event: %s\n", msg.Payload)
		// Handle expiration events (e.g., reload data or clean up sessions)
	}
}

// Simulate user activity across multiple nodes
func simulateUserActivity(username string) {
	sessionKey := fmt.Sprintf("session:%s", username)
	// Query or update session information
	val, err := redisClient.Get(ctx, sessionKey).Result()
	if err == redis.Nil {
		fmt.Println("Session expired or not found, loading from PostgreSQL...")
		var sessionData string
		err := pgConn.QueryRow(ctx, "SELECT username FROM users WHERE username=$1", username).Scan(&sessionData)
		if err != nil {
			log.Printf("Failed to load session from PostgreSQL: %v\n", err)
		}

		// Store the session back in Redis
		err = redisClient.Set(ctx, sessionKey, sessionData, sessionTTL).Err()
		if err != nil {
			log.Printf("Failed to store session in Redis: %v\n", err)
		}
		fmt.Printf("Session for user %s restored to Redis.\n", username)
	} else if err != nil {
		log.Printf("Failed to retrieve session from Redis: %v\n", err)
	} else {
		fmt.Printf("Session for user %s found in Redis: %s\n", username, val)
	}
}

func RedisTransferMoney() {
	var err error
	// Redis single-node configuration
	redisClient, err = initRedis()
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}

	// PostgreSQL configuration
	pgConn, err = initDB()
	if err != nil {
		log.Fatalf("Error initializing PostgreSQL: %v", err)
	}

	if err := checkAndCreateTableMoney(pgConn); err != nil {
		log.Fatalf("Error checking/creating money table: %v", err)
	}

	// Register users
	users := []struct {
		username string
		balance  float64
	}{
		{"alice", 100.00},
		{"bob", 100.00},
	}

	for _, user := range users {
		if err := handleUserRegistrationMoney(ctx, pgConn, user.username); err != nil {
			log.Fatal(err)
		}
	}

	// Start listening to Redis expiration events
	go subscribeToRedis()

	// Simulate transfer operation
	if err = transferFunds("alice", "bob", 50.0); err != nil {
		log.Fatalf("Error transferring funds: %v", err)
	}

	// Simulate user activity
	simulateUserActivity("alice")
	simulateUserActivity("bob")
}

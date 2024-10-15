package api

import (
	"testing"
	"time"
)

func TestRedisClient(t *testing.T) {
	// Initialize Redis client
	rdb := NewRedisClient("localhost:6379", "", 0)

	// Key and value for testing
	key := "test_key"
	value := "test_value"
	expiration := 5 * time.Second

	// 1. Test setting a key-value pair
	err := rdb.SetKey(key, value, expiration)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// 2. Test retrieving the value of a key
	got, err := rdb.GetKey(key)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if got != value {
		t.Errorf("Expected value %s, but got %s", value, got)
	}

	// 3. Test if the key exists
	exists, err := rdb.KeyExists(key)
	if err != nil {
		t.Fatalf("Failed to check if key exists: %v", err)
	}
	if !exists {
		t.Errorf("Expected key %s to exist", key)
	}

	// 4. Test setting expiration time for a key
	err = rdb.ExpireKey(key, 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to set key expiration: %v", err)
	}
	time.Sleep(3 * time.Second) // Wait for the key to expire

	// 5. Test that the key no longer exists after expiration
	exists, err = rdb.KeyExists(key)
	if err != nil {
		t.Fatalf("Failed to check if key exists after expiration: %v", err)
	}
	if exists {
		t.Errorf("Expected key %s to be expired", key)
	}

	// 6. Test deleting the key
	err = rdb.SetKey(key, value, 0) // Set the key again
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}
	err = rdb.DeleteKey(key)
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}
	exists, err = rdb.KeyExists(key)
	if err != nil {
		t.Fatalf("Failed to check if key exists after deletion: %v", err)
	}
	if exists {
		t.Errorf("Expected key %s to be deleted", key)
	}
}

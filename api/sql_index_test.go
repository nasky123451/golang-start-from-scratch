package api_test

import (
	"testing"

	"example.com/m/api"
)

func TestPostgresIndex(t *testing.T) {
	// Database connection string (adjust to your settings)
	connString := "postgres://postgres:henry@localhost:5432/test"

	// Test connection to DB
	db, err := api.PoolConnectToDB(connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	// Ensure db connection is closed after the test
	defer db.Close()

	// Check and create tables
	if err := api.CheckAndCreateTableChat(db); err != nil {
		t.Fatalf("Error checking/creating chat table: %v", err)
	}

	// Now call createIndexIfNotExists to check and create the index
	if err := api.CreateIndexIfNotExists(db, "idx_customer_order_date", "orders"); err != nil {
		t.Fatalf("Error creating index: %v", err)
	}

	// Check the query plan using EXPLAIN ANALYZE
	if err := api.ExplainAnalyze(db, "SELECT * FROM orders WHERE customer_id = 123 AND order_date >= '2024-01-01'"); err != nil {
		t.Fatalf("Error explaining query: %v", err)
	}

	// Drop index
	if err := api.DropIndex(db, "idx_customer_order_date"); err != nil {
		t.Fatalf("Error dropping index: %v", err)
	}

	// Check the query plan again after dropping index
	if err := api.ExplainAnalyze(db, "SELECT * FROM orders WHERE customer_id = 123 AND order_date >= '2024-01-01'"); err != nil {
		t.Fatalf("Error explaining query after dropping index: %v", err)
	}

	// Now call createIndexIfNotExists to check and create the index
	if err := api.CreateIndexIfNotExists(db, "idx_customer_order_date", "orders"); err != nil {
		t.Fatalf("Error creating index after drop: %v", err)
	}

	// Check the query plan after re-creating index
	if err := api.ExplainAnalyze(db, "SELECT * FROM orders WHERE customer_id = 123 AND order_date >= '2024-01-01'"); err != nil {
		t.Fatalf("Error explaining query after re-creating index: %v", err)
	}

	// Re-index the index
	if err := api.ReIndex(db, "idx_customer_order_date"); err != nil {
		t.Fatalf("Error re-indexing: %v", err)
	}
}

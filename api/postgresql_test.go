package api

import (
	"context"
	"testing"
)

func TestPostgresFunctions(t *testing.T) {
	// Database connection string (adjust to your settings)
	connString := "postgres://postgres:henry@localhost:5432/test"

	// Test connection to DB
	conn, err := ConnectToDB(connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Test table creation
	err = CreateTable(conn)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Test inserting a user
	err = InsertUser(conn, "John Doe")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Test querying users
	users, err := QueryUsers(conn)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}
	if len(users) != 1 || users[0] != "User: 1, Name: John Doe" {
		t.Errorf("Expected 1 user, got: %v", users)
	}

	// Test deleting the user
	err = DeleteUser(conn, "John Doe")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Test querying users after deletion
	users, err = QueryUsers(conn)
	if err != nil {
		t.Fatalf("Failed to query users: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected no users, got: %v", users)
	}

	// Clean up by dropping the table
	err = DropTable(conn)
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
}

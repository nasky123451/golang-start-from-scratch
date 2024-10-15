package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ConnectToDB establishes a connection to the PostgreSQL database
func ConnectToDB(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}

// CreateTable creates the users table if it does not exist
func CreateTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT)")
	return err
}

// InsertUser inserts a user into the users table
func InsertUser(conn *pgx.Conn, name string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO users (name) VALUES ($1)", name)
	return err
}

// QueryUsers retrieves all users from the users table
func QueryUsers(conn *pgx.Conn) ([]string, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("error querying data: %v", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("error scanning data: %v", err)
		}
		users = append(users, fmt.Sprintf("User: %d, Name: %s", id, name))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}
	return users, nil
}

// DeleteUser deletes a user from the users table
func DeleteUser(conn *pgx.Conn, name string) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM users WHERE name = $1", name)
	return err
}

// DropTable drops the users table
func DropTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "DROP TABLE IF EXISTS users")
	return err
}

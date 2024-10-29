package server_test

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"

	"example.com/m/server"
)

// TestHandleConnection tests the basic functionality of the TCP server
func TestHandleConnection(t *testing.T) {
	// Start the server in a goroutine so that it runs in the background
	go func() {
		server.TCPIPServer()
	}()

	// Wait a short moment for the server to start
	time.Sleep(100 * time.Millisecond)

	// Connect to the server as a client
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Test case 1: Send "ping" and expect "pong"
	testSendAndReceive(t, conn, "ping", "pong")

	// Test case 2: Send "hello" and expect "You said: hello"
	testSendAndReceive(t, conn, "hello", "You said: hello")

	// Test case 3: Send an empty message and expect "You said: "
	testSendAndReceive(t, conn, "", "You said:")
}

// Helper function to send a message to the server and verify the response
func testSendAndReceive(t *testing.T, conn net.Conn, input string, expectedResponse string) {
	// Write the message to the server
	_, err := conn.Write([]byte(input + "\n"))
	if err != nil {
		t.Fatalf("Failed to send message to server: %v", err)
	}

	// Read the server's response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response from server: %v", err)
	}

	// Trim any extra spaces and compare the response with the expected output
	response = strings.TrimSpace(response)
	if response != expectedResponse {
		t.Errorf("Expected response %q, got %q", expectedResponse, response)
	}
}

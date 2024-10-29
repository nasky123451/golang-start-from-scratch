package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/m/server"
	"github.com/gorilla/websocket"
)

// TestWebsocketHandler tests the WebSocket server's functionality
func TestWebsocketHandler(t *testing.T) {
	// Create a test HTTP server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(server.WsHandler))
	defer server.Close()

	// Replace "http" with "ws" in the server URL to form the WebSocket URL
	wsURL := "ws" + server.URL[len("http"):]

	// Dial the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer ws.Close()

	// Test case 1: Send a message and expect the same message back
	messageToSend := "Hello, WebSocket!"
	if err := ws.WriteMessage(websocket.TextMessage, []byte(messageToSend)); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Read the message back from the server
	_, messageReceived, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Verify the sent and received message are the same
	if string(messageReceived) != messageToSend {
		t.Errorf("Expected message %q, but got %q", messageToSend, messageReceived)
	}
}

// TestWebsocketHandlerTimeout tests the WebSocket connection timeout
func TestWebsocketHandlerTimeout(t *testing.T) {
	// Create a test HTTP server with the WebSocket handler
	server := httptest.NewServer(http.HandlerFunc(server.WsHandler))
	defer server.Close()

	// Replace "http" with "ws" in the server URL to form the WebSocket URL
	wsURL := "ws" + server.URL[len("http"):]

	// Dial the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer ws.Close()

	// Set a short read deadline to simulate a timeout
	ws.SetReadDeadline(time.Now().Add(1 * time.Second))

	// Wait for a timeout, expect an error
	_, _, err = ws.ReadMessage()
	if err == nil {
		t.Fatalf("Expected a timeout error, but got none")
	}
}

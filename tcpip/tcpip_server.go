package tcpip

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Define the Client structure
type Client struct {
	Conn     net.Conn
	Username string
}

// Global variable to store all connected clients
var clients = make(map[string]*Client)
var mutex = &sync.RWMutex{}

// Handle client connection
func handleClient(client *Client) {
	defer client.Conn.Close()

	reader := bufio.NewReader(client.Conn)
	for {
		// Read messages from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			removeClient(client)
			return
		}

		message = strings.TrimSpace(message)
		if strings.HasPrefix(message, "LOGIN:") {
			handleLogin(client, message)
		} else if strings.HasPrefix(message, "MSG_ALL:") {
			broadcastMessage(client, message[len("MSG_ALL:"):])
		} else if strings.HasPrefix(message, "MSG_USER:") {
			handlePrivateMessage(client, message)
		} else {
			client.Conn.Write([]byte("Unknown command\n"))
		}
	}
}

// Handle client login
func handleLogin(client *Client, message string) {
	username := strings.TrimSpace(message[len("LOGIN:"):])
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := clients[username]; exists {
		client.Conn.Write([]byte("Username already taken\n"))
	} else {
		client.Username = username
		clients[username] = client
		client.Conn.Write([]byte("Login successful\n"))
		fmt.Printf("User %s logged in\n", username)
	}
}

// Broadcast message to all clients
func broadcastMessage(sender *Client, message string) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, client := range clients {
		if client != sender {
			client.Conn.Write([]byte(sender.Username + ": " + message + "\n"))
		}
	}
}

// Handle private messaging
func handlePrivateMessage(sender *Client, message string) {
	// Parse the message, format is MSG_USER:username:message
	parts := strings.SplitN(message[len("MSG_USER:"):], ":", 2)
	if len(parts) < 2 {
		sender.Conn.Write([]byte("Invalid private message format\n"))
		return
	}

	recipientUsername := strings.TrimSpace(parts[0])
	messageToSend := strings.TrimSpace(parts[1])

	mutex.RLock()
	defer mutex.RUnlock()

	recipient, exists := clients[recipientUsername]
	if !exists {
		sender.Conn.Write([]byte("User not found\n"))
	} else {
		recipient.Conn.Write([]byte("(Private) " + sender.Username + ": " + messageToSend + "\n"))
	}
}

// Remove client
func removeClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	if client.Username != "" {
		delete(clients, client.Username)
		fmt.Printf("User %s disconnected\n", client.Username)
	}
}

// Monitor system resource usage periodically
func monitorResources() {
	ticker := time.NewTicker(time.Second * 10) // Report every 10 seconds
	defer ticker.Stop()

	for range ticker.C {
		// Get memory statistics
		var m runtime.MemStats
		runtime.ReadMemStats(&m) // Read the memory statistics

		log.Printf("Memory Usage: Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v",
			m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.Sys/1024/1024, m.NumGC) // Log memory usage
		log.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
	}
}

// Start the server
func TcpipServer(enableMonitoring *bool) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080")

	// If the monitoring flag is set, start monitoring resources
	if *enableMonitoring {
		log.Println("Resource monitoring enabled")
		go monitorResources()
	} else {
		log.Println("Resource monitoring disabled")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Create a client object and start a goroutine to handle the connection
		client := &Client{Conn: conn}
		go handleClient(client)
	}
}

package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket URL of the server
var serverURL = "ws://localhost:8080/ws"

// TestClient represents a test client
type TestClient struct {
	id   int
	conn *websocket.Conn
	mu   sync.Mutex // Add a mutex to protect write operations
}

// Create a new test client
func newTestClient(id int) *TestClient {
	return &TestClient{id: id}
}

// Connect the client to the WebSocket server
func (c *TestClient) connect() error {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("client %d failed to connect: %v", c.id, err)
	}
	c.conn = conn
	return nil
}

// Send heartbeat messages to keep the connection alive
func (c *TestClient) sendHeartbeat(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second) // Set the heartbeat interval
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock() // Lock before writing
			err := c.conn.WriteMessage(websocket.PingMessage, []byte("ping"))
			c.mu.Unlock() // Unlock after writing
			if err != nil {
				log.Printf("Client %d heartbeat error: %v", c.id, err)
				return
			}
		}
	}
}

// Randomly send messages to the server
func (c *TestClient) sendMessages(wg *sync.WaitGroup) {
	defer wg.Done()
	defer c.conn.Close() // Ensure the connection is closed when done

	// Start a goroutine for heartbeat messages
	var heartbeatWg sync.WaitGroup
	heartbeatWg.Add(1)
	go c.sendHeartbeat(&heartbeatWg)

	// Set a random number of messages to send, between 5 and 20
	numMessages := rand.Intn(16) + 5
	for i := 0; i < numMessages; i++ {
		// Random delay between 100ms and 1000ms
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(900)+100))

		// Generate a random message
		message := fmt.Sprintf("Client %d message %d", c.id, i+1)

		// Lock before writing to avoid concurrent writes
		c.mu.Lock()
		err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
		c.mu.Unlock() // Unlock after writing
		if err != nil {
			log.Printf("Client %d write error: %v", c.id, err)
			return
		}
		log.Printf("Client %d sent: %s", c.id, message)
	}

	// Wait for a random period before disconnecting, between 1 and 3 seconds
	time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
	log.Printf("Client %d finished sending messages", c.id)

	// Close the connection after sending messages
	c.mu.Lock()
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.mu.Unlock() // Unlock after writing
	if err != nil {
		log.Printf("Client %d close error: %v", c.id, err)
	}

	// Wait for the heartbeat goroutine to finish
	heartbeatWg.Wait()
}

func WebsocketClients() {
	// Number of clients to simulate
	numClients := 1000 // You can change this to a higher number for more load

	var wg sync.WaitGroup
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		client := newTestClient(i)

		// Connect the client
		err := client.connect()
		if err != nil {
			log.Printf("Client %d connection error: %v", i, err)
			wg.Done()
			continue
		}

		// Start sending messages in a goroutine
		go client.sendMessages(&wg)

		// Randomly delay the start of the next client to avoid all clients connecting simultaneously
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	}

	// Wait for all clients to finish
	wg.Wait()

	log.Println("All clients finished")
}

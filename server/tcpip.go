package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close() // Close the connection when the function finishes
	fmt.Println("Client connected:", conn.RemoteAddr())

	// Create a buffered reader to read input from the client
	reader := bufio.NewReader(conn)
	for {
		// Read data sent by the client until a newline character
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		// Trim the newline and any extra spaces
		message = strings.TrimSpace(message)
		fmt.Printf("Received message from %s: %s\n", conn.RemoteAddr(), message)

		// Process the message sent by the client; custom logic can go here
		response := handleMessage(message)

		// Send the response back to the client
		conn.Write([]byte(response + "\n"))
	}
}

func handleMessage(message string) string {
	// Here you can process different messages and generate responses
	if message == "ping" {
		return "pong"
	}
	return "You said: " + message
}

func TCPIPServer() {
	// Listen on a specific TCP address and port
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	// Infinite loop to accept new client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Start a new goroutine to handle each client connection
		go handleConnection(conn)
	}
}

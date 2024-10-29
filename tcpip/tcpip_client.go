package tcpip

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Connect to the server
func connectToServer(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Read messages from the server
func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from server")
			return
		}
		fmt.Print("Server: " + message)
	}
}

// Client main program
func TcpipClient() {
	// Connect to the server
	serverAddress := "localhost:8080"
	conn, err := connectToServer(serverAddress)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	// Start a goroutine to read messages from the server
	go readFromServer(conn)

	// Client login
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	_, err = conn.Write([]byte("LOGIN:" + username + "\n"))
	if err != nil {
		fmt.Println("Error during login:", err)
		return
	}

	// Enter the main loop to handle user input
	for {
		fmt.Print("Enter command (MSG_ALL, MSG_USER, EXIT): ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		// Handle exit command
		if command == "EXIT" {
			fmt.Println("Exiting...")
			break
		}

		// Handle broadcast message command
		if strings.HasPrefix(command, "MSG_ALL:") {
			_, err = conn.Write([]byte(command + "\n"))
			if err != nil {
				fmt.Println("Error sending broadcast message:", err)
			}
		} else if strings.HasPrefix(command, "MSG_USER:") {
			// Handle private message command
			_, err = conn.Write([]byte(command + "\n"))
			if err != nil {
				fmt.Println("Error sending private message:", err)
			}
		} else {
			fmt.Println("Unknown command. Use MSG_ALL, MSG_USER, or EXIT.")
		}
	}
}

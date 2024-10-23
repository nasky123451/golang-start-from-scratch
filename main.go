package main

import (
	"flag"
	"fmt"

	chat "example.com/m/chat"
	gpm "example.com/m/goroutine"
	prometheus "example.com/m/prometheus"
	rdb "example.com/m/redis"
	server "example.com/m/server"
	tracing "example.com/m/tracing"
	ws "example.com/m/websocket"
)

func main() {
	// Define flags
	flags := map[string]*bool{
		"websocketServer":          flag.Bool("websocketServer", false, "Enable resource websocket server disable monitor"),
		"websocketServerMonitor":   flag.Bool("websocketServerMonitor", false, "Enable resource websocket server enable monitor"),
		"websocketClients":         flag.Bool("websocketClients", false, "Enable resource websocket clients"),
		"websocketClient":          flag.Bool("websocketClient", false, "Enable resource websocket client"),
		"monitor":                  flag.Bool("monitor", false, "Enable resource monitoring"),
		"goroutine":                flag.Bool("goroutine", false, "Enable goroutine base"),
		"goroutineMutex":           flag.Bool("goroutineMutex", false, "Enable goroutine mutex"),
		"goroutineChannel":         flag.Bool("goroutineChannel", false, "Enable goroutine channel"),
		"tracingJeager":            flag.Bool("tracingJeager", false, "Enable tracing jeager"),
		"tracingZipkin":            flag.Bool("tracingZipkin", false, "Enable tracing zipkin"),
		"prometheus":               flag.Bool("prometheus", false, "Enable prometheus base"),
		"prometheusApiApplication": flag.Bool("prometheusApiApplication", false, "Enable prometheus api application"),
		"redisbase":                flag.Bool("redisbase", false, "Enable redis base"),
		"redisTransferMoney":       flag.Bool("redisTransferMoney", false, "Enable redis transfer money"),
		"wsServer":                 flag.Bool("wsServer", false, "Enable websocket server"),
		"httpServer":               flag.Bool("httpServer", false, "Enable http server"),
		"chatServer":               flag.Bool("chatServer", false, "Enable chat server"),
		"help":                     flag.Bool("help", false, "Display help information"),
	}

	// Parse command line flags
	flag.Parse()

	// Check if help information needs to be displayed
	if *flags["help"] {
		displayHelp()
		return
	}

	// Check enabled flags
	enabledCount := 0
	for _, enabled := range flags {
		if *enabled {
			enabledCount++
		}
	}

	// Check if more than one flag is enabled
	if enabledCount > 1 {
		fmt.Println("Error: Only one option can be enabled at a time. Please refer to -help for more information.")
		return
	}

	// Start corresponding functionality based on enabled flags
	switch {
	case *flags["websocketServer"]:
		isSecure := false
		ws.WebsocketServer(&isSecure)
	case *flags["websocketServerMonitor"]:
		isSecure := true
		ws.WebsocketServer(&isSecure)
	case *flags["websocketClients"]:
		ws.WebsocketClients()
	case *flags["websocketClient"]:
		ws.WebsocketClient()
	case *flags["goroutine"]:
		gpm.GoroutineBase()
	case *flags["goroutineMutex"]:
		gpm.GoroutineMutex()
	case *flags["goroutineChannel"]:
		gpm.GoroutineChannel()
	case *flags["tracingJeager"]:
		tracing.TracingJeager()
	case *flags["tracingZipkin"]:
		tracing.TracingZipkin()
	case *flags["prometheus"]:
		prometheus.PrometheusBase()
	case *flags["prometheusApiApplication"]:
		prometheus.PrometheusApiApplication()
	case *flags["redisbase"]:
		rdb.RedisBase()
	case *flags["redisTransferMoney"]:
		rdb.RedisTransferMoney()
	case *flags["httpServer"]:
		server.HttpServer()
	case *flags["wsServer"]:
		server.WebsocketServer()
	case *flags["chatServer"]:
		chat.ChatServer()
	default:
		// Display error message if no flags are enabled
		fmt.Println("Error: At least one option must be enabled. Please refer to -help for more information.")
	}
}

// displayHelp prints the help information for available flags
func displayHelp() {
	fmt.Println("Available options:")
	fmt.Println("  -websocketServer   		  This code implements a Go language WebSocket server using the Gorilla WebSocket library, which can handle multiple client connections, message broadcasting, and system resource monitoring.")
	fmt.Println("  -websocketServerMonitor    This code implements a Go language WebSocket server using the Gorilla WebSocket library, which can handle multiple client connections, message broadcasting, and system resource monitoring.")
	fmt.Println("  -websocketClients  		  This program implements a WebSocket client simulation where multiple clients connect to a server, send random messages, and periodically send heartbeat messages to maintain the connection, while ensuring thread-safe operations with mutexes.")
	fmt.Println("  -websocketClient  		  This program implements a WebSocket client that connects to a server, receives messages asynchronously, and sends periodic messages for a limited duration before gracefully closing the connection.")
	fmt.Println("  -goroutine        		  This program simulates multiple customers attempting to purchase a product concurrently, using atomic operations to safely manage the stock level in a thread-safe manner.")
	fmt.Println("  -goroutineMutex   		  This program simulates a bank account with concurrent deposit and withdrawal operations, using a mutex to ensure thread-safe access to the account balance.")
	fmt.Println("  -goroutineChannel 		  This program implements a producer-consumer model in Go, where producers generate prioritized tasks and a consumer processes them concurrently.")
	fmt.Println("  -tracingJeager             This is a Go program that initializes a Jaeger tracer using OpenTelemetry, creating trace spans for multiple operations and recording their execution status and errors.")
	fmt.Println("  -tracingZipkin             This is a Go program that uses OpenTelemetry to initialize the Zipkin tracker and create tracking spans when the operation is executed multiple times to record the execution of the operation and its errors.")
	fmt.Println("  -prometheus                This is a simple Go program that starts an HTTP server, provides a web UI for Prometheus metrics, and handles requests on the root route, logging the request count and duration.")
	fmt.Println("  -prometheusApiApplication  This is a Go program that provides Prometheus monitoring API applications, supports HTTP request processing, resource query, user login, and health check, and has the function of gracefully shutting down services.")
	fmt.Println("  -redisbase  				  This is a Go program used to record user access. It stores access logs through PostgreSQL and uses Redis to cache the user's last access time to improve query efficiency.")
	fmt.Println("  -redisTransferMoney  	  This is a Use Redis distributed locks to securely transfer funds and query and update user balances via PostgreSQL while subscribing to Redis expiration events to manage sessions and user activity.")
	fmt.Println("  -wsServer  	 	  	  	  This is an implementation of a WebSocket server that upgrades HTTP connections and handles messages sent by the client and passes received messages back to the client.")
	fmt.Println("  -httpServer  	 	  	  This is an implements a simple HTTP server that handles different request methods and prints the request details to the console.")
	fmt.Println("  -chatServer  	 	  	  This is a chat server implemented in Go and Gin, supporting user registration, login, real-time chat and WebSocket connections, and integrating Redis and PostgreSQL management data.")
	fmt.Println("  -help              		  Display help information")
}

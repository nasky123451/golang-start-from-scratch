package main

import (
	"flag"
	"fmt"

	gpm "example.com/m/goroutine"
	tracing "example.com/m/tracing"
	ws "example.com/m/websocket"
)

func main() {
	// Define flags
	flags := map[string]*bool{
		"websocketServer":  flag.Bool("websocketServer", false, "Enable resource websocket server"),
		"websocketClients": flag.Bool("websocketClients", false, "Enable resource websocket clients"),
		"websocketClient":  flag.Bool("websocketClient", false, "Enable resource websocket client"),
		"monitor":          flag.Bool("monitor", false, "Enable resource monitoring"),
		"goroutine":        flag.Bool("goroutine", false, "Enable goroutine"),
		"goroutineMutex":   flag.Bool("goroutineMutex", false, "Enable goroutine mutex"),
		"goroutineChannel": flag.Bool("goroutineChannel", false, "Enable goroutine channel"),
		"tracingJeager":    flag.Bool("tracingJeager", false, "Enable tracing jeager"),
		"tracingZipkin":    flag.Bool("tracingZipkin", false, "Enable tracing zipkin"),
		"help":             flag.Bool("help", false, "Display help information"),
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
		ws.WebsocketServer(flags["monitor"])
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
	default:
		// Display error message if no flags are enabled
		fmt.Println("Error: At least one option must be enabled. Please refer to -help for more information.")
	}
}

// displayHelp prints the help information for available flags
func displayHelp() {
	fmt.Println("Available options:")
	fmt.Println("  -websocketServer   Enable websocket server to use 8080 port")
	fmt.Println("  -websocketClients  Used to brute force test websocket server")
	fmt.Println("  -websocketClient   Test connection to websocket server and send messages")
	fmt.Println("  -monitor           Enable websocket server monitoring")
	fmt.Println("  -goroutine         Enable goroutine")
	fmt.Println("  -goroutineMutex    Enable goroutine mutex")
	fmt.Println("  -goroutineChannel  Enable goroutine channel")
	fmt.Println("  -tracingJeager     Enable tracing jeager to use 16686 port")
	fmt.Println("  -tracingZipkin     Enable tracing zipkin to use 9412 port")
	fmt.Println("  -help              Display help information")
}

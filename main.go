package main

import (
	"flag"
	"fmt"

	gpm "example.com/m/goroutine"
	prometheus "example.com/m/prometheus"
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
	default:
		// Display error message if no flags are enabled
		fmt.Println("Error: At least one option must be enabled. Please refer to -help for more information.")
	}
}

// displayHelp prints the help information for available flags
func displayHelp() {
	fmt.Println("Available options:")
	fmt.Println("  -websocketServer   		  Enable websocket server to use 8080 port disable monitor")
	fmt.Println("  -websocketServerMonitor    Enable websocket server to use 8080 port enable monitor")
	fmt.Println("  -websocketClients  		  Used to brute force test websocket server")
	fmt.Println("  -websocketClient  		  Test connection to websocket server and send messages")
	fmt.Println("  -monitor          		  Enable websocket server monitoring")
	fmt.Println("  -goroutine        		  Enable goroutine base")
	fmt.Println("  -goroutineMutex   		  Enable goroutine mutex")
	fmt.Println("  -goroutineChannel 		  Enable goroutine channel")
	fmt.Println("  -tracingJeager             Enable tracing jeager to use 16686 port")
	fmt.Println("  -tracingZipkin             Enable tracing zipkin to use 9412 port")
	fmt.Println("  -prometheus                Enable prometheus base to use 8080 & 9090 port")
	fmt.Println("  -prometheusApiApplication  Enable prometheus api application to use 8080 & 9090 port")
	fmt.Println("  -help              		  Display help information")
}

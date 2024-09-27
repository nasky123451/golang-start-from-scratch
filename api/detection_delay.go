package api

import (
	"fmt" // Standard log package
	"os"
	"time"

	kitlog "github.com/go-kit/kit/log" // Aliasing to avoid conflict
	"github.com/go-kit/kit/log/level"
)

// doOperation simulates a function whose latency you want to measure.
func doOperation() error {
	// Simulate some work with a sleep.
	time.Sleep(10 * time.Millisecond) // Simulating a delay
	return nil
}

// ExampleLatencySimplest measures the latency of doOperation multiple times.
func ExampleLatencySimplest() {
	xTimes := 5 // Define how many times to run the operation
	for i := 0; i < xTimes; i++ {
		start := time.Now()
		err := doOperation()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		elapsed := time.Since(start)

		fmt.Printf("Execution %d: %v ns\n", i+1, elapsed.Nanoseconds())
	}
}

// ExampleLatencyAggregated measures the average latency of doOperation over multiple executions.
func ExampleLatencyAggregated() {
	var count int64
	var sum int64
	xTimes := 5 // Define how many times to run the operation
	for i := 0; i < xTimes; i++ {
		start := time.Now()
		err := doOperation()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		elapsed := time.Since(start)

		sum += elapsed.Nanoseconds()
		count++
	}

	if count > 0 {
		fmt.Printf("Average: %v ns/op\n", sum/count)
	} else {
		fmt.Println("No operations were successfully executed.")
	}
}

// ExampleLatencyLog logs the results of operations using a structured logger.
func ExampleLatencyLog() {
	// Use an alias for the kit log package to avoid conflict
	kitLogger := kitlog.NewLogfmtLogger(os.Stderr)

	xTimes := 5 // Define how many times to run the operation
	for i := 0; i < xTimes; i++ {
		now := time.Now()
		err := doOperation()
		elapsed := time.Since(now)

		level.Info(kitLogger).Log(
			"msg", "finished operation",
			"result", err,
			"elapsed", elapsed.String(),
		)
	}
}

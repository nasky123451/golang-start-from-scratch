package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a task with an ID and a priority
type Task struct {
	ID       int
	taskID   int
	Priority int // Higher number means higher priority
}

// Producer function generates tasks and sends them to the channel
func producer(id int, tasks chan<- Task, stop <-chan struct{}) {
	taskID := 1 // Initialize task ID
	for {
		select {
		case <-stop:
			fmt.Printf("Producer %d stopping\n", id)
			return
		default:
			// Generate a new task with a random priority
			task := Task{ID: id, taskID: taskID, Priority: rand.Intn(5)} // Random priority between 0 and 4
			tasks <- task                                                // Send task to the channel
			currentTime := time.Now().Format("15:04:05")                 // Get current time
			fmt.Printf("Producer %d produced task: %v at %s\n", id, task, currentTime)
			taskID++ // Increment the task ID for the next task

			// Sleep for a random duration to simulate sporadic task generation
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000))) // Random sleep between 0 and 1000 ms
		}
	}
}

// Consumer function processes tasks received from the channel
func consumer(id int, tasks <-chan Task, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()           // Indicate that this goroutine is done when it returns
	for task := range tasks { // Read from the channel until it's closed
		go processTask(task, id, results)
	}
}

// processTask function simulates processing of a single task
func processTask(task Task, consumerID int, results chan<- string) {
	// Simulate processing based on priority
	time.Sleep(time.Millisecond * time.Duration(2000-task.Priority*500)) // Higher priority takes less time
	result := fmt.Sprintf("Consumer %d processed task %d and taskID %d with priority %d", consumerID, task.ID, task.taskID, task.Priority)
	currentTime := time.Now().Format("15:04:05")           // Get current time
	results <- fmt.Sprintf("%s - %s", currentTime, result) // Include time in results
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator
	tasks := make(chan Task, 10)     // Buffered channel for tasks
	results := make(chan string)     // Buffered channel for results

	var wg sync.WaitGroup

	// Create a stop channel for producers
	stop := make(chan struct{})

	// Start producer goroutines
	numProducers := 3
	for i := 1; i <= numProducers; i++ {
		go producer(i, tasks, stop)
	}

	// Start consumer goroutines
	numConsumers := 1
	for i := 1; i <= numConsumers; i++ {
		wg.Add(1) // Add to wait group for each consumer
		go consumer(i, tasks, results, &wg)
	}

	// Stop task production after 3 seconds
	go func() {
		time.Sleep(3 * time.Second) // Wait for 3 seconds
		close(stop)                 // Close the stop channel
	}()

	// Close the tasks channel when all tasks are processed
	go func() {
		wg.Wait()    // Wait for all consumers to finish
		close(tasks) // Close tasks channel after all consumers finish
	}()

	// Close the results channel when all tasks are processed
	go func() {
		wg.Wait()      // Wait for all consumers to finish
		close(results) // Close results channel after all tasks are processed
	}()

	// Print results as they come in
	for result := range results {
		fmt.Println(result) // Print results
	}

	fmt.Println("All tasks processed.")
}

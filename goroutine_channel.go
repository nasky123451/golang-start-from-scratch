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
	Priority int // Higher number means higher priority
}

// Producer function generates tasks and sends them to the channel
func producer(id int, tasks chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ { // Each producer generates 5 tasks
		task := Task{ID: i + 1 + (id-1)*5, Priority: rand.Intn(5)} // Random priority between 0 and 4
		tasks <- task                                              // Send task to the channel
		currentTime := time.Now().Format("15:04:05")               // Get current time
		fmt.Printf("Producer %d produced task: %v at %s\n", id, task, currentTime)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) // Simulate work
	}
}

// Consumer function processes tasks received from the channel
func consumer(id int, tasks <-chan Task, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks { // Read from the channel until it's closed
		// For each task, start a new goroutine to process it
		wg.Add(1)                             // Increment WaitGroup counter
		go processTask(task, id, results, wg) // Process the task in a new goroutine
	}
}

// processTask function simulates processing of a single task
func processTask(task Task, consumerID int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement counter when function completes
	// Simulate processing based on priority
	time.Sleep(time.Millisecond * time.Duration(2000-task.Priority*500)) // Higher priority takes less time
	result := fmt.Sprintf("Consumer %d processed task %d with priority %d", consumerID, task.ID, task.Priority)
	currentTime := time.Now().Format("15:04:05")           // Get current time
	results <- fmt.Sprintf("%s - %s", currentTime, result) // Include time in results
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed random number generator
	tasks := make(chan Task, 15)     // Buffered channel for tasks
	results := make(chan string, 15) // Buffered channel for results

	var wg sync.WaitGroup

	// Start producer goroutines
	numProducers := 3
	for i := 1; i <= numProducers; i++ {
		wg.Add(1)
		go producer(i, tasks, &wg)
	}

	// Wait for all producers to finish
	wg.Wait()
	close(tasks) // Close the tasks channel after producers are done

	// Start consumer goroutines
	numConsumers := 1
	for i := 1; i <= numConsumers; i++ {
		wg.Add(1)
		go consumer(i, tasks, results, &wg)
	}

	// Wait for consumers to finish processing tasks
	go func() {
		wg.Wait()
		close(results) // Close the results channel after consumers are done
	}()

	// Display results
	for result := range results {
		fmt.Println(result) // Print results as they come in
	}

	fmt.Println("All tasks processed.")
}

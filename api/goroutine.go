package api

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Seeding the random number generator once globally
func init() {
	rand.Seed(time.Now().UnixNano())
}

// randInt64 generates a random int64 number
func randInt64() int64 {
	return rand.Int63n(100) // random number between 0 and 99
}

func sharingWithAtomic() (sum int64) {
	var wg sync.WaitGroup

	concurrentFn := func() {
		defer wg.Done()
		atomic.AddInt64(&sum, randInt64())
	}

	wg.Add(3) // Adding three goroutines to the WaitGroup
	go concurrentFn()
	go concurrentFn()
	go concurrentFn()

	wg.Wait() // Waiting for all goroutines to finish
	return sum
}

func sharingWithMutex() (sum int64) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	concurrentFn := func() {
		mu.Lock()
		sum += randInt64()
		mu.Unlock()
		wg.Done()
	}

	wg.Add(3) // Adding three goroutines to the WaitGroup
	go concurrentFn()
	go concurrentFn()
	go concurrentFn()

	wg.Wait() // Waiting for all goroutines to finish
	return sum
}

func sharingWithChannel() (sum int64) {
	result := make(chan int64)

	// This function will send random values to the channel
	concurrentFn := func() {
		result <- randInt64()
	}

	go concurrentFn()
	go concurrentFn()
	go concurrentFn()

	// Reading values from the channel
	for i := 0; i < 3; i++ {
		sum += <-result
	}
	close(result)
	return sum
}

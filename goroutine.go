package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Product represents a product in the inventory
type Product struct {
	ID    int64
	Stock int64
}

func main() {
	product := &Product{ID: 1, Stock: 10}
	var wg sync.WaitGroup

	// Simulate customers trying to purchase the product
	customerPurchase := func(customerID int) {
		defer wg.Done()
		time.Sleep(time.Millisecond * time.Duration(100)) // Simulate some delay
		for {
			if atomic.LoadInt64(&product.Stock) > 0 {
				atomic.AddInt64(&product.Stock, -1)
				fmt.Printf("Customer %d purchased product %d. Remaining stock: %d\n", customerID, product.ID, product.Stock)
				break
			} else {
				fmt.Printf("Customer %d could not purchase product %d. Out of stock.\n", customerID, product.ID)
				break
			}
		}
	}

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go customerPurchase(i)
	}

	wg.Wait()
	fmt.Printf("Final stock of product %d: %d\n", product.ID, product.Stock)
}

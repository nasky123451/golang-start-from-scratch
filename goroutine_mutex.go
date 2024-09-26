package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// BankAccount represents a bank account with a balance
type BankAccount struct {
	mu      sync.Mutex // Mutex to protect access to the balance
	balance int64      // Current balance
}

// Deposit adds amount to the account balance
func (account *BankAccount) Deposit(amount int64) {
	account.mu.Lock()         // Lock the mutex
	defer account.mu.Unlock() // Ensure to unlock the mutex
	account.balance += amount
	fmt.Printf("Deposited %d. New balance: %d\n", amount, account.balance)
}

// Withdraw deducts amount from the account balance
func (account *BankAccount) Withdraw(amount int64) bool {
	account.mu.Lock()         // Lock the mutex
	defer account.mu.Unlock() // Ensure to unlock the mutex
	if account.balance >= amount {
		account.balance -= amount
		fmt.Printf("Withdrew %d. New balance: %d\n", amount, account.balance)
		return true
	}
	fmt.Printf("Withdrawal of %d failed. Insufficient funds: %d\n", amount, account.balance)
	return false
}

// DisplayBalance shows the current balance
func (account *BankAccount) DisplayBalance() int64 {
	account.mu.Lock()         // Lock the mutex
	defer account.mu.Unlock() // Ensure to unlock the mutex
	return account.balance
}

func main() {
	rand.Seed(time.Now().UnixNano())        // Set seed for random number generator
	account := &BankAccount{balance: 10000} // Initial balance of 10000
	var wg sync.WaitGroup

	// Define a larger number of deposit and withdrawal operations
	const numOperations = 50
	operations := make([]struct {
		opType string
		amount int64
	}, numOperations)

	// Randomly generate operations
	for i := 0; i < numOperations; i++ {
		if rand.Intn(2) == 0 {
			operations[i] = struct {
				opType string
				amount int64
			}{"deposit", rand.Int63n(1000) + 1} // Random deposit between 1 and 1000
		} else {
			operations[i] = struct {
				opType string
				amount int64
			}{"withdraw", rand.Int63n(1000) + 1} // Random withdrawal between 1 and 1000
		}
	}

	for _, operation := range operations {
		wg.Add(1)
		go func(opType string, amount int64) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) // Random delay
			if opType == "deposit" {
				account.Deposit(amount)
			} else if opType == "withdraw" {
				account.Withdraw(amount)
			}
		}(operation.opType, operation.amount)
	}

	wg.Wait()                                                   // Wait for all goroutines to finish
	fmt.Printf("Final balance: %d\n", account.DisplayBalance()) // Display final balance
}

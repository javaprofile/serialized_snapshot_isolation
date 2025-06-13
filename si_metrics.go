package main

import (
	"fmt"
	"sync"
	"time"
)

type Transaction struct {
	id          int
	waitingFor  *Transaction
	lockedItems map[string]bool
}

var (
	mutex          sync.Mutex
	transactionDB  map[int]*Transaction
	deadlockCount  int
	resourceLocks  map[string]*Transaction
)

func init() {
	transactionDB = make(map[int]*Transaction)
	resourceLocks = make(map[string]*Transaction)
	deadlockCount = 0
}

func detectDeadlock(transaction *Transaction) bool {
	
	visited := make(map[int]bool)
	return checkCycle(transaction, visited)
}

func checkCycle(transaction *Transaction, visited map[int]bool) bool {
	if visited[transaction.id] {
		
		return true
	}
	visited[transaction.id] = true

	if transaction.waitingFor != nil {
		return checkCycle(transaction.waitingFor, visited)
	}

	return false
}

func tryLockTransaction(transaction *Transaction, resource string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if owner, exists := resourceLocks[resource]; exists {
		transaction.waitingFor = owner
		return false
	}
	resourceLocks[resource] = transaction
	transaction.lockedItems[resource] = true
	return true
}

func releaseLocks(transaction *Transaction) {
	mutex.Lock()
	defer mutex.Unlock()
	for resource := range transaction.lockedItems {
		delete(resourceLocks, resource)
	}
}

func simulateTransaction(id int, resource string) {
	transaction := &Transaction{
		id:          id,
		lockedItems: make(map[string]bool),
	}

	transactionDB[id] = transaction

	if !tryLockTransaction(transaction, resource) {
		if detectDeadlock(transaction) {
			fmt.Printf("Deadlock detected for transaction %d!\n", id)
			deadlockCount++
		} else {
			fmt.Printf("Transaction %d is waiting for resource %s\n", id, resource)
		}
	}

	time.Sleep(2 * time.Second)
	releaseLocks(transaction)
}

func main() {
	
	go simulateTransaction(1, "resource1")
	go simulateTransaction(2, "resource1")
	go simulateTransaction(3, "resource2")
	go simulateTransaction(4, "resource2")	
	time.Sleep(5 * time.Second)
	fmt.Printf("Total deadlocks detected: %d\n", deadlockCount)
}

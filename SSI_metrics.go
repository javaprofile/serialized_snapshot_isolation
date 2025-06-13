package main

import (
	"fmt"
	"sync"
	"time"
)

type Transaction struct {
	ID       int
	ReadSet  map[string]bool
	WriteSet map[string]bool
	Status   string
}

type Database struct {
	mu          sync.Mutex
	transactions map[int]*Transaction
	deadlocks    int
}

func (db *Database) startTransaction(id int) *Transaction {
	db.mu.Lock()
	defer db.mu.Unlock()

	txn := &Transaction{
		ID:       id,
		ReadSet:  make(map[string]bool),
		WriteSet: make(map[string]bool),
		Status:   "active",
	}
	db.transactions[id] = txn
	return txn
}

func (db *Database) detectDeadlock(txn1, txn2 *Transaction) bool {
	for key := range txn1.WriteSet {
		if _, exists := txn2.ReadSet[key]; exists {
			return true
		}
	}
	for key := range txn2.WriteSet {
		if _, exists := txn1.ReadSet[key]; exists {
			return true
		}
	}
	return false
}

func (db *Database) commitTransaction(txn *Transaction) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, otherTxn := range db.transactions {
		if otherTxn.ID != txn.ID && otherTxn.Status == "active" {
			if db.detectDeadlock(txn, otherTxn) {
				db.deadlocks++
				fmt.Printf("Deadlock detected between transaction %d and %d\n", txn.ID, otherTxn.ID)
				txn.Status = "aborted"
				return
			}
		}
	}

	txn.Status = "committed"
}

func main() {
	db := &Database{
		transactions: make(map[int]*Transaction),
	}

	txn1 := db.startTransaction(1)
	txn2 := db.startTransaction(2)
	txn3 := db.startTransaction(3)

	txn1.WriteSet["item1"] = true
	txn2.ReadSet["item1"] = true
	txn2.WriteSet["item2"] = true
	txn3.WriteSet["item1"] = true

	go db.commitTransaction(txn1)
	go db.commitTransaction(txn2)
	go db.commitTransaction(txn3)

	time.Sleep(1 * time.Second)

	fmt.Printf("Total deadlocks: %d\n", db.deadlocks)
}

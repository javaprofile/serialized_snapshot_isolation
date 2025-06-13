package main

import (
	"fmt"
	"sync"
)

type Transaction struct {
	ID        int
	Timestamp int
	ReadSet   map[string]int
	WriteSet  map[string]int
}

type Database struct {
	mu      sync.Mutex
	data    map[string]int
	txnLog  []Transaction
}

func NewDatabase() *Database {
	return &Database{data: make(map[string]int)}
}

func (db *Database) StartTransaction(txnID int) *Transaction {
	db.mu.Lock()
	defer db.mu.Unlock()

	txn := Transaction{ID: txnID, Timestamp: len(db.txnLog) + 1, ReadSet: make(map[string]int), WriteSet: make(map[string]int)}
	db.txnLog = append(db.txnLog, txn)
	return &txn
}

func (txn *Transaction) Read(db *Database, key string) int {
	db.mu.Lock()
	defer db.mu.Unlock()

	txn.ReadSet[key] = db.data[key]
	return db.data[key]
}

func (txn *Transaction) Write(db *Database, key string, value int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	txn.WriteSet[key] = value
	db.data[key] = value
}

func (txn *Transaction) Commit(db *Database) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, otherTxn := range db.txnLog {
		if txn.Timestamp < otherTxn.Timestamp {
			for key := range txn.WriteSet {
				if _, exists := otherTxn.ReadSet[key]; exists {
					return false
				}
			}
		}
	}

	for key, value := range txn.WriteSet {
		db.data[key] = value
	}
	return true
}

func main() {
	db := NewDatabase()

	txn1 := db.StartTransaction(1)
	txn1.Read(db, "A")
	txn1.Write(db, "A", 5)
	txn1.Commit(db)

	txn2 := db.StartTransaction(2)
	txn2.Read(db, "A")
	txn2.Write(db, "B", 10)
	txn2.Commit(db)

	fmt.Println(db.data)
}

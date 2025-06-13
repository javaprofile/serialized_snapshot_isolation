package main
import (
	"fmt"
	"sync"
)
type Transaction struct {
	id       int
	snapshot map[string]int
	committed bool
}
type Database struct {
	data map[string]int
	mu   sync.Mutex
}
func NewDatabase() *Database {
	return &Database{
		data: make(map[string]int),
	}
}
func (db *Database) BeginTransaction(id int) *Transaction {
	db.mu.Lock()
	defer db.mu.Unlock()
	snapshot := make(map[string]int)
	for key, value := range db.data {
		snapshot[key] = value
	}
	return &Transaction{id: id, snapshot: snapshot, committed: false}
}
func (t *Transaction) Read(db *Database, key string) (int, bool) {
	if value, exists := t.snapshot[key]; exists {
		return value, true
	}
	return 0, false
}
func (t *Transaction) Write(db *Database, key string, value int) bool {
	if currentValue, exists := t.snapshot[key]; exists {
		if currentValue != db.data[key] {
			return false
		}
		db.data[key] = value
		t.snapshot[key] = value
		return true
	}
	return false
}
func (t *Transaction) Commit(db *Database) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	for key, value := range t.snapshot {
		if currentValue, exists := db.data[key]; exists && currentValue != value {
			return false
		}
	}
	for key, value := range t.snapshot {
		db.data[key] = value
	}
	t.committed = true
	return true
}
func main() {
	db := NewDatabase()
	t1 := db.BeginTransaction(1)
	t2 := db.BeginTransaction(2)
	t1.Write(db, "a", 10)
	t2.Write(db, "a", 20)
	if t1.Commit(db) {
		fmt.Println("Transaction 1 committed successfully")
	} else {
		fmt.Println("Transaction 1 failed due to conflict")
	}
	if t2.Commit(db) {
		fmt.Println("Transaction 2 committed successfully")
	} else {
		fmt.Println("Transaction 2 failed due to conflict")
	}
	fmt.Println("Database state:", db.data)
}

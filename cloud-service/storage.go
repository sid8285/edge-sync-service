package main

import (
	"fmt"
	"sync"
)

type CloudStore struct {
	mu           sync.RWMutex
	transactions map[string]Transaction
}

func NewCloudStore() *CloudStore {
	return &CloudStore{
		transactions: make(map[string]Transaction), //we use make to create new maps bc Go panics at a nil map
	}
}

//this function belongs to CloudStore, s is the specific store instance
func (s *CloudStore) SaveSyncedTransaction(t Transaction) (Transaction, bool, error) {
	if t.ID == "" {
		return Transaction{}, false, fmt.Errorf("transaction ID cannot be empty")
	}
	s.mu.Lock() //only one writer at a time
	defer s.mu.Unlock() //unlock regarless of how i return

	existing, alreadyExists := s.transactions[t.ID] //Go's two-value map index. existing maps to whatever was stored under the key, already exists returns a bool - good to know
	if alreadyExists {
		return existing, true, nil
	}

	s.transactions[t.ID] = t
	return t, false, nil
}

func (s *CloudStore) ListTransactions() []Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	transactions := make([]Transaction, 0, len(s.transactions))

	for _, tx := range s.transactions {
		transactions = append(transactions, tx)
	}

	return transactions
}

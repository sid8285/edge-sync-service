package main

import (
	"fmt"
	"sync"
)

type EdgeStore struct {
	mu           sync.RWMutex
	transactions map[string]Transaction
}

func newEdgeStore() *EdgeStore {
	return &EdgeStore{
		transactions: make(map[string]Transaction),
	}
}

// SaveLocal stores t under t.ID (idempotent: same ID returns the existing row).
func (s *EdgeStore) SaveLocal(t Transaction) (Transaction, bool, error) {
	if t.ID == "" {
		return Transaction{}, false, fmt.Errorf("transaction id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.transactions[t.ID]; ok {
		return existing, true, nil
	}

	s.transactions[t.ID] = t
	return t, false, nil
}

func (s *EdgeStore) ListAll() []Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Transaction, 0, len(s.transactions))
	for _, tx := range s.transactions {
		out = append(out, tx)
	}
	return out
}

func (s *EdgeStore) ListUnsynced() []Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Transaction, 0, len(s.transactions))
	for _, tx := range s.transactions {
		if !tx.Synced {
			out = append(out, tx)
		}
	}
	return out
}

func (s *EdgeStore) MarkSynced(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" {
		return fmt.Errorf("There is no id provided")
	}

	tx, ok := s.transactions[id]
	if !ok {
		return fmt.Errorf("unknown id: %s", id)
	}

	// Cannot write s.transactions[id].Synced = true: values in maps are not
	// assignable fields. Copy-out, mutate, assign back.
	tx.Synced = true
	s.transactions[id] = tx
	return nil
}

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type server struct {
	store *EdgeStore
}

// CreateTransactionRequest is the JSON body for POST /transactions (amount + items only).
// ID, store_id, created_at, and synced are set in handleCreateTransaction.
type CreateTransactionRequest struct {
	Amount float64  `json:"amount"`
	Items  []string `json:"items"`
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "edge-service",
	})
}

func (s *server) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	txs := s.store.ListAll()
	_ = json.NewEncoder(w).Encode(txs)
}

func (s *server) handleListUnsynced(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	txs := s.store.ListUnsynced()
	_ = json.NewEncoder(w).Encode(txs)
}

func (s *server) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	storeID := getenvDefault("STORE_ID", "store_001")

	tx := Transaction{
		ID:        uuid.NewString(),
		StoreID:   storeID,
		Amount:    req.Amount,
		Items:     req.Items,
		CreatedAt: time.Now().UTC(),
		Synced:    false,
	}

	saved, _, err := s.store.SaveLocal(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(saved)
}

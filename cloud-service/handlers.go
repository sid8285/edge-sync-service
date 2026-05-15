package main

import (
	"encoding/json"
	"net/http"
)

type server struct {
	store *CloudStore
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method isnt allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"status":  "ok",
		"service": "cloud-service",
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func (s *server) listTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this method isnt allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	txs := s.store.ListTransactions()
	_ = json.NewEncoder(w).Encode(txs)
}

func (s *server) handleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		// Malformed JSON is a client mistake, not "wrong HTTP method".
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	saved, dup, err := s.store.SaveSyncedTransaction(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	msg := "stored"
	if dup {
		msg = "already synced"
	}

	resp := struct {
		Message     string      `json:"message"`
		Transaction Transaction `json:"transaction"`
	}{
		Message:     msg,
		Transaction: saved,
	}

	_ = json.NewEncoder(w).Encode(resp)

}

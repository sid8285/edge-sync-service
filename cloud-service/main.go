package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	addr := ":" + port

	store := NewCloudStore()
	srv := &server{store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", srv.handleHealth)
	mux.HandleFunc("GET /transactions", srv.listTransactions)
	mux.HandleFunc("POST /transactions/sync", srv.handleSync)

	log.Printf("cloud-service listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

package main

import "time"


//Transaction here represents a single retail sale recorded at a store and to the cloud.
//Both services will use a shape like this so that transactions can travel end to end without needing to be converted.
type Transaction struct {
	ID        string    `json:"id"`
	StoreID   string    `json:"store_id"`
	Amount    float64   `json:"amount"`
	Items     []string  `json:"items"`
	CreatedAt time.Time `json:"created_at"`
	Synced    bool      `json:"synced"`
}

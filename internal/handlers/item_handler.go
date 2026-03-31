package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"findit/internal/models"
)

var (
	items []models.Item
	mu    sync.Mutex
)

// Create Item
func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	items = append(items, item)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// Get All Items
func GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
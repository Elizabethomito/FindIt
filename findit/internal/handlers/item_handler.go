package handlers

import (
	"encoding/json"
	"net/http"

	"findit/internal/db"
	"findit/internal/models"
	"findit/pkg/utils"
)

// CreateItem handles POST /items
func CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if item.ID == "" || item.Name == "" || item.Type == "" || item.UserID == "" {
		utils.JSONError(w, http.StatusBadRequest, "id, name, type, and user_id are required")
		return
	}
	if item.Type != "lost" && item.Type != "found" {
		utils.JSONError(w, http.StatusBadRequest, "type must be 'lost' or 'found'")
		return
	}

	_, err := db.DB.Exec(
		"INSERT INTO items (id, user_id, type, name, description, location, date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		item.ID, item.UserID, item.Type, item.Name, item.Description, item.Location, item.Date,
	)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, item)
}

// GetItems handles GET /items
func GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rows, err := db.DB.Query("SELECT id, user_id, type, name, description, location, date FROM items")
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "failed to fetch items")
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.UserID, &item.Type, &item.Name, &item.Description, &item.Location, &item.Date); err != nil {
			continue
		}
		items = append(items, item)
	}

	if items == nil {
		items = []models.Item{} // return empty array, not null
	}
	utils.JSONResponse(w, http.StatusOK, items)
}

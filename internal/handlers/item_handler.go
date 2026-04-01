package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"findit/internal/db"
	"findit/internal/models"
	"findit/pkg/utils"
)

// generateUUID generates a random UUID
func generateUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	// Set version (4) and variant (2)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return hex.EncodeToString(uuid)
}

// CreateItem handles POST /items
func CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Printf("Error decoding request body: %v", err)
		utils.JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Generate auto ID if not provided
	if item.ID == "" {
		item.ID = generateUUID()
	}

	log.Printf("Received item data: ID=%s, UserID=%s, Type=%s, Name=%s", item.ID, item.UserID, item.Type, item.Name)

	if item.Name == "" || item.Type == "" || item.UserID == "" {
		utils.JSONError(w, http.StatusBadRequest, "name, type, and user_id are required")
		return
	}
	if item.Type != "lost" && item.Type != "found" {
		utils.JSONError(w, http.StatusBadRequest, "type must be 'lost' or 'found'")
		return
	}

	// Verify user exists
	var userExists bool
	checkErr := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", item.UserID).Scan(&userExists)
	if checkErr != nil {
		log.Printf("Error checking user existence: %v", checkErr)
		utils.JSONError(w, http.StatusInternalServerError, "failed to verify user")
		return
	}
	if !userExists {
		log.Printf("User not found: %s", item.UserID)
		utils.JSONError(w, http.StatusBadRequest, "user not found")
		return
	}

	_, err := db.DB.Exec(
		"INSERT INTO items (id, user_id, type, name, description, location, date) VALUES (?, ?, ?, ?, ?, ?, ?)",
		item.ID, item.UserID, item.Type, item.Name, item.Description, item.Location, item.Date,
	)
	if err != nil {
		log.Printf("Error creating item: %v", err)
		log.Printf("Item data: ID=%s, UserID=%s, Type=%s, Name=%s", item.ID, item.UserID, item.Type, item.Name)
		utils.JSONError(w, http.StatusInternalServerError, "failed to create item: "+err.Error())
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

// DeleteItem handles DELETE /items/{id}
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract item ID from URL path
	itemID := r.URL.Path[len("/items/"):]
	if itemID == "" {
		utils.JSONError(w, http.StatusBadRequest, "item id is required")
		return
	}

	log.Printf("Deleting item: %s", itemID)

	// Check if item exists
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM items WHERE id = ?)", itemID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking item existence: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "failed to check item")
		return
	}
	if !exists {
		utils.JSONError(w, http.StatusNotFound, "item not found")
		return
	}

	// Delete item
	_, err = db.DB.Exec("DELETE FROM items WHERE id = ?", itemID)
	if err != nil {
		log.Printf("Error deleting item: %v", err)
		utils.JSONError(w, http.StatusInternalServerError, "failed to delete item")
		return
	}

	log.Printf("Item deleted successfully: %s", itemID)
	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "item deleted successfully"})
}

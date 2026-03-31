package routes

import (
	"net/http"

	"findit/internal/handlers"
)

// Register sets up all API routes and returns the mux.
func Register() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("/signup", handlers.Signup)
	mux.HandleFunc("/login", handlers.Login)

	// Items — single handler dispatches by method internally
	mux.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateItem(w, r)
		case http.MethodGet:
			handlers.GetItems(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

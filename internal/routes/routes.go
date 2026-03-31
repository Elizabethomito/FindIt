package routes

import (
	"net/http"
	"findit/internal/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateItem(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetItems(w, r)
		}
	})
}
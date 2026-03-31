package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse writes a JSON response with the given status code.
func JSONResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// JSONError writes a JSON error response.
func JSONError(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}

package main

import (
	"log"
	"net/http"

	"findit/internal/db"
	"findit/internal/routes"
)

func main() {
	db.Init()

	mux := routes.Register()

	log.Println("FindIt server starting on http://localhost:8080/")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

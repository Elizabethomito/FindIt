package main

import (
	"log"
	"net/http"
	"os"

	"findit/internal/db"
	"findit/internal/routes"
)

func main() {
	db.Init()

	mux := routes.Register()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	http.ListenAndServe(":"+port, mux)
}

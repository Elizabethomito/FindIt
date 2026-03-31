package main

import (
	"fmt"
	"net/http"

	"findit/internal/routes"
)

func main() {
	routes.RegisterRoutes()

	fmt.Println("Server running on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
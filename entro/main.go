package main

import (
	"fmt"
	"github.com/ireuven89/entro/service"
	"net/http"
)

func scanHandler(w http.ResponseWriter, r *http.Request) {
	// Accept only POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract query parameters from URL
	query := r.URL.Query()
	repository := query.Get("repository")
	organization := query.Get("organization")
	token := query.Get("token")

	service.StartScan(organization, repository, token)
}

func main() {

	http.HandleFunc("/scan", scanHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Quote struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/quote", func(w http.ResponseWriter, r *http.Request) {
		// ENABLE CORS:
		w.Header().Set("Access-Control-Allow-Origin", "*") // In production, replace "*" with your Vue domain
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "OPTIONS" {
			return
		}

		quote := Quote{
			Message: "Go and Vue are a match made in heaven.",
			Author:  "The Gopher",
		}

		json.NewEncoder(w).Encode(quote)
	})

	log.Printf("Backend starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

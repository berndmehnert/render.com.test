package main

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"log"
	"net/http"
)

type Quote struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

func main() {
	cfg := config.Load()
	port := cfg.Port

	http.HandleFunc("/stream", middleware.CORS(handler.SSE(cfg), cfg.UiUrl))

	log.Printf("Backend starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

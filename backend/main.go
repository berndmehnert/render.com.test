package main

import (
	"backend/config"
	"backend/handler"
	"backend/middleware"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	port := cfg.Port

	http.HandleFunc("/stream", middleware.CORS(handler.SSE(cfg), cfg.UiUrl))

	log.Printf("Backend starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

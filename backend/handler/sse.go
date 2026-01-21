package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"backend/client"
	"backend/config"
	"backend/model"
)

func SSE(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify streaming is supported
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		ollama := client.NewOllama(cfg)

		resp, err := ollama.Chat(ctx, "Hi, please give me a short Go gopher quote.")
		if err != nil {
			log.Printf("Internal server error: %v", err)
			return
		}
		defer resp.Body.Close()

		// NOW set SSE headers and start streaming
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		// Stream response
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				log.Println("Client disconnected")
				return
			default:
				var response model.StreamResponse
				if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
					log.Printf("Error unmarshaling JSON: %v", err)
					continue
				}

				// Send SSE formatted data
				if response.Message.Content != "" {
					fmt.Fprintf(w, "data: %s\n\n", response.Message.Content)
					flusher.Flush()
				}

				// Send done event
				if response.Done {
					fmt.Fprintf(w, "event: done\ndata: stream complete\n\n")
					flusher.Flush()
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Scanner error: %v", err)
		}
	}
}

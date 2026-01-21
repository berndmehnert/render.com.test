package client

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"backend/config"
	"backend/model"
)

type OllamaClient struct {
	httpClient *http.Client
	endpoint   string
	apiKey     string
	model      string
}

func NewOllama(cfg *config.Config) *OllamaClient {
	return &OllamaClient{
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		endpoint:   cfg.OllamaEndpoint,
		apiKey:     cfg.OllamaAPIKey,
		model:      cfg.OllamaModel,
	}
}

func (c *OllamaClient) Chat(ctx context.Context, prompt string) (*http.Response, error) {
	requestBody := model.RequestBody{
		Model: c.model,
		Messages: []model.Message{
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return nil, err
	}

	// Create context-aware request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}

	// Check upstream status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Ollama returned status: %d", resp.StatusCode)
		return nil, err
	}

	return resp, nil
}

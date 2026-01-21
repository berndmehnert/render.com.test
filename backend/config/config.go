package config

import "os"

type Config struct {
	Port           string
	OllamaEndpoint string
	OllamaModel    string
	OllamaAPIKey   string
	UiUrl          string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		OllamaEndpoint: getEnv("OLLAMA_API_ENDPOINT", "http://localhost:11434/api/chat"),
		OllamaModel:    getEnv("OLLAMA_API_MODEL", "llama3.2"),
		OllamaAPIKey:   os.Getenv("OLLAMA_API_KEY"),
		UiUrl:          getEnv("UI_URL", "http://localhost:5173"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

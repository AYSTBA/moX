package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const mimoBase = "https://api.xiaomimimo.com/v1"

type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content,omitempty"`
}

type ChatRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	Stream    bool          `json:"stream"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

type ChatResponse struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Choices []Choice      `json:"choices"`
	Usage   *Usage        `json:"usage,omitempty"`
}

type Choice struct {
	Index        int            `json:"index"`
	Delta        *MessageDelta  `json:"delta,omitempty"`
	Message      *ChatMessage   `json:"message,omitempty"`
	FinishReason string         `json:"finish_reason,omitempty"`
}

type MessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

var configDir string

func getConfigDir() string {
	if configDir != "" {
		return configDir
	}
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".mox")
	os.MkdirAll(configDir, 0755)
	return configDir
}

func getSettingsPath() string {
	return filepath.Join(getConfigDir(), "settings.json")
}

type Settings struct {
	APIKey string `json:"api_key"`
	Model  string `json:"model"`
}

func loadSettings() *Settings {
	s := &Settings{Model: "mimo-v2.5"}
	data, _ := os.ReadFile(getSettingsPath())
	if len(data) > 0 {
		json.Unmarshal(data, s)
	}
	return s
}

func saveSettings(s *Settings) error {
	data, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(getSettingsPath(), data, 0644)
}

func main() {
	http.HandleFunc("/v1/chat/completions", corsMiddleware(handleChat))
	http.HandleFunc("/v1/models", corsMiddleware(handleModels))
	http.HandleFunc("/api/settings", corsMiddleware(handleSettings))

	addr := ":3099"
	log.Printf("MOX backend listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next(w, r)
	}
}

func handleModels(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"object": "list",
		"data": []map[string]string{
			{"id": "mimo-v2.5-pro", "object": "model", "owned_by": "mimo"},
			{"id": "mimo-v2.5", "object": "model", "owned_by": "mimo"},
			{"id": "mimo-v2-pro", "object": "model", "owned_by": "mimo"},
			{"id": "mimo-v2-flash", "object": "model", "owned_by": "mimo"},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s := loadSettings()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	case "POST":
		var s Settings
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := saveSettings(&s); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", 405)
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	settings := loadSettings()
	apiKey := r.Header.Get("X-Mimo-Key")
	if apiKey == "" {
		apiKey = settings.APIKey
	}
	if apiKey == "" {
		http.Error(w, `{"error":"API Key not configured"}`, 401)
		return
	}

	mimoReq := map[string]interface{}{
		"model":    req.Model,
		"messages": req.Messages,
		"stream":   true,
	}

	body, _ := json.Marshal(mimoReq)
	mimoURL := mimoBase + "/chat/completions"

	httpReq, _ := http.NewRequestWithContext(r.Context(), "POST", mimoURL, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	mimoResp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 502)
		return
	}
	defer mimoResp.Body.Close()

	if mimoResp.StatusCode != 200 {
		errBody, _ := io.ReadAll(mimoResp.Body)
		http.Error(w, string(errBody), mimoResp.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)

	id := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
	scanner := bufio.NewScanner(mimoResp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var fullContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			// Send final usage if available
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}

		var mimoEvent struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
			Usage *Usage `json:"usage,omitempty"`
		}
		if err := json.Unmarshal([]byte(data), &mimoEvent); err != nil {
			continue
		}

		if len(mimoEvent.Choices) == 0 {
			continue
		}
		delta := mimoEvent.Choices[0].Delta
		finish := mimoEvent.Choices[0].FinishReason

		fullContent.WriteString(delta.Content)

		openaiEvent := ChatResponse{
			ID:     id,
			Object: "chat.completion.chunk",
			Choices: []Choice{{
				Index: 0,
				Delta: &MessageDelta{Content: delta.Content},
			}},
		}

		if finish != "" {
			openaiEvent.Choices[0].FinishReason = finish
		}

		evt, _ := json.Marshal(openaiEvent)
		fmt.Fprintf(w, "data: %s\n\n", evt)
		flusher.Flush()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}
}


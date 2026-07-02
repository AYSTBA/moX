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

type StreamChoice struct {
	Index int          `json:"index"`
	Delta *MessageDelta `json:"delta,omitempty"`
}

type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Choices []StreamChoice `json:"choices"`
}

type MessageDelta struct {
	Content string `json:"content,omitempty"`
}

type Settings struct {
	APIKey       string `json:"api_key"`
	Model        string `json:"model"`
	SystemPrompt string `json:"system_prompt"`
}

var configDir string

func getConfigDir() string {
	if configDir != "" { return configDir }
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".mox")
	os.MkdirAll(configDir, 0755)
	return configDir
}
func getSettingsPath() string { return filepath.Join(getConfigDir(), "settings.json") }

func loadSettings() *Settings {
	s := &Settings{Model: "mimo-v2.5", SystemPrompt: "你是MiMo，由小米公司研发的AI智能助手。你擅长中文对话、逻辑推理和代码编写。回答时保持简洁、准确。"}
	data, _ := os.ReadFile(getSettingsPath())
	if len(data) > 0 { json.Unmarshal(data, s) }
	return s
}
func saveSettings(s *Settings) error {
	data, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(getSettingsPath(), data, 0644)
}

var defaultSystem = "你是MiMo，由小米公司研发的AI智能助手。你擅长中文对话、逻辑推理和代码编写。回答时保持简洁、准确。"

func buildSystemPrompt(s *Settings) string {
	var b strings.Builder
	p := s.SystemPrompt
	if p == "" { p = defaultSystem }
	b.WriteString(p)
	b.WriteString(fmt.Sprintf("\n当前时间：%s", time.Now().Format("2006-01-02 15:04:05 Mon")))
	return b.String()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", cors(handleChat))
	mux.HandleFunc("/v1/models", cors(handleModels))
	mux.HandleFunc("/api/settings", cors(handleSettings))

	// Determine port from env or default
	port := os.Getenv("MOX_PORT")
	if port == "" { port = "3099" }

	log.Printf("MOX backend starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Mimo-Key")
		if r.Method == "OPTIONS" { w.WriteHeader(204); return }
		next(w, r)
	}
}

func handleModels(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"object": "list",
		"data": []map[string]string{
			{"id": "mimo-v2.5-pro", "object": "model", "owned_by": "mimo"},
			{"id": "mimo-v2.5", "object": "model", "owned_by": "mimo"},
			{"id": "mimo-v2-pro", "object": "model", "owned_by": "mimo"},
			{"id": "mimo-v2-flash", "object": "model", "owned_by": "mimo"},
		},
	})
}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(loadSettings())
	case "POST":
		var s Settings
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), 400); return
		}
		if err := saveSettings(&s); err != nil {
			http.Error(w, err.Error(), 500); return
		}
		w.WriteHeader(200)
	}
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" { http.Error(w, "method not allowed", 405); return }

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400); return
	}

	settings := loadSettings()
	apiKey := r.Header.Get("X-Mimo-Key")
	if apiKey == "" { apiKey = settings.APIKey }
	if apiKey == "" { http.Error(w, `{"error":"API Key not configured"}`, 401); return }

	// Prepend system prompt with time awareness
	var messages []ChatMessage
	sysMsg := ChatMessage{Role: "system", Content: buildSystemPrompt(settings)}
	messages = append(messages, sysMsg)
	messages = append(messages, req.Messages...)

	mimoReq := map[string]interface{}{
		"model":    req.Model,
		"messages": messages,
		"stream":   true,
	}
	body, _ := json.Marshal(mimoReq)

	httpReq, _ := http.NewRequestWithContext(r.Context(), "POST", mimoBase+"/chat/completions", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	mimoResp, err := client.Do(httpReq)
	if err != nil { http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), 502); return }
	defer mimoResp.Body.Close()

	if mimoResp.StatusCode != 200 {
		errBody, _ := io.ReadAll(mimoResp.Body)
		// Return API error as OpenAI-compatible error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mimoResp.StatusCode)
		fmt.Fprintf(w, `{"error":{"message":%s}}`, errBody)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, _ := w.(http.Flusher)

	id := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
	scanner := bufio.NewScanner(mimoResp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") { continue }
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}

		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil { continue }
		if len(event.Choices) == 0 { continue }

		evt, _ := json.Marshal(StreamResponse{
			ID:     id,
			Object: "chat.completion.chunk",
			Choices: []StreamChoice{{Index: 0, Delta: &MessageDelta{Content: event.Choices[0].Delta.Content}}},
		})
		fmt.Fprintf(w, "data: %s\n\n", evt)
		flusher.Flush()
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Stream error: %v", err)
	}
}
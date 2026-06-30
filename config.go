package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	APIKey      string `json:"api_key"`
	Model       string `json:"model"`
	Theme       string `json:"theme"`
	SystemPrompt string `json:"system_prompt"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	MaxTokens   int    `json:"max_tokens"`
	ThinkingEnabled bool `json:"thinking_enabled"`
	WebSearchEnabled bool `json:"web_search_enabled"`
	ExternalSearchAPIKey string `json:"external_search_api_key"`
	ExternalSearchEnabled bool `json:"external_search_enabled"`
	TimeAwareness bool `json:"time_awareness"`
	PersonalizationEnabled bool   `json:"personalization_enabled"`
	BlurIntensity           int    `json:"blur_intensity"`
	BackgroundImage        string `json:"background_image"`
}

type Session struct {
	Key       string    `json:"key"`
	Label     string    `json:"label"`
	Messages  []Message `json:"messages"`
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"updated_at"`
}

type Message struct {
	ID               string       `json:"id"`
	Role             string       `json:"role"`
	Content          string       `json:"content"`
	ReasoningContent string       `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall   `json:"tool_calls,omitempty"`
	Annotations      []Annotation `json:"annotations,omitempty"`
	Attachments      []AttachmentItem `json:"attachments,omitempty"`
	Usage            *Usage           `json:"usage,omitempty"`
	Timestamp        int64        `json:"timestamp"`
}

type AttachmentItem struct {
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Size int64  `json:"size,omitempty"`
	Data string `json:"data,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
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

func getSessionsDir() string {
	dir := filepath.Join(getConfigDir(), "sessions")
	os.MkdirAll(dir, 0755)
	return dir
}

func LoadSettings() *Settings {
	s := &Settings{
		Model:       "mimo-v2.5",
		Theme:       "dark",
		SystemPrompt: "你是MiMo，是小米公司研发的AI智能助手。",
		Temperature: 1.0,
		TopP:        0.95,
		MaxTokens:   4096,
		ThinkingEnabled: true,
	}
	data, err := os.ReadFile(getSettingsPath())
	if err == nil {
		json.Unmarshal(data, s)
	}
	return s
}

func SaveSettings(s *Settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getSettingsPath(), data, 0644)
}

func LoadSessions() []Session {
	dir := getSessionsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var sessions []Session
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var s Session
		if json.Unmarshal(data, &s) == nil {
			sessions = append(sessions, s)
		}
	}
	return sessions
}

func SaveSession(s *Session) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(getSessionsDir(), s.Key+".json"), data, 0644)
}

func DeleteSession(key string) error {
	path := filepath.Join(getSessionsDir(), key+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

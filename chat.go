package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiBase = "https://api.xiaomimimo.com/v1"

type ChatMessage struct {
	Role               string      `json:"role"`
	Content            interface{} `json:"content,omitempty"`
	ReasoningContent   string      `json:"reasoning_content,omitempty"`
	ToolCalls          []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID         string      `json:"tool_call_id,omitempty"`
}

type ChatRequest struct {
	Model              string        `json:"model"`
	Messages           []ChatMessage `json:"messages"`
	Stream             bool          `json:"stream"`
	MaxCompletionTokens int          `json:"max_completion_tokens,omitempty"`
	Temperature        float64       `json:"temperature,omitempty"`
	TopP               float64       `json:"top_p,omitempty"`
	FrequencyPenalty   float64       `json:"frequency_penalty,omitempty"`
	PresencePenalty    float64       `json:"presence_penalty,omitempty"`
	Thinking           *Thinking     `json:"thinking,omitempty"`
}

type Thinking struct {
	Type string `json:"type"`
}

type StreamDelta struct {
	Role             string     `json:"role"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content"`
	ToolCalls        []ToolCall `json:"tool_calls"`
}

type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason string      `json:"finish_reason"`
}

type StreamResponse struct {
	ID      string         `json:"id"`
	Choices []StreamChoice `json:"choices"`
	Usage   *Usage         `json:"usage,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamEvent struct {
	Type     string // "token", "thinking", "toolcall", "done", "error"
	Content  string
	Reasoning string
	ToolCalls []ToolCall
	Finish   string
	Error    string
	Usage    *Usage
}

func SendChatMessage(ctx context.Context, apiKey string, req ChatRequest, events chan<- StreamEvent) error {
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiBase+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(errBody))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			events <- StreamEvent{Type: "done"}
			return nil
		}

		var sr StreamResponse
		if err := json.Unmarshal([]byte(data), &sr); err != nil {
			continue
		}

		if len(sr.Choices) == 0 {
			continue
		}

		delta := sr.Choices[0].Delta
		finish := sr.Choices[0].FinishReason

		if delta.ReasoningContent != "" {
			events <- StreamEvent{Type: "thinking", Reasoning: delta.ReasoningContent}
		}

		if delta.Content != "" {
			events <- StreamEvent{Type: "token", Content: delta.Content}
		}

		if len(delta.ToolCalls) > 0 {
			events <- StreamEvent{Type: "toolcall", ToolCalls: delta.ToolCalls}
		}

		if finish != "" {
			events <- StreamEvent{Type: "done", Finish: finish, Usage: sr.Usage}
			return nil
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	events <- StreamEvent{Type: "done"}
	return nil
}

func GenerateTitle(ctx context.Context, apiKey string, messages []ChatMessage) (string, error) {
	trimmed := messages
	if len(trimmed) > 6 {
		trimmed = trimmed[len(trimmed)-6:]
	}

	prompt := []ChatMessage{
		{Role: "system", Content: "根据以下对话内容，生成一个简短的中文标题（不超过15个字）。只输出标题，不要任何解释、引号或标点。"},
	}
	prompt = append(prompt, trimmed...)

	body, _ := json.Marshal(ChatRequest{
		Model:              "mimo-v2-flash",
		Messages:           prompt,
		Stream:             false,
		MaxCompletionTokens: 50,
		Temperature:        0.3,
	})

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiBase+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := (&http.Client{}).Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("title API error %d: %s", resp.StatusCode, string(errBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	title := strings.TrimSpace(result.Choices[0].Message.Content)
	title = strings.Trim(title, "\"'")
	return title, nil
}

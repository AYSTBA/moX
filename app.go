package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetSettings() *Settings {
	return LoadSettings()
}

func (a *App) SaveSettings(s *Settings) error {
	return SaveSettings(s)
}

func (a *App) GetSessions() []Session {
	return LoadSessions()
}

func (a *App) CreateSession(label string) *Session {
	s := &Session{
		Key:       uuid.New().String()[:8],
		Label:     label,
		Messages:  []Message{},
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	SaveSession(s)
	return s
}

func (a *App) SaveSession(s *Session) error {
	s.UpdatedAt = time.Now().UnixMilli()
	return SaveSession(s)
}

func (a *App) DeleteSession(key string) error {
	return DeleteSession(key)
}

func (a *App) SendMessage(sessionKey string, userContent string, model string, thinking bool) {
	settings := LoadSettings()
	apiKey := settings.APIKey
	if apiKey == "" {
		runtime.EventsEmit(a.ctx, "chat:error", "请先在设置中配置 API Key")
		return
	}

	sessions := LoadSessions()
	var session *Session
	for i := range sessions {
		if sessions[i].Key == sessionKey {
			session = &sessions[i]
			break
		}
	}
	if session == nil {
		runtime.EventsEmit(a.ctx, "chat:error", "会话不存在")
		return
	}

	userMsg := Message{
		ID:        uuid.New().String(),
		Role:      "user",
		Content:   userContent,
		Timestamp: time.Now().UnixMilli(),
	}
	session.Messages = append(session.Messages, userMsg)
	SaveSession(session)

	runtime.EventsEmit(a.ctx, "chat:userMessage", userMsg)

	apiMessages := make([]ChatMessage, 0)
	if settings.SystemPrompt != "" || settings.TimeAwareness {
		prompt := settings.SystemPrompt
		if settings.TimeAwareness {
			now := time.Now()
			weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
			timeInfo := fmt.Sprintf("\n当前时间：%s 星期%s %02d:%02d", now.Format("2006年01月02日"), weekdays[now.Weekday()], now.Hour(), now.Minute())
			prompt = prompt + timeInfo
		}
		apiMessages = append(apiMessages, ChatMessage{
			Role:    "system",
			Content: prompt,
		})
	}

	if settings.ExternalSearchEnabled && settings.ExternalSearchAPIKey != "" && userContent != "" {
		results, err := ExternalSearch(a.ctx, settings.ExternalSearchAPIKey, userContent)
		if err != nil {
			runtime.EventsEmit(a.ctx, "chat:toast", "搜索失败: "+err.Error())
		} else if len(results) > 0 {
			var sb strings.Builder
			sb.WriteString("你具备联网搜索能力。以下是针对用户问题搜索到的实时网络资料：\n")
			for i, r := range results {
				sb.WriteString(fmt.Sprintf("\n\n[%d] %s\n来源: %s\n%s", i+1, r.Title, r.URL, r.Content))
			}
			sb.WriteString("\n\n请基于以上搜索结果回答用户问题。直接引用其中的信息，不要说「我无法访问互联网」或「我的知识有截止日期」。如果搜索结果不足以回答，如实说明。")
			apiMessages = append(apiMessages, ChatMessage{
				Role:    "system",
				Content: sb.String(),
			})
		}
	}

	for _, m := range session.Messages {
		cm := ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		}
		if m.ReasoningContent != "" {
			cm.ReasoningContent = m.ReasoningContent
		}
		if len(m.ToolCalls) > 0 {
			cm.ToolCalls = m.ToolCalls
		}
		apiMessages = append(apiMessages, cm)
	}

	req := ChatRequest{
		Model:              model,
		Messages:           apiMessages,
		MaxCompletionTokens: settings.MaxTokens,
		Temperature:        settings.Temperature,
		TopP:               settings.TopP,
	}

	if thinking {
		req.Thinking = &Thinking{Type: "enabled"}
	}

	if settings.WebSearchEnabled {
		req.Tools = []interface{}{
			WebSearchTool{
				Type:        "web_search",
				ForceSearch: false,
				MaxKeyword:  3,
				Limit:       5,
			},
		}
	}

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel

	events := make(chan StreamEvent, 100)

	go func() {
		err := SendChatMessage(ctx, apiKey, req, events)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "webSearchEnabled is false") {
				errMsg = "MiMo 联网搜索插件未开启。请前往 platform.xiaomimimo.com 控制台激活搜索插件，或关闭搜索开关。"
			}
			events <- StreamEvent{Type: "error", Error: errMsg}
		}
	}()

	assistantMsg := Message{
		ID:        uuid.New().String(),
		Role:      "assistant",
		Timestamp: time.Now().UnixMilli(),
	}

	var reasoningBuilder strings.Builder
	var contentBuilder strings.Builder
	var toolCallsBuilder []ToolCall
	var annotationsBuilder []Annotation

	for event := range events {
		switch event.Type {
		case "thinking":
			reasoningBuilder.WriteString(event.Reasoning)
			runtime.EventsEmit(a.ctx, "chat:thinking", reasoningBuilder.String())
		case "token":
			contentBuilder.WriteString(event.Content)
			runtime.EventsEmit(a.ctx, "chat:token", contentBuilder.String())
		case "toolcall":
			toolCallsBuilder = append(toolCallsBuilder, event.ToolCalls...)
			runtime.EventsEmit(a.ctx, "chat:toolcall", toolCallsBuilder)
		case "annotations":
			annotationsBuilder = append(annotationsBuilder, event.Annotations...)
			runtime.EventsEmit(a.ctx, "chat:annotations", annotationsBuilder)
		case "done":
			assistantMsg.Content = contentBuilder.String()
			assistantMsg.ReasoningContent = reasoningBuilder.String()
			assistantMsg.ToolCalls = toolCallsBuilder
			assistantMsg.Annotations = annotationsBuilder
			session.Messages = append(session.Messages, assistantMsg)
			SaveSession(session)
			runtime.EventsEmit(a.ctx, "chat:done", assistantMsg)

			if len(session.Messages) >= 4 && session.Label == "新对话" {
				go a.generateTitle(apiKey, session)
			}
		case "error":
			runtime.EventsEmit(a.ctx, "chat:error", event.Error)
		}
	}
}

func (a *App) StopGeneration() {
	if a.cancelFunc != nil {
		a.cancelFunc()
		a.cancelFunc = nil
	}
}

func (a *App) GetModels() []map[string]string {
	return []map[string]string{
		{"id": "mimo-v2.5-pro", "name": "MiMo-V2.5 Pro", "desc": "旗舰模型，深度思考，1M上下文"},
		{"id": "mimo-v2.5", "name": "MiMo-V2.5", "desc": "全模态理解（文本/图片/音频/视频）"},
		{"id": "mimo-v2-pro", "name": "MiMo-V2 Pro", "desc": "上一代旗舰"},
		{"id": "mimo-v2-flash", "name": "MiMo-V2 Flash", "desc": "快速响应，低成本"},
		{"id": "mimo-v2-omni", "name": "MiMo-V2 Omni", "desc": "多模态理解"},
	}
}

func (a *App) generateTitle(apiKey string, session *Session) {
	apiMessages := make([]ChatMessage, 0, len(session.Messages))
	for _, m := range session.Messages {
		apiMessages = append(apiMessages, ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	title, err := GenerateTitle(a.ctx, apiKey, apiMessages)
	if err != nil || title == "" {
		return
	}

	session.Label = title
	SaveSession(session)
	runtime.EventsEmit(a.ctx, "chat:titleUpdated", map[string]string{
		"key":   session.Key,
		"label": title,
	})
}

func (a *App) TestAPIKey(apiKey string) string {
	req, err := http.NewRequest("GET", apiBase+"/models", nil)
	if err != nil {
		return fmt.Sprintf("请求创建失败: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("连接失败: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return "ok"
	}
	return fmt.Sprintf("API 返回错误: %d", resp.StatusCode)
}

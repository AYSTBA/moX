package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	mu         sync.Mutex
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

func (a *App) SendMessage(sessionKey string, userContent string, model string, thinking bool, attachments []FileAttachment) {
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
		Attachments: toAttachmentItems(attachments),
		Timestamp: time.Now().UnixMilli(),
	}
	session.Messages = append(session.Messages, userMsg)
	SaveSession(session)

	runtime.EventsEmit(a.ctx, "chat:userMessage", userMsg)

	apiMessages := make([]ChatMessage, 0)
	if settings.SystemPrompt != "" || settings.TimeAwareness || settings.ExternalSearchEnabled {
		prompt := settings.SystemPrompt
		if settings.TimeAwareness {
			now := time.Now()
			weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
			timeInfo := fmt.Sprintf("\n当前时间：%s 星期%s %02d:%02d", now.Format("2006年01月02日"), weekdays[now.Weekday()], now.Hour(), now.Minute())
			prompt = prompt + timeInfo
		}
		if settings.ExternalSearchEnabled {
			prompt = prompt + "\n\n你拥有联网搜索能力。当用户提问涉及实时信息、最新资讯、特定URL内容、当前事件等时，系统会自动搜索并将结果传入。你应当直接使用搜索结果回答，不要说「我无法访问互联网」或「我的知识有截止日期」。"
		}
		apiMessages = append(apiMessages, ChatMessage{
			Role:    "system",
			Content: prompt,
		})
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

	// If there are file attachments, convert the last user message to multimodal content parts
	if len(attachments) > 0 {
		for i := len(apiMessages) - 1; i >= 0; i-- {
			if apiMessages[i].Role == "user" {
				var parts []ContentPart
				if text, ok := apiMessages[i].Content.(string); ok && text != "" {
					parts = append(parts, ContentPart{Type: "text", Text: text})
				}
				for _, att := range attachments {
					if strings.HasPrefix(att.MimeType, "image/") {
						parts = append(parts, ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: "data:" + att.MimeType + ";base64," + att.Data}})
					} else {
						parts = append(parts, ContentPart{Type: "text", Text: fmt.Sprintf("[附件: %s]", att.Name)})
					}
				}
				apiMessages[i].Content = parts
				break
			}
		}
	}

	ctx, cancel := context.WithCancel(a.ctx)
	a.mu.Lock()
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
	a.cancelFunc = cancel
	a.mu.Unlock()

	if settings.ExternalSearchEnabled && settings.ExternalSearchAPIKey != "" {
		a.agentLoop(ctx, apiKey, model, thinking, settings, apiMessages, session)
	} else {
		a.directChat(ctx, apiKey, model, thinking, settings, apiMessages, session, nil)
	}
}

func (a *App) agentLoop(ctx context.Context, apiKey, model string, thinking bool, settings *Settings, history []ChatMessage, session *Session) {
	runtime.EventsEmit(a.ctx, "chat:status", "planning")

	lastUserMsg := ""
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			switch c := history[i].Content.(type) {
			case string:
				lastUserMsg = c
			case []ContentPart:
				for _, part := range c {
					if part.Type == "text" {
						lastUserMsg = part.Text
						break
					}
				}
			}
			break
		}
	}

	planMessages := []ChatMessage{
		{Role: "system", Content: `你是一个规划助手。判断用户问题并返回JSON格式结果。

格式: {"action":"search"|"think"|"direct","query":"搜索关键词"}

规则:
- "search": 需要实时信息、最新资讯、特定URL/仓库/文件内容、当前事件
- "think": 问题复杂需要深入推理
- "direct": 可以直接回答

只输出JSON，不要解释。`},
		{Role: "user", Content: lastUserMsg},
	}

	planReq := ChatRequest{
		Model:              model,
		Messages:           planMessages,
		Stream:             false,
		MaxCompletionTokens: 30,
		Temperature:        0.1,
	}

	planResp, err := SendChatMessageSync(ctx, apiKey, planReq)
	if err != nil {
		runtime.EventsEmit(a.ctx, "chat:status", "")
		a.directChat(ctx, apiKey, model, thinking, settings, history, session, nil)
		return
	}

	runtime.EventsEmit(a.ctx, "chat:toast", "Agent: "+planResp)

	// Strip markdown code fences
	planResp = strings.TrimSpace(planResp)
	planResp = strings.TrimPrefix(planResp, "```json")
	planResp = strings.TrimPrefix(planResp, "```")
	planResp = strings.TrimSuffix(planResp, "```")
	planResp = strings.TrimSpace(planResp)

	type planResponse struct {
		Action string `json:"action"`
		Query  string `json:"query"`
	}
	var plan planResponse
	if err := json.Unmarshal([]byte(planResp), &plan); err != nil {
		// Non-JSON response means LLM misunderstood — fallback to direct
		runtime.EventsEmit(a.ctx, "chat:toast", "Agent 回复格式异常，跳过搜索")
		a.directChat(ctx, apiKey, model, thinking, settings, history, session, nil)
		return
	}

	if plan.Action == "" {
		runtime.EventsEmit(a.ctx, "chat:toast", "Agent 响应为空，跳过搜索")
		a.directChat(ctx, apiKey, model, thinking, settings, history, session, nil)
		return
	}

	var searchResults []TavilyResult
	if plan.Action == "search" && plan.Query != "" {
		runtime.EventsEmit(a.ctx, "chat:status", "searching")
		searchResults, err = ExternalSearch(ctx, settings.ExternalSearchAPIKey, plan.Query)
		if err != nil {
			runtime.EventsEmit(a.ctx, "chat:toast", "搜索失败: "+err.Error())
		}
		runtime.EventsEmit(a.ctx, "chat:status", "")
	}

	a.directChat(ctx, apiKey, model, thinking || plan.Action == "think", settings, history, session, searchResults)
}

func (a *App) directChat(ctx context.Context, apiKey, model string, thinking bool, settings *Settings, history []ChatMessage, session *Session, searchResults []TavilyResult) {
	apiMessages := append([]ChatMessage{}, history...)

	if len(searchResults) > 0 {
		var sb strings.Builder
		sb.WriteString("<search_results>\n")
		for i, r := range searchResults {
			sb.WriteString(fmt.Sprintf("[%d] %s\nURL: %s\n内容: %s\n\n", i+1, r.Title, r.URL, r.Content))
		}
		sb.WriteString("</search_results>")
		apiMessages = append(apiMessages, ChatMessage{
			Role:    "system",
			Content: sb.String(),
		})
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
			cleanContent := cleanFunctionCalls(contentBuilder.String())
			assistantMsg.Content = cleanContent
			assistantMsg.ReasoningContent = reasoningBuilder.String()
			assistantMsg.ToolCalls = toolCallsBuilder
			assistantMsg.Annotations = annotationsBuilder
			if event.Usage != nil {
				assistantMsg.Usage = &Usage{
					PromptTokens:     event.Usage.PromptTokens,
					CompletionTokens: event.Usage.CompletionTokens,
					TotalTokens:      event.Usage.TotalTokens,
				}
			}
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

func cleanFunctionCalls(content string) string {
	// 仅匹配工具调用 JSON 模式 {"name":"...","arguments":"..."}，避免误删含 name 字段的普通 JSON
	re := regexp.MustCompile(`\{\s*"name"\s*:\s*"[^"]+"\s*,\s*"arguments"\s*:\s*"[^"]*"\s*\}`)
	content = re.ReplaceAllString(content, "")
	content = strings.TrimSpace(content)
	return content
}

func (a *App) StopGeneration() {
	a.mu.Lock()
	defer a.mu.Unlock()
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
	// 仅用调用时刻的消息快照构建生成标题的 prompt
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

	// 从磁盘重新加载会话，仅更新 Label 字段后保存，避免覆盖此期间新增的消息
	sessions := LoadSessions()
	for i := range sessions {
		if sessions[i].Key == session.Key {
			sessions[i].Label = title
			if err := SaveSession(&sessions[i]); err == nil {
				runtime.EventsEmit(a.ctx, "chat:titleUpdated", map[string]string{
					"key":   session.Key,
					"label": title,
				})
			}
			return
		}
	}
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

func toAttachmentItems(attachments []FileAttachment) []AttachmentItem {
	items := make([]AttachmentItem, 0, len(attachments))
	for _, att := range attachments {
		dataURI := "data:" + att.MimeType + ";base64," + att.Data
		items = append(items, AttachmentItem{
			Name: att.Name,
			Type: att.MimeType,
			Data: dataURI,
		})
	}
	return items
}

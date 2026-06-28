# MOX

一个基于 Wails 构建的小米 MiMo AI 桌面客户端。

## 为什么做这个

小米 MiMo 模型能力很强，但官方一直没有面向普通用户（C 端）的客户端产品——既没有电脑端，也没有手机端。普通用户想要体验 MiMo，只能在 [开发者后台用 API 调试的方式手动测试，门槛很高，或者用网页，体验也很糟糕。

MOX 是想解决这个问题的一次尝试：

- 让普通用户也能像用 ChatGPT、豆包那样，直接打开客户端就能和 MiMo 对话
- 接入小米账号登录（理想情况下由官方支持），免去手动填 API Key 的麻烦
- 哪怕只是一个雏形，也希望推动小米官方重视 C 端产品，尽快推出手机端 App（电脑端可以晚一点，但手机端真的不能再等了）

如果你也觉得小米应该出一个面向 C 端的 MiMo 客户端，欢迎在 Issue 区发声、Star 本项目，把需求传递出去。

## 功能

- 多模型支持：MiMo-V2.5 Pro、MiMo-V2.5、MiMo-V2 Pro、MiMo-V2 Flash、MiMo-V2 Omni
- 思考模式：可视化展示模型推理过程
- 联网搜索：集成 Tavily API，获取实时信息
- 会话管理：本地持久化存储对话历史
- 时间感知：自动注入当前时间信息
- 双主题：深色 / 浅色切换
- 流式响应：实时输出对话内容

## 技术栈

- 后端：Go 1.23 + Wails v2.12
- 前端：Vue 3 + Vite + Pinia
- 渲染：WebView2

## 快速开始

### 前置要求

- Go >= 1.23
- Node.js >= 18
- Wails CLI >= 2.12

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发

```bash
git clone https://github.com/AYSTBA/mox.git
cd mox/mox

npm install --prefix frontend

wails dev
```

### 构建

```bash
wails build                     # Windows
wails build -platform darwin    # macOS
wails build -platform linux     # Linux
```

## 配置

API Key 获取：

- MiMo API Key：前往 [platform.xiaomimimo.com](https://platform.xiaomimimo.com)
- Tavily API Key：前往 [tavily.com](https://tavily.com)（用于联网搜索）

配置文件位于 `~/.mox/settings.json`：

```json
{
  "api_key": "your-mimo-api-key",
  "model": "mimo-v2.5",
  "theme": "dark",
  "system_prompt": "你是MiMo，是小米公司研发的AI智能助手。",
  "temperature": 1.0,
  "top_p": 0.95,
  "max_tokens": 4096,
  "thinking_enabled": true,
  "web_search_enabled": false,
  "external_search_api_key": "your-tavily-api-key",
  "external_search_enabled": true,
  "time_awareness": true
}
```

## 项目结构

```
mox/
├── main.go                  # 应用入口
├── app.go                   # 主应用逻辑
├── chat.go                  # AI 聊天核心逻辑
├── config.go                # 配置与数据持久化
├── go.mod                   # Go 依赖
├── wails.json               # Wails 配置
├── frontend/                # Vue 前端
│   └── src/
│       ├── App.vue
│       ├── main.js
│       ├── style.css
│       ├── components/
│       │   ├── ChatInput.vue
│       │   ├── ChatView.vue
│       │   ├── MessageBubble.vue
│       │   ├── ModelSelector.vue
│       │   ├── Settings.vue
│       │   └── Sidebar.vue
│       └── stores/
│           ├── chat.js
│           └── settings.js
└── build/                   # 构建资源
```

## License

MIT

## 致谢

- [Wails](https://wails.io/)
- [MiMo](https://platform.xiaomimimo.com)
- [Tavily](https://tavily.com)
- [Vue.js](https://vuejs.org/)

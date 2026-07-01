# MOX

<div align="center">

MiMo AI 桌面客户端 - 基于 Wails

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![Wails](https://img.shields.io/badge/Wails-2.12-DF4A2F?style=flat&logo=wails)](https://wails.io/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat)](LICENSE)

[功能](#功能) | [技术栈](#技术栈) | [快速开始](#快速开始) | [项目结构](#项目结构) | [配置](#配置) | [更新日志](./CHANGELOG.md)

</div>

---

## 项目简介

MOX 是一个基于 Wails 构建的小米 MiMo AI 桌面客户端。MiMo 模型能力很强，但官方一直没有面向普通用户的客户端产品。MOX 通过在 Wails（Go + Vue 3 + WebView2）之上封装 MiMo API，为普通用户提供开箱即用的桌面端 MiMo 体验。

---

## 功能

### 对话
- 多模型支持：MiMo-V2.5 Pro、MiMo-V2.5、MiMo-V2 Pro、MiMo-V2 Flash、MiMo-V2 Omni
- 流式输出：实时展示 Token、思考过程、工具调用、搜索来源
- 思考时间：推理耗时实时显示，括号中标注用时
- 会话管理：本地持久化存储，自动按时间排序
- 一键复制：完整复制输出内容（不含思考过程）
- 信息面板：鼠标悬浮查看 Token 用量、输出长度、响应时间

### 搜索
- Agent 规划：自动判断是否需要联网搜索
- Tavily 集成：多搜索结果注入上下文
- MiMo 内置搜索：官方搜索插件可选

### 个性化
- 深色 / 浅色双主题
- 自定义背景图片 + 磨砂玻璃效果（雾面强度可调）
- 配色方案（默认 / 蓝色 / 绿色 / 紫色 / 暖色）

### 附件
- 图片 / 视频内嵌显示
- Lightbox 图片放大
- 多模态 content parts 发送

### 系统
- 时间感知：自动注入当前日期时间
- 自定义系统提示词
- Temperature / Top P / Max Tokens 参数调节

---

## 技术栈

### 前端
| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.4 | 渐进式框架 |
| Pinia | 2.x | 状态管理 |
| Vite | 3.2 | 构建工具 |
| marked + highlight.js | - | Markdown 渲染与代码高亮 |

### 后端
| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.23 | 编译语言 |
| Wails | 2.12 | 桌面框架 |
| WebView2 | - | 渲染引擎 |

### API
| 服务 | 说明 |
|------|------|
| MiMo API | 对话与推理 |
| Tavily API | 联网搜索结果 |

---

## 快速开始

### 前置要求

- Go >= 1.23
- Node.js >= 18
- Wails CLI >= 2.12

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 开发模式

```bash
git clone https://github.com/AYSTBA/mox.git
cd mox
npm install --prefix frontend
wails dev
```

### 构建

```bash
wails build                 # Windows（默认）
wails build -platform darwin  # macOS
wails build -platform linux   # Linux
```

构建产物位于 `build/bin/` 目录。

---

## 项目结构

```
mox/
├── main.go                  # 应用入口
├── app.go                   # 主应用逻辑（消息、Agent、事件）
├── chat.go                  # 聊天核心（流式请求、搜索）
├── config.go                # 配置与数据持久化
├── go.mod / go.sum          # Go 依赖
├── wails.json               # Wails 配置
│
├── frontend/                # Vue 前端
│   ├── src/
│   │   ├── App.vue          # 根组件
│   │   ├── main.js          # 入口
│   │   ├── style.css        # 全局样式
│   │   ├── components/
│   │   │   ├── ChatInput.vue
│   │   │   ├── ChatView.vue
│   │   │   ├── MessageBubble.vue
│   │   │   ├── ModelSelector.vue
│   │   │   ├── Settings.vue
│   │   │   └── Sidebar.vue
│   │   └── stores/
│   │       ├── chat.js
│   │       └── settings.js
│   └── wailsjs/             # Wails 自动生成的绑定
│
├── build/                   # 构建资源（图标、安装包配置）
├── .github/workflows/       # CI 工作流
├── CHANGELOG.md
├── LICENSE
└── README.md
```

---

## 配置

### API Key 获取

- MiMo API Key：[platform.xiaomimimo.com](https://platform.xiaomimimo.com)
- Tavily API Key：[tavily.com](https://tavily.com)（可选，用于联网搜索）

### 配置文件

配置文件位于 `~/.mox/settings.json`，程序首次启动时自动创建。

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
  "external_search_api_key": "",
  "external_search_enabled": false,
  "time_awareness": false
}
```

> **注意：** API Key 仅存储在本地配置文件中，不会随程序分发或上传。详见 [安全说明](#安全说明)。

### 安全说明

- API Key 仅存储在 `~/.mox/settings.json`，请勿将此文件分享给他人
- 程序不收集任何使用数据
- 所有对话记录仅保存在本地 `~/.mox/sessions/` 目录

---

## 更新日志

详见 [CHANGELOG.md](./CHANGELOG.md)

---

## 致谢

- [Wails](https://wails.io/) - 桌面应用框架
- [MiMo](https://platform.xiaomimimo.com) - 小米 AI 模型
- [Tavily](https://tavily.com) - 搜索引擎
- [Vue.js](https://vuejs.org/) - 前端框架

---

## License

MIT License

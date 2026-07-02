# MOX

<div align="center">

**MiMo AI 桌面客户端**

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.4-4FC08D?style=flat&logo=vue.js)](https://vuejs.org/)
[![Tauri](https://img.shields.io/badge/Tauri-2.11-FFC131?style=flat&logo=tauri)](https://tauri.app/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat)](LICENSE)

</div>

MOX 是一个基于 **Tauri v2 + Go + Vue 3** 构建的小米 MiMo AI 桌面客户端。MiMo 模型能力很强，但官方没有面向普通用户的桌面产品。MOX 通过封装 MiMo API，提供开箱即用的桌面端 AI 对话体验，支持多对话管理、Markdown 渲染、代码高亮、联网搜索。

---

## 截图

![MOX 主界面](frontend/screenshot2.png)

---

## 功能

### 对话
- 多模型支持：MiMo-V2.5 Pro / V2.5 / V2 Pro / V2 Flash
- 流式输出（SSE）+ 加载动画
- 多对话管理：创建、切换、重命名、删除会话
- 本地持久化（localStorage），对话历史不丢失
- Markdown 渲染 + 代码语法高亮（highlight.js + github-dark 主题）
- 消息气泡间距精调，AI 消息宽度控制在 65%

### 界面
- 左侧边栏：对话列表、新建对话、设置入口
- 侧边栏可折叠，主区域顶部 60px 留白视觉平衡
- 附件按钮（placeholder，位于发送按钮上方）
- 深色主题（#111 / #1a1a1a），全界面适配

### 设置
- MiMo API Key 管理
- 自定义系统提示词（默认："你是mimo"）
- 磨砂玻璃效果（backdrop-filter: blur）对话框

### 系统
- Go 后端独立进程，Tauri sidecar 自动拉起/关闭
- 启动前清理残留进程，端口轮询确保就绪
- 无命令行窗口（Windows GUI 子系统）
- 时间感知：每次请求自动注入当前日期时间

---

## 快速开始

### 前置要求

- Go 1.26+
- Node.js 18+
- Rust 1.70+（Tauri 编译需要）
- Windows WebView2 Runtime（Win10+ 自带）

### 开发模式

```bash
# 1. 启动 Go 后端
cd backend
go run . &

# 2. 启动 Vite 前端
cd frontend
npm install
npm run dev
```

浏览器访问 `http://localhost:5173`

### 生产构建

```bash
# 1. 编译 Go 后端
cd backend
go build -ldflags="-s -w" -o mox-server.exe .

# 2. 复制到 Tauri sidecar 目录
cp mox-server.exe ../frontend/src-tauri/binaries/mox-server-x86_64-pc-windows-msvc.exe

# 3. Tauri 构建（自动编译前端 + Rust 壳 + 打包 sidecar）
cd ../frontend
npm install
npm run tauri build
```

构建产物位于 `frontend/src-tauri/target/release/mox.exe`（单文件，双击运行）。

---

## 项目结构

```
mox/
├── backend/                  # Go 后端
│   ├── main.go               # 服务入口 /v1/chat/completions + /api/settings
│   ├── go.mod / go.sum
│
├── frontend/                 # Vue 3 前端
│   ├── src/
│   │   ├── App.vue           # 根组件（侧边栏 + 聊天区 + 设置对话框）
│   │   ├── main.js           # 入口
│   │   ├── style.css         # 全局样式
│   │   ├── api.js            # API 调用 + 设置管理
│   │   └── stores/
│   │       └── chat.js       # Pinia 会话状态
│   ├── src-tauri/            # Tauri Rust 壳
│   │   ├── src/
│   │   │   ├── main.rs       # 入口（windows_subsystem = "windows"）
│   │   │   └── lib.rs        # 后端进程管理 + 清除残留
│   │   ├── icons/            # 应用图标
│   │   ├── binaries/         # sidecar 二进制
│   │   └── tauri.conf.json
│   ├── index.html
│   ├── package.json
│   └── vite.config.js
│
├── CHANGELOG.md
├── LICENSE
└── README.md
```

---

## 配置

### API Key

前往 [platform.xiaomimimo.com](https://platform.xiaomimimo.com) 获取 MiMo API Key，在设置面板中输入保存。

### 系统提示词

设置面板中可自定义系统提示词，默认值为 `"你是mimo"`。提示词通过 Go 后端 `/api/settings` 接口持久化到 `~/.mox/settings.json`。

> **安全说明：** API Key 仅存储在本地配置文件中，不会随程序分发或上传。程序不收集任何使用数据，所有对话记录仅保存在浏览器本地存储中。

---

## 更新日志

详见 [CHANGELOG.md](./CHANGELOG.md)

---

## 技术栈

| 层 | 技术 |
|---|---|
| 桌面壳 | [Tauri v2](https://tauri.app/)（Rust） |
| 后端 | [Go](https://go.dev/) + net/http |
| 前端 | [Vue 3](https://vuejs.org/) + Pinia + Vite |
| API | MiMo OpenAI 兼容 API |
| 图标 | Pillow (Python) 生成 |

---

## License

MIT License

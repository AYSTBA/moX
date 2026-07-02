# Changelog

## [0.3.0] - 2026-07-02

### Added
- Markdown 渲染 + 代码语法高亮（marked + highlight.js + markedHighlight）
- 代码块 dark 主题（github-dark.css）、表格、引用块、列表、标题样式
- 设置面板系统提示词编辑（textarea，默认 "你是mimo"）
- 附件按钮（placeholder，位于发送按钮上方纵向排列）
- 磨砂玻璃效果（backdrop-filter: blur）对话框

### Changed
- AI 消息气泡 max-width 75% → 65%
- 消息顶部留白 80px → 60px，气泡间距 24px → 10px
- 输入框最小高度 64px → 76px，底部留白 20px → 28px
- v-text → v-html 渲染 Markdown 内容
- 附件按钮从输入框左侧移至发送按钮上方（input-actions 纵向排列）
- 设置对话框取消按钮改为直接关闭（showSettings = false）

### Fixed
- 应用图标右下角圆点缩小（r=48→32）并右移（cx=352→420）
- AI 头像 SVG 从 16px 放大到 22px
- 对话框按钮重置操作消除 saveKey 覆盖问题

## [0.2.0] - 2026-07-02

### Added
- Tauri v2 桌面壳（单一 exe 启动）
- Go 后端作为 Tauri sidecar 自动拉起/关闭
- 多对话管理（localStorage 持久化，新建/切换/重命名/删除）
- 应用图标（MOX Logo，Pillow 生成，256px .ico + 多尺寸 .png）
- 后端启动前 kill_stale_backend() 清理残留进程
- wait_for_backend_ready() 轮询 port 3099 直到就绪

### Changed
- Wails → Tauri 架构迁移
- 前端 UI 重构（侧边栏 + 主聊天区 + 设置弹窗，AI Studio 风格）
- 气泡发送前自动清理残留进程
- Cargo.toml release profile: lto=false + codegen-units=1 修复 LTO 链接错误
- macOS/Linux 侧使用 cfg(windows) 条件编译

### Fixed
- 后端残留进程锁住 3099 端口导致新进程静默退出
- Rust LTO 链接错误（58 个未解析外部符号）
- 命令行窗口拖尾（#![windows_subsystem = "windows"]）
- find_server() 路径检查（优先同目录 mox-server.exe）

### Removed
- 旧版 Wails 代码和构建产物
- 冲突副本文件（chat(冲突副本).js + 4 个 gen/ 构建产物）

## [0.1.0] - 2026-07-01

### Added
- 基础聊天界面（侧边栏、对话列表、消息气泡）
- 多模型支持（MiMo-V2.5 Pro / V2.5 / V2 Pro / Flash / Omni）
- 流式输出（思考过程、Token、工具调用、搜索来源）
- 自定义系统提示词、Temperature / Top P / Max Tokens
- 主题切换（深色/浅色）
- 个性化设置（背景图片 + 磨砂玻璃 + 雾面强度滑块）
- 一键复制回复内容、Token / 用时信息 tooltip
- 思考时间显示 [X.Xs]

### Fixed
- 附件上传后实际发送到 API（base64 → multimodal content parts）
- 消息气泡中图片/视频内嵌显示
- 雾面强度滑块即时响应（@input 直接更新 CSS）
- 个人信息 tooltip 显示真实 Token 和用时
- 思考过程改为无框文字 + hover 变色
- 流式响应添加 30s 心跳保活
- 移除冗余 sb.WriteString 调用

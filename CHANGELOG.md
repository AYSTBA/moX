# Changelog

## [0.1.0] - 2026-07-01

### Added
- 基础聊天界面（侧边栏、对话列表、消息气泡）
- 多模型支持（MiMo-V2.5 Pro / V2.5 / V2 Pro / Flash / Omni）
- 流式输出（思考过程、Token、工具调用、搜索来源）
- Agent 规划 + Tavily 联网搜索
- 消息附件（图片/视频内嵌显示、lightbox 放大）
- 自定义系统提示词、Temperature / Top P / Max Tokens
- 主题切换（深色/浅色）
- 个性化设置（背景图片 + 磨砂玻璃 + 雾面强度滑块）
- 一键复制回复内容、Token / 用时信息 tooltip
- 思考时间显示 `[X.Xs]`

### Fixed
- 附件上传后实际发送到 API（base64 → multimodal content parts）
- 消息气泡中图片/视频内嵌显示
- 雾面强度滑块即时响应（`@input` 直接更新 CSS）
- 个人信息 tooltip 显示真实 Token 和用时
- 思考过程改为无框文字 + hover 变色
- cancelFunc 加 `sync.Mutex` 防止 goroutine 泄漏
- Agent 多模态消息 Content 类型断言容错
- Agent 决策改为 JSON 结构化输出
- 会话列表按 `UpdatedAt` 降序排列
- `cleanFunctionCalls` 改用正则避免 JSON 碎片残留
- 流式响应添加 30s 心跳保活
- 移除冗余 `sb.WriteString` 调用

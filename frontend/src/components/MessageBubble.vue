<script setup>
import {computed, ref, onMounted} from 'vue'
import {marked} from 'marked'
import hljs from 'highlight.js'

const props = defineProps({
  message: Object,
  streaming: Boolean,
})

const showThinking = ref(false)
const showAnnotations = ref(false)
const lightboxSrc = ref(null)

function isImage(att) {
  return att.type && att.type.startsWith('image/')
}

function isVideo(att) {
  return att.type && att.type.startsWith('video/')
}

marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(code, {language: lang}).value
      } catch {}
    }
    return code
  },
  breaks: true,
  gfm: true,
})

const renderedContent = computed(() => {
  if (!props.message.content) return ''
  return marked.parse(props.message.content)
})

const hasThinking = computed(() => {
  return props.message.reasoning_content && props.message.reasoning_content.trim()
})

const hasToolCalls = computed(() => {
  return props.message.tool_calls && props.message.tool_calls.length > 0
})

const hasAnnotations = computed(() => {
  return props.message.annotations && props.message.annotations.length > 0
})

function copyCode(e) {
  const btn = e.target.closest('.code-copy-btn')
  const code = btn?.parentElement?.querySelector('code')
  if (code) {
    navigator.clipboard.writeText(code.textContent)
    btn.textContent = 'OK'
    setTimeout(() => btn.textContent = '复制', 1500)
  }
}

onMounted(() => {
  showThinking.value = !hasThinking.value ? false : showThinking.value
})
</script>

<template>
  <div class="message" :class="[message.role, {streaming}]">
    <div class="message-avatar">
      <span v-if="message.role === 'user'" class="avatar-user">U</span>
      <span v-else class="avatar-bot">M</span>
    </div>
    <div class="message-body">
      <div v-if="hasThinking" class="thinking-block">
        <button class="thinking-toggle" @click="showThinking = !showThinking">
          <svg class="thinking-icon" :class="{open: showThinking}" width="10" height="10" viewBox="0 0 24 24" fill="currentColor"><polygon points="5 3 19 12 5 21 5 3"/></svg>
          <span>思考过程</span>
          <span v-if="streaming" class="thinking-dots">...</span>
        </button>
        <div v-if="showThinking" class="thinking-content">
          {{ message.reasoning_content }}
        </div>
      </div>

      <div v-if="hasToolCalls" class="toolcalls-block">
        <div v-for="tc in message.tool_calls" :key="tc.id || tc.function?.name" class="toolcall-item">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
          <span class="toolcall-name">{{ tc.function?.name }}</span>
          <code class="toolcall-args">{{ tc.function?.arguments }}</code>
        </div>
      </div>

      <div v-if="hasAnnotations" class="annotations-block">
        <div class="annotations-header" @click="showAnnotations = !showAnnotations">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/></svg>
          <span>搜索来源 ({{ message.annotations.length }})</span>
          <svg class="annotations-arrow" :class="{open: showAnnotations}" width="10" height="10" viewBox="0 0 24 24" fill="currentColor"><polygon points="5 3 19 12 5 21 5 3"/></svg>
        </div>
        <div v-if="showAnnotations" class="annotations-list">
          <a v-for="(ann, i) in message.annotations" :key="i" :href="ann.url" target="_blank" class="annotation-item">
            <span class="ann-title">{{ ann.title }}</span>
            <span class="ann-site">{{ ann.site_name }}</span>
          </a>
        </div>
      </div>

      <div v-if="message.attachments && message.attachments.length" class="attachments-block">
        <div v-for="(att, i) in message.attachments" :key="i" class="attachment-item">
          <img v-if="isImage(att)" :src="att.data" class="attachment-media" @click="lightboxSrc = att.data" />
          <video v-else-if="isVideo(att)" :src="att.data" controls preload="metadata" class="attachment-media"></video>
          <div v-else class="attachment-file">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
            <span>{{ att.name }}</span>
          </div>
        </div>
      </div>

      <div v-if="message.content" class="markdown-body" v-html="renderedContent" @click="copyCode"></div>

      <div v-if="streaming && !message.content && !hasThinking" class="typing-indicator">
        <span></span><span></span><span></span>
      </div>
    </div>
  <div v-if="lightboxSrc" class="lightbox-overlay" @click="lightboxSrc = null">
    <img :src="lightboxSrc" class="lightbox-img" @click.stop />
    <button class="lightbox-close" @click="lightboxSrc = null">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
    </button>
  </div>
  </div>
</template>

<style scoped>
.message {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.message.user .message-avatar {
  background: var(--text-primary);
  color: var(--bg-primary);
  border: none;
}

.avatar-user, .avatar-bot {
  font-family: 'Smiley Sans', sans-serif;
}

.message-body {
  max-width: 85%;
  min-width: 0;
}

.message.user .message-body {
  background: var(--bg-message-user);
  padding: 10px 16px;
  border-radius: 12px 12px 2px 12px;
  border: 1px solid var(--border-color);
}

.message.assistant .message-body {
  padding: 4px 0;
}

.thinking-block {
  margin-bottom: 8px;
}

.thinking-toggle {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: var(--thinking-bg);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  padding: 4px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.thinking-toggle:hover {
  border-color: var(--text-muted);
  color: var(--text-primary);
}

.thinking-icon {
  transition: transform 0.2s;
}

.thinking-icon.open {
  transform: rotate(90deg);
}

.thinking-dots {
  animation: blink 1s infinite;
}

@keyframes blink {
  50% { opacity: 0.3; }
}

.thinking-content {
  margin-top: 8px;
  padding: 12px 16px;
  background: var(--thinking-bg);
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
  max-height: 400px;
  overflow-y: auto;
  border-left: 2px solid var(--border-color);
}

.toolcalls-block {
  margin-bottom: 8px;
}

.toolcall-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: var(--thinking-bg);
  border-radius: 4px;
  margin-bottom: 4px;
  font-size: 12px;
  border: 1px solid var(--border-color);
}

.toolcall-name {
  font-weight: 600;
  color: var(--text-primary);
  font-family: monospace;
}

.toolcall-args {
  color: var(--text-muted);
  font-size: 11px;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.annotations-block {
  margin-bottom: 8px;
}

.annotations-header {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  background: var(--thinking-bg);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.annotations-header:hover {
  border-color: var(--text-muted);
  color: var(--text-primary);
}

.annotations-arrow {
  transition: transform 0.2s;
}

.annotations-arrow.open {
  transform: rotate(90deg);
}

.annotations-list {
  margin-top: 6px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.annotation-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  background: var(--thinking-bg);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  text-decoration: none;
  color: var(--text-primary);
  font-size: 12px;
  transition: background 0.15s;
}

.annotation-item:hover {
  background: var(--bg-hover);
}

.ann-title {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ann-site {
  color: var(--text-muted);
  font-size: 11px;
  flex-shrink: 0;
}

.attachments-block {
  margin-bottom: 8px;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: var(--thinking-bg);
  border-radius: 4px;
  margin-bottom: 4px;
  font-size: 12px;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 8px 0;
}

.typing-indicator span {
  width: 5px;
  height: 5px;
  background: var(--text-muted);
  border-radius: 50%;
  animation: bounce 1.4s infinite;
}

.typing-indicator span:nth-child(2) { animation-delay: 0.2s; }
.typing-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}
.attachments-block {
  margin-bottom: 8px;
}

.attachment-item {
  margin-bottom: 6px;
}

.attachment-item:last-child {
  margin-bottom: 0;
}

.attachment-media {
  max-width: 100%;
  max-height: 400px;
  border-radius: 8px;
  cursor: pointer;
  display: block;
  background: var(--bg-tertiary);
  object-fit: cover;
}

.attachment-media:hover {
  opacity: 0.95;
}

video.attachment-media {
  max-height: 350px;
  cursor: default;
  background: #000;
}

.attachment-file {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: var(--thinking-bg);
  border-radius: 4px;
  font-size: 12px;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}
</style>

<style>
.markdown-body {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-primary);
  word-wrap: break-word;
}

.markdown-body p {
  margin-bottom: 12px;
}

.markdown-body p:last-child {
  margin-bottom: 0;
}

.markdown-body pre {
  background: var(--code-bg);
  border-radius: 6px;
  padding: 16px;
  overflow-x: auto;
  margin: 12px 0;
  position: relative;
  border: 1px solid var(--border-color);
}

.markdown-body pre code {
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.5;
}

.markdown-body code {
  background: var(--code-bg);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  border: 1px solid var(--border-color);
}

.markdown-body pre code {
  background: none;
  padding: 0;
  border: none;
}

.markdown-body .code-copy-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  padding: 2px 8px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 11px;
  opacity: 0;
  transition: opacity 0.2s;
}

.markdown-body pre:hover .code-copy-btn {
  opacity: 1;
}

.markdown-body .code-copy-btn:hover {
  background: var(--accent);
  color: var(--bg-primary);
}

.markdown-body ul, .markdown-body ol {
  padding-left: 24px;
  margin: 8px 0;
}

.markdown-body li {
  margin-bottom: 4px;
}

.markdown-body blockquote {
  border-left: 2px solid var(--border-color);
  padding-left: 12px;
  color: var(--text-secondary);
  margin: 12px 0;
}

.markdown-body table {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}

.markdown-body th, .markdown-body td {
  border: 1px solid var(--border-color);
  padding: 8px 12px;
  text-align: left;
}

.markdown-body th {
  background: var(--bg-tertiary);
  font-weight: 600;
}

.markdown-body h1, .markdown-body h2, .markdown-body h3 {
  margin: 16px 0 8px;
  font-weight: 600;
  font-family: 'Smiley Sans', sans-serif;
}

.markdown-body a {
  color: var(--text-primary);
  text-decoration: underline;
}

.markdown-body a:hover {
  opacity: 0.7;
}

.markdown-body hr {
  border: none;
  border-top: 1px solid var(--border-color);
  margin: 16px 0;
}
.attachments-block {
  margin-bottom: 8px;
}

.attachment-item {
  margin-bottom: 6px;
}

.attachment-item:last-child {
  margin-bottom: 0;
}

.attachment-media {
  max-width: 100%;
  max-height: 400px;
  border-radius: 8px;
  cursor: pointer;
  display: block;
  background: var(--bg-tertiary);
  object-fit: cover;
}

.attachment-media:hover {
  opacity: 0.95;
}

video.attachment-media {
  max-height: 350px;
  cursor: default;
  background: #000;
}

.attachment-file {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: var(--thinking-bg);
  border-radius: 4px;
  font-size: 12px;
  color: var(--text-secondary);
  border: 1px solid var(--border-color);
}

.lightbox-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  cursor: zoom-out;
}

.lightbox-img {
  max-width: 90%;
  max-height: 90%;
  object-fit: contain;
  border-radius: 4px;
  cursor: default;
}

.lightbox-close {
  position: fixed;
  top: 16px;
  right: 16px;
  background: rgba(255,255,255,0.15);
  border: none;
  color: #fff;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.lightbox-close:hover {
  background: rgba(255,255,255,0.3);
}
</style>



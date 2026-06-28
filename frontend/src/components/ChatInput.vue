<script setup>
import {ref, nextTick} from 'vue'

const props = defineProps({
  disabled: Boolean,
  isStreaming: Boolean,
})

const emit = defineEmits(['send', 'stop'])

const input = ref('')
const textareaRef = ref(null)
const attachments = ref([])
const showAttachMenu = ref(false)
const fileInputRef = ref(null)

function handleSend() {
  if (!input.value.trim() && attachments.value.length === 0) return
  if (props.disabled) return
  emit('send', input.value.trim(), attachments.value.length > 0 ? [...attachments.value] : undefined)
  input.value = ''
  attachments.value = []
  nextTick(() => {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
    }
  })
}

function handleKeydown(e) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSend()
  }
}

function autoResize() {
  const el = textareaRef.value
  if (!el) return
  el.style.height = 'auto'
  el.style.height = Math.min(el.scrollHeight, 300) + 'px'
}

function triggerFileInput(type) {
  showAttachMenu.value = false
  if (type === 'file') {
    fileInputRef.value?.click()
  }
}

function onFileSelected(e) {
  const files = e.target.files
  if (!files) return
  for (const f of files) {
    attachments.value.push({
      name: f.name,
      size: f.size,
      type: f.type,
      file: f,
    })
  }
  e.target.value = ''
}

function removeAttachment(index) {
  attachments.value.splice(index, 1)
}

function formatSize(bytes) {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1048576).toFixed(1) + ' MB'
}

function getFileIcon(type) {
  if (type.startsWith('image/')) return 'IMG'
  if (type.startsWith('video/')) return 'VID'
  if (type.startsWith('audio/')) return 'AUD'
  return 'FILE'
}
</script>

<template>
  <div class="chat-input-area">
    <div class="input-wrapper">
      <div v-if="attachments.length > 0" class="attachments-preview">
        <div v-for="(att, i) in attachments" :key="i" class="att-chip">
          <span class="att-icon">{{ getFileIcon(att.type) }}</span>
          <span class="att-name">{{ att.name }}</span>
          <span class="att-size">{{ formatSize(att.size) }}</span>
          <button class="att-remove" @click="removeAttachment(i)">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
      </div>

      <div class="input-container">
        <div class="attach-area">
          <button class="btn-attach" @click="showAttachMenu = !showAttachMenu" title="添加附件">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
          </button>
          <div v-if="showAttachMenu" class="attach-menu">
            <button class="attach-option" @click="triggerFileInput('file')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="13 2 13 9 20 9"/></svg>
              <span>文件</span>
            </button>
            <button class="attach-option" @click="triggerFileInput('file')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
              <span>图片</span>
            </button>
            <button class="attach-option" @click="triggerFileInput('file')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>
              <span>视频</span>
            </button>
            <button class="attach-option" @click="triggerFileInput('file')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 18V5l12-2v13"/><circle cx="6" cy="18" r="3"/><circle cx="18" cy="16" r="3"/></svg>
              <span>音频</span>
            </button>
          </div>
        </div>

        <textarea
          ref="textareaRef"
          v-model="input"
          :disabled="disabled"
          placeholder="输入消息... (Shift+Enter 换行)"
          rows="3"
          @keydown="handleKeydown"
          @input="autoResize"
          @blur="showAttachMenu = false"
        ></textarea>

        <div class="input-actions">
          <button
            v-if="isStreaming"
            class="btn-stop"
            @click="emit('stop')"
            title="停止生成"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
              <rect x="6" y="6" width="12" height="12" rx="1"/>
            </svg>
          </button>
          <button
            v-else
            class="btn-send"
            :disabled="!input.trim() && attachments.length === 0"
            @click="handleSend"
            title="发送"
          >
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="22" y1="2" x2="11" y2="13"/>
              <polygon points="22 2 15 22 11 13 2 9 22 2"/>
            </svg>
          </button>
        </div>
      </div>

      <input
        ref="fileInputRef"
        type="file"
        multiple
        accept="image/*,video/*,audio/*,.pdf,.doc,.docx,.txt,.csv,.json,.xml,.md"
        style="display:none"
        @change="onFileSelected"
      />
    </div>
  </div>
</template>

<style scoped>
.chat-input-area {
  padding: 16px 20px;
  background: var(--bg-primary);
}

.input-wrapper {
  max-width: 800px;
  margin: 0 auto;
}

.attachments-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 8px;
}

.att-chip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 12px;
}

.att-icon {
  font-family: monospace;
  font-size: 10px;
  font-weight: 700;
  color: var(--text-muted);
  padding: 1px 4px;
  border: 1px solid var(--border-color);
  border-radius: 2px;
}

.att-name {
  color: var(--text-primary);
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.att-size {
  color: var(--text-muted);
  font-size: 10px;
}

.att-remove {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0;
  display: flex;
  align-items: center;
}

.att-remove:hover {
  color: var(--text-primary);
}

.input-container {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 10px 12px;
  transition: border-color 0.2s;
}

.input-container:focus-within {
  border-color: var(--text-muted);
}

.attach-area {
  position: relative;
  flex-shrink: 0;
}

.btn-attach {
  width: 36px;
  height: 36px;
  background: none;
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.btn-attach:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.attach-menu {
  position: absolute;
  bottom: calc(100% + 6px);
  left: 0;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 4px;
  min-width: 140px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.2);
  z-index: 100;
}

.attach-option {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 12px;
  background: none;
  border: none;
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
  border-radius: 4px;
  text-align: left;
}

.attach-option:hover {
  background: var(--bg-hover);
}

textarea {
  flex: 1;
  background: none;
  border: none;
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  outline: none;
  min-height: 48px;
  max-height: 300px;
  font-family: inherit;
  padding: 6px 0;
}

textarea::placeholder {
  color: var(--text-muted);
}

textarea:disabled {
  opacity: 0.5;
}

.input-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.btn-send, .btn-stop {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.btn-send {
  background: var(--text-primary);
  color: var(--bg-primary);
}

.btn-send:hover:not(:disabled) {
  opacity: 0.8;
  transform: scale(1.05);
}

.btn-send:disabled {
  opacity: 0.2;
  cursor: not-allowed;
}

.btn-stop {
  background: var(--text-primary);
  color: var(--bg-primary);
  animation: pulse 1.5s infinite;
}

.btn-stop:hover {
  opacity: 0.8;
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(255,255,255,0.3); }
  50% { box-shadow: 0 0 0 6px rgba(255,255,255,0); }
}
</style>

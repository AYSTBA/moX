<script setup>
import {ref} from 'vue'
import {useChatStore} from '../stores/chat'
import {useSettingsStore} from '../stores/settings'

const chat = useChatStore()
const settings = useSettingsStore()
const editingKey = ref('')
const editLabel = ref('')

function startRename(s) {
  editingKey.value = s.key
  editLabel.value = s.label
}

function finishRename(s) {
  if (editLabel.value.trim()) {
    chat.renameSession(s.key, editLabel.value.trim())
  }
  editingKey.value = ''
}

function formatDate(ts) {
  const d = new Date(ts)
  return `${d.getMonth()+1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2,'0')}`
}
</script>

<template>
  <div class="sidebar">
    <div class="sidebar-header">
      <div class="logo">
        <span class="logo-text">MOX</span>
      </div>
      <button class="btn-icon" @click="chat.newSession()" title="新对话">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
      </button>
    </div>

    <div class="session-list">
      <div
        v-for="s in chat.sessions"
        :key="s.key"
        class="session-item"
        :class="{active: s.key === chat.activeSessionKey}"
        @click="chat.selectSession(s.key)"
      >
        <div class="session-content" v-if="editingKey !== s.key">
          <div class="session-label">{{ s.label }}</div>
          <div class="session-time">{{ formatDate(s.updated_at) }}</div>
        </div>
        <input
          v-else
          v-model="editLabel"
          class="session-edit"
          @blur="finishRename(s)"
          @keydown.enter="finishRename(s)"
          @keydown.escape="editingKey = ''"
          autofocus
        />
        <div class="session-actions">
          <button class="btn-tiny" @click.stop="startRename(s)" title="重命名">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
          </button>
          <button class="btn-tiny" @click.stop="chat.removeSession(s.key)" title="删除">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
          </button>
        </div>
      </div>
      <div v-if="chat.sessions.length === 0" class="empty-hint">
        点击 + 开始新对话
      </div>
    </div>

    <div class="sidebar-footer">
      <button class="btn-settings" @click="settings.showSettings = true">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
        <span>设置</span>
      </button>
      <button class="btn-theme" @click="settings.toggleTheme()" :title="settings.settings.theme === 'dark' ? '切换亮色' : '切换暗色'">
        <svg v-if="settings.settings.theme === 'dark'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.sidebar {
  width: 260px;
  min-width: 260px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.logo-text {
  font-family: 'Smiley Sans', sans-serif;
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 1px;
}

.btn-icon {
  background: none;
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  cursor: pointer;
  padding: 6px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.btn-icon:hover {
  background: var(--bg-hover);
}

.session-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.session-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  margin-bottom: 2px;
  transition: background 0.15s;
}

.session-item:hover {
  background: var(--bg-hover);
}

.session-item.active {
  background: var(--bg-tertiary);
}

.session-content {
  flex: 1;
  min-width: 0;
}

.session-label {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.session-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
}

.session-edit {
  flex: 1;
  background: var(--bg-input);
  border: 1px solid var(--accent);
  color: var(--text-primary);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 13px;
  outline: none;
}

.session-actions {
  display: none;
  gap: 2px;
  margin-left: 4px;
}

.session-item:hover .session-actions {
  display: flex;
}

.btn-tiny {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s;
}

.btn-tiny:hover {
  color: var(--text-primary);
}

.empty-hint {
  text-align: center;
  color: var(--text-muted);
  padding: 40px 20px;
  font-size: 13px;
}

.sidebar-footer {
  padding: 12px;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 8px;
}

.btn-settings {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: background 0.2s;
}

.btn-settings:hover {
  background: var(--bg-hover);
}

.btn-theme {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-theme:hover {
  background: var(--bg-hover);
}
</style>

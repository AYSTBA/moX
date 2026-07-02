<script setup>
import {ref, onMounted, watch, nextTick} from "vue"
import {useChatStore} from "./stores/chat.js"
import {getKey, setKey} from "./api.js"

const chat = useChatStore()
const input = ref("")
const apiKey = ref(getKey())
const showSettings = ref(false)
const streaming = ref(false)
const sidebarOpen = ref(true)
const messagesContainer = ref(null)
const editingTitle = ref(null)
const editTitleValue = ref("")
const welcomeModel = ref("mimo-v2.5")

onMounted(async () => {
  chat.loadConversations()
  await chat.loadModels()
  if (!apiKey.value) showSettings.value = true
  if (!chat.currentId && chat.conversations.length === 0) {
    chat.newConversation()
  }
})

watch(() => chat.currentMessages.length, () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
})

function getFirstMessagePreview(conv) {
  const userMsg = conv.messages.find((m) => m.role === "user")
  if (userMsg) return userMsg.content.slice(0, 40) + (userMsg.content.length > 40 ? "..." : "")
  return conv.title
}

function formatTime(ts) {
  const d = new Date(ts)
  const now = new Date()
  const diff = now - d
  if (diff < 86400000) return d.toLocaleTimeString("zh-CN", {hour: "2-digit", minute: "2-digit"})
  if (diff < 604800000) {
    const days = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"]
    return days[d.getDay()]
  }
  return d.toLocaleDateString("zh-CN", {month: "2-digit", day: "2-digit"})
}

async function handleSend() {
  const text = input.value.trim()
  if (!text) return
  input.value = ""

  if (!chat.currentId) chat.newConversation()
  const conv = chat.currentConversation
  if (!conv) return

  conv.messages.push({role: "user", content: text})
  if (conv.messages.filter((m) => m.role === "user").length === 1) {
    conv.title = text.slice(0, 24) + (text.length > 24 ? "..." : "")
  }
  conv.updatedAt = Date.now()
  chat.persistConversations()

  streaming.value = true
  const msg = {role: "assistant", content: "", streaming: true}
  conv.messages.push(msg)

  try {
    const history = conv.messages.filter((m) => !m.streaming).map((m) => ({role: m.role, content: m.content}))
    msg.content = await chat.send(history, conv.model || "mimo-v2.5")
    msg.streaming = false
    chat.persistConversations()
  } catch (e) {
    msg.content = "Error: " + e.message
    msg.streaming = false
  }
  streaming.value = false
}

function newChat() {
  const id = chat.newConversation()
  const conv = chat.conversations.find((c) => c.id === id)
  if (conv) conv.model = welcomeModel.value
}

function selectChat(id) {
  chat.selectConversation(id)
}

function deleteChat(id, e) {
  e.stopPropagation()
  chat.removeConversation(id)
}

function saveKey() {
  setKey(apiKey.value)
  showSettings.value = false
}

function startEditTitle(conv, e) {
  e.stopPropagation()
  editingTitle.value = conv.id
  editTitleValue.value = conv.title
}

function saveTitle(conv) {
  conv.title = editTitleValue.value || "新对话"
  conv.updatedAt = Date.now()
  chat.persistConversations()
  editingTitle.value = null
}

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value
}
</script>

<template>
  <div class="app-layout" :class="{sidebarCollapsed: !sidebarOpen}">
    <!-- Left Sidebar -->
    <aside class="sidebar" :class="{open: sidebarOpen}">
      <div class="sidebar-header">
        <div class="logo-area">
          <svg class="logo-icon" width="22" height="22" viewBox="0 0 48 48" fill="none">
            <rect x="6" y="6" width="36" height="36" rx="8" fill="#e8e8e8"/>
            <path d="M16 20h16M16 28h12" stroke="#111" stroke-width="3" stroke-linecap="round"/>
            <circle cx="34" cy="32" r="4" fill="#111"/>
          </svg>
          <span class="logo-text">MOX</span>
        </div>
      </div>

      <button class="btn-new-chat" @click="newChat">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        新对话
      </button>

      <div class="conversation-list">
        <div
          v-for="conv in chat.conversations"
          :key="conv.id"
          class="conversation-item"
          :class="{active: conv.id === chat.currentId}"
          @click="selectChat(conv.id)"
        >
          <div class="conv-content">
            <div class="conv-title-row">
              <template v-if="editingTitle === conv.id">
                <input
                  class="conv-title-input"
                  v-model="editTitleValue"
                  @blur="saveTitle(conv)"
                  @keydown.enter="saveTitle(conv)"
                  @click.stop
                  autofocus
                />
              </template>
              <template v-else>
                <span class="conv-title" @dblclick="startEditTitle(conv, $event)">
                  {{ getFirstMessagePreview(conv) }}
                </span>
              </template>
            </div>
            <div class="conv-meta">{{ formatTime(conv.updatedAt) }}</div>
          </div>
          <button class="btn-delete" @click="deleteChat(conv.id, $event)" title="删除">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
          </button>
        </div>
        <div v-if="chat.conversations.length === 0" class="conv-empty">
          暂无对话
        </div>
      </div>

      <div class="sidebar-footer">
        <button class="btn-sidebar-footer" @click="showSettings = true">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
          </svg>
          设置
        </button>
      </div>
    </aside>

    <!-- Sidebar toggle -->
    <button class="btn-toggle-sidebar" @click="toggleSidebar" :title="sidebarOpen ? '收起' : '展开'">
      <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline :points="sidebarOpen ? '15 18 9 12 15 6' : '9 18 15 12 9 6'"/>
      </svg>
    </button>

    <!-- Main Area -->
    <main class="main-area">
      <div class="messages" ref="messagesContainer">
        <!-- Welcome -->
        <div v-if="!chat.currentConversation || chat.currentMessages.length === 0" class="welcome">
          <div class="welcome-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <rect x="6" y="6" width="36" height="36" rx="8" fill="#2a2a2a"/>
              <path d="M16 20h16M16 28h12" stroke="#e8e8e8" stroke-width="3" stroke-linecap="round"/>
              <circle cx="34" cy="32" r="4" fill="#e8e8e8"/>
            </svg>
          </div>
          <h1>MOX</h1>
          <p class="welcome-subtitle">MiMo AI Client</p>
          <div class="welcome-actions">
            <select v-model="welcomeModel" class="model-select">
              <option v-for="m in chat.models" :key="m.id" :value="m.id">{{ m.id }}</option>
            </select>
          </div>
        </div>

        <!-- Messages -->
        <div v-for="(msg, i) in chat.currentMessages" :key="i" class="message" :class="msg.role">
          <div class="avatar" :class="msg.role">
            <svg v-if="msg.role === 'user'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/>
            </svg>
            <svg v-else width="16" height="16" viewBox="0 0 48 48" fill="none">
              <rect x="6" y="6" width="36" height="36" rx="8" fill="#e8e8e8"/>
              <path d="M16 20h16M16 28h12" stroke="#111" stroke-width="3" stroke-linecap="round"/>
              <circle cx="34" cy="32" r="4" fill="#111"/>
            </svg>
          </div>
          <div class="bubble">
            <div v-if="msg.streaming && !msg.content" class="loading-dots">
              <span></span><span></span><span></span>
            </div>
            <div v-else class="bubble-content" v-text="msg.content"></div>
          </div>
        </div>
      </div>

      <!-- Input -->
      <div class="input-area">
        <div class="model-bar">
          <select v-if="chat.currentConversation" v-model="chat.currentConversation.model" class="model-select">
            <option v-for="m in chat.models" :key="m.id" :value="m.id">{{ m.id }}</option>
          </select>
        </div>
        <div class="input-wrap">
          <textarea
            v-model="input"
            @keydown.enter.exact="handleSend"
            rows="1"
            placeholder="输入消息... (Enter 发送)"
            :disabled="streaming"
          ></textarea>
          <button class="btn-send" @click="handleSend" :disabled="!input.trim() || streaming">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/>
            </svg>
          </button>
        </div>
      </div>
    </main>

    <!-- Settings Dialog -->
    <div v-if="showSettings" class="dialog-overlay" @click.self="saveKey">
      <div class="dialog-panel">
        <div class="dialog-header">
          <h3>设置</h3>
          <button class="btn-close" @click="saveKey">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>MiMo API Key</label>
            <input type="password" v-model="apiKey" placeholder="输入 API Key" />
            <p class="form-hint">API Key 存储在浏览器本地，不会上传到服务器</p>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-cancel" @click="saveKey">取消</button>
          <button class="btn-primary" @click="saveKey">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>


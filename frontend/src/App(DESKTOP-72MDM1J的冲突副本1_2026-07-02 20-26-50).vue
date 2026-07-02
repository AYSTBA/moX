<script setup>
import {ref, onMounted} from 'vue'
import {useChatStore} from './stores/chat.js'
import {getKey, setKey} from './api.js'

const chat = useChatStore()
const messages = ref([])
const input = ref('')
const model = ref('mimo-v2.5')
const apiKey = ref(getKey())
const showSettings = ref(false)
const streaming = ref(false)

onMounted(async () => {
  await chat.loadModels()
  if (!apiKey.value) showSettings.value = true
})

async function handleSend() {
  const text = input.value.trim()
  if (!text) return
  input.value = ''
  messages.value.push({role: 'user', content: text})
  streaming.value = true
  messages.value.push({role: 'assistant', content: '', streaming: true})
  try {
    const lastMsg = messages.value[messages.value.length - 1]
    await chat.send(
      messages.value.filter(m => !m.streaming).map(m => ({role: m.role, content: m.content})),
      model.value
    )
    lastMsg.content = chat.streamingContent
    lastMsg.streaming = false
  } catch (e) {
    messages.value[messages.value.length - 1] = {role: 'assistant', content: 'Error: ' + e.message}
  }
  streaming.value = false
}

function saveKey() {
  setKey(apiKey.value)
  showSettings.value = false
}
</script>

<template>
  <div class="app-layout">
    <div class="sidebar">
      <div style="padding:16px;font-weight:600;font-size:16px;border-bottom:1px solid #333">MOX</div>
      <div style="flex:1"></div>
    </div>

    <div class="chat-area">
      <div class="model-bar">
        <select v-model="model">
          <option v-for="m in chat.models" :key="m.id" :value="m.id">{{ m.id }}</option>
        </select>
        <button class="btn-settings" @click="showSettings = true">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
        </button>
      </div>

      <div class="messages">
        <div v-if="messages.length === 0" class="welcome">
          <h1>MOX</h1>
          <p>MiMo AI Client</p>
        </div>

        <div v-for="(msg, i) in messages" :key="i" class="message" :class="msg.role">
          <div class="avatar">{{ msg.role === 'user' ? 'U' : 'M' }}</div>
          <div class="bubble">
            <div v-if="msg.streaming && !msg.content" class="loading-dots">
              <span></span><span></span><span></span>
            </div>
            <div v-else class="bubble-content" v-html="msg.content"></div>
          </div>
        </div>
      </div>

      <div class="input-area">
        <div class="input-wrap">
          <textarea v-model="input" @keydown.enter.exact="handleSend" rows="1" placeholder="输入消息... (Enter 发送)" :disabled="streaming"></textarea>
          <button class="btn-send" @click="handleSend" :disabled="!input.trim() || streaming">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
          </button>
        </div>
      </div>
    </div>

    <div v-if="showSettings" class="settings-overlay" @click.self="saveKey">
      <div class="settings-panel">
        <button class="btn-close" @click="saveKey">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
        <h3>设置</h3>
        <label>MiMo API Key</label>
        <input type="password" v-model="apiKey" placeholder="输入 API Key" />
        <button class="btn-save" @click="saveKey">保存</button>
      </div>
    </div>
  </div>
</template>

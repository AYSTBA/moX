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
  const msg = {role: 'assistant', content: '', streaming: true}
  messages.value.push(msg)
  try {
    const history = messages.value.filter(m => !m.streaming).map(m => ({role: m.role, content: m.content}))
    await chat.send(history, model.value)
    msg.content = chat.streamingContent
    msg.streaming = false
  } catch (e) {
    msg.content = 'Error: ' + e.message
    msg.streaming = false
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
    <div class="messages">
      <div v-if="messages.length === 0" class="welcome">
        <h1>MOX</h1>
        <p style="color:#666;font-size:14px">MiMo AI Client</p>
        <select v-model="model" style="margin-top:12px;background:#1a1a1a;border:1px solid #333;color:#e8e8e8;padding:4px 10px;border-radius:4px;font-size:13px;outline:none">
          <option v-for="m in chat.models" :key="m.id" :value="m.id">{{ m.id }}</option>
        </select>
      </div>

      <div v-for="(msg, i) in messages" :key="i" class="message" :class="msg.role">
        <div class="avatar">{{ msg.role === 'user' ? 'U' : 'M' }}</div>
        <div class="bubble">
          <div v-if="msg.streaming && !msg.content" class="loading-dots">
            <span></span><span></span><span></span>
          </div>
          <div v-else class="bubble-content">{{ msg.content }}</div>
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
      <div style="max-width:800px;margin:4px auto 0;display:flex;justify-content:space-between">
        <span style="font-size:11px;color:#666">MiMo {{ model }}</span>
        <button style="font-size:11px;color:#666;background:none;border:none;cursor:pointer" @click="showSettings = true">设置</button>
      </div>
    </div>

    <div v-if="showSettings" class="settings-overlay" @click.self="saveKey">
      <div class="settings-panel">
        <h3>设置</h3>
        <label>MiMo API Key</label>
        <input type="password" v-model="apiKey" placeholder="输入 API Key" />
        <button class="btn-save" @click="saveKey">保存</button>
      </div>
    </div>
  </div>
</template>
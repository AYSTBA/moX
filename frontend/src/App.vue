<script setup>
import {onMounted} from 'vue'
import {useSettingsStore} from './stores/settings'
import {useChatStore} from './stores/chat'
import Sidebar from './components/Sidebar.vue'
import ChatView from './components/ChatView.vue'
import Settings from './components/Settings.vue'

const settings = useSettingsStore()
const chat = useChatStore()

onMounted(async () => {
  await settings.load()
  await chat.loadSessions()
  await chat.loadModels()
  chat.setupEvents()
})
</script>

<template>
  <div class="app-layout" :class="settings.settings.theme">
    <Sidebar />
    <ChatView />
    <Settings v-if="settings.showSettings" />
    <transition name="toast">
      <div v-if="chat.showToast" class="toast">
        {{ chat.toastMessage }}
      </div>
    </transition>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Smiley+Sans&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

::-webkit-scrollbar {
  width: 4px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(180, 180, 180, 0.3);
  border-radius: 2px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(200, 200, 200, 0.5);
  box-shadow: 0 0 6px rgba(200, 200, 200, 0.3);
}

.dark ::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
  box-shadow: 0 0 6px rgba(255, 255, 255, 0.2);
}

.app-layout {
  display: flex;
  height: 100vh;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.dark {
  --bg-primary: #111111;
  --bg-secondary: #1a1a1a;
  --bg-tertiary: #222222;
  --bg-hover: #2a2a2a;
  --bg-input: #1a1a1a;
  --bg-message-user: #2a2a2a;
  --bg-message-assistant: #111111;
  --text-primary: #e8e8e8;
  --text-secondary: #a0a0a0;
  --text-muted: #666666;
  --border-color: #333333;
  --accent: #ffffff;
  --accent-hover: #cccccc;
  --thinking-bg: #1e1e1e;
  --code-bg: #0a0a0a;
  --scrollbar-thumb: #444444;
}

.light {
  --bg-primary: #ffffff;
  --bg-secondary: #f5f5f5;
  --bg-tertiary: #e8e8e8;
  --bg-hover: #eeeeee;
  --bg-input: #f5f5f5;
  --bg-message-user: #f0f0f0;
  --bg-message-assistant: #ffffff;
  --text-primary: #111111;
  --text-secondary: #555555;
  --text-muted: #999999;
  --border-color: #e0e0e0;
  --accent: #111111;
  --accent-hover: #333333;
  --thinking-bg: #f5f5f5;
  --code-bg: #f6f6f6;
  --scrollbar-thumb: #cccccc;
}

.personalized {
  position: relative;
  background: transparent !important;
}

.personalized::before {
  content: '';
  position: fixed;
  inset: 0;
  background: var(--bg-image, none) center/cover no-repeat fixed;
  filter: blur(24px) brightness(0.55);
  z-index: -1;
  pointer-events: none;
  transition: background 0.5s ease;
}

/* Frosted glass - override component scoped backgrounds */
.personalized .sidebar,
.personalized .chat-header,
.personalized .chat-input-area {
  background: rgba(26, 26, 26, 0.55) !important;
  backdrop-filter: blur(8px) !important;
  -webkit-backdrop-filter: blur(8px) !important;
}

.personalized .message.user .message-body {
  background: rgba(42, 42, 42, 0.55) !important;
  backdrop-filter: blur(6px) !important;
  -webkit-backdrop-filter: blur(6px) !important;
}

.personalized .settings-panel {
  background: rgba(26, 26, 26, 0.7) !important;
  backdrop-filter: blur(12px) !important;
  -webkit-backdrop-filter: blur(12px) !important;
}

/* Accent color on toggle switches */
.personalized .toggle-switch.on {
  background: var(--toggle-on-bg, var(--text-primary)) !important;
  border-color: var(--toggle-on-bg, var(--text-primary)) !important;
}

.personalized .toggle-switch.on .toggle-knob {
  left: 20px !important;
  background: var(--bg-primary) !important;
}

/* Accent color on primary buttons */
.personalized .btn-send,
.personalized .btn-primary {
  background: var(--accent, var(--text-primary)) !important;
}

/* Assure consistent border radius and spacing */
.personalized .settings-overlay {
  background: rgba(0, 0, 0, 0.5) !important;
}

.toast {
  position: fixed;
  bottom: 80px;
  left: 50%;
  transform: translateX(-50%);
  padding: 10px 20px;
  background: #333;
  color: #fff;
  border-radius: 6px;
  font-size: 13px;
  z-index: 9999;
  box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  max-width: 500px;
  text-align: center;
}

.light .toast {
  background: #222;
  color: #fff;
}

.toast-enter-active, .toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(10px);
}
</style>



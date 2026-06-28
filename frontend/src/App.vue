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

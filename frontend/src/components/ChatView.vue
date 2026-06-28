<script setup>
import {ref, watch, nextTick, computed} from 'vue'
import {useChatStore} from '../stores/chat'
import {useSettingsStore} from '../stores/settings'
import MessageBubble from './MessageBubble.vue'
import ChatInput from './ChatInput.vue'
import ModelSelector from './ModelSelector.vue'

const chat = useChatStore()
const settings = useSettingsStore()
const messagesEnd = ref(null)
const selectedModel = ref(settings.settings.model)

const displayMessages = computed(() => {
  if (!chat.activeSession) return []
  return chat.activeSession.messages || []
})

const hasStream = computed(() => chat.isStreaming && (chat.streamingContent || chat.streamingThinking))

watch(
  () => [displayMessages.value.length, chat.streamingContent, chat.streamingThinking],
  () => {
    nextTick(() => {
      messagesEnd.value?.scrollIntoView({behavior: 'smooth'})
    })
  },
  {deep: true}
)

watch(selectedModel, (v) => {
  settings.settings.model = v
  settings.save()
})

function handleSend(content, attachments) {
  chat.send(content, selectedModel.value, settings.settings.thinking_enabled, attachments)
}
</script>

<template>
  <div class="chat-view">
    <div class="chat-header">
      <ModelSelector v-model="selectedModel" :models="chat.models" />
      <div class="header-actions">
        <label class="thinking-toggle" title="深度思考">
          <input type="checkbox" v-model="settings.settings.thinking_enabled" @change="settings.save()" />
          <span>思考</span>
        </label>
      </div>
    </div>

    <div class="messages-area">
      <div v-if="!chat.activeSession || displayMessages.length === 0" class="welcome">
        <h2 class="welcome-title">MOX</h2>
        <p class="welcome-sub">MiMo AI Client</p>
        <div class="model-badge">{{ selectedModel }}</div>
      </div>

      <div v-else class="messages-list">
        <MessageBubble
          v-for="msg in displayMessages"
          :key="msg.id"
          :message="msg"
        />

        <div v-if="hasStream" class="streaming-message">
          <MessageBubble
            :message="{
              id: 'streaming',
              role: 'assistant',
              content: chat.streamingContent,
              reasoning_content: chat.streamingThinking,
              tool_calls: chat.streamingToolCalls,
              timestamp: Date.now()
            }"
            :streaming="true"
          />
        </div>

        <div ref="messagesEnd"></div>
      </div>
    </div>

    <ChatInput
      @send="handleSend"
      :disabled="chat.isStreaming"
      @stop="chat.stop()"
      :is-streaming="chat.isStreaming"
    />
  </div>
</template>

<style scoped>
.chat-view {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100vh;
  min-width: 0;
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-primary);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.thinking-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  cursor: pointer;
  user-select: none;
}

.thinking-toggle input {
  accent-color: var(--accent);
}

.messages-area {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.welcome {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.welcome-title {
  font-family: 'Smiley Sans', sans-serif;
  font-size: 48px;
  color: var(--text-primary);
  margin-bottom: 8px;
  letter-spacing: 2px;
}

.welcome-sub {
  font-size: 14px;
  margin-bottom: 16px;
  color: var(--text-muted);
}

.model-badge {
  padding: 4px 14px;
  background: var(--bg-tertiary);
  border-radius: 4px;
  font-size: 12px;
  color: var(--text-secondary);
  font-family: monospace;
}

.messages-list {
  max-width: 800px;
  width: 100%;
  margin: 0 auto;
  padding: 20px 20px 0;
}
</style>

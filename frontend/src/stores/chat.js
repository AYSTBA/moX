import {defineStore} from 'pinia'
import {ref, computed} from 'vue'
import {
  GetSessions,
  CreateSession,
  SaveSession,
  DeleteSession,
  SendMessage,
  StopGeneration,
  GetModels,
} from '../../wailsjs/go/main/App'
import {EventsOn} from '../../wailsjs/runtime/runtime'

export const useChatStore = defineStore('chat', () => {
  const sessions = ref([])
  const activeSessionKey = ref('')
  const isStreaming = ref(false)
  const streamingContent = ref('')
  const streamingThinking = ref('')
  const streamingToolCalls = ref([])
  const streamingAnnotations = ref([])
  const models = ref([])
  const toastMessage = ref('')
  const showToast = ref(false)

  const activeSession = computed(() =>
    sessions.value.find(s => s.key === activeSessionKey.value)
  )

  async function loadSessions() {
    sessions.value = await GetSessions() || []
    if (sessions.value.length > 0 && !activeSessionKey.value) {
      activeSessionKey.value = sessions.value[0].key
    }
  }

  async function loadModels() {
    models.value = await GetModels() || []
  }

  async function newSession() {
    const s = await CreateSession('新对话')
    sessions.value.unshift(s)
    activeSessionKey.value = s.key
  }

  async function selectSession(key) {
    activeSessionKey.value = key
  }

  async function removeSession(key) {
    await DeleteSession(key)
    sessions.value = sessions.value.filter(s => s.key !== key)
    if (activeSessionKey.value === key) {
      activeSessionKey.value = sessions.value[0]?.key || ''
    }
  }

  async function renameSession(key, label) {
    const s = sessions.value.find(s => s.key === key)
    if (s) {
      s.label = label
      await SaveSession(s)
    }
  }

  async function send(content, model, thinking, attachments) {
    if ((!content || !content.trim()) && (!attachments || attachments.length === 0)) return
    if (isStreaming.value) return

    if (!activeSessionKey.value || !sessions.value.find(s => s.key === activeSessionKey.value)) {
      await newSession()
    }

    isStreaming.value = true
    streamingContent.value = ''
    streamingThinking.value = ''
    streamingToolCalls.value = []
    streamingAnnotations.value = []

    SendMessage(activeSessionKey.value, content || '', model, thinking)
  }

  function stop() {
    StopGeneration()
    isStreaming.value = false
  }

  function setupEvents() {
    EventsOn('chat:userMessage', (msg) => {
      const s = sessions.value.find(s => s.key === activeSessionKey.value)
      if (s) s.messages.push(msg)
    })

    EventsOn('chat:token', (content) => {
      streamingContent.value = content
    })

    EventsOn('chat:thinking', (content) => {
      streamingThinking.value = content
    })

    EventsOn('chat:toolcall', (calls) => {
      streamingToolCalls.value = calls
    })

    EventsOn('chat:annotations', (anns) => {
      streamingAnnotations.value = anns
    })

    EventsOn('chat:done', (msg) => {
      isStreaming.value = false
      streamingContent.value = ''
      streamingThinking.value = ''
      streamingToolCalls.value = []
      streamingAnnotations.value = []
      const s = sessions.value.find(s => s.key === activeSessionKey.value)
      if (s) s.messages.push(msg)
    })

    EventsOn('chat:error', (err) => {
      isStreaming.value = false
      streamingContent.value = ''
      streamingThinking.value = ''
      streamingAnnotations.value = []
      const s = sessions.value.find(s => s.key === activeSessionKey.value)
      if (s) {
        s.messages.push({
          id: 'error-' + Date.now(),
          role: 'assistant',
          content: err,
          timestamp: Date.now(),
        })
      }
    })

    EventsOn('chat:titleUpdated', ({key, label}) => {
      const s = sessions.value.find(s => s.key === key)
      if (s) s.label = label
    })

    EventsOn('chat:toast', (msg) => {
      toastMessage.value = msg
      showToast.value = true
      setTimeout(() => { showToast.value = false }, 4000)
    })
  }

  return {
    sessions,
    activeSessionKey,
    activeSession,
    isStreaming,
    streamingContent,
    streamingThinking,
    streamingToolCalls,
    streamingAnnotations,
    models,
    toastMessage,
    showToast,
    loadSessions,
    loadModels,
    newSession,
    selectSession,
    removeSession,
    renameSession,
    send,
    stop,
    setupEvents,
  }
})

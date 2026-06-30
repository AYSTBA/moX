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
const thinkingStartTime = ref(0)
const thinkingDuration = ref(0)
const sendTime = ref(0)
  const models = ref([])
  const toastMessage = ref('')
  const showToast = ref(false)
  const agentStatus = ref('')

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

    let processedAttachments = []
    if (attachments && attachments.length > 0) {
      processedAttachments = await Promise.all(attachments.map(readFileAtt))
    }

    SendMessage(activeSessionKey.value, content || '', model, thinking, processedAttachments)
  }

  function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => {
        const dataUrl = reader.result
        resolve(dataUrl.split(',')[1])
      }
      reader.onerror = reject
      reader.readAsDataURL(file)
    })
  }

  function readFileAtt(att) {
    return readFileAsBase64(att.file).then(data => ({
      name: att.name,
      mimeType: att.type,
      data: data,
    }))
  }

  function stop() {
    StopGeneration()
    isStreaming.value = false
    thinkingStartTime.value = 0
    thinkingDuration.value = 0
  }

  function setupEvents() {
    EventsOn('chat:userMessage', (msg) => {
      msg.thinking_duration = thinkingDuration.value
      msg.total_duration = Date.now() - sendTime.value
      sendTime.value = 0
      const s = sessions.value.find(s => s.key === activeSessionKey.value)
      if (s) s.messages.push(msg)
    })

    EventsOn('chat:token', (content) => {
      streamingContent.value = content
    })

    EventsOn('chat:thinking', (content) => {
      if (thinkingStartTime.value === 0) {
        thinkingStartTime.value = Date.now()
      }
      thinkingDuration.value = Date.now() - thinkingStartTime.value
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
      agentStatus.value = ''
      streamingContent.value = ''
      streamingThinking.value = ''
      streamingToolCalls.value = []
      streamingAnnotations.value = []
      msg.thinking_duration = thinkingDuration.value
      msg.total_duration = Date.now() - sendTime.value
      sendTime.value = 0
      const s = sessions.value.find(s => s.key === activeSessionKey.value)
      if (s) s.messages.push(msg)
    })

    EventsOn('chat:error', (err) => {
      isStreaming.value = false
      agentStatus.value = ''
      streamingContent.value = ''
      streamingThinking.value = ''
      streamingAnnotations.value = []
      thinkingStartTime.value = 0
      thinkingDuration.value = 0
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

    EventsOn('chat:status', (status) => {
      agentStatus.value = status
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
    thinkingDuration,
    models,
    toastMessage,
    showToast,
    agentStatus,
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





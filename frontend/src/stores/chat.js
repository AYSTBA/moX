import {defineStore} from "pinia"
import {ref, computed} from "vue"
import {chat, getModels, getConversations, saveConversation, deleteConversation} from "../api.js"

export const useChatStore = defineStore("chat", () => {
  const models = ref([])
  const streamingContent = ref("")
  const isStreaming = ref(false)
  const conversations = ref([])
  const currentId = ref(null)

  const currentConversation = computed(() => {
    return conversations.value.find((c) => c.id === currentId.value) || null
  })

  const currentMessages = computed({
    get() {
      return currentConversation.value?.messages || []
    },
    set(messages) {
      const conv = conversations.value.find((c) => c.id === currentId.value)
      if (conv) conv.messages = messages
    },
  })

  async function loadModels() {
    models.value = await getModels()
  }

  function loadConversations() {
    const saved = getConversations()
    conversations.value = saved
    if (saved.length > 0 && !currentId.value) {
      currentId.value = saved[0].id
    }
  }

  function newConversation() {
    const id = Date.now().toString(36) + Math.random().toString(36).slice(2, 6)
    const conv = {
      id,
      title: "新对话",
      messages: [],
      model: "mimo-v2.5",
      createdAt: Date.now(),
      updatedAt: Date.now(),
    }
    conversations.value.unshift(conv)
    currentId.value = id
    persistConversations()
    return id
  }

  function selectConversation(id) {
    currentId.value = id
  }

  function removeConversation(id) {
    const idx = conversations.value.findIndex((c) => c.id === id)
    if (idx === -1) return
    conversations.value.splice(idx, 1)
    if (currentId.value === id) {
      currentId.value = conversations.value.length > 0 ? conversations.value[0].id : null
    }
    persistConversations()
  }

  function persistConversations() {
    saveConversation(conversations.value)
  }

  async function send(messages, model) {
    if (!currentId.value) newConversation()
    streamingContent.value = ""
    isStreaming.value = true
    try {
      await chat(messages, model, (token) => {
        streamingContent.value += token
      })
      const result = streamingContent.value
      streamingContent.value = ""
      return result
    } finally {
      isStreaming.value = false
    }
  }

  return {
    models,
    streamingContent,
    isStreaming,
    conversations,
    currentId,
    currentConversation,
    currentMessages,
    loadModels,
    loadConversations,
    newConversation,
    selectConversation,
    removeConversation,
    persistConversations,
    send,
  }
})

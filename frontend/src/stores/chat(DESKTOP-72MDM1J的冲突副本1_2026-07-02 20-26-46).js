import {defineStore} from 'pinia'
import {ref, computed} from 'vue'
import {chat, getModels} from '../api.js'

export const useChatStore = defineStore('chat', () => {
  const models = ref([])
  const streamingContent = ref('')
  const isStreaming = ref(false)

  async function loadModels() {
    models.value = await getModels()
  }

  async function send(messages, model, onFirstToken) {
    isStreaming.value = true
    streamingContent.value = ''
    await chat(messages, model, (token) => {
      streamingContent.value += token
      if (onFirstToken && streamingContent.value.length === token.length) onFirstToken()
    })
    isStreaming.value = false
    const result = streamingContent.value
    streamingContent.value = ''
    return result
  }

  return {models, streamingContent, isStreaming, loadModels, send}
})

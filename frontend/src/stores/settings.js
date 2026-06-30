import {defineStore} from 'pinia'
import {ref} from 'vue'
import {GetSettings, SaveSettings} from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref({
    api_key: '',
    model: 'mimo-v2.5',
    theme: 'dark',
    system_prompt: '你是MiMo，是小米公司研发的AI智能助手。',
    temperature: 1.0,
    top_p: 0.95,
    max_tokens: 4096,
    thinking_enabled: true,
    web_search_enabled: false,
    external_search_api_key: '',
    external_search_enabled: false,
    time_awareness: false,
    personalization_enabled: false,
    background_image: '',
  })

  const showSettings = ref(false)

  async function load() {
    const s = await GetSettings()
    if (s) settings.value = s
    applyTheme(settings.value.theme)
    applyPersonalization()
  }

  async function save() {
    await SaveSettings(settings.value)
    applyPersonalization()
  }

  function applyTheme(theme) {
    if (theme === 'dark') {
      document.documentElement.classList.add('dark')
      document.documentElement.classList.remove('light')
    } else {
      document.documentElement.classList.add('light')
      document.documentElement.classList.remove('dark')
    }
  }

  function applyPersonalization() {
    const s = settings.value
    const root = document.documentElement

    root.classList.toggle('personalized', s.personalization_enabled && !!s.background_image)

    if (s.personalization_enabled && s.background_image) {
      root.style.setProperty('--bg-image', 'url(' + s.background_image + ')')
    } else {
      root.style.setProperty('--bg-image', 'none')
    }
  }

  function toggleTheme() {
    settings.value.theme = settings.value.theme === 'dark' ? 'light' : 'dark'
    applyTheme(settings.value.theme)
    save()
  }

  return {settings, showSettings, load, save, toggleTheme}
})

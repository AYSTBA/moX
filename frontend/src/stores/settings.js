import {defineStore} from 'pinia'
import {ref} from 'vue'
import {GetSettings, SaveSettings} from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref({
    api_key: '',
    model: 'mimo-v2.5-pro',
    theme: 'dark',
    system_prompt: '你是MiMo，是小米公司研发的AI智能助手。',
    temperature: 1.0,
    top_p: 0.95,
    max_tokens: 4096,
    thinking_enabled: true,
  })

  const showSettings = ref(false)

  async function load() {
    const s = await GetSettings()
    if (s) settings.value = s
    applyTheme(settings.value.theme)
  }

  async function save() {
    await SaveSettings(settings.value)
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

  function toggleTheme() {
    settings.value.theme = settings.value.theme === 'dark' ? 'light' : 'dark'
    applyTheme(settings.value.theme)
    save()
  }

  return {settings, showSettings, load, save, toggleTheme}
})

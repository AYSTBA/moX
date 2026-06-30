import {defineStore} from 'pinia'
import {ref} from 'vue'
import {GetSettings, SaveSettings} from '../../wailsjs/go/main/App'

function hexToRgba(hex, alpha) {
  if (!hex || hex.length < 7) return hex
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

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
    accent_color: '',
    bg_color: '',
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
    const enabled = s.personalization_enabled

    root.classList.toggle('personalized', enabled && (s.accent_color || s.background_image))

    if (enabled && s.accent_color) {
      const c = s.accent_color
      root.style.setProperty('--accent', c)
      root.style.setProperty('--accent-hover', c + 'cc')
      root.style.setProperty('--toggle-on-bg', c)
      root.style.setProperty('--input-focus-border', c)
      root.style.setProperty('--avatar-bg', c)
      root.style.setProperty('--thinking-accent', hexToRgba(c, 0.3))
      root.style.setProperty('--message-user-border', hexToRgba(c, 0.4))
    } else {
      const vars = ['--accent','--accent-hover','--toggle-on-bg','--input-focus-border','--avatar-bg','--thinking-accent','--message-user-border']
      vars.forEach(v => root.style.removeProperty(v))
    }

    if (enabled && s.background_image) {
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

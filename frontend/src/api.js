const API = "http://localhost:3099"

let apiKey = localStorage.getItem("mox_api_key") || ""

export function setKey(key) {
  apiKey = key
  localStorage.setItem("mox_api_key", key)
}

export function getKey() {
  return apiKey
}

// === Conversation persistence ===

const STORAGE_KEY = "mox_conversations_v2"

export function getConversations() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : []
  } catch {
    return []
  }
}

export function saveConversation(conversations) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(conversations))
  } catch (e) {
    console.error("Failed to save conversations:", e)
  }
}

export function deleteConversation(id) {
  const conversations = getConversations()
  const idx = conversations.findIndex((c) => c.id === id)
  if (idx !== -1) {
    conversations.splice(idx, 1)
    saveConversation(conversations)
  }
}

// === Settings API ===

export async function getSettings() {
  const resp = await fetch(API + "/api/settings")
  if (!resp.ok) return {}
  return resp.json()
}

export async function saveSettings(apiKey, systemPrompt) {
  await fetch(API + "/api/settings", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({api_key: apiKey, system_prompt: systemPrompt}),
  })
}

// === API calls ===

export async function chat(messages, model, onToken) {
  const resp = await fetch(API + "/v1/chat/completions", {
    method: "POST",
    headers: {"Content-Type": "application/json", "X-Mimo-Key": apiKey},
    body: JSON.stringify({model: model || "mimo-v2.5", messages, stream: true}),
  })
  if (!resp.ok) {
    const err = await resp.text()
    throw new Error(err)
  }
  const reader = resp.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ""
  while (true) {
    const {done, value} = await reader.read()
    if (done) break
    buffer += decoder.decode(value, {stream: true})
    const lines = buffer.split("\n")
    buffer = lines.pop() || ""
    for (const line of lines) {
      if (!line.startsWith("data: ")) continue
      const data = line.slice(6).trim()
      if (data === "[DONE]") return
      try {
        const json = JSON.parse(data)
        const content = json.choices?.[0]?.delta?.content || ""
        if (content) onToken(content)
      } catch {}
    }
  }
}

export async function getModels() {
  const resp = await fetch(API + "/v1/models")
  return (await resp.json()).data || []
}


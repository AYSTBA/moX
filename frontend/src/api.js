const API = 'http://localhost:3099'

let apiKey = localStorage.getItem('mox_api_key') || ''

export function setKey(key) {
  apiKey = key
  localStorage.setItem('mox_api_key', key)
}

export function getKey() {
  return apiKey
}

export async function chat(messages, model, onToken) {
  const resp = await fetch(API + '/v1/chat/completions', {
    method: 'POST',
    headers: {'Content-Type': 'application/json', 'X-Mimo-Key': apiKey},
    body: JSON.stringify({model: model || 'mimo-v2.5', messages, stream: true}),
  })
  if (!resp.ok) {
    const err = await resp.text()
    throw new Error(err)
  }
  const reader = resp.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  while (true) {
    const {done, value} = await reader.read()
    if (done) break
    buffer += decoder.decode(value, {stream: true})
    const lines = buffer.split('\n')
    buffer = lines.pop() || ''
    for (const line of lines) {
      if (!line.startsWith('data: ')) continue
      const data = line.slice(6).trim()
      if (data === '[DONE]') return
      try {
        const json = JSON.parse(data)
        const content = json.choices?.[0]?.delta?.content || ''
        if (content) onToken(content)
      } catch {}
    }
  }
}

export async function getModels() {
  const resp = await fetch(API + '/v1/models')
  return (await resp.json()).data || []
}

<script setup>
import {ref} from 'vue'
import {useSettingsStore} from '../stores/settings'
import {TestAPIKey} from '../../wailsjs/go/main/App'

const settings = useSettingsStore()
const testing = ref(false)
const testResult = ref('')

async function testKey() {
  testing.value = true
  testResult.value = ''
  const result = await TestAPIKey(settings.settings.api_key)
  testResult.value = result === 'ok' ? 'API Key 有效' : result
  testing.value = false
}

function close() {
  settings.save()
  settings.showSettings = false
}
</script>

<template>
  <div class="settings-overlay" @click.self="close">
    <div class="settings-panel">
      <div class="settings-header">
        <h3>设置</h3>
        <button class="btn-close" @click="close">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <div class="settings-body">
        <div class="setting-group">
          <label>API Key</label>
          <div class="api-key-row">
            <input
              type="password"
              v-model="settings.settings.api_key"
              placeholder="输入 MiMo API Key"
              class="input"
            />
            <button class="btn-test" @click="testKey" :disabled="testing || !settings.settings.api_key">
              {{ testing ? '测试中...' : '测试' }}
            </button>
          </div>
          <div v-if="testResult" class="test-result" :class="{ok: testResult === 'API Key 有效'}">{{ testResult }}</div>
          <div class="setting-hint">
            从 <a href="https://platform.xiaomimimo.com/#/console/api-keys" target="_blank">MiMo 开放平台</a> 获取 API Key
          </div>
        </div>

        <div class="setting-group">
          <label>系统提示词</label>
          <textarea
            v-model="settings.settings.system_prompt"
            class="input textarea"
            rows="3"
            placeholder="系统提示词"
          ></textarea>
        </div>

        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-group half">
              <label>时间感知</label>
              <div class="toggle-row">
                <button
                  class="toggle-switch"
                  :class="{on: settings.settings.time_awareness}"
                  @click="settings.settings.time_awareness = !settings.settings.time_awareness"
                >
                  <span class="toggle-knob"></span>
                </button>
                <span class="toggle-desc">向模型提供当前日期时间</span>
              </div>
            </div>
          </div>
        </div>

        <div class="setting-group">
          <label>外部搜索 API (Tavily)</label>
          <div class="api-key-row">
            <input
              type="password"
              v-model="settings.settings.external_search_api_key"
              placeholder="Tavily API Key (可选)"
              class="input"
            />
            <button
              class="toggle-switch"
              :class="{on: settings.settings.external_search_enabled}"
              @click="settings.settings.external_search_enabled = !settings.settings.external_search_enabled"
            >
              <span class="toggle-knob"></span>
            </button>
          </div>
          <div class="setting-hint">
            使用 Tavily API 进行联网搜索，需先获取 <a href="https://tavily.com" target="_blank">Tavily API Key</a>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-group half">
            <label>Temperature</label>
            <div class="slider-row">
              <input type="range" v-model.number="settings.settings.temperature" min="0" max="1.5" step="0.1" />
              <span class="slider-value">{{ settings.settings.temperature }}</span>
            </div>
          </div>
          <div class="setting-group half">
            <label>Top P</label>
            <div class="slider-row">
              <input type="range" v-model.number="settings.settings.top_p" min="0.01" max="1" step="0.01" />
              <span class="slider-value">{{ settings.settings.top_p }}</span>
            </div>
          </div>
        </div>

        <div class="setting-group">
          <label>最大输出 Token</label>
          <input
            type="number"
            v-model.number="settings.settings.max_tokens"
            class="input"
            min="256"
            max="131072"
          />
        </div>
      </div>

      <div class="settings-footer">
        <button class="btn-primary" @click="close">保存</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.settings-panel {
  width: 460px;
  max-height: 80vh;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  box-shadow: 0 16px 48px rgba(0,0,0,0.4);
  display: flex;
  flex-direction: column;
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.settings-header h3 {
  font-family: 'Smiley Sans', sans-serif;
  font-size: 18px;
  font-weight: 600;
}

.btn-close {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
}

.btn-close:hover {
  color: var(--text-primary);
}

.settings-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.setting-group {
  margin-bottom: 20px;
}

.setting-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 6px;
  color: var(--text-secondary);
}

.input {
  width: 100%;
  padding: 10px 14px;
  background: var(--bg-input);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.input:focus {
  border-color: var(--text-muted);
}

.textarea {
  resize: vertical;
  min-height: 60px;
  font-family: inherit;
}

.api-key-row {
  display: flex;
  gap: 8px;
}

.api-key-row .input {
  flex: 1;
}

.btn-test {
  padding: 8px 16px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  white-space: nowrap;
  transition: background 0.2s;
}

.btn-test:hover:not(:disabled) {
  background: var(--bg-hover);
}

.btn-test:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.test-result {
  margin-top: 6px;
  font-size: 12px;
  color: #cc4444;
}

.test-result.ok {
  color: #44aa44;
}

.setting-hint {
  margin-top: 4px;
  font-size: 11px;
  color: var(--text-muted);
}

.setting-hint a {
  color: var(--text-primary);
  text-decoration: underline;
}

.toggle-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.toggle-desc {
  font-size: 12px;
  color: var(--text-muted);
}

.toggle-switch {
  width: 40px;
  height: 22px;
  border-radius: 11px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  cursor: pointer;
  position: relative;
  transition: all 0.25s ease;
  padding: 0;
  flex-shrink: 0;
}

.toggle-switch.on {
  background: var(--text-primary);
  border-color: var(--text-primary);
}

.toggle-knob {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--text-muted);
  transition: all 0.25s ease;
}

.toggle-switch.on .toggle-knob {
  left: 20px;
  background: var(--bg-primary);
}

.setting-row {
  display: flex;
  gap: 16px;
}

.setting-group.half {
  flex: 1;
}

.slider-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.slider-row input[type="range"] {
  flex: 1;
  accent-color: var(--text-primary);
}

.slider-value {
  font-size: 13px;
  font-weight: 500;
  min-width: 40px;
  text-align: right;
  font-family: monospace;
}

.settings-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
}

.btn-primary {
  width: 100%;
  padding: 10px;
  background: var(--text-primary);
  color: var(--bg-primary);
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn-primary:hover {
  opacity: 0.85;
}
</style>

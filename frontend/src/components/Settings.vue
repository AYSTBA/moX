<script setup>
import {ref} from 'vue'
import {useSettingsStore} from '../stores/settings'
import {TestAPIKey} from '../../wailsjs/go/main/App'

const settings = useSettingsStore()
const testing = ref(false)
const testResult = ref('')
const bgInput = ref(null)

function onBgSelected(e) {
  const file = e.target.files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    settings.settings.background_image = reader.result
    settings.save()
  }
  reader.readAsDataURL(file)
}

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
          <label>联网搜索</label>
          <div class="search-options">
            <div class="search-option">
              <div class="search-option-header">
                <span class="search-option-name">MiMo 内置搜索</span>
                <button
                  class="toggle-switch"
                  :class="{on: settings.settings.web_search_enabled}"
                  @click="settings.settings.web_search_enabled = !settings.settings.web_search_enabled"
                >
                  <span class="toggle-knob"></span>
                </button>
              </div>
              <div class="search-option-desc">使用 MiMo 官方搜索插件，按次计费 (约 ¥16/千次)</div>
            </div>
            <div class="search-option">
              <div class="search-option-header">
                <span class="search-option-name">Tavily 外部搜索</span>
                <button
                  class="toggle-switch"
                  :class="{on: settings.settings.external_search_enabled}"
                  @click="settings.settings.external_search_enabled = !settings.settings.external_search_enabled"
                >
                  <span class="toggle-knob"></span>
                </button>
              </div>
              <div class="search-option-desc">使用 Tavily API，需自行申请 Key</div>
              <input
                v-if="settings.settings.external_search_enabled"
                type="password"
                v-model="settings.settings.external_search_api_key"
                placeholder="Tavily API Key"
                class="input mt-6"
              />
            </div>
          </div>
          <div class="setting-hint">两个搜索可同时开启，也可只开其一。关闭 MiMo 搜索可节省费用。</div>
        <div class="setting-group">
          <div class="personalization-header">
            <span class="personalization-title">个性化</span>
            <button
              class="toggle-switch"
              :class="{on: settings.settings.personalization_enabled}"
              @click="settings.settings.personalization_enabled = !settings.settings.personalization_enabled; settings.save()"
            >
              <span class="toggle-knob"></span>
            </button>
          </div>

          <div v-if="settings.settings.personalization_enabled" class="personalization-body">

            <label>背景图片</label>
            <div class="bg-image-row">
              <button class="btn-upload" @click="bgInput.click()">
                {{ settings.settings.background_image ? "更换图片" : "选择图片" }}
              </button>
              <button v-if="settings.settings.background_image" class="btn-clear" @click="settings.settings.background_image = ''; settings.save()">
                清除
              </button>
            </div>
            <div v-if="settings.settings.background_image" class="bg-preview">
              <img :src="settings.settings.background_image" class="bg-thumb" />
            </div>
            <input type="file" ref="bgInput" accept="image/*" style="display:none" @change="onBgSelected" />

            <label>雾面强度</label>
            <div class="slider-row">
              <input type="range" v-model.number="settings.settings.blur_intensity" min="0" max="100" step="1" @input="settings.save()" />
              <span class="slider-value">{{ settings.settings.blur_intensity }}px</span>
            </div>
            <div class="setting-hint">建议图片: 1920×1080 或更高分辨率</div>
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

.search-options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.search-option {
  padding: 10px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.search-option-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-option-name {
  font-size: 13px;
  font-weight: 500;
}

.search-option-desc {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}

.mt-6 {
  margin-top: 8px;
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
.section-title {
  font-family: 'Smiley Sans', sans-serif;
  font-size: 16px;
  font-weight: 600;
}

.personalization-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.personalization-title {
  font-family: 'Smiley Sans', sans-serif;
  font-size: 16px;
  font-weight: 600;
}

.personalization-body {
  padding: 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.personalization-body label {
  display: block;
  font-size: 12px;
  margin-top: 12px;
  margin-bottom: 6px;
  color: var(--text-secondary);
}

.personalization-body label:first-child {
  margin-top: 0;
}


.color-option:hover {
  border-color: var(--text-muted);
}

.color-option.active {
  border-color: var(--text-primary);
}

.color-swatch {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 1px solid rgba(255,255,255,0.1);
}

.bg-image-row {
  display: flex;
  gap: 8px;
}

.btn-upload, .btn-clear {
  padding: 8px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: background 0.2s;
}

.btn-upload:hover, .btn-clear:hover {
  background: var(--bg-hover);
}

.bg-preview {
  margin-top: 8px;
}

.bg-thumb {
  max-width: 100%;
  max-height: 120px;
  border-radius: 4px;
  object-fit: cover;
  border: 1px solid var(--border-color);
}
</style>

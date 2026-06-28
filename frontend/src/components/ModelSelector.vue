<script setup>
import {ref} from 'vue'

const props = defineProps({
  modelValue: String,
  models: Array,
})

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false)

function select(id) {
  emit('update:modelValue', id)
  isOpen.value = false
}

function getCurrentModel() {
  const m = props.models?.find(m => m.id === props.modelValue)
  return m?.name || props.modelValue
}
</script>

<template>
  <div class="model-selector" @click="isOpen = !isOpen" @blur="isOpen = false" tabindex="0">
    <span class="model-current">{{ getCurrentModel() }}</span>
    <svg class="dropdown-icon" :class="{open: isOpen}" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="6 9 12 15 18 9"/>
    </svg>
    <div v-if="isOpen" class="dropdown-menu">
      <div
        v-for="m in models"
        :key="m.id"
        class="dropdown-item"
        :class="{active: m.id === modelValue}"
        @click.stop="select(m.id)"
      >
        <div class="item-name">{{ m.name }}</div>
        <div class="item-desc">{{ m.desc }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.model-selector {
  position: relative;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  user-select: none;
  transition: border-color 0.2s;
}

.model-selector:hover {
  border-color: var(--text-muted);
}

.dropdown-icon {
  transition: transform 0.2s;
}

.dropdown-icon.open {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  min-width: 280px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 4px;
  z-index: 100;
  box-shadow: 0 8px 24px rgba(0,0,0,0.3);
}

.dropdown-item {
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.15s;
}

.dropdown-item:hover {
  background: var(--bg-hover);
}

.dropdown-item.active {
  background: var(--bg-tertiary);
}

.item-name {
  font-size: 13px;
  font-weight: 500;
}

.item-desc {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
}
</style>

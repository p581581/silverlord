<template>
  <div class="qty-control">
    <button class="qty-btn" :disabled="disabled || localVal <= min" @click.prevent="step(-1)">
      <i class="pi pi-minus" />
    </button>
    <input
      type="number"
      class="qty-input"
      :value="localVal"
      :min="min"
      :max="max"
      :disabled="disabled"
      @input="(e) => localVal = parseInt((e.target as HTMLInputElement).value)"
      @blur="commit"
    />
    <button class="qty-btn" :disabled="disabled || (max != null && localVal >= max)" @click.prevent="step(1)">
      <i class="pi pi-plus" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: number
  min?: number
  max?: number
  disabled?: boolean
}>(), {
  min: 1,
  max: undefined,
  disabled: false,
})

const emit = defineEmits<{ 'update:modelValue': [value: number] }>()

const localVal = ref(props.modelValue)

watch(() => props.modelValue, (v) => { localVal.value = v })

function clamp(v: number) {
  const lo = props.min ?? 1
  const hi = props.max ?? Infinity
  return isNaN(v) ? lo : Math.min(Math.max(v, lo), hi)
}

function commit() {
  const v = clamp(localVal.value)
  localVal.value = v
  emit('update:modelValue', v)
}

function step(delta: number) {
  const v = clamp(localVal.value + delta)
  localVal.value = v
  emit('update:modelValue', v)
}
</script>

<style scoped>
.qty-control {
  display: flex;
  align-items: center;
  border: 1.5px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
  background: #f8fafc;
  flex-shrink: 0;
}

.qty-btn {
  background: none;
  border: none;
  width: 2rem;
  height: 2.2rem;
  color: #475569;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  transition: background 0.15s, color 0.15s;
}
.qty-btn:hover:not(:disabled) { background: #e2e8f0; color: #0f172a; }
.qty-btn:disabled { color: #cbd5e1; cursor: not-allowed; }

.qty-input {
  width: 2.5rem;
  min-width: 2.5rem;
  text-align: center;
  font-weight: 700;
  font-size: 0.9rem;
  color: #0f172a;
  background: transparent;
  border: none;
  outline: none;
  box-shadow: none;
  padding: 0;
  height: 2.2rem;
  -moz-appearance: textfield;
}
.qty-input::-webkit-outer-spin-button,
.qty-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
}
.qty-input:disabled { color: #cbd5e1; }
</style>

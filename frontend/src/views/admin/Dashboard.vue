<template>
  <div>
    <h1 class="page-title">{{ t('admin.dashboard') }}</h1>

    <!-- 即時報價監控 -->
    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-chart-line" /> {{ t('admin.spotPrice') }}
      </h2>
      <div v-if="spot" class="spot-grid">
        <div class="spot-item">
          <div class="spot-label">XAG / USD</div>
          <div class="spot-value">${{ spot.xag_usd?.toFixed(4) }}</div>
        </div>
        <div class="spot-item">
          <div class="spot-label">USD / TWD</div>
          <div class="spot-value">{{ spot.usd_twd?.toFixed(2) }}</div>
        </div>
        <div class="spot-item">
          <div class="spot-label">{{ t('admin.lastUpdated') }}</div>
          <div class="spot-value updated-at">{{ lastUpdated }}</div>
        </div>
      </div>
      <p v-else class="not-ready">{{ t('errors.PRICE_NOT_READY') }}</p>
    </div>

    <!-- 控盤中心 -->
    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-sliders-h" /> {{ t('admin.dashboard') }}
      </h2>
      <div class="control-grid">
        <div class="field">
          <label>{{ t('admin.manualOffset') }}</label>
          <InputNumber
            v-model="form.manual_offset"
            :min-fraction-digits="2"
            :max-fraction-digits="4"
            :step="0.1"
            show-buttons
            prefix="± " />
          <small>目前: {{ form.manual_offset >= 0 ? '+' : '' }}{{ form.manual_offset }}</small>
        </div>
        <div class="field">
          <label>{{ t('admin.manualExRate') }}</label>
          <InputNumber
            v-model="form.manual_ex_rate"
            :min-fraction-digits="2"
            :max-fraction-digits="4"
            :min="0"
            :step="0.1"
            show-buttons />
          <small v-if="form.manual_ex_rate > 0" class="warn">
            ⚠ 固定匯率啟用中: {{ form.manual_ex_rate }}
          </small>
          <small v-else>停用（使用即時匯率）</small>
        </div>
      </div>
      <Button
        :label="t('admin.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="saveConfig" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import { usePriceStore } from '@/stores/price'
import api from '@/api'

const { t } = useI18n()
const toast = useToast()
const priceStore = usePriceStore()

const spot = computed(() => priceStore.snapshot)
const lastUpdated = computed(() =>
  spot.value?.updated_at ? new Date(spot.value.updated_at).toLocaleTimeString() : '',
)

const form = ref({ manual_offset: 0, manual_ex_rate: 0 })
const saving = ref(false)

onMounted(async () => {
  priceStore.startAutoRefresh()
  try {
    const res = await api.get('/api/admin/config')
    form.value.manual_offset = res.data.manual_offset
    form.value.manual_ex_rate = res.data.manual_ex_rate
  } catch { }
})

async function saveConfig() {
  saving.value = true
  try {
    await api.put('/api/admin/config', form.value)
    toast.add({ severity: 'success', summary: t('admin.save'), detail: 'OK', life: 2000 })
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
}

.section-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 1.25rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.spot-grid {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}

.spot-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.spot-label {
  font-size: 0.8rem;
  color: #999;
}

.spot-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1a1a2e;
}

.updated-at {
  font-size: 1rem;
  color: #27ae60;
}

.not-ready {
  color: #e74c3c;
}

.control-grid {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
  margin-bottom: 1.25rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  min-width: 200px;
}

.field label {
  font-weight: 600;
  font-size: 0.9rem;
}

.field small {
  color: #999;
  font-size: 0.8rem;
}

.warn {
  color: #e67e22 !important;
}
</style>

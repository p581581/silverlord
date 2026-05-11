<template>
  <div>
    <h1 class="page-title">{{ t('admin.shipping') }}</h1>

    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-truck" /> {{ t('admin.shipFeeSettings') }}
      </h2>

      <div class="fee-grid">
        <!-- 郵寄 -->
        <div class="fee-card">
          <div class="fee-icon"><i class="pi pi-envelope" /></div>
          <div class="fee-label">{{ t('checkout.shipMail') }}</div>
          <div class="fee-input">
            <span class="currency">NT$</span>
            <InputNumber
              v-model="form.ship_fee_mail"
              :use-grouping="false"
              :min="0"
              placeholder="0"
              class="fee-number" />
          </div>
          <div class="fee-hint">{{ form.ship_fee_mail === 0 ? t('checkout.shipFreeText') : `NT$ ${form.ship_fee_mail}`
            }}</div>
        </div>

        <!-- 7-11 -->
        <div class="fee-card">
          <div class="fee-icon"><i class="pi pi-shop" /></div>
          <div class="fee-label">{{ t('checkout.ship711') }}</div>
          <div class="fee-input">
            <span class="currency">NT$</span>
            <InputNumber
              v-model="form.ship_fee_711"
              :use-grouping="false"
              :min="0"
              placeholder="0"
              class="fee-number" />
          </div>
          <div class="fee-hint">{{ form.ship_fee_711 === 0 ? t('checkout.shipFreeText') : `NT$ ${form.ship_fee_711}` }}
          </div>
        </div>

        <!-- 自取 -->
        <div class="fee-card">
          <div class="fee-icon"><i class="pi pi-home" /></div>
          <div class="fee-label">{{ t('checkout.shipPickup') }}</div>
          <div class="fee-input">
            <span class="currency">NT$</span>
            <InputNumber
              v-model="form.ship_fee_pickup"
              :use-grouping="false"
              :min="0"
              placeholder="0"
              class="fee-number" />
          </div>
          <div class="fee-hint">{{ form.ship_fee_pickup === 0 ? t('checkout.shipFreeText') : `NT$ ${form.ship_fee_pickup}` }}</div>
        </div>
      </div>

      <Button
        :label="t('admin.save')"
        icon="pi pi-save"
        :loading="saving"
        class="mt-4"
        @click="save" />
    </div>

    <!-- 說明 -->
    <div class="section-card info-card">
      <h2 class="section-title"><i class="pi pi-info-circle" /> 說明</h2>
      <ul class="info-list">
        <li>運費設為 <strong>0</strong> 表示免運費</li>
        <li>運費將於結帳時加入訂單總額</li>
        <li>修改後立即生效，不影響已建立的訂單</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import InputNumber from 'primevue/inputnumber'
import Button from 'primevue/button'
import api from '@/api'

const { t } = useI18n()
const toast = useToast()

const form = ref({ ship_fee_mail: 0, ship_fee_711: 0, ship_fee_pickup: 0 })
const saving = ref(false)

onMounted(async () => {
  try {
    const res = await api.get('/api/admin/config')
    form.value.ship_fee_mail = res.data.ship_fee_mail ?? 0
    form.value.ship_fee_711 = res.data.ship_fee_711 ?? 0
    form.value.ship_fee_pickup = res.data.ship_fee_pickup ?? 0
  } catch { }
})

async function save() {
  saving.value = true
  try {
    // 只更新運費欄位，先讀取完整 config 再合併儲存
    const res = await api.get('/api/admin/config')
    await api.put('/api/admin/config', {
      ...res.data,
      ship_fee_mail: form.value.ship_fee_mail,
      ship_fee_711: form.value.ship_fee_711,
      ship_fee_pickup: form.value.ship_fee_pickup,
    })
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
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.fee-grid {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
}

.fee-card {
  flex: 1;
  min-width: 200px;
  border: 2px solid #eee;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  transition: border-color 0.2s;
}

.fee-card:hover {
  border-color: #c0c0c0;
}

.fee-icon {
  font-size: 2rem;
  color: #1a1a2e;
}

.fee-label {
  font-weight: 600;
  font-size: 1rem;
}

.fee-input {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.currency {
  font-weight: 600;
  color: #666;
}

.fee-number {
  width: 100px;
}
.fee-number :deep(.p-inputnumber-input) {
  width: 100%;
}

.fee-hint {
  font-size: 0.85rem;
  color: #27ae60;
  font-weight: 600;
}

.info-card {
  background: #f8f9fa;
}

.info-list {
  margin: 0;
  padding-left: 1.25rem;
  color: #555;
  line-height: 2;
}
</style>

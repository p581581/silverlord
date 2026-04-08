<template>
  <div class="page">
    <Navbar />

    <main class="container">
      <div class="payment-card">
        <h1 class="page-title">{{ t('payment.title') }}</h1>

        <!-- 成功 -->
        <div v-if="success" class="success-box">
          <i class="pi pi-check-circle success-icon" />
          <p>{{ t('payment.success') }}</p>
          <Button :label="t('nav.products')" @click="router.push('/')" />
        </div>

        <!-- 表單 -->
        <form v-else @submit.prevent="submit">
          <div class="field">
            <label>{{ t('payment.orderNo') }}</label>
            <InputText
              v-model="form.orderNo"
              :placeholder="'SN20250407XXXX'"
              class="w-full"
              required
            />
          </div>

          <div class="field">
            <label>{{ t('payment.remitInfo') }}</label>
            <InputMask
              v-model="form.remitInfo"
              mask="99999"
              :placeholder="t('payment.remitPlaceholder')"
              class="w-full"
              required
            />
            <small class="field-hint">{{ t('payment.remitPlaceholder') }}</small>
          </div>

          <Button
            type="submit"
            :label="t('payment.submit')"
            icon="pi pi-send"
            :loading="submitting"
            class="w-full mt-2"
          />
        </form>
      </div>
    </main>

    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import InputMask from 'primevue/inputmask'
import Button from 'primevue/button'
import Navbar from '@/components/shared/Navbar.vue'
import Footer from '@/components/shared/Footer.vue'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const toast = useToast()

const form = ref({ orderNo: '', remitInfo: '' })
const submitting = ref(false)
const success = ref(false)

async function submit() {
  submitting.value = true
  try {
    await api.post('/api/v1/report-payment', {
      order_no: form.value.orderNo,
      remit_info: form.value.remitInfo,
    })
    success.value = true
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 4000 })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page { min-height: 100vh; background: #f8fafc; display: flex; flex-direction: column; }
.container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 3rem 1rem;
}
.payment-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}
.page-title { font-size: 1.4rem; font-weight: 700; margin-bottom: 1.5rem; }
.field { display: flex; flex-direction: column; gap: 0.4rem; margin-bottom: 1.25rem; }
.field label { font-weight: 600; font-size: 0.9rem; }
.field-hint { color: #999; font-size: 0.8rem; }
.success-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem 0;
}
.success-icon { font-size: 3.5rem; color: #27ae60; }
</style>

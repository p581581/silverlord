<template>
  <div class="page">
    <Navbar />

    <main class="container">
      <div class="query-card">
        <h1 class="page-title">{{ t('orderQuery.title') }}</h1>

        <!-- 搜尋列 -->
        <div class="search-row">
          <InputText
            v-model="orderNo"
            :placeholder="t('orderQuery.placeholder')"
            class="search-input"
            @keyup.enter="search" />
          <Button
            :label="t('orderQuery.search')"
            icon="pi pi-search"
            :loading="loading"
            class="search-btn"
            @click="search" />
        </div>

        <!-- 本機歷史訂單 -->
        <div v-if="myOrders.length > 0 && !result" class="history-section">
          <div class="history-title">
            <i class="pi pi-history" /> 本裝置的訂單紀錄
          </div>
          <div class="history-list">
            <div
              v-for="no in myOrders"
              :key="no"
              class="history-item"
              @click="quickSearch(no)">
              <code>{{ no }}</code>
              <i class="pi pi-chevron-right" />
            </div>
          </div>
        </div>

        <!-- 查無訂單 -->
        <div v-if="notFound" class="not-found">
          <i class="pi pi-exclamation-circle" />
          {{ t('orderQuery.notFound') }}
        </div>

        <!-- 查詢結果 -->
        <div v-if="result" class="result-box">
          <div class="result-header">
            <div class="meta-item">
              <span class="meta-label">{{ t('admin.orderNo') }}</span>
              <code class="order-no">{{ result.order_no }}</code>
            </div>
            <div class="meta-item">
              <span class="meta-label">{{ t('admin.status') }}</span>
              <Tag
                :value="t(`admin.statusList.${result.status}`)"
                :severity="statusSeverity(result.status)" />
            </div>
            <div class="meta-item">
              <span class="meta-label">{{ t('admin.totalLockedPrice') }}</span>
              <span class="locked-price">NT$ {{ result.total_locked_price.toLocaleString() }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-label">{{ t('admin.createdAt') }}</span>
              <span>{{ formatDate(result.created_at) }}</span>
            </div>
          </div>

          <Divider />

          <h3 class="items-title">{{ t('orderQuery.items') }}</h3>
          <DataTable :value="result.items" size="small">
            <Column style="width: 80px">
              <template #body="{ data }">
                <RouterLink :to="`/products/${data.product_id}`">
                  <img
                    :src="resolveImage(data.image_url)"
                    class="item-thumb"
                    :alt="data.product_name" />
                </RouterLink>
              </template>
            </Column>
            <Column :header="t('admin.name')">
              <template #body="{ data }">
                <RouterLink :to="`/products/${data.product_id}`" class="product-link">
                  {{ data.product_name }}
                </RouterLink>
              </template>
            </Column>
            <Column :header="t('product.weight')">
              <template #body="{ data }">{{ data.weight }} oz t</template>
            </Column>
            <Column :header="t('cart.quantity')" field="quantity" />
            <Column :header="t('cart.unitPrice')">
              <template #body="{ data }">NT$ {{ data.locked_unit_price.toLocaleString() }}</template>
            </Column>
            <Column :header="t('cart.subtotal')">
              <template #body="{ data }">
                NT$ {{ (data.locked_unit_price * data.quantity).toLocaleString() }}
              </template>
            </Column>
          </DataTable>

          <div class="price-summary">
            <div class="price-row">
              <span>{{ t('checkout.subtotal') }}</span>
              <span>NT$ {{ (result.total_locked_price - (result.ship_fee || 0)).toLocaleString() }}</span>
            </div>
            <div class="price-row">
              <span>{{ t('checkout.shipFee') }}</span>
              <span>{{ result.ship_fee > 0 ? 'NT$ ' + result.ship_fee.toLocaleString() : t('checkout.shipFreeText')
              }}</span>
            </div>
            <div class="price-row total-row">
              <span>{{ t('checkout.grandTotal') }}</span>
              <strong>NT$ {{ result.total_locked_price.toLocaleString() }}</strong>
            </div>
          </div>

          <!-- 匯款回報（status=0 才顯示） -->
          <div v-if="result.status === 0" class="remit-section">
            <div class="remit-header">
              <i class="pi pi-send" />
              {{ t('payment.title') }}
            </div>
            <p class="remit-hint">請於 30 分鐘內完成匯款，並輸入末五碼通知我們確認。</p>
            <div v-if="!remitSuccess" class="remit-form">
              <InputMask v-model="remitInfo" mask="99999" :placeholder="t('payment.remitPlaceholder')"
                class="remit-input" />
              <Button :label="t('payment.submit')" icon="pi pi-send" :loading="remitting"
                :disabled="remitInfo.replace(/[_]/g, '').length < 5" @click="submitRemit" />
            </div>
            <div v-else class="remit-success">
              <i class="pi pi-check-circle" />
              {{ t('payment.success') }}
            </div>
          </div>

          <div class="result-footer">
            <Button label="← 返回" text @click="result = null; notFound = false; remitInfo = ''; remitSuccess = false" />
            <Button v-if="result.email" :label="t('orderQuery.resendEmail')" icon="pi pi-envelope" severity="secondary"
              :loading="resending" @click="resendEmail" />
          </div>
        </div>
      </div>
    </main>

    <Footer />

    <!-- Scroll to top button -->
    <button v-show="showScrollTop" class="scroll-top-btn" @click="scrollToTop" aria-label="回到頂部">
      <i class="pi pi-chevron-up" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useScrollTop } from '@/composables/useScrollTop'
import { useI18n } from 'vue-i18n'
import { RouterLink } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import InputMask from 'primevue/inputmask'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Divider from 'primevue/divider'
import Navbar from '@/components/shared/Navbar.vue'
import Footer from '@/components/shared/Footer.vue'
import { resolveImage } from '@/api/image'
import api from '@/api'

const { t } = useI18n()
const toast = useToast()
const { showScrollTop, scrollToTop } = useScrollTop()

const orderNo = ref('')
const loading = ref(false)
const resending = ref(false)
const notFound = ref(false)
const result = ref<any>(null)
const myOrders = ref<string[]>(JSON.parse(localStorage.getItem('my_orders') || '[]'))

// 匯款回報
const remitInfo = ref('')
const remitting = ref(false)
const remitSuccess = ref(false)

async function search() {
  const no = orderNo.value.trim()
  if (!no) return
  loading.value = true
  notFound.value = false
  result.value = null

  try {
    const res = await api.get(`/api/v1/orders/${no}`)
    result.value = res.data
  } catch {
    notFound.value = true
  } finally {
    loading.value = false
  }
}

function quickSearch(no: string) {
  orderNo.value = no
  search()
}

async function submitRemit() {
  if (!result.value) return
  remitting.value = true
  try {
    await api.post('/api/v1/report-payment', {
      order_no: result.value.order_no,
      remit_info: remitInfo.value,
    })
    remitSuccess.value = true
    result.value.status = 1
    toast.add({ severity: 'success', summary: t('payment.title'), detail: t('payment.success'), life: 3000 })
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 4000 })
  } finally {
    remitting.value = false
  }
}

async function resendEmail() {
  if (!result.value) return
  resending.value = true
  try {
    await api.post(`/api/v1/orders/${result.value.order_no}/resend-email`)
    toast.add({ severity: 'success', summary: 'Email', detail: t('orderQuery.resendSuccess'), life: 3000 })
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    resending.value = false
  }
}

function statusSeverity(status: number) {
  if (status === 0) return 'warn'
  if (status === 1) return 'info'
  if (status === 2) return 'success'
  if (status === 3) return 'danger'
  if (status === 4) return 'warn'
  if (status === 5) return 'success'
  return 'secondary'
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString()
}

</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
}

.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2.5rem 1.5rem;
  flex: 1;
}

.query-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.page-title {
  font-size: 1.4rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
}

.search-row {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.search-input {
  flex: 1;
}

@media (max-width: 480px) {
  .search-btn :deep(.p-button-label) {
    display: none;
  }
  .search-btn :deep(.p-button-icon) {
    margin: 0;
  }
}

.history-section {
  margin-bottom: 1.5rem;
}

.history-title {
  font-size: 0.85rem;
  color: #999;
  display: flex;
  align-items: center;
  gap: 0.4rem;
  margin-bottom: 0.75rem;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.6rem 1rem;
  border: 1px solid #eee;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.history-item:hover {
  background: #f0f7ff;
  border-color: #90caf9;
}

.history-item code {
  font-size: 0.95rem;
  color: #2c3e50;
}

.history-item .pi-chevron-right {
  font-size: 0.75rem;
  color: #bbb;
}

.not-found {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #e74c3c;
  padding: 1rem;
  background: #fdf3f3;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.result-box {
  margin-top: 0.5rem;
}

.result-header {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  margin-bottom: 0.5rem;
}

.meta-item {
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
}

.meta-label {
  font-size: 0.8rem;
  color: #999;
}

.order-no {
  font-size: 1rem;
  font-weight: 700;
  color: #2c3e50;
  user-select: all;
}

.locked-price {
  font-size: 1.2rem;
  font-weight: 700;
  color: #c0392b;
}

.items-title {
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
}

.price-summary {
  margin-top: 1rem;
  border-top: 1px solid #eee;
  padding-top: 0.75rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.price-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
  color: #666;
}

.total-row {
  font-size: 1rem;
  font-weight: 600;
  color: #1a1a2e;
  border-top: 1px solid #eee;
  padding-top: 0.5rem;
  margin-top: 0.25rem;
}

.total-row strong {
  color: #c0392b;
  font-size: 1.1rem;
}

.item-thumb {
  width: 52px;
  aspect-ratio: 4 / 5;
  object-fit: cover;
  border-radius: 6px;
  display: block;
}

.product-link {
  color: #1a1a2e;
  font-weight: 500;
  text-decoration: none;
}

.product-link:hover {
  color: #0f3460;
  text-decoration: underline;
}

/* 匯款回報 */
.remit-section {
  margin-top: 1.25rem;
  padding: 1.25rem 1.5rem;
  background: #faf9f7;
  border: 1px solid #e8e3dc;
  border-radius: 8px;
}

.remit-header {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1rem;
  font-weight: 600;
  color: #0f172a;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.4rem;
}

.remit-hint {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 0.95rem;
  color: #64748b;
  margin-bottom: 1rem;
}

.remit-form {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  flex-wrap: wrap;
}

.remit-input {
  width: 160px;
}

.remit-success {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #15803d;
  font-weight: 600;
  font-size: 0.95rem;
}

.result-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1rem;
  flex-wrap: wrap;
  gap: 0.5rem;
}
</style>

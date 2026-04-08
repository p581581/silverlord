<template>
  <div>
    <h1 class="page-title">{{ t('admin.orders') }}</h1>

    <!-- 搜尋列 -->
    <div class="search-bar">
      <InputText
        v-model="searchQ"
        :placeholder="t('admin.orderNo') + '...'"
        class="search-input"
        @keyup.enter="fetchOrders"
      />
      <Select
        v-model="searchStatus"
        :options="statusOptions"
        option-label="label"
        option-value="value"
        :placeholder="t('admin.status')"
        show-clear
        class="status-select"
        @change="fetchOrders"
      />
      <Button :label="t('orderQuery.search')" icon="pi pi-search" @click="fetchOrders" />
      <Button icon="pi pi-refresh" severity="secondary" text @click="resetSearch" />
    </div>

    <div class="section-card">
      <DataTable
        :value="orders"
        :loading="loading"
        striped-rows
        :row-class="rowClass"
        @row-click="openDetail"
        paginator
        :rows="20"
      >
        <Column :header="t('admin.orderNo')" field="order_no">
          <template #body="{ data }">
            <code class="order-no">{{ data.order_no }}</code>
          </template>
        </Column>
        <Column :header="t('admin.totalLockedPrice')">
          <template #body="{ data }">
            <span class="price-tag">NT$ {{ data.total_locked_price.toLocaleString() }}</span>
          </template>
        </Column>
        <Column header="Email">
          <template #body="{ data }">{{ data.email || '—' }}</template>
        </Column>
        <Column :header="t('admin.remitInfo')">
          <template #body="{ data }">{{ data.remit_info || '—' }}</template>
        </Column>
        <Column :header="t('admin.status')">
          <template #body="{ data }">
            <Tag
              :value="t(`admin.statusList.${data.status}`)"
              :severity="statusSeverity(data.status)"
            />
          </template>
        </Column>
        <Column :header="t('admin.createdAt')">
          <template #body="{ data }">{{ formatDate(data.created_at) }}</template>
        </Column>
        <Column style="width: 260px">
          <template #body="{ data }">
            <div class="action-btns">
              <Button
                v-if="data.status === 1"
                :label="t('admin.confirmOrder')"
                size="small"
                icon="pi pi-check"
                outlined
                @click.stop="confirmPayment(data)"
              />
              <Button
                v-if="data.status === 1 || data.status === 4"
                :label="t('admin.rejectOrder')"
                size="small"
                severity="warning"
                icon="pi pi-undo"
                outlined
                @click.stop="rejectPayment(data)"
              />
              <Button
                v-if="data.status === 2"
                :label="t('admin.completeOrder')"
                size="small"
                severity="success"
                icon="pi pi-check-circle"
                outlined
                @click.stop="markComplete(data)"
              />
              <Button
                v-if="[0, 1, 2].includes(data.status)"
                :label="t('admin.abnormalOrder')"
                size="small"
                severity="warning"
                icon="pi pi-exclamation-triangle"
                outlined
                @click.stop="markAbnormal(data)"
              />
              <Button
                v-if="data.status !== 3 && data.status !== 1 && data.status !== 5"
                :label="t('admin.invalidateOrder')"
                size="small"
                severity="danger"
                icon="pi pi-ban"
                outlined
                @click.stop="confirmInvalidate(data)"
              />
              <Button
                icon="pi pi-trash"
                size="small"
                outlined
                severity="danger"
                @click.stop="confirmDelete(data)"
              />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 訂單詳情 Dialog -->
    <Dialog
      v-model:visible="detailVisible"
      :header="t('admin.orderNo') + ': ' + selectedOrder?.order_no"
      :style="{ width: '600px' }"
      modal
      :draggable="false"
    >
      <div v-if="selectedOrder" class="detail-content">
        <div class="detail-meta">
          <div>
            {{ t('admin.totalLockedPrice') }}:
            <strong>NT$ {{ selectedOrder.total_locked_price.toLocaleString() }}</strong>
          </div>
          <div>{{ t('admin.name') }}：<strong>{{ selectedOrder.name || '—' }}</strong></div>
          <div>電話：<code>{{ selectedOrder.phone || '—' }}</code></div>
          <div>Email: <code>{{ selectedOrder.email || '—' }}</code></div>
          <div>
            出貨方式：
            <strong>{{
              selectedOrder.ship_method === '711' ? '7-11 店到店'
              : selectedOrder.ship_method === 'pickup' ? '自取'
              : '郵寄'
            }}</strong>
            <span v-if="selectedOrder.ship_method === '711'">
              — 門市：<code>{{ selectedOrder.ship_store || '—' }}</code>
            </span>
            <span v-else-if="selectedOrder.ship_method === 'mail'">
              — 地址：<code>{{ selectedOrder.ship_address || '—' }}</code>
            </span>
          </div>
          <div>{{ t('admin.remitInfo') }}: <code>{{ selectedOrder.remit_info || '—' }}</code></div>
          <div>
            {{ t('admin.status') }}:
            <Tag
              :value="t(`admin.statusList.${selectedOrder.status}`)"
              :severity="statusSeverity(selectedOrder.status)"
            />
          </div>
        </div>

        <DataTable :value="orderItems" class="mt-3" size="small">
          <Column :header="t('admin.name')" field="product_name" />
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
      </div>

      <template #footer>
        <Button
          v-if="selectedOrder?.status === 1"
          :label="t('admin.confirmOrder')"
          icon="pi pi-check"
          outlined
          @click="confirmPayment(selectedOrder)"
        />
        <Button
          v-if="selectedOrder?.status === 1 || selectedOrder?.status === 4"
          :label="t('admin.rejectOrder')"
          icon="pi pi-undo"
          severity="warning"
          outlined
          @click="rejectPayment(selectedOrder)"
        />
        <Button
          v-if="selectedOrder?.status === 2"
          :label="t('admin.completeOrder')"
          icon="pi pi-check-circle"
          severity="success"
          outlined
          @click="markComplete(selectedOrder)"
        />
        <Button
          v-if="[0, 1, 2].includes(selectedOrder?.status)"
          :label="t('admin.abnormalOrder')"
          icon="pi pi-exclamation-triangle"
          severity="warning"
          outlined
          @click="markAbnormal(selectedOrder)"
        />
        <Button
          v-if="selectedOrder?.status !== 3 && selectedOrder?.status !== 1 && selectedOrder?.status !== 5"
          :label="t('admin.invalidateOrder')"
          icon="pi pi-ban"
          severity="danger"
          outlined
          @click="confirmInvalidate(selectedOrder)"
        />
        <Button
          :label="t('admin.deleteProduct')"
          icon="pi pi-trash"
          severity="danger"
          outlined
          @click="confirmDelete(selectedOrder)"
        />
      </template>
    </Dialog>

    <ConfirmDialog />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import ConfirmDialog from 'primevue/confirmdialog'
import api from '@/api'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const orders = ref<any[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const selectedOrder = ref<any>(null)
const orderItems = ref<any[]>([])

const searchQ = ref('')
const searchStatus = ref<number | null>(null)

const statusOptions = [
  { label: t('admin.statusList.0'), value: 0 },
  { label: t('admin.statusList.1'), value: 1 },
  { label: t('admin.statusList.2'), value: 2 },
  { label: t('admin.statusList.3'), value: 3 },
  { label: t('admin.statusList.4'), value: 4 },
  { label: t('admin.statusList.5'), value: 5 },
]

async function fetchOrders() {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (searchQ.value) params.q = searchQ.value
    if (searchStatus.value !== null) params.status = searchStatus.value
    const res = await api.get('/api/admin/orders', { params })
    orders.value = res.data
  } finally {
    loading.value = false
  }
}

function resetSearch() {
  searchQ.value = ''
  searchStatus.value = null
  fetchOrders()
}

async function openDetail(event: any) {
  const order = event.data
  selectedOrder.value = order
  detailVisible.value = true
  try {
    const res = await api.get(`/api/admin/orders/${order.id}`)
    selectedOrder.value = res.data.order
    orderItems.value = res.data.items
  } catch {}
}

async function confirmPayment(order: any) {
  try {
    await api.patch(`/api/admin/orders/${order.id}/status`, { status: 2 })
    toast.add({ severity: 'success', summary: t('admin.confirmOrder'), detail: 'OK', life: 2000 })
    detailVisible.value = false
    await fetchOrders()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

async function rejectPayment(order: any) {
  try {
    await api.patch(`/api/admin/orders/${order.id}/status`, { status: 0 })
    toast.add({ severity: 'warn', summary: t('admin.rejectOrder'), detail: 'OK', life: 2000 })
    detailVisible.value = false
    await fetchOrders()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

function confirmInvalidate(order: any) {
  confirm.require({
    message: `將訂單「${order.order_no}」標記為不成立？庫存將自動退回。`,
    header: t('admin.invalidateOrder'),
    icon: 'pi pi-exclamation-triangle',
    acceptProps: { severity: 'danger' },
    accept: () => invalidateOrder(order),
  })
}

async function invalidateOrder(order: any) {
  try {
    await api.patch(`/api/admin/orders/${order.id}/status`, { status: 3 })
    toast.add({ severity: 'error', summary: t('admin.invalidateOrder'), detail: 'OK', life: 2000 })
    detailVisible.value = false
    await fetchOrders()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

async function markAbnormal(order: any) {
  try {
    await api.patch(`/api/admin/orders/${order.id}/status`, { status: 4 })
    toast.add({ severity: 'warn', summary: t('admin.abnormalOrder'), detail: 'OK', life: 2000 })
    detailVisible.value = false
    await fetchOrders()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

async function markComplete(order: any) {
  try {
    await api.patch(`/api/admin/orders/${order.id}/status`, { status: 5 })
    toast.add({ severity: 'success', summary: t('admin.completeOrder'), detail: 'OK', life: 2000 })
    detailVisible.value = false
    await fetchOrders()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

function confirmDelete(order: any) {
  confirm.require({
    message: `刪除訂單「${order.order_no}」？此操作無法還原。`,
    header: '確認刪除',
    icon: 'pi pi-exclamation-triangle',
    acceptProps: { severity: 'danger' },
    accept: () => deleteOrder(order.id),
  })
}

async function deleteOrder(id: number) {
  try {
    await api.delete(`/api/admin/orders/${id}`)
    toast.add({ severity: 'success', summary: '已刪除', detail: 'OK', life: 2000 })
    detailVisible.value = false
    await fetchOrders()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

function statusSeverity(status: number) {
  if (status === 0) return 'warn'
  if (status === 1) return 'info'
  if (status === 2) return 'success'
  if (status === 3) return 'danger'
  if (status === 4) return 'warn'   // 異常：橘色
  if (status === 5) return 'success' // 已完成：綠色
  return 'secondary'
}

function rowClass(data: any) {
  if (data.status === 1) return 'row-notified'
  if (data.status === 3) return 'row-invalid'
  if (data.status === 4) return 'row-abnormal'
  if (data.status === 5) return 'row-completed'
  return ''
}

function formatDate(iso: string) {
  return new Date(iso).toLocaleString()
}

onMounted(fetchOrders)
</script>

<style scoped>
.page-title { font-size: 1.5rem; font-weight: 700; margin-bottom: 1.5rem; }
.search-bar {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}
.search-input { flex: 1; min-width: 200px; }
.status-select { width: 150px; }
.section-card {
  background: white;
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}
.order-no { font-size: 0.9rem; color: #2c3e50; }
.price-tag { color: #c0392b; font-weight: 600; }
.action-btns { display: flex; gap: 0.25rem; align-items: center; }
.detail-content { display: flex; flex-direction: column; gap: 0.75rem; }
.detail-meta { display: flex; flex-direction: column; gap: 0.5rem; font-size: 0.95rem; }
:deep(.row-notified) { background: #fff8e1 !important; }
:deep(.row-invalid) { background: #fdf2f2 !important; opacity: 0.7; }
:deep(.row-abnormal) { background: #fff3e0 !important; }
:deep(.row-completed) { background: #f0fdf4 !important; }
</style>

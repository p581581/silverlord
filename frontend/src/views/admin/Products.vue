<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">{{ t('admin.products') }}</h1>
      <Button :label="t('admin.addProduct')" icon="pi pi-plus" @click="openDialog()" />
    </div>

    <div class="section-card">
      <DataTable :value="products" :loading="loading" striped-rows>
        <Column :header="t('admin.image')" style="width: 80px">
          <template #body="{ data }">
            <img :src="resolveImage(data.image_url)" class="product-thumb" :alt="data.name" />
          </template>
        </Column>
        <Column :header="t('admin.name')" field="name" />
        <Column :header="t('admin.productType')">
          <template #body="{ data }">
            <Tag
              :value="data.product_type === 'pre_order' ? t('product.preOrder') : t('product.inStock')"
              :severity="data.product_type === 'pre_order' ? 'info' : 'success'"
            />
          </template>
        </Column>
        <Column :header="t('admin.stockQty')">
          <template #body="{ data }">
            <span v-if="data.product_type === 'in_stock'">
              <span :class="data.stock_qty === 0 ? 'stock-zero' : ''">{{ data.stock_qty }}</span>
            </span>
            <span v-else class="text-muted">—</span>
          </template>
        </Column>
        <Column :header="t('admin.weight')">
          <template #body="{ data }">{{ data.weight }} oz t</template>
        </Column>
        <Column :header="t('admin.premiumPercent')">
          <template #body="{ data }">{{ data.premium_percent }}%</template>
        </Column>
        <Column :header="t('admin.laborCost')">
          <template #body="{ data }">NT$ {{ data.labor_cost }}</template>
        </Column>
        <Column header="Price">
          <template #body="{ data }">
            <span class="price-tag">NT$ {{ data.price?.toLocaleString() ?? '-' }}</span>
          </template>
        </Column>
        <Column :header="t('admin.isActive')">
          <template #body="{ data }">
            <Tag
              :value="data.is_active ? t('admin.isActive') : t('product.outOfStock')"
              :severity="data.is_active ? 'success' : 'secondary'"
            />
          </template>
        </Column>
        <Column :header="t('admin.config')" style="width: 180px">
          <template #body="{ data }">
            <div class="action-btns">
              <Button icon="pi pi-pencil" size="small" text @click="openDialog(data)" />
              <Button
                v-if="data.product_type === 'in_stock'"
                icon="pi pi-warehouse"
                size="small"
                text
                severity="info"
                :title="t('admin.stockAdjust')"
                @click="openStockModal(data)"
              />
              <Button
                :icon="data.is_active ? 'pi pi-eye-slash' : 'pi pi-eye'"
                size="small"
                text
                :severity="data.is_active ? 'warn' : 'success'"
                @click="toggleProduct(data)"
              />
              <Button
                icon="pi pi-trash"
                size="small"
                text
                severity="danger"
                @click="confirmDelete(data)"
              />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- 新增/編輯 Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editTarget ? t('admin.editProduct') : t('admin.addProduct')"
      :style="{ width: '560px' }"
      modal
      :draggable="false"
    >
      <form @submit.prevent="saveProduct" class="dialog-form">
        <div class="field">
          <label>{{ t('admin.name') }} *</label>
          <InputText v-model="form.name" class="w-full" required />
        </div>
        <div class="field">
          <label>{{ t('admin.description') }}</label>
          <Textarea v-model="form.description" class="w-full" rows="3" />
        </div>

        <!-- 商品性質 -->
        <div class="field">
          <label>{{ t('admin.productType') }} *</label>
          <div class="type-options">
            <div
              class="type-option"
              :class="{ active: form.product_type === 'in_stock' }"
              @click="form.product_type = 'in_stock'"
            >
              <i class="pi pi-box" />
              <span>{{ t('product.inStock') }}</span>
            </div>
            <div
              class="type-option"
              :class="{ active: form.product_type === 'pre_order' }"
              @click="form.product_type = 'pre_order'"
            >
              <i class="pi pi-clock" />
              <span>{{ t('product.preOrder') }}</span>
            </div>
          </div>
        </div>

        <!-- 庫存品：庫存數量 -->
        <div v-if="form.product_type === 'in_stock'" class="field">
          <label>{{ t('admin.stockQty') }}</label>
          <InputNumber v-model="form.stock_qty" :min="0" class="w-full" />
        </div>

        <!-- 預購：預計到貨 -->
        <div v-if="form.product_type === 'pre_order'" class="field">
          <label>{{ t('admin.expectedArrival') }}</label>
          <DatePicker v-model="expectedArrivalDate" date-format="yy-mm-dd" class="w-full" />
        </div>

        <div class="field-row">
          <div class="field">
            <label>{{ t('admin.weight') }} *</label>
            <InputNumber v-model="form.weight" :min="0.001" :min-fraction-digits="3" suffix=" oz t" />
          </div>
          <div class="field">
            <label>{{ t('admin.laborCost') }}</label>
            <InputNumber v-model="form.labor_cost" :min="0" prefix="NT$ " />
          </div>
          <div class="field">
            <label>{{ t('admin.premiumPercent') }}</label>
            <InputNumber v-model="form.premium_percent" :min="0" :min-fraction-digits="1" suffix=" %" />
          </div>
        </div>
        <div class="field">
          <label>{{ t('admin.image') }} <span class="field-hint">（最多 3 張，第 1 張為列表主圖）</span></label>
          <div class="image-slots">
            <div v-for="i in 3" :key="i" class="image-slot">
              <div class="slot-preview">
                <img
                  v-if="slotUrl(i - 1)"
                  :src="resolveImage(slotUrl(i - 1))"
                  class="slot-thumb"
                  :alt="`圖片 ${i}`"
                />
                <div v-else class="slot-empty">
                  <i class="pi pi-image" />
                </div>
                <Button
                  v-if="slotUrl(i - 1)"
                  icon="pi pi-times"
                  class="slot-remove"
                  rounded
                  text
                  severity="danger"
                  size="small"
                  @click="clearSlot(i - 1)"
                />
              </div>
              <FileUpload
                mode="basic"
                accept="image/*"
                :max-file-size="5000000"
                :auto="true"
                :custom-upload="true"
                :choose-label="i === 1 ? t('admin.upload') : '+'"
                class="slot-upload"
                @uploader="(e: any) => onUpload(e, i - 1)"
              />
              <span class="slot-label">{{ i === 1 ? '主圖' : `圖 ${i}` }}</span>
            </div>
          </div>
        </div>

        <!-- 裁切 Dialog -->
        <Dialog
          v-model:visible="cropVisible"
          header="框選圖片範圍"
          :style="{ width: '680px' }"
          modal
          :draggable="false"
          append-to="body"
        >
          <div class="crop-wrap">
            <Cropper
              ref="cropperRef"
              :src="cropSrc"
              :stencil-props="{ aspectRatio: 4 / 5 }"
              class="cropper"
            />
          </div>
          <template #footer>
            <Button label="取消" severity="secondary" outlined @click="cropVisible = false" />
            <Button label="確認裁切" icon="pi pi-check" :loading="cropUploading" @click="confirmCrop" />
          </template>
        </Dialog>
        <div class="field-check">
          <Checkbox v-model="form.is_active" binary input-id="is_active" />
          <label for="is_active">{{ t('admin.isActive') }}</label>
        </div>
      </form>

      <template #footer>
        <Button :label="t('admin.save')" icon="pi pi-save" :loading="saving" @click="saveProduct" />
      </template>
    </Dialog>

    <!-- 庫存調整 Modal -->
    <Dialog
      v-model:visible="stockModalVisible"
      :header="t('admin.stockAdjust')"
      :style="{ width: '360px' }"
      modal
      :draggable="false"
    >
      <div class="stock-modal">
        <p class="stock-product-name">{{ stockTarget?.name }}</p>
        <div class="field">
          <label>{{ t('admin.stockQty') }}</label>
          <InputNumber v-model="stockQtyInput" :min="0" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button :label="t('admin.save')" icon="pi pi-save" :loading="stockSaving" @click="saveStock" />
      </template>
    </Dialog>

    <ConfirmDialog />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import FileUpload from 'primevue/fileupload'
import Checkbox from 'primevue/checkbox'
import Tag from 'primevue/tag'
import DatePicker from 'primevue/datepicker'
import ConfirmDialog from 'primevue/confirmdialog'
import { Cropper } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import api from '@/api'
import { resolveImage } from '@/api/image'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const products = ref<any[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const saving = ref(false)
const editTarget = ref<any>(null)

// 日期 picker 綁定（DatePicker 用 Date 物件）
const expectedArrivalDate = ref<Date | null>(null)

const emptyForm = () => ({
  name: '',
  description: '',
  weight: 1,
  labor_cost: 0,
  premium_percent: 5.0,
  image_url: '',
  image_url_2: '',
  image_url_3: '',
  is_active: true,
  product_type: 'in_stock',
  stock_qty: 0,
  expected_arrival: '',
})
const form = ref(emptyForm())

// 同步 DatePicker → form.expected_arrival
watch(expectedArrivalDate, (d) => {
  if (d) {
    const y = d.getFullYear()
    const m = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    form.value.expected_arrival = `${y}-${m}-${day}`
  } else {
    form.value.expected_arrival = ''
  }
})

// 庫存調整
const stockModalVisible = ref(false)
const stockTarget = ref<any>(null)
const stockQtyInput = ref(0)
const stockSaving = ref(false)

function openStockModal(product: any) {
  stockTarget.value = product
  stockQtyInput.value = product.stock_qty ?? 0
  stockModalVisible.value = true
}

async function saveStock() {
  stockSaving.value = true
  try {
    await api.patch(`/api/admin/products/${stockTarget.value.id}/stock`, {
      stock_qty: stockQtyInput.value,
    })
    toast.add({ severity: 'success', summary: t('admin.stockAdjust'), detail: 'OK', life: 2000 })
    stockModalVisible.value = false
    await fetchProducts()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    stockSaving.value = false
  }
}

async function fetchProducts() {
  loading.value = true
  try {
    const res = await api.get('/api/admin/products')
    products.value = res.data
  } finally {
    loading.value = false
  }
}

function openDialog(product?: any) {
  editTarget.value = product || null
  if (product) {
    form.value = { ...product }
    expectedArrivalDate.value = product.expected_arrival
      ? new Date(product.expected_arrival)
      : null
  } else {
    form.value = emptyForm()
    expectedArrivalDate.value = null
  }
  dialogVisible.value = true
}

async function saveProduct() {
  saving.value = true
  try {
    const payload = { ...form.value }
    if (editTarget.value) {
      await api.put(`/api/admin/products/${editTarget.value.id}`, payload)
    } else {
      await api.post('/api/admin/products', payload)
    }
    toast.add({ severity: 'success', summary: t('admin.save'), detail: 'OK', life: 2000 })
    dialogVisible.value = false
    await fetchProducts()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    saving.value = false
  }
}

async function toggleProduct(product: any) {
  try {
    await api.patch(`/api/admin/products/${product.id}/toggle`)
    await fetchProducts()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

function confirmDelete(product: any) {
  confirm.require({
    message: `刪除「${product.name}」？`,
    header: t('admin.deleteProduct'),
    icon: 'pi pi-exclamation-triangle',
    acceptProps: { severity: 'danger' },
    accept: () => deleteProduct(product.id),
  })
}

async function deleteProduct(id: number) {
  try {
    await api.delete(`/api/admin/products/${id}`)
    await fetchProducts()
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  }
}

// ── 圖片裁切 ──────────────────────────────────────────────
const cropVisible = ref(false)
const cropSrc = ref('')
const cropperRef = ref<any>(null)
const cropUploading = ref(false)
const cropTargetSlot = ref(0)
let pendingFileName = 'crop.jpg'

function slotUrl(slot: number): string {
  if (slot === 0) return form.value.image_url
  if (slot === 1) return form.value.image_url_2
  return form.value.image_url_3
}

function setSlotUrl(slot: number, url: string) {
  if (slot === 0) form.value.image_url = url
  else if (slot === 1) form.value.image_url_2 = url
  else form.value.image_url_3 = url
}

function clearSlot(slot: number) {
  setSlotUrl(slot, '')
}

function onUpload(event: any, slot: number) {
  cropTargetSlot.value = slot
  const file = event.files[0]
  pendingFileName = file.name
  const reader = new FileReader()
  reader.onload = (e) => {
    cropSrc.value = e.target?.result as string
    cropVisible.value = true
  }
  reader.readAsDataURL(file)
}

async function confirmCrop() {
  cropUploading.value = true
  try {
    const { canvas } = cropperRef.value.getResult()
    canvas.toBlob(async (blob: Blob | null) => {
      if (!blob) return
      const formData = new FormData()
      formData.append('file', blob, pendingFileName)
      const res = await api.post('/api/admin/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      setSlotUrl(cropTargetSlot.value, res.data.url)
      toast.add({ severity: 'success', summary: t('admin.upload'), detail: 'OK', life: 2000 })
      cropVisible.value = false
      cropUploading.value = false
    }, 'image/jpeg', 0.92)
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
    cropUploading.value = false
  }
}

onMounted(fetchProducts)
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; }
.page-title { font-size: 1.5rem; font-weight: 700; }
.section-card {
  background: white;
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}
.product-thumb { width: 48px; height: 48px; object-fit: cover; border-radius: 6px; }
.price-tag { color: #c0392b; font-weight: 600; }
.stock-zero { color: #e74c3c; font-weight: 700; }
.text-muted { color: #aaa; }
.action-btns { display: flex; gap: 0.25rem; }
.dialog-form { display: flex; flex-direction: column; gap: 1rem; }
.field { display: flex; flex-direction: column; gap: 0.4rem; }
.field label { font-weight: 600; font-size: 0.85rem; }
.field-row { display: flex; gap: 1rem; flex-wrap: wrap; }
.field-check { display: flex; align-items: center; gap: 0.5rem; }
.field-hint { font-size: 0.78rem; color: #999; font-weight: 400; margin-left: 0.4rem; }
.image-slots { display: flex; gap: 1rem; flex-wrap: wrap; }
.image-slot { display: flex; flex-direction: column; align-items: center; gap: 0.4rem; }
.slot-preview {
  position: relative;
  width: 90px;
  height: 90px;
  border: 2px dashed #e0e0e0;
  border-radius: 8px;
  overflow: visible;
  background: #fafafa;
  display: flex;
  align-items: center;
  justify-content: center;
}
.slot-thumb { width: 90px; height: 90px; object-fit: cover; border-radius: 6px; display: block; }
.slot-empty { display: flex; align-items: center; justify-content: center; color: #ccc; font-size: 1.8rem; width: 100%; height: 100%; }
.slot-remove {
  position: absolute !important;
  top: -8px !important;
  right: -8px !important;
  width: 1.6rem !important;
  height: 1.6rem !important;
  padding: 0 !important;
  background: white !important;
  box-shadow: 0 1px 4px rgba(0,0,0,0.2) !important;
  z-index: 1;
}
.slot-label { font-size: 0.72rem; color: #999; }

.type-options { display: flex; gap: 0.75rem; }
.type-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem;
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  color: #666;
  background: white;
  transition: all 0.15s;
}
.type-option:hover { border-color: #aaa; color: #333; }
.type-option.active { border-color: #1a1a2e; background: #1a1a2e; color: white; }

.stock-modal { display: flex; flex-direction: column; gap: 1rem; padding: 0.5rem 0; }
.stock-product-name { font-weight: 600; font-size: 1rem; color: #1a1a2e; margin: 0; }

.crop-wrap {
  width: 100%;
  height: 420px;
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
}
.cropper {
  width: 100%;
  height: 100%;
}
</style>

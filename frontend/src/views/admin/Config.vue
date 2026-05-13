<template>
  <div>
    <h1 class="page-title">{{ t('admin.config') }}</h1>

    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-send" /> Telegram
      </h2>

      <div class="field">
        <label>Bot Token</label>
        <InputText v-model="form.tg_bot_token" class="w-full" placeholder="123456:ABC-..." />
      </div>

      <div class="field">
        <label>Chat ID <span class="optional">（可新增多組管理者）</span></label>
        <div v-for="(_, idx) in chatIds" :key="idx" class="chat-id-row">
          <InputText
            v-model="chatIds[idx]"
            placeholder="-100123456789"
            class="flex-1" />
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            @click="chatIds.splice(idx, 1)" />
        </div>
        <Button
          label="新增 Chat ID"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          style="align-self: flex-start"
          @click="chatIds.push('')" />
      </div>

      <div class="action-row">
        <Button
          :label="t('admin.save')"
          icon="pi pi-save"
          :loading="saving"
          @click="saveConfig" />
        <Button
          :label="t('admin.testTelegram')"
          icon="pi pi-send"
          severity="secondary"
          :loading="testing"
          @click="testTelegram" />
      </div>
    </div>

    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-building-columns" /> 匯款資訊
      </h2>
      <div class="field-row">
        <div class="field">
          <label>銀行名稱</label>
          <InputText v-model="form.bank_name" class="w-full" placeholder="玉山銀行（808）" />
        </div>
        <div class="field">
          <label>戶名</label>
          <InputText v-model="form.bank_account_name" class="w-full" placeholder="王小明" />
        </div>
      </div>
      <div class="field">
        <label>帳號</label>
        <InputText v-model="form.bank_account" class="w-full" placeholder="0123456789012" />
      </div>
      <Button
        :label="t('admin.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="saveConfig" />
    </div>

    <!-- Banner 上傳 -->
    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-image" /> 首頁 Banner
      </h2>
      <!-- 預覽：按優先順序顯示 -->
      <div class="banner-preview-wrap">
        <img
          v-if="form.banner_url || form.banner_custom_url"
          :src="resolveImage(form.banner_url || form.banner_custom_url)"
          class="banner-preview"
          alt="banner"
        />
        <div v-else class="banner-empty">
          <i class="pi pi-image" />
          <span>預設 Banner（Unsplash）</span>
        </div>
      </div>

      <!-- 優先順序說明 -->
      <div class="banner-priority">
        <span :class="['priority-tag', form.banner_url ? 'active' : '']">① 上傳圖片</span>
        <i class="pi pi-angle-right" />
        <span :class="['priority-tag', !form.banner_url && form.banner_custom_url ? 'active' : '']">② 自訂 URL</span>
        <i class="pi pi-angle-right" />
        <span :class="['priority-tag', !form.banner_url && !form.banner_custom_url ? 'active' : '']">③ 預設</span>
      </div>

      <!-- 上傳區 -->
      <div class="field" style="margin-top: 1rem">
        <label>① 上傳圖片 <span class="optional">（優先顯示）</span></label>
        <div class="action-row">
          <FileUpload
            mode="basic"
            accept="image/*"
            :max-file-size="10000000"
            :auto="true"
            :custom-upload="true"
            choose-label="上傳 Banner"
            @uploader="onBannerUpload"
          />
          <Button
            v-if="form.banner_url"
            label="移除上傳圖片"
            icon="pi pi-trash"
            severity="danger"
            outlined
            size="small"
            :loading="saving"
            @click="removeBanner"
          />
        </div>
        <small v-if="form.banner_url" class="hint">目前：{{ form.banner_url }}</small>
      </div>

      <!-- 自訂 URL 區 -->
      <div class="field">
        <label>② 自訂 URL <span class="optional">（無上傳圖片時使用）</span></label>
        <InputText v-model="form.banner_custom_url" class="w-full" placeholder="https://example.com/banner.jpg" />
      </div>
      <Button
        :label="t('admin.save')"
        icon="pi pi-save"
        :loading="saving"
        @click="saveConfig"
      />

      <!-- 裁切 Dialog -->
      <Dialog
        v-model:visible="bannerCropVisible"
        header="框選 Banner 範圍"
        :style="{ width: '760px' }"
        modal
        :draggable="false"
      >
        <div class="banner-crop-wrap">
          <Cropper
            ref="bannerCropperRef"
            :src="bannerCropSrc"
            :stencil-props="{ aspectRatio: 16 / 9 }"
            class="cropper"
          />
        </div>
        <template #footer>
          <Button label="取消" severity="secondary" outlined @click="bannerCropVisible = false" />
          <Button label="確認裁切" icon="pi pi-check" :loading="bannerCropUploading" @click="confirmBannerCrop" />
        </template>
      </Dialog>
    </div>

    <div class="section-card">
      <h2 class="section-title">
        <i class="pi pi-info-circle" /> 系統資訊
      </h2>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">{{ t('admin.manualOffset') }}</span>
          <span class="info-value">{{ form.manual_offset >= 0 ? '+' : '' }}{{ form.manual_offset }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">{{ t('admin.manualExRate') }}</span>
          <span class="info-value" :class="{ active: form.manual_ex_rate > 0 }">
            {{ form.manual_ex_rate > 0 ? form.manual_ex_rate : '停用' }}
          </span>
        </div>
      </div>
      <small class="hint">在「控盤中心」調整報價修正與匯率設定。</small>
      <small class="hint smtp-hint">SMTP 信箱設定請至伺服器 <code>.env</code> 檔案修改。</small>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import FileUpload from 'primevue/fileupload'
import { Cropper } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import api from '@/api'
import { resolveImage } from '@/api/image'

const { t } = useI18n()
const toast = useToast()

const form = ref({
  tg_bot_token: '',
  tg_chat_ids: '',
  manual_offset: 0,
  manual_ex_rate: 0,
  bank_name: '',
  bank_account: '',
  bank_account_name: '',
  banner_url: '',
  banner_custom_url: '',
})
// 將逗號分隔字串拆成陣列供 UI 編輯
const chatIds = ref<string[]>([''])
const saving = ref(false)
const testing = ref(false)

onMounted(async () => {
  try {
    const res = await api.get('/api/admin/config')
    form.value = res.data
    chatIds.value = res.data.tg_chat_ids
      ? res.data.tg_chat_ids.split(',').map((s: string) => s.trim())
      : ['']
  } catch { }
})

async function saveConfig() {
  saving.value = true
  form.value.tg_chat_ids = chatIds.value.filter(id => id.trim() !== '').join(',')
  try {
    const { data: current } = await api.get('/api/admin/config')
    await api.put('/api/admin/config', {
      ...current,
      tg_bot_token: form.value.tg_bot_token,
      tg_chat_ids: form.value.tg_chat_ids,
      bank_name: form.value.bank_name,
      bank_account: form.value.bank_account,
      bank_account_name: form.value.bank_account_name,
      banner_url: form.value.banner_url,
      banner_custom_url: form.value.banner_custom_url,
    })
    toast.add({ severity: 'success', summary: t('admin.save'), detail: 'OK', life: 2000 })
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    saving.value = false
  }
}

// ── Banner 裁切上傳 ───────────────────────────────────────
const bannerCropVisible = ref(false)
const bannerCropSrc = ref('')
const bannerCropperRef = ref<any>(null)
const bannerCropUploading = ref(false)
let bannerPendingFileName = 'banner.jpg'

function onBannerUpload(event: any) {
  const file = event.files[0]
  bannerPendingFileName = file.name
  const reader = new FileReader()
  reader.onload = (e) => {
    bannerCropSrc.value = e.target?.result as string
    bannerCropVisible.value = true
  }
  reader.readAsDataURL(file)
}

async function confirmBannerCrop() {
  bannerCropUploading.value = true
  try {
    const { canvas } = bannerCropperRef.value.getResult()
    canvas.toBlob(async (blob: Blob | null) => {
      if (!blob) return
      const fd = new FormData()
      fd.append('file', blob, bannerPendingFileName)
      const res = await api.post('/api/admin/upload', fd, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      form.value.banner_url = res.data.url
      await saveConfig()
      bannerCropVisible.value = false
      bannerCropUploading.value = false
    }, 'image/jpeg', 0.92)
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
    bannerCropUploading.value = false
  }
}

async function removeBanner() {
  form.value.banner_url = ''
  await saveConfig()
}

async function testTelegram() {
  testing.value = true
  try {
    await api.post('/api/admin/config/test-tg')
    toast.add({ severity: 'success', summary: 'Telegram', detail: '✅ 傳送成功', life: 3000 })
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    testing.value = false
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

.field {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  margin-bottom: 1.25rem;
}

.field label {
  font-weight: 600;
  font-size: 0.9rem;
}

.action-row {
  display: flex;
  gap: 1rem;
}

.info-grid {
  display: flex;
  gap: 2rem;
  margin-bottom: 0.75rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-label {
  font-size: 0.8rem;
  color: #999;
}

.info-value {
  font-size: 1.1rem;
  font-weight: 600;
}

.info-value.active {
  color: #e67e22;
}

.hint {
  color: #999;
  font-size: 0.8rem;
  display: block;
}

.smtp-hint {
  margin-top: 0.4rem;
}

.smtp-hint code {
  background: #f0f0f0;
  padding: 0 4px;
  border-radius: 3px;
}

.field-row { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 0; }
.field-row .field { flex: 1; min-width: 160px; }

.banner-preview-wrap {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  background: #f0f0f0;
  border: 2px dashed #e0e0e0;
}
.banner-preview {
  width: 100%;
  aspect-ratio: 16 / 9;
  object-fit: cover;
  display: block;
}
.banner-empty {
  width: 100%;
  aspect-ratio: 16 / 9;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: #bbb;
  font-size: 0.9rem;
}
.banner-empty i { font-size: 2rem; }

.banner-priority {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  margin-top: 0.75rem;
  font-size: 0.82rem;
}
.priority-tag {
  padding: 0.2rem 0.6rem;
  border-radius: 999px;
  background: #f0f0f0;
  color: #999;
  font-weight: 500;
}
.priority-tag.active {
  background: #1a1a2e;
  color: white;
}

.banner-crop-wrap {
  width: 100%;
  height: 400px;
  background: #1a1a1a;
  border-radius: 8px;
  overflow: hidden;
}
.cropper { width: 100%; height: 100%; }
.chat-id-row {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  margin-bottom: 0.5rem;
}

.optional {
  font-weight: 400;
  color: #999;
  font-size: 0.82rem;
}
</style>

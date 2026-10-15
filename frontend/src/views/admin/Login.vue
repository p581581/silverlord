<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <img src="@/assets/img/logo.png" alt="Silver Lord" class="login-logo" />
        <h1>Silver Lord</h1>
        <p>{{ t('admin.login') }}</p>
      </div>

      <form class="login-form" @submit.prevent="submit">
        <div class="field">
          <Password
            v-model="password"
            :feedback="false"
            :placeholder="t('admin.password')"
            :aria-label="t('admin.password')"
            toggle-mask
            class="login-password"
            input-class="login-password-input"
            required />
        </div>

        <Button
          type="submit"
          :label="t('admin.loginBtn')"
          icon="pi pi-sign-in"
          :loading="loading" />
      </form>

      <div class="locale-row">
        <Select
          v-model="currentLocale"
          :options="localeOptions"
          option-label="label"
          option-value="value"
          size="small"
          @change="onLocaleChange" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Select from 'primevue/select'
import { useAuthStore } from '@/stores/auth'
import { setLocale, type Locale } from '@/i18n'

const { t, locale } = useI18n()
const router = useRouter()
const toast = useToast()
const auth = useAuthStore()

const password = ref('')
const loading = ref(false)
const currentLocale = ref<Locale>(locale.value as Locale)

const localeOptions = [
  { label: '中文', value: 'zh-TW' },
  { label: 'English', value: 'en' },
  { label: '日本語', value: 'ja' },
]

function onLocaleChange() {
  setLocale(currentLocale.value)
}

async function submit() {
  loading.value = true
  try {
    await auth.login(password.value)
    router.push('/admin/dashboard')
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 3000 })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 60%, #0f3460 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  background: white;
  border-radius: 16px;
  padding: 1.25rem 2.5rem 3.25rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.login-header {
  text-align: center;
  margin-bottom: 1.5rem;
}

.login-logo {
  display: block;
  height: 7.5rem;
  width: auto;
  margin: 0 auto;
}

.login-header h1 {
  font-size: 1.125rem;
  font-weight: 700;
  line-height: 1.25;
  margin: 0.5rem 0 0;
}

.login-header p {
  color: #666;
  font-size: 1.125rem;
  line-height: 1.25;
  margin: 0.625rem 0 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.field {
  width: 15rem;
  max-width: calc(100% - 2.5rem);
  margin-bottom: 1.25rem;
}

.login-password {
  position: relative;
  width: 100%;
}

.login-password :deep(.login-password-input) {
  width: 100%;
  height: 2.25rem;
  font-size: 0.9rem;
  padding: 0.5rem 0.75rem;
}

.login-password :deep(.login-password-input::placeholder) {
  color: #111;
  font-weight: 600;
  opacity: 1;
}

.login-password :deep(.p-password-toggle-mask-icon) {
  position: absolute;
  right: -1.75rem;
  top: 50%;
  margin-top: 0;
  transform: translateY(-50%);
}

.locale-row {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

.locale-row :deep(.p-select-dropdown) {
  width: 2rem;
}
</style>

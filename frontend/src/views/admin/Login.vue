<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <span class="login-logo">🥈</span>
        <h1>Silver Lord</h1>
        <p>{{ t('admin.login') }}</p>
      </div>

      <form @submit.prevent="submit">
        <div class="field">
          <label>{{ t('admin.password') }}</label>
          <Password
            v-model="password"
            :feedback="false"
            toggle-mask
            class="w-full"
            input-class="w-full"
            required />
        </div>

        <Button
          type="submit"
          :label="t('admin.loginBtn')"
          icon="pi pi-sign-in"
          :loading="loading"
          class="w-full" />
      </form>

      <div class="locale-row">
        <Select
          v-model="currentLocale"
          :options="localeOptions"
          option-label="label"
          option-value="value"
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
  padding: 2.5rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.4);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-logo {
  font-size: 3rem;
}

.login-header h1 {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0.5rem 0 0.25rem;
}

.login-header p {
  color: #666;
  font-size: 0.9rem;
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

.locale-row {
  display: flex;
  justify-content: center;
  margin-top: 1.5rem;
}
</style>

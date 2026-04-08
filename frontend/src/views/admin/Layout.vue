<template>
  <div class="admin-layout">
    <!-- 側邊欄 -->
    <aside class="sidebar">
      <div class="sidebar-brand">🥈 Silver Lord</div>

      <nav class="sidebar-nav">
        <RouterLink to="/admin/dashboard" class="nav-item">
          <i class="pi pi-chart-line" /> {{ t('admin.dashboard') }}
        </RouterLink>
        <RouterLink to="/admin/products" class="nav-item">
          <i class="pi pi-box" /> {{ t('admin.products') }}
        </RouterLink>
        <RouterLink to="/admin/orders" class="nav-item">
          <i class="pi pi-list" /> {{ t('admin.orders') }}
        </RouterLink>
        <RouterLink to="/admin/shipping" class="nav-item">
          <i class="pi pi-truck" /> {{ t('admin.shipping') }}
        </RouterLink>
        <RouterLink to="/admin/config" class="nav-item">
          <i class="pi pi-cog" /> {{ t('admin.config') }}
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <Select
          v-model="currentLocale"
          :options="localeOptions"
          option-label="label"
          option-value="value"
          class="w-full mb-2"
          @change="onLocaleChange"
        />
        <Button
          :label="t('admin.logout')"
          icon="pi pi-sign-out"
          severity="secondary"
          text
          class="w-full"
          @click="logout"
        />
      </div>
    </aside>

    <!-- 主內容 -->
    <main class="admin-main">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import Select from 'primevue/select'
import { useAuthStore } from '@/stores/auth'
import { setLocale, type Locale } from '@/i18n'

const { t, locale } = useI18n()
const router = useRouter()
const auth = useAuthStore()

const currentLocale = ref<Locale>(locale.value as Locale)
const localeOptions = [
  { label: '中文', value: 'zh-TW' },
  { label: 'English', value: 'en' },
  { label: '日本語', value: 'ja' },
]

function onLocaleChange() {
  setLocale(currentLocale.value)
}

function logout() {
  auth.logout()
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}
.sidebar {
  width: 220px;
  background: #1a1a2e;
  color: white;
  display: flex;
  flex-direction: column;
  padding: 1.5rem 1rem;
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  z-index: 50;
}
.sidebar-brand {
  font-size: 1.1rem;
  font-weight: 700;
  padding: 0.5rem 0 1.5rem;
  border-bottom: 1px solid #333;
  margin-bottom: 1rem;
}
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}
.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.65rem 0.75rem;
  border-radius: 8px;
  color: #aaa;
  text-decoration: none;
  font-size: 0.9rem;
  transition: background 0.2s, color 0.2s;
}
.nav-item:hover,
.nav-item.router-link-active {
  background: #0f3460;
  color: white;
}
.sidebar-footer {
  border-top: 1px solid #333;
  padding-top: 1rem;
}
.admin-main {
  margin-left: 220px;
  flex: 1;
  background: #f5f5f5;
  min-height: 100vh;
  padding: 2rem;
}
</style>

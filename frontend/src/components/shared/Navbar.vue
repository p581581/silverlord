<template>
  <div class="nav-wrapper">
    <!-- 即時報價 Ticker -->
    <div class="ticker-bar">
      <div class="ticker-inner">
        <span class="ticker-live">
          <i class="pi pi-circle-fill ticker-dot" />
          <span class="ticker-live-text">即時市場報價</span>
        </span>
        <span class="ticker-sep">|</span>
        <span>
          XAG/USD：
          <b class="ticker-val">{{ xagUsd }}</b>
        </span>
        <span class="ticker-sep">|</span>
        <span>
          USD/TWD：
          <b class="ticker-val">{{ usdTwd }}</b>
        </span>
      </div>
    </div>

    <!-- 主導覽列 -->
    <nav class="navbar">
      <RouterLink to="/" class="navbar-brand">
        <img src="@/assets/img/logo.png" alt="Silver Lord" class="brand-logo" />
        <span class="brand-text">SILVER LORD</span>
      </RouterLink>

      <!-- 桌面導覽 -->
      <div class="navbar-links">
        <RouterLink to="/" class="nav-link nav-hide-xs" :class="{ active: route.path === '/' }">
          {{ t('nav.products') }}
        </RouterLink>
        <RouterLink to="/order-query" class="nav-link nav-hide-xs" :class="{ active: route.path === '/order-query' }">
          {{ t('orderQuery.title') }}
        </RouterLink>

        <RouterLink to="/cart" class="cart-btn">
          <i class="pi pi-shopping-cart" />
          <span v-if="cartCount > 0" class="cart-badge">{{ cartCount }}</span>
        </RouterLink>

        <Select
          v-model="currentLocale"
          :options="localeOptions"
          option-label="label"
          option-value="value"
          class="locale-select nav-hide-xs"
          @change="onLocaleChange" />

        <!-- 漢堡按鈕：xs 才顯示 -->
        <button class="hamburger" aria-label="選單" @click="drawerOpen = true">
          <i class="pi pi-bars" />
        </button>
      </div>
    </nav>

    <!-- Drawer -->
    <Drawer v-model:visible="drawerOpen" position="right" class="nav-drawer"
      :pt="{ root: { style: 'width: 280px; background: #0f172a;' } }">
      <template #header>
        <div class="drawer-header">
          <img src="@/assets/img/logo.png" alt="Silver Lord" class="drawer-logo" />
          <span class="drawer-brand">SILVER LORD</span>
        </div>
      </template>

      <nav class="drawer-nav">
        <RouterLink to="/" class="drawer-link" :class="{ active: route.path === '/' }" @click="drawerOpen = false">
          <i class="pi pi-home" />
          {{ t('nav.products') }}
        </RouterLink>
        <RouterLink to="/order-query" class="drawer-link" :class="{ active: route.path === '/order-query' }"
          @click="drawerOpen = false">
          <i class="pi pi-search" />
          {{ t('orderQuery.title') }}
        </RouterLink>
        <RouterLink to="/cart" class="drawer-link" :class="{ active: route.path === '/cart' }"
          @click="drawerOpen = false">
          <i class="pi pi-shopping-cart" />
          {{ t('nav.cart') }}
          <span v-if="cartCount > 0" class="drawer-badge">{{ cartCount }}</span>
        </RouterLink>
      </nav>

      <div class="drawer-footer">
        <div class="drawer-label">語言 / Language</div>
        <Select v-model="currentLocale" :options="localeOptions" option-label="label" option-value="value"
          class="w-full"
          @change="onLocaleChange" />
      </div>
    </Drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { RouterLink, useRoute } from 'vue-router'
import Select from 'primevue/select'
import Drawer from 'primevue/drawer'
import { useCartStore } from '@/stores/cart'
import { usePriceStore } from '@/stores/price'
import { setLocale, type Locale } from '@/i18n'

const { t, locale } = useI18n()
const route = useRoute()
const cart = useCartStore()
const priceStore = usePriceStore()

const drawerOpen = ref(false)

const cartCount = computed(() => cart.items.reduce((sum, i) => sum + i.quantity, 0))

const xagUsd = computed(() =>
  priceStore.snapshot?.xag_usd ? priceStore.snapshot.xag_usd.toFixed(2) : '—'
)
const usdTwd = computed(() =>
  priceStore.snapshot?.usd_twd ? priceStore.snapshot.usd_twd.toFixed(2) : '—'
)

const currentLocale = ref<Locale>(locale.value as Locale)
const localeOptions = [
  { label: '中文', value: 'zh-TW' },
  { label: 'English', value: 'en' },
  { label: '日本語', value: 'ja' },
]
function onLocaleChange() {
  setLocale(currentLocale.value)
}

onMounted(() => {
  if (!priceStore.snapshot) priceStore.fetchPrice()
})
</script>

<style scoped>
.nav-wrapper {
  position: sticky;
  top: 0;
  z-index: 100;
}

/* ── Ticker ── */
.ticker-bar {
  background: #020617;
  color: #94a3b8;
  font-size: 0.72rem;
  font-family: 'Courier New', monospace;
  letter-spacing: 0.02em;
  padding: 0.4rem 0;
  overflow: hidden;
}

.ticker-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: nowrap;
  overflow-x: auto;
  scrollbar-width: none;
}

.ticker-inner::-webkit-scrollbar {
  display: none;
}

.ticker-live {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: #e2e8f0;
  font-weight: 600;
  font-family: Inter, system-ui, sans-serif;
  white-space: nowrap;
}

.ticker-dot {
  font-size: 0.5rem;
  color: #4ade80;
  animation: pulse-dot 2s infinite;
}

@keyframes pulse-dot {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.3;
  }
}

.ticker-sep {
  color: #334155;
  user-select: none;
}

.ticker-val {
  color: #f8fafc;
  font-weight: 700;
}

/* ── Navbar ── */
.navbar {
  background: #0f172a;
  height: 72px;
  display: flex;
  align-items: center;
  padding: 0 1.5rem;
  box-shadow: 0 1px 0 rgba(255, 255, 255, 0.06);
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  margin-right: auto;
  color: #f1f5f9;
  text-decoration: none;
}

.brand-logo {
  height: 2rem;
  width: auto;
  filter: brightness(0) invert(1);
  opacity: 0.85;
}

.brand-text {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 0.1em;
  line-height: 1.8;
  color: #f1f5f9;
}

.navbar-links {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.nav-link {
  color: #94a3b8;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 0.4rem 0.75rem;
  border-radius: 6px;
  transition: color 0.15s, background 0.15s;
  text-decoration: none;
  white-space: nowrap;
  position: relative;
}

.nav-link:hover {
  color: #f1f5f9;
  background: rgba(255, 255, 255, 0.06);
}

.nav-link.active {
  color: #fff;
  font-weight: 600;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0.75rem;
  right: 0.75rem;
  height: 2px;
  background: white;
  border-radius: 2px;
}

.cart-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  color: #94a3b8;
  border-radius: 8px;
  transition: color 0.15s, background 0.15s;
  margin: 0 0.25rem;
  font-size: 1.15rem;
}

.cart-btn:hover {
  color: white;
  background: rgba(255, 255, 255, 0.08);
}

.cart-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  background: #ef4444;
  color: white;
  font-size: 0.6rem;
  font-weight: 700;
  border-radius: 9999px;
  min-width: 1.1rem;
  height: 1.1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 3px;
  line-height: 1;
}

.locale-select {
  width: 120px;
  margin-left: 0.5rem;
}

.locale-select :deep(.p-select) {
  --p-select-background: rgba(255, 255, 255, 0.07);
  --p-select-border-color: rgba(255, 255, 255, 0.12);
  --p-select-hover-border-color: rgba(255, 255, 255, 0.25);
  --p-select-color: #94a3b8;
  --p-select-dropdown-color: #64748b;
  font-size: 0.85rem;
}

/* ── 漢堡 ── */
.hamburger {
  display: none;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 1.2rem;
  border-radius: 8px;
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
  margin-left: 0.25rem;
}

.hamburger:hover {
  color: white;
  background: rgba(255, 255, 255, 0.08);
}

/* ── Drawer 內容 ── */
.drawer-header {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.drawer-logo {
  height: 1.6rem;
  width: auto;
  filter: brightness(0) invert(1);
  opacity: 0.8;
}

.drawer-brand {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.1rem;
  font-weight: 700;
  letter-spacing: 0.18em;
  color: #f1f5f9;
}

.drawer-nav {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 1rem 0;
}

.drawer-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.85rem 1rem;
  border-radius: 8px;
  color: #94a3b8;
  font-size: 0.95rem;
  font-weight: 500;
  text-decoration: none;
  transition: color 0.15s, background 0.15s;
}

.drawer-link:hover {
  color: #f1f5f9;
  background: rgba(255, 255, 255, 0.07);
}

.drawer-link.active {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  font-weight: 600;
}

.drawer-badge {
  margin-left: auto;
  background: #ef4444;
  color: white;
  font-size: 0.7rem;
  font-weight: 700;
  border-radius: 9999px;
  min-width: 1.3rem;
  height: 1.3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
}

.drawer-footer {
  position: absolute;
  bottom: 1.5rem;
  left: 1.25rem;
  right: 1.25rem;
}

.drawer-label {
  font-size: 0.78rem;
  color: #475569;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
}

/* ── RWD ── */
@media (max-width: 640px) {
  .ticker-live-text {
    display: none;
  }

  .nav-hide-xs {
    display: none;
  }

  .hamburger {
    display: flex;
  }

  .brand-text {
    font-size: 1.05rem;
  }
}
</style>

<!-- Drawer 深色主題覆寫（non-scoped） -->
<style>
.nav-drawer .p-drawer-header {
  background: #0f172a;
  border-bottom: 1px solid #1e293b;
  padding: 1.1rem 1.25rem;
}

.nav-drawer .p-drawer-content {
  background: #0f172a;
  padding: 0 1.25rem;
  flex: 1;
  overflow-y: auto;
}

.nav-drawer .p-drawer-close-button {
  color: #64748b !important;
}

.nav-drawer .p-drawer-close-button:hover {
  color: #f1f5f9 !important;
  background: rgba(255, 255, 255, 0.08) !important;
}

/* locale select panel */
.p-select-overlay {
  --p-select-overlay-background: #1e293b;
  --p-select-overlay-border-color: #334155;
  --p-select-overlay-color: #94a3b8;
  --p-select-option-color: #94a3b8;
  --p-select-option-focus-background: rgba(255, 255, 255, 0.08);
  --p-select-option-focus-color: #f1f5f9;
  --p-select-option-selected-background: rgba(255, 255, 255, 0.1);
  --p-select-option-selected-color: #f1f5f9;
  --p-select-option-selected-focus-background: rgba(255, 255, 255, 0.12);
  font-size: 0.85rem;
}
</style>

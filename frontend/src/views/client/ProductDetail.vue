<template>
  <div class="page">
    <Navbar />

    <main class="container">

      <!-- 返回 -->
      <Button icon="pi pi-arrow-left" text class="mb-2" :label="t('nav.back')" @click="router.back()" />

      <div v-if="loading" class="center-msg">
        <ProgressSpinner />
      </div>

      <div v-else-if="!product" class="center-msg error-msg">
        {{ t('errors.PRODUCT_NOT_FOUND') }}
      </div>

      <div v-else class="detail-layout">
        <!-- 圖片 -->
        <div class="image-wrap">
          <img
            :src="resolveImage(allImages[activeImg])"
            :alt="product.name"
            class="product-image"
            @error="(e: any) => { e.target.src = '/placeholder.svg' }"
          />
          <span class="type-badge" :class="product.product_type === 'pre_order' ? 'pre-order' : 'in-stock'">
            {{ product.product_type === 'pre_order' ? t('product.preOrder') : t('product.inStock') }}
          </span>
          <div v-if="allImages.length > 1" class="thumb-strip">
            <img
              v-for="(img, i) in allImages"
              :key="i"
              :src="resolveImage(img)"
              class="thumb-dot"
              :class="{ active: activeImg === i }"
              @click="activeImg = i"
            />
          </div>
        </div>

        <!-- 資訊 -->
        <div class="info-panel">
          <h1 class="product-name">{{ product.name }}</h1>
          <p class="product-desc">{{ product.description }}</p>

          <div class="meta-list">
            <div class="meta-row">
              <span class="meta-label">{{ t('product.weight') }}</span>
              <span>{{ product.weight }} oz t</span>
            </div>
            <div v-if="product.product_type === 'in_stock'" class="meta-row">
              <span class="meta-label">{{ t('product.stockQty') }}</span>
              <span :class="product.stock_qty === 0 ? 'stock-zero' : ''">
                {{ product.stock_qty > 0 ? product.stock_qty : t('product.insufficientStock') }}
              </span>
            </div>
            <div v-if="product.product_type === 'pre_order' && product.expected_arrival" class="meta-row">
              <span class="meta-label">{{ t('product.expectedArrival') }}</span>
              <span>{{ product.expected_arrival }}</span>
            </div>
          </div>

          <div class="price-block" :class="{ flash: flashing }">
            <span class="price-label">{{ t('product.price') }}</span>
            <span class="price-value">NT$ {{ product.price.toLocaleString() }}</span>
          </div>

          <div class="countdown-bar">
            <span class="countdown-label">
              {{ t('product.nextUpdate') }} {{ priceStore.countdown }}{{ t('product.seconds') }}
            </span>
            <ProgressBar :value="countdownPercent" :show-value="false" class="countdown-progress" />
          </div>

          <Button
            :label="t('product.addToCart')"
            icon="pi pi-shopping-cart"
            size="large"
            class="w-full mt-4"
            :disabled="product.product_type === 'in_stock' && product.stock_qty === 0"
            @click="addToCart" />
        </div>
      </div>
    </main>

    <Footer />
    <button v-show="showScrollTop" class="scroll-top-btn" @click="scrollToTop" aria-label="回到頂部">
      <i class="pi pi-chevron-up" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useScrollTop } from '@/composables/useScrollTop'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import ProgressBar from 'primevue/progressbar'
import ProgressSpinner from 'primevue/progressspinner'
import Navbar from '@/components/shared/Navbar.vue'
import Footer from '@/components/shared/Footer.vue'
import { usePriceStore } from '@/stores/price'
import { useCartStore } from '@/stores/cart'
import { resolveImage } from '@/api/image'
import api from '@/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToast()
const priceStore = usePriceStore()
const cartStore = useCartStore()

const { showScrollTop, scrollToTop } = useScrollTop()

const product = ref<any>(null)
const loading = ref(true)
const flashing = ref(false)
const activeImg = ref(0)

const allImages = computed(() => {
  if (!product.value) return []
  return [product.value.image_url, product.value.image_url_2, product.value.image_url_3].filter(Boolean)
})

const countdownPercent = computed(() => (priceStore.countdown / 60) * 100)

watch(() => product.value?.price, () => {
  flashing.value = true
  setTimeout(() => (flashing.value = false), 800)
})

watch(() => product.value?.id, () => { activeImg.value = 0 })

async function fetchProduct() {
  try {
    const res = await api.get(`/api/v1/products/${route.params.id}`)
    product.value = res.data
  } catch {
    product.value = null
  } finally {
    loading.value = false
  }
}

// 每 60 秒更新報價
let timer: ReturnType<typeof setInterval>

onMounted(() => {
  priceStore.startAutoRefresh()
  fetchProduct()
  timer = setInterval(fetchProduct, 60_000)
})

onUnmounted(() => {
  priceStore.stopAutoRefresh()
  clearInterval(timer)
})

function addToCart() {
  if (!product.value) return
  cartStore.addItem({
    product_id: product.value.id,
    product_name: product.value.name,
    quantity: 1,
    unit_price: product.value.price,
    image_url: product.value.image_url,
    stock_qty: product.value.product_type === 'in_stock' ? product.value.stock_qty : undefined,
  })
  toast.add({ severity: 'success', summary: product.value.name, detail: t('product.addToCart'), life: 2000 })
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
  max-width: 1000px;
  margin: 0 auto;
  padding: 2.5rem 1.5rem;
  flex: 1;
}

.center-msg {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4rem;
  color: #999;
}

.error-msg {
  color: #e74c3c;
}

.detail-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  background: white;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.09);
  border: 1px solid #e2e8f0;
}

.image-wrap {
  position: relative;
}

.product-image {
  width: 100%;
  aspect-ratio: 4 / 5;
  object-fit: cover;
  display: block;
}

.thumb-strip {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 6px;
  background: rgba(0,0,0,0.35);
  border-radius: 8px;
  padding: 5px 8px;
}

.thumb-dot {
  width: 44px;
  height: 44px;
  object-fit: cover;
  border-radius: 5px;
  cursor: pointer;
  opacity: 0.6;
  border: 2px solid transparent;
  transition: opacity 0.15s, border-color 0.15s;
}

.thumb-dot:hover { opacity: 0.85; }
.thumb-dot.active { opacity: 1; border-color: white; }

.type-badge {
  position: absolute;
  top: 14px;
  left: 14px;
  font-size: 0.95rem;
  font-weight: 700;
  padding: 5px 14px;
  border-radius: 20px;
  pointer-events: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}
.type-badge.in-stock { background: #27ae60; color: white; }
.type-badge.pre-order { background: #2980b9; color: white; }

.stock-zero { color: #e74c3c; font-weight: 700; }

.info-panel {
  padding: 2rem;
  display: flex;
  flex-direction: column;
}

.product-name {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.55rem;
  font-weight: 700;
  letter-spacing: 0.02em;
  margin: 0 0 0.75rem;
  color: #0f172a;
}

.product-desc {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 1.05rem;
  color: #64748b;
  line-height: 1.85;
  margin-bottom: 1.5rem;
}

.meta-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
}

.meta-label {
  color: #999;
}

.price-block {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
  padding: 1rem;
  background: #fdf3f3;
  border-radius: 10px;
  margin-bottom: 1rem;
}

.price-label {
  font-size: 0.9rem;
  color: #999;
}

.price-value {
  font-size: 2rem;
  font-weight: 800;
  color: #0f172a;
}

.price-block.flash {
  animation: flash-price 0.8s ease;
}

@keyframes flash-price {
  30% {
    background: #fff8e1;
  }

  100% {
    background: #fdf3f3;
  }
}

.countdown-bar {
  margin-bottom: 0.5rem;
}

.countdown-label {
  font-size: 0.75rem;
  color: #999;
}

.countdown-progress {
  height: 4px;
  margin-top: 4px;
}

@media (max-width: 640px) {
  .detail-layout {
    grid-template-columns: 1fr;
  }

  .product-image {
    min-height: 240px;
  }
}
</style>

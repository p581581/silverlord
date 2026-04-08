<template>
  <div class="page">
    <Navbar />

    <!-- Hero Banner -->
    <section class="hero">
      <div class="hero-overlay" />
      <div class="hero-content">
        <p class="hero-sub">您的專屬實體白銀金庫</p>
        <h1 class="hero-title">SILVER LORD</h1>
        <div class="hero-divider" />
        <p class="hero-desc">
          我們提供國際認證、最高純度的白銀實體投資產品。<br />透明報價、安全交易，是您避險保值與財富傳承的最佳選擇。
        </p>
        <button class="hero-cta" @click="scrollToProducts">選購白銀</button>
      </div>
    </section>

    <!-- About Section -->
    <!-- <section class="about-section">
      <div class="section-inner">
        <div class="about-inner">
          <p class="about-label">關於我們</p>
          <blockquote class="about-quote">
            源自全球頂尖鑄幣廠，純粹價值，如實呈現。
          </blockquote>
          <div class="about-divider" />
          <p class="about-body">
            我們精選來自國際知名鑄幣廠與美國權威通路的銀條、銀幣，承諾 99.9% 以上的絕對純度。價格同步國際白銀走勢，讓您在下單瞬間即刻鎖定行情，並由我們以最高效率配送，確保這份實體價值，安全、完整地交付予您。
          </p>
          <div class="about-stats">
            <div class="about-stat">
              <span class="stat-value">99.9%</span>
              <span class="stat-label">純度保證</span>
            </div>
            <div class="stat-sep" />
            <div class="about-stat">
              <span class="stat-value">即時</span>
              <span class="stat-label">鎖定行情</span>
            </div>
            <div class="stat-sep" />
            <div class="about-stat">
              <span class="stat-value">迅速</span>
              <span class="stat-label">效率配送</span>
            </div>
          </div>
        </div>
      </div>
    </section> -->

    <!-- Products Section -->
    <section id="products" class="products-section">
      <div class="section-inner">
        <div class="products-header">
          <div>
            <h2 class="products-title">熱銷白銀商品</h2>
            <p class="products-sub">實體白銀，掌握在手的真實財富</p>
          </div>
          <div v-if="priceStore.snapshot" class="price-badge">
            <span class="price-badge-label">下次更新</span>
            <span class="price-badge-countdown">{{ priceStore.countdown }}s</span>
          </div>
        </div>

        <div v-if="loading" class="center-msg">
          <ProgressSpinner />
        </div>
        <div v-else-if="error" class="center-msg error-msg">{{ error }}</div>
        <div v-else-if="products.length === 0" class="center-msg">目前沒有上架商品</div>
        <div v-else class="product-grid">
          <ProductCard v-for="product in products" :key="product.id" :product="product" />
        </div>
      </div>
    </section>

    <Footer />

    <!-- Scroll to top button -->
    <button v-show="showScrollTop" class="scroll-top-btn" @click="scrollToTop" aria-label="回到頂部">
      <i class="pi pi-chevron-up" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import ProgressSpinner from 'primevue/progressspinner'
import Navbar from '@/components/shared/Navbar.vue'
import Footer from '@/components/shared/Footer.vue'
import ProductCard from '@/components/client/ProductCard.vue'
import { usePriceStore } from '@/stores/price'
import { useCartStore } from '@/stores/cart'
import api from '@/api'

const { t } = useI18n()
const priceStore = usePriceStore()
const cartStore = useCartStore()

const products = ref<any[]>([])
const loading = ref(false)
const error = ref('')

async function fetchProducts() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('/api/v1/products')
    products.value = res.data
    const priceMap: Record<number, number> = {}
    res.data.forEach((p: any) => (priceMap[p.id] = p.price))
    cartStore.updatePrices(priceMap)
  } catch (err: any) {
    error.value = err.userMessage || t('errors.INTERNAL_ERROR')
  } finally {
    loading.value = false
  }
}

function scrollToProducts() {
  document.getElementById('products')?.scrollIntoView({ behavior: 'smooth' })
}

let productTimer: ReturnType<typeof setInterval>

// Scroll to top
const showScrollTop = ref(false)

function onScroll() {
  const scrollY = window.scrollY ?? document.documentElement.scrollTop
  showScrollTop.value = scrollY > 200
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.scrollTo(0, 0)
  window.addEventListener('scroll', onScroll, { passive: true })
  priceStore.startAutoRefresh()
  fetchProducts()
  productTimer = setInterval(fetchProducts, 60_000)
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  priceStore.stopAutoRefresh()
  clearInterval(productTimer)
})
</script>

<style scoped>
.scroll-top-btn {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  border: 1px solid #aaaaaa;
  background: rgba(30, 30, 30, 0.8);
  color: #cccccc;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 10px rgba(0,0,0,0.25);
  transition: background 0.2s, transform 0.2s;
  z-index: 999;
}
.scroll-top-btn:hover {
  background: rgba(180, 180, 180, 0.15);
  transform: translateY(-2px);
}

.page {
  --section-pad: 5rem;
  min-height: 100vh;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
}

/* ── Hero ── */
.hero {
  position: relative;
  background-image: url('https://images.unsplash.com/photo-1610375461246-83df859d849d?auto=format&fit=crop&q=80&w=1920&sat=-100');
  background-size: cover;
  background-position: center;
  min-height: 520px;
  display: flex;
  align-items: center;
  overflow: hidden;
}

.hero-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to right, rgba(10, 14, 28, 0.92) 38%, rgba(10, 14, 28, 0.55) 70%, rgba(10, 14, 28, 0.35));
}

.hero-content {
  position: relative;
  z-index: 1;
  max-width: 1280px;
  margin: 0 auto;
  padding: 5.5rem 1.5rem;
  width: 100%;
}

.hero-sub {
  font-family: 'EB Garamond', Georgia, serif;
  font-style: italic;
  color: #c8bfa8;
  font-size: 1.5rem;
  letter-spacing: 0.04em;
  margin-bottom: 1rem;
}

.hero-title {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(2.8rem, 6.5vw, 4.5rem);
  font-weight: 700;
  letter-spacing: 0.15em;
  color: #f5f0e8;
  margin-bottom: 1.4rem;
  margin-top: 0;
}

.hero-divider {
  width: 3rem;
  height: 1px;
  background: linear-gradient(to right, #c8bfa8, transparent);
  margin-bottom: 1.4rem;
}

.hero-desc {
  font-family: 'EB Garamond', Georgia, serif;
  color: #a89f8c;
  font-size: 1.15rem;
  line-height: 1.9;
  max-width: 520px;
  margin-bottom: 2.25rem;
}

@media (max-width: 640px) {
  .page {
    --section-pad: 3rem;
  }

  .hero-desc>br {
    display: none;
  }
}

@media (max-width: 480px) {
  .page {
    --section-pad: 1.5rem;
  }

  .hero {
    /* ticker ≈ 32px + navbar 72px = 104px */
    min-height: 0;
    height: calc(100dvh - 104px);
  }

  .hero-sub {
    font-size: 1.15rem;
    margin-bottom: 0.5rem;
  }

  .hero-title {
    font-size: 2rem;
    white-space: nowrap;
    margin-bottom: 0.75rem;
  }

  .hero-divider {
    margin-bottom: 0.75rem;
  }

  .hero-desc {
    font-size: 1.05rem;
    line-height: 1.7;
    margin-bottom: 1.25rem;
  }

  .hero-cta {
    font-size: 0.85rem;
    padding: 0.6rem 1.75rem;
  }
}

.hero-cta {
  font-family: 'Playfair Display', Georgia, serif;
  background: transparent;
  color: #f5f0e8;
  border: 1px solid rgba(200, 191, 168, 0.6);
  padding: 0.8rem 2.5rem;
  border-radius: 2px;
  font-size: 0.95rem;
  font-weight: 600;
  letter-spacing: 0.15em;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s, color 0.2s;
}

.hero-cta:hover {
  background: rgba(200, 191, 168, 0.15);
  border-color: #c8bfa8;
  color: #fff;
}

/* ── About ── */
.about-section {
  background: var(--color-surface, #fff);
  padding: var(--section-pad) 0;
  border-bottom: 1px solid #e8e3dc;
}

.section-inner {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.5rem;
}

.about-inner {
  max-width: 680px;
  margin: 0 auto;
  text-align: center;
}

.about-label {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 0.8rem;
  letter-spacing: 0.2em;
  text-transform: uppercase;
  color: #94a3b8;
  margin-bottom: 1.5rem;
}

.about-quote {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: clamp(1.3rem, 2.5vw, 1.75rem);
  font-style: italic;
  font-weight: 500;
  color: #0f172a;
  line-height: 1.55;
  margin: 0 0 1.5rem;
}

.about-divider {
  width: 2.5rem;
  height: 1px;
  background: linear-gradient(to right, transparent, #c8bfa8, transparent);
  margin: 0 auto 1.5rem;
}

.about-body {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 1.1rem;
  color: #475569;
  line-height: 2;
  margin-bottom: 2.5rem;
}

.about-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 2rem;
}

.about-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
}

.stat-value {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.4rem;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: 0.04em;
}

.stat-label {
  font-family: 'EB Garamond', Georgia, serif;
  font-size: 0.82rem;
  color: #94a3b8;
  letter-spacing: 0.08em;
}

.stat-sep {
  width: 1px;
  height: 2rem;
  background: #e2e8f0;
}

/* ── Products ── */
.products-section {
  padding: var(--section-pad) 0;
  flex: 1;
}

.products-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.products-title {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.9rem;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 0.35rem;
  margin-top: .5rem;
  letter-spacing: 0.02em;
}

.products-sub {
  font-family: 'EB Garamond', Georgia, serif;
  font-style: italic;
  color: #64748b;
  font-size: 1.05rem;
}

.price-badge {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: #0f172a;
  color: white;
  padding: 0.4rem 0.9rem;
  border-radius: 9999px;
  font-size: 0.8rem;
}

.price-badge-label {
  color: #94a3b8;
}

.price-badge-countdown {
  font-weight: 700;
  font-family: monospace;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 1.75rem;
}

.center-msg {
  text-align: center;
  padding: 4rem;
  color: #94a3b8;
}

.error-msg {
  color: #ef4444;
}
</style>

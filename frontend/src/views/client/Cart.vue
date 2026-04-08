<template>
  <div class="page">
    <Navbar />

    <main class="container">
      <h1 class="page-title">
        <i class="pi pi-shopping-cart" />
        {{ t('cart.title') }}
      </h1>

      <!-- 逾時警告 -->
      <Message v-if="cartStore.isCheckoutExpired" severity="error" class="mb-4">
        {{ t('cart.sessionExpired') }}
        <Button label="Refresh" size="small" text @click="refreshCart" />
      </Message>
      <Message
        v-else-if="cartStore.checkoutCountdown > 0 && cartStore.checkoutCountdown <= 60"
        severity="warn"
        class="mb-4">
        {{ t('cart.sessionWarning', { seconds: cartStore.checkoutCountdown }) }}
      </Message>

      <!-- 空購物車 -->
      <div v-if="cartStore.items.length === 0" class="empty-cart">
        <i class="pi pi-shopping-cart empty-icon" />
        <h3>{{ t('cart.empty') }}</h3>
        <p>快去商品區挑選具備投資價值的白銀吧！</p>
        <Button :label="t('nav.products')" icon="pi pi-arrow-left" @click="router.push('/')" />
      </div>

      <!-- 購物車內容 -->
      <div v-else class="cart-layout">

        <!-- 商品列表 -->
        <div class="cart-items">
          <div class="items-card">
            <div
              v-for="item in cartStore.items"
              :key="item.product_id"
              class="cart-item">
              <RouterLink :to="`/products/${item.product_id}`">
                <img :src="resolveImage(item.image_url)" class="item-thumb" :alt="item.product_name"
                  @error="(e) => ((e.target as HTMLImageElement).src = '/placeholder.svg')" />
              </RouterLink>

              <div class="item-info">
                <div class="item-top">
                  <div>
                    <RouterLink :to="`/products/${item.product_id}`" class="item-name">
                      {{ item.product_name }}
                    </RouterLink>
                    <div class="item-unit">NT$ {{ item.unit_price.toLocaleString() }} / 件</div>
                  </div>
                  <button class="remove-btn" @click="cartStore.removeItem(item.product_id)">
                    <i class="pi pi-trash" />
                  </button>
                </div>

                <div class="item-bottom">
                  <div class="qty-control">
                    <button class="qty-btn"
                      @click="cartStore.updateQuantity(item.product_id, Math.max(1, item.quantity - 1))">
                      <i class="pi pi-minus" />
                    </button>
                    <span class="qty-value">{{ item.quantity }}</span>
                    <button
                      class="qty-btn"
                      :disabled="item.stock_qty != null && item.quantity >= item.stock_qty"
                      @click="cartStore.updateQuantity(item.product_id, item.quantity + 1)">
                      <i class="pi pi-plus" />
                    </button>
                  </div>
                  <span class="item-subtotal">
                    NT$ {{ (item.unit_price * item.quantity).toLocaleString() }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 訂單摘要 -->
        <div class="cart-summary">
          <div class="summary-card">
            <h3 class="summary-title">訂單摘要</h3>

            <div class="summary-rows">
              <div class="summary-row">
                <span>商品小計</span>
                <span>NT$ {{ cartStore.total.toLocaleString() }}</span>
              </div>
              <div class="summary-row muted">
                <span>運費</span>
                <span>結帳時確認</span>
              </div>
            </div>

            <div class="summary-total">
              <span>合計金額</span>
              <span class="total-value">NT$ {{ cartStore.total.toLocaleString() }}</span>
            </div>

            <Button
              :label="t('cart.checkout')"
              icon="pi pi-lock"
              class="checkout-btn w-full"
              :disabled="cartStore.isCheckoutExpired"
              @click="goCheckout" />

            <div class="secure-note">
              <i class="pi pi-shield" />
              安全加密結帳環境
            </div>
          </div>
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
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, RouterLink } from 'vue-router'
import Button from 'primevue/button'
import Message from 'primevue/message'
import Navbar from '@/components/shared/Navbar.vue'
import Footer from '@/components/shared/Footer.vue'
import { useCartStore } from '@/stores/cart'
import { resolveImage } from '@/api/image'

const { t } = useI18n()
const router = useRouter()
const cartStore = useCartStore()

onMounted(() => {
  window.scrollTo(0, 0)
  window.addEventListener('scroll', onScroll, { passive: true })
  cartStore.startCheckoutTimer()
})
onUnmounted(() => {
  cartStore.stopCheckoutTimer()
  window.removeEventListener('scroll', onScroll)
})

const showScrollTop = ref(false)
function onScroll() {
  showScrollTop.value = (window.scrollY ?? document.documentElement.scrollTop) > 200
}
function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function goCheckout() {
  if (cartStore.isCheckoutExpired) return
  router.push('/checkout')
}
function refreshCart() {
  router.push('/')
}
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
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
  transition: background 0.2s, transform 0.2s;
  z-index: 999;
}

.scroll-top-btn:hover {
  background: rgba(180, 180, 180, 0.15);
  transform: translateY(-2px);
}

.page {
  min-height: 100vh;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 3rem 1.5rem;
  flex: 1;
}

.page-title {
  font-size: 1.75rem;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 2rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

/* Empty */
.empty-cart {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 5rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  border: 1px solid #e2e8f0;
  text-align: center;
}

.empty-icon {
  font-size: 3.5rem;
  color: #cbd5e1;
}

.empty-cart h3 {
  font-size: 1.2rem;
  font-weight: 700;
  color: #475569;
}

.empty-cart p {
  color: #94a3b8;
  margin-bottom: 0.5rem;
}

/* Layout */
.cart-layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 1.5rem;
  align-items: start;
}

@media (max-width: 900px) {
  .cart-layout {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .cart-item {
    gap: 0.75rem;
    padding: 1rem;
  }

  .item-name {
    font-size: 0.88rem;
  }

  .item-bottom {
    flex-wrap: wrap;
    gap: .5rem;
  }

  .qty-control {
    margin-left: auto;
  }

  .item-subtotal {
    font-size: 1rem;
    margin-left: auto;
  }
}

/* Items */
.items-card {
  background: white;
  border-radius: 1rem;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.cart-item {
  display: flex;
  align-items: stretch;
  gap: .8rem;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #f1f5f9;
  transition: background 0.15s;
}

.cart-item:last-child {
  border-bottom: none;
}

.cart-item:hover {
  background: #fafafa;
}

.item-thumb {
  width: 100px;
  height: 100%;
  min-height: 100px;
  object-fit: cover;
  border-radius: 8px;
  flex-shrink: 0;
  background: #e2e8f0;
}

.item-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  justify-content: space-between;
}

.item-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: .5rem;
  margin-top: 0.2rem;
}

.item-name {
  font-size: 0.95rem;
  font-weight: 700;
  color: #0f172a;
  text-decoration: none;
  transition: color 0.15s;
  display: block;
  margin-bottom: 0.75rem;

}

.item-name:hover {
  color: #1d4ed8;
}

.item-unit {
  font-size: 0.82rem;
  color: #94a3b8;
  margin-top: .25rem
}

.remove-btn {
  background: none;
  border: none;
  padding: 0.3rem;
  color: #cbd5e1;
  cursor: pointer;
  border-radius: 6px;
  transition: color 0.15s, background 0.15s;
  line-height: 1;
}

.remove-btn:hover {
  color: #ef4444;
  background: #fef2f2;
}

.item-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.qty-control {
  display: flex;
  align-items: center;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
  background: #f8fafc;
}

.qty-btn {
  background: none;
  border: none;
  width: 2rem;
  height: 2rem;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  transition: background 0.15s;
}

.qty-btn:hover:not(:disabled) {
  background: #e2e8f0;
}

.qty-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.qty-value {
  min-width: 2.5rem;
  text-align: center;
  font-weight: 700;
  font-size: 0.9rem;
  color: #0f172a;
}

.item-subtotal {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

/* Summary */
.summary-card {
  background: #0f172a;
  color: white;
  border-radius: 1rem;
  padding: 1.75rem;
  position: sticky;
  top: calc(36px + 72px + 1rem);
  /* ticker + navbar + gap */
}

.summary-title {
  font-size: 1.1rem;
  font-weight: 700;
  padding-bottom: 1rem;
  border-bottom: 1px solid #1e293b;
  margin-bottom: 1.25rem;
}

.summary-rows {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.25rem;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
  color: #94a3b8;
}

.summary-row.muted {
  color: #475569;
}

.summary-total {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  border-top: 1px solid #1e293b;
  padding-top: 1rem;
  margin-bottom: 1.5rem;
  font-size: 1rem;
  font-weight: 600;
}

.total-value {
  font-size: 1.75rem;
  font-weight: 800;
  color: white;
}

.checkout-btn {
  width: 100%;
  justify-content: center;
  border-radius: 10px !important;
  padding: 0.85rem !important;
  font-size: 1rem !important;
  font-weight: 700 !important;
}

.secure-note {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  font-size: 0.78rem;
  color: #475569;
  margin-top: 1rem;
}
</style>

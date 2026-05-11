<template>
  <div class="product-card" :class="{ updating: priceStore.loading }">

    <!-- 圖片區 4:5 -->
    <RouterLink :to="`/products/${product.id}`" class="image-wrap">
      <img
        :src="resolveImage(product.image_url)"
        :alt="product.name"
        class="product-image"
        @error="(e: any) => { e.target.src = '/placeholder.svg' }"
      />
      <div class="image-overlay" />
      <span class="type-badge" :class="product.product_type === 'pre_order' ? 'pre-order' : (product.stock_qty === 0 ? 'restocking' : 'in-stock')">
        {{ product.product_type === 'pre_order' ? t('product.preOrder') : (product.stock_qty === 0 ? t('product.insufficientStock') : t('product.inStock')) }}
      </span>
    </RouterLink>

    <!-- 資訊區 -->
    <div class="product-body">
      <RouterLink :to="`/products/${product.id}`" class="product-name-link">
        <h3 class="product-name">{{ product.name }}</h3>
      </RouterLink>
      <p class="product-desc">{{ product.description }}</p>

      <div class="product-meta">
        <span>{{ t('product.weight') }}: {{ product.weight }} oz t</span>
        <span v-if="product.product_type === 'in_stock' && product.stock_qty === 0" class="out-of-stock">
          {{ t('product.insufficientStock') }}
        </span>
        <span v-else-if="product.product_type === 'pre_order' && product.expected_arrival" class="arrival">
          {{ t('product.expectedArrival') }}: {{ product.expected_arrival }}
        </span>
      </div>

      <div class="price-row" :class="{ flash: flashing }">
        <span class="price-label">NT$</span>
        <span class="price-value">{{ product.price.toLocaleString() }}</span>
      </div>

      <!-- 數量 + 加入購物車 -->
      <div class="cart-row" :class="{ disabled: isSoldOut || isCartFull }">
        <QtyInput
          v-model="qty"
          :max="product.product_type === 'in_stock' ? remaining : undefined"
          :disabled="isSoldOut || isCartFull"
        />
        <button class="add-btn" :disabled="isSoldOut || isCartFull" @click.prevent="addToCart">
          <i class="pi pi-shopping-cart" />
          {{ isCartFull ? t('product.cartFull') : t('product.addToCart') }}
        </button>
      </div>

      <div class="countdown-bar">
        <ProgressBar :value="countdownPercent" :show-value="false" class="countdown-progress" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { RouterLink } from 'vue-router'
import ProgressBar from 'primevue/progressbar'
import { usePriceStore } from '@/stores/price'
import { useCartStore } from '@/stores/cart'
import { resolveImage } from '@/api/image'
import QtyInput from '@/components/shared/QtyInput.vue'

interface Product {
  id: number
  name: string
  description: string
  image_url: string
  image_url_2?: string
  image_url_3?: string
  weight: number
  price: number
  product_type: string
  stock_qty: number
  expected_arrival: string
}

const props = defineProps<{ product: Product }>()

const { t } = useI18n()
const toast = useToast()
const priceStore = usePriceStore()
const cartStore = useCartStore()

const qty = ref(1)
const flashing = ref(false)

const countdownPercent = computed(() => (priceStore.countdown / 60) * 100)
const isSoldOut = computed(() => props.product.product_type === 'in_stock' && props.product.stock_qty === 0)

const cartQty = computed(() => cartStore.items.find(i => i.product_id === props.product.id)?.quantity ?? 0)
const remaining = computed(() =>
  props.product.product_type === 'in_stock' ? props.product.stock_qty - cartQty.value : Infinity
)
const isCartFull = computed(() => props.product.product_type === 'in_stock' && remaining.value <= 0)

watch(
  () => props.product.price,
  () => {
    flashing.value = true
    setTimeout(() => (flashing.value = false), 800)
  },
)

// 報價更新時，若數量超過最新庫存則自動修正
watch(
  () => props.product.stock_qty,
  (newQty) => {
    if (props.product.product_type === 'in_stock' && qty.value > newQty) {
      qty.value = Math.max(1, newQty)
    }
  },
)

function addToCart() {
  if (isSoldOut.value || isCartFull.value) return
  const addQty = isFinite(remaining.value) ? Math.min(qty.value, remaining.value) : qty.value
  cartStore.addItem({
    product_id: props.product.id,
    product_name: props.product.name,
    quantity: addQty,
    unit_price: props.product.price,
    image_url: props.product.image_url,
    stock_qty: props.product.product_type === 'in_stock' ? props.product.stock_qty : undefined,
  })
  toast.add({ severity: 'success', summary: props.product.name, detail: `× ${addQty} ${t('product.addToCart')}`, life: 2000 })
  qty.value = 1
}
</script>

<style scoped>
.product-card {
  background: white;
  border-radius: 2px;
  overflow: hidden;
  box-shadow: 0 1px 6px rgba(0,0,0,0.06);
  border: 1px solid #e8e3dc;
  transition: box-shadow 0.2s, transform 0.2s;
  display: flex;
  flex-direction: column;
}
.product-card:hover {
  box-shadow: 0 10px 30px rgba(0,0,0,0.11);
  transform: translateY(-4px);
}

/* ── Image 4:5 ── */
.image-wrap {
  position: relative;
  display: block;
  overflow: hidden;
  background: #e2e8f0;
}
.product-image {
  width: 100%;
  aspect-ratio: 4 / 5;
  object-fit: cover;
  display: block;
  transition: transform 0.5s ease;
}
.product-card:hover .product-image { transform: scale(1.04); }
.image-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0,0,0,0.12) 0%, transparent 45%);
  pointer-events: none;
}
.type-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  font-size: 0.78rem;
  font-weight: 700;
  padding: 4px 12px;
  border-radius: 9999px;
  pointer-events: none;
  letter-spacing: 0.03em;
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
}
.type-badge.in-stock   { background: #15803d; color: white; }
.type-badge.pre-order  { background: #1d4ed8; color: white; }
.type-badge.restocking { background: #f87171; color: white; }

/* ── Body ── */
.product-body {
  padding: 1.1rem 1.1rem 1rem;
  display: flex;
  flex-direction: column;
  flex: 1;
}

.product-name-link { text-decoration: none; color: inherit; }
.product-name {
  font-family: 'Playfair Display', Georgia, serif;
  font-size: 1.05rem;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 0.35rem;
  line-height: 1.4;
  transition: color 0.15s;
  letter-spacing: 0.01em;
}
.product-name-link:hover .product-name { color: #3d3529; }

.product-desc {
  font-family: 'EB Garamond', Georgia, serif;
  color: #64748b;
  font-size: 0.95rem;
  line-height: 1.65;
  margin-bottom: 0.65rem;
  min-height: 2.8em;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.product-meta {
  font-size: 0.78rem;
  color: #94a3b8;
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  margin-bottom: 0.65rem;
}
.out-of-stock { color: #ef4444; font-weight: 600; }
.arrival { color: #3b82f6; }

/* Price */
.price-row {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
  margin-bottom: 0.9rem;
}
.price-label { font-family: var(--font-ui, sans-serif); font-size: 0.82rem; color: #94a3b8; font-weight: 500; letter-spacing: 0.03em; }
.price-value { font-family: 'Playfair Display', Georgia, serif; font-size: 1.55rem; font-weight: 700; color: #0f172a; line-height: 1; }
.price-row.flash { animation: flash-price 0.8s ease; }
@keyframes flash-price {
  30% { background: #fef9c3; border-radius: 6px; }
  100% { background: transparent; }
}

/* ── Cart Row ── */
.cart-row {
  display: flex;
  gap: 0.6rem;
  align-items: center;
  margin-bottom: 0.75rem;
}
.cart-row.disabled { opacity: 0.45; pointer-events: none; }


.add-btn {
  flex: 1;
  background: #0f172a;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0 0.85rem;
  height: 2.2rem;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  transition: background 0.15s;
  white-space: nowrap;
}
.add-btn:hover:not(:disabled) { background: #1e293b; }
.add-btn:disabled { background: #94a3b8; cursor: not-allowed; }

/* Countdown */
.countdown-bar { margin-top: auto; }
.countdown-progress { height: 3px !important; }
</style>

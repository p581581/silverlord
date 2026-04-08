<template>
  <div class="page">
    <Navbar />

    <main class="container">
      <h1 class="page-title">{{ t('checkout.title') }}</h1>

      <!-- 逾時 -->
      <Message v-if="cartStore.isCheckoutExpired" severity="error">
        {{ t('cart.sessionExpired') }}
        <Button label="Back" size="small" text @click="router.push('/')" />
      </Message>

      <!-- 訂單成功 -->
      <div v-else-if="orderResult" class="success-box">
        <i class="pi pi-check-circle success-icon" />
        <h2>{{ t('checkout.success') }}</h2>
        <div class="order-no-box">
          {{ t('checkout.orderNo') }}: <code>{{ orderResult.order_no }}</code>
        </div>
        <div class="locked-price">
          {{ t('checkout.lockedPrice') }}:
          <strong>NT$ {{ orderResult.total_locked_price.toLocaleString() }}</strong>
        </div>
        <p class="hint">{{ t('checkout.hint') }}</p>
        <Button
          :label="t('payment.title')"
          icon="pi pi-send"
          @click="router.push('/payment')"
        />
      </div>

      <!-- 確認訂單 -->
      <div v-else>
        <!-- Email -->
        <div class="form-field">
          <label>Email <span class="required">*</span></label>
          <InputText v-model="email" type="email" placeholder="your@email.com" class="w-full" />
        </div>

        <!-- 姓名 / 電話 -->
        <div class="form-row">
          <div class="form-field">
            <label>{{ t('checkout.name') }} <span class="required">*</span></label>
            <InputText v-model="name" :placeholder="t('checkout.name')" class="w-full" />
          </div>
          <div class="form-field">
            <label>{{ t('checkout.phone') }} <span class="required">*</span></label>
            <InputText v-model="phone" placeholder="0912345678" class="w-full" />
          </div>
        </div>

        <!-- 出貨方式 -->
        <div class="form-field">
          <label>{{ t('checkout.shipMethod') }} <span class="required">*</span></label>
          <div class="ship-options">
            <div
              class="ship-option"
              :class="{ active: shipMethod === 'mail' }"
              @click="shipMethod = 'mail'"
            >
              <i class="pi pi-envelope" />
              <span>{{ t('checkout.shipMail') }}</span>
            </div>
            <div
              class="ship-option"
              :class="{ active: shipMethod === '711' }"
              @click="shipMethod = '711'"
            >
              <i class="pi pi-shop" />
              <span class="label-full">{{ t('checkout.ship711') }}</span>
              <span class="label-xs">7-11</span>
            </div>
            <div
              class="ship-option"
              :class="{ active: shipMethod === 'pickup' }"
              @click="shipMethod = 'pickup'"
            >
              <i class="pi pi-home" />
              <span>{{ t('checkout.shipPickup') }}</span>
            </div>
          </div>
        </div>

        <!-- 郵寄地址 -->
        <div v-if="shipMethod === 'mail'" class="form-field">
          <label>{{ t('checkout.shipAddress') }} <span class="required">*</span></label>
          <InputText
            v-model="shipAddress"
            :placeholder="t('checkout.shipAddressPlaceholder')"
            class="w-full"
          />
        </div>

        <!-- 7-11 店名 -->
        <div v-if="shipMethod === '711'" class="form-field">
          <label>{{ t('checkout.shipStore') }} <span class="required">*</span></label>
          <InputText
            v-model="shipStore"
            :placeholder="t('checkout.shipStorePlaceholder')"
            class="w-full"
          />
        </div>

        <DataTable :value="cartStore.items" class="mb-4">
          <Column :header="t('admin.name')">
            <template #body="{ data }">
              <RouterLink :to="`/products/${data.product_id}`" class="product-link">
                {{ data.product_name }}
              </RouterLink>
            </template>
          </Column>
          <Column :header="t('cart.quantity')">
            <template #body="{ data }">{{ data.quantity }}</template>
          </Column>
          <Column :header="t('cart.unitPrice')">
            <template #body="{ data }">NT$ {{ data.unit_price.toLocaleString() }}</template>
          </Column>
          <Column :header="t('cart.subtotal')">
            <template #body="{ data }">
              NT$ {{ (data.unit_price * data.quantity).toLocaleString() }}
            </template>
          </Column>
        </DataTable>

        <div class="checkout-footer">
          <div class="price-row">
            <span>{{ t('checkout.subtotal') }}</span>
            <span>NT$ {{ cartStore.total.toLocaleString() }}</span>
          </div>
          <div class="price-row">
            <span>{{ t('checkout.shipFee') }}</span>
            <span>{{ currentShipFee > 0 ? 'NT$ ' + currentShipFee.toLocaleString() : t('checkout.shipFreeText') }}</span>
          </div>
          <div class="price-row total-display">
            <span>{{ t('checkout.grandTotal') }}</span>
            <strong>NT$ {{ grandTotal.toLocaleString() }}</strong>
          </div>
          <div class="countdown-hint">
            {{ t('cart.sessionWarning', { seconds: cartStore.checkoutCountdown }) }}
          </div>
          <div class="actions">
            <Button
              :label="t('cart.title')"
              icon="pi pi-arrow-left"
              text
              @click="router.push('/cart')"
            />
            <Button
              :label="t('checkout.confirm')"
              icon="pi pi-lock"
              :loading="submitting"
              @click="placeOrder"
            />
          </div>
        </div>
      </div>
    </main>
    <button v-show="showScrollTop" class="scroll-top-btn" @click="scrollToTop" aria-label="回到頂部">
      <i class="pi pi-chevron-up" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useScrollTop } from '@/composables/useScrollTop'
import { useI18n } from 'vue-i18n'
import { useRouter, RouterLink } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Message from 'primevue/message'
import InputText from 'primevue/inputtext'
import Navbar from '@/components/shared/Navbar.vue'
import { useCartStore } from '@/stores/cart'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const toast = useToast()
const cartStore = useCartStore()

const { showScrollTop, scrollToTop } = useScrollTop()

const submitting = ref(false)
const email = ref('')
const name = ref('')
const phone = ref('')
const shipMethod = ref<'mail' | '711' | 'pickup' | ''>('')
const shipAddress = ref('')
const shipStore = ref('')
const shipFees = ref<{ mail: number; '711': number; pickup: number }>({ mail: 0, '711': 0, pickup: 0 })
const orderResult = ref<{ order_no: string; total_locked_price: number } | null>(null)

const currentShipFee = computed(() =>
  shipMethod.value === 'mail' ? shipFees.value.mail
  : shipMethod.value === '711' ? shipFees.value['711']
  : shipMethod.value === 'pickup' ? shipFees.value.pickup
  : 0
)
const grandTotal = computed(() => cartStore.total + currentShipFee.value)

onMounted(async () => {
  cartStore.startCheckoutTimer()
  try {
    const res = await api.get('/api/v1/ship-fees')
    shipFees.value = res.data
  } catch {}
})
onUnmounted(() => {
  cartStore.stopCheckoutTimer()
})

function isValidEmail(e: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e)
}

async function placeOrder() {
  if (cartStore.isCheckoutExpired || submitting.value) return

  if (!email.value || !isValidEmail(email.value)) {
    toast.add({ severity: 'warn', summary: 'Email', detail: t('checkout.emailInvalid'), life: 3000 })
    return
  }
  if (!name.value.trim() || !phone.value.trim()) {
    toast.add({ severity: 'warn', summary: t('checkout.name'), detail: t('checkout.shipRequired'), life: 3000 })
    return
  }

  if (!shipMethod.value) {
    toast.add({ severity: 'warn', summary: t('checkout.shipMethod'), detail: t('checkout.shipRequired'), life: 3000 })
    return
  }
  if (shipMethod.value === 'mail' && !shipAddress.value.trim()) {
    toast.add({ severity: 'warn', summary: t('checkout.shipAddress'), detail: t('checkout.shipRequired'), life: 3000 })
    return
  }
  if (shipMethod.value === '711' && !shipStore.value.trim()) {
    toast.add({ severity: 'warn', summary: t('checkout.shipStore'), detail: t('checkout.shipRequired'), life: 3000 })
    return
  }
  // pickup: no additional fields required

  submitting.value = true
  try {
    const items = cartStore.items.map((i) => ({
      product_id: i.product_id,
      quantity: i.quantity,
    }))
    const res = await api.post('/api/v1/checkout', {
      email: email.value,
      name: name.value,
      phone: phone.value,
      ship_method: shipMethod.value,
      ship_address: shipAddress.value,
      ship_store: shipStore.value,
      items,
    })
    orderResult.value = res.data
    cartStore.clearCart()

    const saved = JSON.parse(localStorage.getItem('my_orders') || '[]')
    saved.unshift(res.data.order_no)
    localStorage.setItem('my_orders', JSON.stringify(saved.slice(0, 20)))
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.userMessage, life: 4000 })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page { min-height: 100vh; background: #f8fafc; }
.container { max-width: 800px; margin: 0 auto; padding: 2rem 1rem; }
.page-title { font-size: 1.5rem; font-weight: 700; margin-bottom: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.4rem; margin-bottom: 1rem; }
.form-field label { font-weight: 600; }
.form-row { display: flex; gap: 1rem; }
.form-row .form-field { flex: 1; }
.required { color: #e74c3c; margin-left: 2px; }
.product-link { color: #1a1a2e; font-weight: 500; text-decoration: none; }
.product-link:hover { color: #0f3460; text-decoration: underline; }

.ship-options { display: flex; gap: 0.75rem; }
.ship-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.9rem;
  border: 2px solid #e0e0e0;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  color: #666;
  background: white;
  transition: all 0.15s;
}
.ship-option:hover { border-color: #aaa; color: #333; }
.ship-option.active { border-color: #1a1a2e; background: #1a1a2e; color: white; }
.ship-option .pi { font-size: 1.1rem; }

.checkout-footer {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.price-row { display: flex; justify-content: space-between; align-items: center; font-size: 0.95rem; color: #555; }
.total-display { font-size: 1.15rem; font-weight: 600; color: #1a1a2e; border-top: 1px solid #eee; padding-top: 0.5rem; margin-top: 0.25rem; }
.total-display strong { color: #c0392b; font-size: 1.4rem; }
.countdown-hint { font-size: 0.85rem; color: #e67e22; }
.actions { display: flex; justify-content: flex-end; gap: 1rem; }

.success-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  background: white;
  border-radius: 12px;
  padding: 3rem;
  text-align: center;
}
.success-icon { font-size: 4rem; color: #27ae60; }
.order-no-box {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
}
.order-no-box code { font-size: 1.1rem; font-weight: 700; color: #2c3e50; user-select: all; }
.locked-price { font-size: 1.1rem; }
.locked-price strong { color: #c0392b; font-size: 1.4rem; }
.hint { color: #666; font-size: 0.9rem; }

@media (max-width: 640px) {
  .form-row { flex-direction: column; gap: 0; }
}

@media (max-width: 480px) {
  .label-full { display: none; }
}

.label-xs { display: none; }
@media (max-width: 480px) {
  .label-xs { display: inline; }
}
</style>

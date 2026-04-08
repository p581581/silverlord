import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CartItem {
  product_id: number
  product_name: string
  quantity: number
  unit_price: number // 加入購物車當下的報價（顯示用，非鎖價）
  image_url?: string
  stock_qty?: number  // 庫存品上限（pre_order 為 undefined）
}

const CHECKOUT_TIMEOUT = 300 // 秒

export const useCartStore = defineStore('cart', () => {
  // 只持久化 product_id + quantity，價格由 computed 動態算
  const items = ref<CartItem[]>([])

  // 結帳計時器
  const checkoutCountdown = ref(0)
  const isCheckoutExpired = computed(() => checkoutCountdown.value <= 0 && items.value.length > 0)
  let checkoutTimer: ReturnType<typeof setInterval> | null = null

  const total = computed(() =>
    items.value.reduce((sum, item) => sum + item.unit_price * item.quantity, 0),
  )

  function addItem(item: CartItem) {
    const existing = items.value.find((i) => i.product_id === item.product_id)
    if (existing) {
      existing.quantity += item.quantity
      existing.unit_price = item.unit_price // 更新最新報價
    } else {
      items.value.push({ ...item })
    }
  }

  function removeItem(productId: number) {
    items.value = items.value.filter((i) => i.product_id !== productId)
  }

  function updateQuantity(productId: number, quantity: number) {
    const item = items.value.find((i) => i.product_id === productId)
    if (item) item.quantity = quantity
  }

  function updatePrices(priceMap: Record<number, number>) {
    items.value.forEach((item) => {
      if (priceMap[item.product_id] !== undefined) {
        item.unit_price = priceMap[item.product_id]
      }
    })
  }

  function clearCart() {
    items.value = []
    stopCheckoutTimer()
  }

  // 進入結帳頁面時啟動 5 分鐘計時
  function startCheckoutTimer() {
    checkoutCountdown.value = CHECKOUT_TIMEOUT
    stopCheckoutTimer()
    checkoutTimer = setInterval(() => {
      if (checkoutCountdown.value > 0) {
        checkoutCountdown.value--
      } else {
        stopCheckoutTimer()
      }
    }, 1000)
  }

  function stopCheckoutTimer() {
    if (checkoutTimer) {
      clearInterval(checkoutTimer)
      checkoutTimer = null
    }
  }

  return {
    items,
    total,
    checkoutCountdown,
    isCheckoutExpired,
    addItem,
    removeItem,
    updateQuantity,
    updatePrices,
    clearCart,
    startCheckoutTimer,
    stopCheckoutTimer,
  }
})

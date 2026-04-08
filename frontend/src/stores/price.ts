import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export interface SpotPrice {
  xag_usd: number
  usd_twd: number
  updated_at: string
}

export const usePriceStore = defineStore('price', () => {
  const snapshot = ref<SpotPrice | null>(null)
  const countdown = ref(60)
  const loading = ref(false)

  let countdownTimer: ReturnType<typeof setInterval> | null = null

  async function fetchPrice() {
    loading.value = true
    try {
      const res = await api.get<SpotPrice>('/api/v1/spot-price')
      snapshot.value = res.data
      countdown.value = 60
    } finally {
      loading.value = false
    }
  }

  function startAutoRefresh() {
    fetchPrice()

    // 倒數計時
    countdownTimer = setInterval(() => {
      if (countdown.value <= 1) {
        fetchPrice()
      } else {
        countdown.value--
      }
    }, 1000)
  }

  function stopAutoRefresh() {
    if (countdownTimer) {
      clearInterval(countdownTimer)
      countdownTimer = null
    }
  }

  return { snapshot, countdown, loading, fetchPrice, startAutoRefresh, stopAutoRefresh }
})

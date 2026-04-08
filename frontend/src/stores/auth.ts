import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('admin_token'))
  const isLoggedIn = computed(() => !!token.value)

  async function login(password: string) {
    const res = await api.post<{ token: string }>('/api/v1/admin/login', { password })
    token.value = res.data.token
    localStorage.setItem('admin_token', res.data.token)
  }

  function logout() {
    token.value = null
    localStorage.removeItem('admin_token')
  }

  return { token, isLoggedIn, login, logout }
})

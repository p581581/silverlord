import axios from 'axios'
import { i18n } from '@/i18n'

const api = axios.create({
  baseURL: import.meta.env.BASE_URL,
  timeout: 10000,
})

// 自動附加 JWT Token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// 統一錯誤處理：將後端 code 翻譯為 i18n 訊息
api.interceptors.response.use(
  (res) => res,
  (error) => {
    const t = i18n.global.t
    let message: string

    if (!error.response) {
      message = t('errors.NETWORK_ERROR')
    } else {
      const code: string = error.response.data?.code
      message = code && t(`errors.${code}`) !== `errors.${code}`
        ? t(`errors.${code}`)
        : error.response.data?.message || t('errors.INTERNAL_ERROR')
    }

    return Promise.reject({ ...error, userMessage: message })
  },
)

export default api

const BASE: string = import.meta.env.VITE_API_BASE_URL || ''

/**
 * 將後端回傳的相對路徑（/uploads/xxx）
 * 補上 API base URL，讓瀏覽器能正確載入圖片
 */
export function resolveImage(url: string | null | undefined): string {
  if (!url) return '/placeholder.svg'
  if (url.startsWith('http')) return url
  return `${BASE}${url}`
}

import { createI18n } from 'vue-i18n'
import zhTW from './locales/zh-TW'
import en from './locales/en'
import ja from './locales/ja'

export type Locale = 'zh-TW' | 'en' | 'ja'

const STORAGE_KEY = 'locale'

function detectLocale(): Locale {
  const saved = localStorage.getItem(STORAGE_KEY) as Locale | null
  if (saved && ['zh-TW', 'en', 'ja'].includes(saved)) return saved

  const browser = navigator.language
  if (browser.startsWith('zh')) return 'zh-TW'
  if (browser.startsWith('ja')) return 'ja'
  return 'en'
}

export const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'en',
  messages: {
    'zh-TW': zhTW,
    en,
    ja,
  },
})

export function setLocale(locale: Locale) {
  i18n.global.locale.value = locale
  localStorage.setItem(STORAGE_KEY, locale)
  document.documentElement.lang = locale
}

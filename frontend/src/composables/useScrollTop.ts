import { ref, onMounted, onUnmounted } from 'vue'

export function useScrollTop() {
  const showScrollTop = ref(false)

  function onScroll() {
    showScrollTop.value = (window.scrollY ?? document.documentElement.scrollTop) > 200
  }

  function scrollToTop() {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  onMounted(() => {
    window.scrollTo(0, 0)
    window.addEventListener('scroll', onScroll, { passive: true })
  })

  onUnmounted(() => {
    window.removeEventListener('scroll', onScroll)
  })

  return { showScrollTop, scrollToTop }
}

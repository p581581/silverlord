import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // ── 前台 ────────────────────────────────────────────────────
    {
      path: '/',
      component: () => import('@/views/client/ProductList.vue'),
    },
    {
      path: '/products/:id',
      component: () => import('@/views/client/ProductDetail.vue'),
    },
    {
      path: '/cart',
      component: () => import('@/views/client/Cart.vue'),
    },
    {
      path: '/checkout',
      component: () => import('@/views/client/Checkout.vue'),
    },
    {
      path: '/payment',
      component: () => import('@/views/client/Payment.vue'),
    },
    {
      path: '/order-query',
      component: () => import('@/views/client/OrderQuery.vue'),
    },
    // ── 後台 ────────────────────────────────────────────────────
    {
      path: '/admin/login',
      component: () => import('@/views/admin/Login.vue'),
    },
    {
      path: '/admin',
      component: () => import('@/views/admin/Layout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/admin/dashboard' },
        { path: 'dashboard', component: () => import('@/views/admin/Dashboard.vue') },
        { path: 'products', component: () => import('@/views/admin/Products.vue') },
        { path: 'orders', component: () => import('@/views/admin/Orders.vue') },
        { path: 'config', component: () => import('@/views/admin/Config.vue') },
        { path: 'shipping', component: () => import('@/views/admin/Shipping.vue') },
      ],
    },
  ],
})

// 後台路由守衛
router.beforeEach((to) => {
  if (to.meta.requiresAuth) {
    const auth = useAuthStore()
    if (!auth.isLoggedIn) return '/admin/login'
  }
})

export default router

import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/index.vue')
    },
    {
      path: '/',
      name: 'layout',
      component: () => import('@/components/Layout/index.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/index.vue')
        },
        {
          path: 'stores',
          name: 'stores',
          component: () => import('@/views/stores/index.vue')
        },
        {
          path: 'stores/:id/desktop',
          name: 'store-desktop',
          component: () => import('@/views/desktop/index.vue')
        },
        {
          path: 'stores/:id/command',
          name: 'store-command',
          component: () => import('@/views/command/index.vue')
        },
        {
          path: 'stores/:id/process',
          name: 'store-process',
          component: () => import('@/views/process/index.vue')
        },
        {
          path: 'stores/:id/resource',
          name: 'store-resource',
          component: () => import('@/views/resource/index.vue')
        },
        {
          path: 'audit-logs',
          name: 'audit-logs',
          component: () => import('@/views/audit/index.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  authStore.init()

  if (to.path !== '/login' && !authStore.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && authStore.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
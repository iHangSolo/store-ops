import { createRouter, createWebHistory } from 'vue-router'

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
        }
      ]
    }
  ]
})

export default router
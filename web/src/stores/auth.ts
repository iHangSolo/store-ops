import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '@/api/auth'
import type { LoginRequest, LoginResult } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const adminInfo = ref<LoginResult['admin'] | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(data: LoginRequest) {
    const result = await authApi.login(data)
    token.value = result.token
    adminInfo.value = result.admin
    localStorage.setItem('token', result.token)
    return result
  }

  async function logout() {
    try {
      await authApi.logout()
    } finally {
      token.value = null
      adminInfo.value = null
      localStorage.removeItem('token')
    }
  }

  function init() {
    const savedToken = localStorage.getItem('token')
    if (savedToken) {
      token.value = savedToken
    }
  }

  return {
    token,
    adminInfo,
    isLoggedIn,
    login,
    logout,
    init
  }
})
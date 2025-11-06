import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { userApi } from '@/api/modules/user'
import type { UserInfo, LoginParams, RegisterParams } from '@/types/api'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>('')
  const userInfo = ref<UserInfo | null>(null)
  
  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => userInfo.value?.role === 1)
  
  async function login(params: LoginParams) {
    const data = await userApi.login(params)
    token.value = data.token
    userInfo.value = data.user
  }
  
  async function register(params: RegisterParams) {
    const data = await userApi.register(params)
    token.value = data.token
    userInfo.value = data.user
  }
  
  function logout() {
    token.value = ''
    userInfo.value = null
  }
  
  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    login,
    register,
    logout,
  }
}, {
  persist: {
    key: 'user-store',
    storage: localStorage,
    paths: ['token', 'userInfo'],
  },
})

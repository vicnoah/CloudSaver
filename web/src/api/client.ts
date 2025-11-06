import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios'
import { useUserStore } from '@/stores/user'
import type { ApiResponse } from '@/types/api'

// 创建自定义 axios 实例类型
interface CustomAxiosInstance extends AxiosInstance {
  get<T = any>(url: string, config?): Promise<T>
  post<T = any>(url: string, data?, config?): Promise<T>
  put<T = any>(url: string, data?, config?): Promise<T>
  delete<T = any>(url: string, config?): Promise<T>
}

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    const { success, data, error } = response.data as ApiResponse
    if (!success) {
      console.error('API Error:', error)
      return Promise.reject(new Error(error || '请求失败'))
    }
    return data
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default apiClient as CustomAxiosInstance

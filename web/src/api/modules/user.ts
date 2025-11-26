import apiClient from '../client'
import type { LoginParams, RegisterParams, LoginResponse } from '@/types/api'

export const userApi = {
  login(params: LoginParams) {
    return apiClient.post<LoginResponse>('/user/login', params)
  },
  
  register(params: RegisterParams) {
    return apiClient.post<LoginResponse>('/user/register', params)
  },
}

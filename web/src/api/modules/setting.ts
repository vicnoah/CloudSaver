import apiClient from '../client'
import type { UserSetting } from '@/types/api'

export const settingApi = {
  get() {
    return apiClient.get<UserSetting>('/setting/get')
  },
  
  save(data: Partial<UserSetting>) {
    return apiClient.post('/setting/save', data)
  },
}

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { settingApi } from '@/api/modules/setting'
import type { UserSetting } from '@/types/api'

export const useSettingStore = defineStore('setting', () => {
  const setting = ref<UserSetting>({
    cloud115Cookie: '',
    quarkCookie: '',
  })
  
  async function fetchSetting() {
    const data = await settingApi.get()
    setting.value = data
  }
  
  async function saveSetting(data: Partial<UserSetting>) {
    await settingApi.save(data)
    await fetchSetting()
  }
  
  return {
    setting,
    fetchSetting,
    saveSetting,
  }
})

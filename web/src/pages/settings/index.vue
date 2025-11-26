<template>
  <div class="settings-container p-6">
    <div class="header flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">设置</h1>
      <router-link to="/" class="btn-secondary">返回首页</router-link>
    </div>

    <div class="content card max-w-2xl">
      <form @submit.prevent="handleSave" class="space-y-6">
        <div>
          <label class="block text-sm font-medium mb-2">115 网盘 Cookie</label>
          <textarea 
            v-model="form.cloud115Cookie"
            class="input-base w-full min-h-32"
            placeholder="请输入 115 网盘的 Cookie"
          />
        </div>

        <div>
          <label class="block text-sm font-medium mb-2">夸克网盘 Cookie</label>
          <textarea 
            v-model="form.quarkCookie"
            class="input-base w-full min-h-32"
            placeholder="请输入夸克网盘的 Cookie"
          />
        </div>

        <div v-if="message" :class="['text-sm', messageType === 'success' ? 'text-success' : 'text-danger']">
          {{ message }}
        </div>

        <div class="flex gap-4">
          <button type="submit" class="btn-primary" :disabled="loading">
            {{ loading ? '保存中...' : '保存设置' }}
          </button>
          <button type="button" @click="fetchSettings" class="btn-secondary">
            刷新
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSettingStore } from '@/stores/setting'

const settingStore = useSettingStore()

const loading = ref(false)
const message = ref('')
const messageType = ref<'success' | 'error'>('success')

const form = ref({
  cloud115Cookie: '',
  quarkCookie: '',
})

async function fetchSettings() {
  try {
    await settingStore.fetchSetting()
    form.value.cloud115Cookie = settingStore.setting.cloud115Cookie || ''
    form.value.quarkCookie = settingStore.setting.quarkCookie || ''
  } catch (err: any) {
    showMessage(err.message || '获取设置失败', 'error')
  }
}

async function handleSave() {
  loading.value = true
  message.value = ''
  
  try {
    await settingStore.saveSetting(form.value)
    showMessage('保存成功', 'success')
  } catch (err: any) {
    showMessage(err.message || '保存失败', 'error')
  } finally {
    loading.value = false
  }
}

function showMessage(msg: string, type: 'success' | 'error') {
  message.value = msg
  messageType.value = type
  setTimeout(() => {
    message.value = ''
  }, 3000)
}

onMounted(() => {
  fetchSettings()
})
</script>

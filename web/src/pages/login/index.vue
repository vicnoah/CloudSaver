<template>
  <div class="login-container flex items-center justify-center min-h-screen bg-gray-100">
    <div class="login-box bg-white rounded-lg shadow-lg p-8 w-96">
      <h1 class="text-2xl font-bold text-center mb-6">CloudSaver</h1>
      
      <div class="tabs flex mb-6">
        <button 
          :class="['flex-1 py-2', isLogin ? 'border-b-2 border-primary text-primary' : 'text-gray-500']"
          @click="isLogin = true"
        >
          登录
        </button>
        <button 
          :class="['flex-1 py-2', !isLogin ? 'border-b-2 border-primary text-primary' : 'text-gray-500']"
          @click="isLogin = false"
        >
          注册
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <input 
            v-model="form.username"
            type="text"
            placeholder="用户名"
            class="input-base w-full"
            required
          />
        </div>
        
        <div>
          <input 
            v-model="form.password"
            type="password"
            placeholder="密码"
            class="input-base w-full"
            required
          />
        </div>

        <div v-if="!isLogin">
          <input 
            v-model="form.registerCode"
            type="number"
            placeholder="注册码"
            class="input-base w-full"
            required
          />
        </div>

        <div v-if="error" class="text-danger text-sm">{{ error }}</div>

        <button 
          type="submit"
          class="btn-primary w-full"
          :disabled="loading"
        >
          {{ loading ? '处理中...' : (isLogin ? '登录' : '注册') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const form = ref({
  username: '',
  password: '',
  registerCode: 0,
})

async function handleSubmit() {
  error.value = ''
  loading.value = true
  
  try {
    if (isLogin.value) {
      await userStore.login({
        username: form.value.username,
        password: form.value.password,
      })
    } else {
      await userStore.register({
        username: form.value.username,
        password: form.value.password,
        registerCode: form.value.registerCode,
      })
    }
    router.push('/')
  } catch (err: any) {
    error.value = err.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

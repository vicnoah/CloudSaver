<template>
  <div class="home-container p-6">
    <div class="header flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold">CloudSaver</h1>
      <div class="flex gap-4">
        <router-link to="/settings" class="btn-secondary">设置</router-link>
        <button @click="handleLogout" class="btn-secondary">退出</button>
      </div>
    </div>

    <div class="content card">
      <h2 class="text-xl font-bold mb-4">欢迎使用 CloudSaver</h2>
      <p class="text-gray-600">
        这是一个云盘资源管理系统，支持 115 网盘和夸克网盘。
      </p>
      <p class="text-gray-600 mt-2">
        用户: {{ userStore.userInfo?.username }} 
        ({{ userStore.isAdmin ? '管理员' : '普通用户' }})
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

function handleLogout() {
  userStore.logout()
  router.push('/login')
}
</script>

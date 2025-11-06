<template>
  <div class="home-container min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center gap-8">
            <h1 class="text-2xl font-bold text-primary">CloudSaver</h1>
            <nav class="hidden md:flex gap-6">
              <button 
                :class="['nav-link', currentTab === 'search' && 'active']"
                @click="currentTab = 'search'"
              >
                <i class="i-carbon-search mr-1"></i>
                资源搜索
              </button>
              <button 
                :class="['nav-link', currentTab === 'douban' && 'active']"
                @click="currentTab = 'douban'"
              >
                <i class="i-carbon-movie mr-1"></i>
                豆瓣榜单
              </button>
              <button 
                :class="['nav-link', currentTab === 'transfer' && 'active']"
                @click="currentTab = 'transfer'"
              >
                <i class="i-carbon-cloud-upload mr-1"></i>
                快速转存
              </button>
            </nav>
          </div>
          <div class="flex items-center gap-4">
            <span class="text-sm text-gray-600">
              {{ userStore.userInfo?.username }}
              <span class="text-xs text-primary">
                ({{ userStore.isAdmin ? '管理员' : '普通用户' }})
              </span>
            </span>
            <router-link to="/settings" class="btn-secondary">
              <i class="i-carbon-settings mr-1"></i>
              设置
            </router-link>
            <button @click="handleLogout" class="btn-secondary">
              <i class="i-carbon-logout mr-1"></i>
              退出
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- 搜索功能 -->
      <div v-if="currentTab === 'search'" class="space-y-6">
        <div class="card">
          <h2 class="text-xl font-bold mb-4">资源搜索</h2>
          <div class="flex gap-4">
            <input 
              v-model="searchKeyword"
              type="text"
              placeholder="输入关键词搜索资源..."
              class="input-base flex-1"
              @keypress.enter="handleSearch"
            />
            <button @click="handleSearch" class="btn-primary" :disabled="searching">
              <i class="i-carbon-search mr-1"></i>
              {{ searching ? '搜索中...' : '搜索' }}
            </button>
          </div>
        </div>

        <!-- 搜索结果 -->
        <div v-if="searchResults.length > 0" class="space-y-4">
          <div v-for="channel in searchResults" :key="channel.id" class="card">
            <div class="flex items-center gap-3 mb-4">
              <img 
                v-if="channel.channelInfo?.channelLogo" 
                :src="channel.channelInfo.channelLogo"
                class="w-10 h-10 rounded-full"
                alt="频道头像"
              />
              <h3 class="text-lg font-bold">{{ channel.channelInfo?.name }}</h3>
              <span class="text-sm text-gray-500">({{ channel.list.length }} 条结果)</span>
            </div>

            <div class="space-y-3">
              <div 
                v-for="item in channel.list" 
                :key="item.messageId"
                class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
              >
                <div class="flex gap-4">
                  <img 
                    v-if="item.image"
                    :src="item.image"
                    class="w-24 h-24 object-cover rounded"
                    alt="封面"
                  />
                  <div class="flex-1">
                    <h4 class="font-bold text-lg mb-2">{{ item.title }}</h4>
                    <p class="text-sm text-gray-600 mb-2">{{ item.content }}</p>
                    <div class="flex items-center gap-2 mb-2">
                      <span 
                        v-for="tag in item.tags" 
                        :key="tag"
                        class="text-xs bg-gray-100 px-2 py-1 rounded"
                      >
                        {{ tag }}
                      </span>
                    </div>
                    <div class="flex items-center justify-between">
                      <div class="flex items-center gap-2">
                        <span :class="['text-xs px-2 py-1 rounded', getCloudTypeClass(item.cloudType)]">
                          {{ getCloudTypeName(item.cloudType) }}
                        </span>
                        <span class="text-xs text-gray-500">{{ formatDate(item.pubDate) }}</span>
                      </div>
                      <button 
                        @click="handleTransfer(item.cloudLinks[0])"
                        class="btn-primary btn-sm"
                      >
                        <i class="i-carbon-cloud-upload mr-1"></i>
                        转存
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-else-if="searched && !searching" class="card text-center text-gray-500 py-12">
          <i class="i-carbon-search text-4xl mb-2"></i>
          <p>暂无搜索结果</p>
        </div>
      </div>

      <!-- 豆瓣榜单 -->
      <div v-if="currentTab === 'douban'" class="space-y-6">
        <div class="card">
          <h2 class="text-xl font-bold mb-4">豆瓣榜单</h2>
          <div class="flex gap-4 flex-wrap">
            <select v-model="doubanType" class="input-base">
              <option value="movie">电影</option>
              <option value="tv">电视剧</option>
            </select>
            <select v-model="doubanTag" class="input-base">
              <option value="热门">热门</option>
              <option value="最新">最新</option>
              <option value="经典">经典</option>
              <option value="豆瓣高分">豆瓣高分</option>
            </select>
            <button @click="fetchDoubanList" class="btn-primary" :disabled="loadingDouban">
              <i class="i-carbon-renew mr-1"></i>
              {{ loadingDouban ? '加载中...' : '刷新' }}
            </button>
          </div>
        </div>

        <div v-if="doubanList.length > 0" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
          <div 
            v-for="item in doubanList"
            :key="item.id"
            class="card p-0 overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
            @click="openDoubanLink(item.url)"
          >
            <img :src="item.cover" class="w-full h-64 object-cover" :alt="item.title" />
            <div class="p-3">
              <h3 class="font-bold text-sm mb-1 truncate" :title="item.title">{{ item.title }}</h3>
              <div class="flex items-center justify-between">
                <span class="text-warning text-sm">⭐ {{ item.rate }}</span>
                <span v-if="item.is_new" class="text-xs bg-danger text-white px-2 py-0.5 rounded">NEW</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 快速转存 -->
      <div v-if="currentTab === 'transfer'" class="space-y-6">
        <div class="card">
          <h2 class="text-xl font-bold mb-4">快速转存</h2>
          <p class="text-sm text-gray-600 mb-4">
            支持 115 网盘和夸克网盘分享链接，粘贴链接后自动识别类型
          </p>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">分享链接</label>
              <input 
                v-model="transferLink"
                type="text"
                placeholder="粘贴分享链接，如: https://115.com/s/xxx?password=xxx"
                class="input-base w-full"
              />
            </div>
            
            <div>
              <label class="block text-sm font-medium mb-2">提取码（如有）</label>
              <input 
                v-model="transferPasscode"
                type="text"
                placeholder="提取码"
                class="input-base w-full"
              />
            </div>

            <button @click="parseShareLink" class="btn-primary" :disabled="parsingLink">
              <i class="i-carbon-link mr-1"></i>
              {{ parsingLink ? '解析中...' : '解析链接' }}
            </button>
          </div>
        </div>

        <!-- 解析结果 -->
        <div v-if="shareInfo" class="card">
          <h3 class="text-lg font-bold mb-4">文件列表</h3>
          <div class="space-y-2 mb-4 max-h-96 overflow-y-auto">
            <div 
              v-for="file in shareInfo.list"
              :key="file.file_id"
              class="flex items-center gap-2 p-2 border border-gray-200 rounded hover:bg-gray-50"
            >
              <input 
                type="checkbox"
                :value="file.file_id"
                v-model="selectedFiles"
                class="w-4 h-4"
              />
              <i :class="[file.file_type === 1 ? 'i-carbon-folder' : 'i-carbon-document']"></i>
              <span class="flex-1">{{ file.file_name }}</span>
            </div>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-2">选择目标文件夹</label>
              <button @click="showFolderSelectorDialog()" class="btn-secondary w-full">
                <i class="i-carbon-folder mr-1"></i>
                {{ selectedFolder ? selectedFolder.name : '点击选择文件夹' }}
              </button>
            </div>

            <button 
              @click="handleSaveFiles"
              class="btn-primary w-full"
              :disabled="selectedFiles.length === 0 || !selectedFolder || saving"
            >
              <i class="i-carbon-cloud-upload mr-1"></i>
              {{ saving ? '转存中...' : `转存已选文件 (${selectedFiles.length})` }}
            </button>
          </div>
        </div>

        <!-- 文件夹选择器弹窗 -->
        <div v-if="showFolderSelector" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showFolderSelector = false">
          <div class="bg-white rounded-lg p-6 w-full max-w-md max-h-[80vh] overflow-y-auto">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-bold">选择文件夹</h3>
              <button @click="showFolderSelector = false" class="text-gray-500 hover:text-gray-700">
                <i class="i-carbon-close text-xl"></i>
              </button>
            </div>
            
            <div v-if="loadingFolders" class="text-center py-8 text-gray-500">
              <i class="i-carbon-circle-dash animate-spin text-2xl"></i>
              <p class="mt-2">加载中...</p>
            </div>
            
            <div v-else class="space-y-2">
              <div 
                v-for="folder in folders"
                :key="folder.cid"
                class="p-3 border border-gray-200 rounded hover:bg-gray-50 cursor-pointer"
                :class="{ 'bg-primary-50 border-primary': selectedFolder?.cid === folder.cid }"
                @click="selectFolder(folder)"
              >
                <i class="i-carbon-folder mr-2"></i>
                {{ folder.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { cloud115Api } from '@/api/modules/cloud115'
import { quarkApi } from '@/api/modules/quark'
import { searchApi } from '@/api/modules/search'
import { doubanApi } from '@/api/modules/douban'
import type { ShareInfoResponse, FolderItem, Save115FileParams, SaveQuarkFileParams } from '@/types/api'

const router = useRouter()
const userStore = useUserStore()

// 当前选项卡
const currentTab = ref<'search' | 'douban' | 'transfer'>('search')

// 搜索相关
const searchKeyword = ref('')
const searching = ref(false)
const searched = ref(false)
const searchResults = ref<any[]>([])

// 豆瓣相关
const doubanType = ref('movie')
const doubanTag = ref('热门')
const loadingDouban = ref(false)
const doubanList = ref<any[]>([])

// 转存相关
const transferLink = ref('')
const transferPasscode = ref('')
const parsingLink = ref(false)
const shareInfo = ref<ShareInfoResponse | null>(null)
const selectedFiles = ref<string[]>([])
const showFolderSelector = ref(false)
const loadingFolders = ref(false)
const folders = ref<FolderItem[]>([])
const selectedFolder = ref<FolderItem | null>(null)
const saving = ref(false)
const currentCloudType = ref<'115' | 'quark' | ''>('')

function handleLogout() {
  userStore.logout()
  router.push('/login')
}

async function handleSearch() {
  if (!searchKeyword.value.trim()) return
  
  searching.value = true
  searched.value = true
  try {
    searchResults.value = await searchApi.search({ keyword: searchKeyword.value })
  } catch (error: any) {
    console.error('搜索失败:', error)
    alert(error.message || '搜索失败')
  } finally {
    searching.value = false
  }
}

async function fetchDoubanList() {
  loadingDouban.value = true
  try {
    doubanList.value = await doubanApi.getHotList({
      type: doubanType.value,
      tag: doubanTag.value,
      page_limit: '20',
      page_start: '0',
    })
  } catch (error: any) {
    console.error('获取豆瓣榜单失败:', error)
    alert(error.message || '获取榜单失败')
  } finally {
    loadingDouban.value = false
  }
}

async function parseShareLink() {
  if (!transferLink.value.trim()) return
  
  parsingLink.value = true
  try {
    // 识别链接类型
    if (transferLink.value.includes('115.com')) {
      currentCloudType.value = '115'
      const match = transferLink.value.match(/s\/([a-zA-Z0-9]+)/)
      if (match) {
        const shareCode = match[1]
        shareInfo.value = await cloud115Api.getShareInfo(shareCode, transferPasscode.value)
        selectedFiles.value = []
      }
    } else if (transferLink.value.includes('pan.quark.cn')) {
      currentCloudType.value = 'quark'
      const match = transferLink.value.match(/s\/([a-zA-Z0-9]+)/)
      if (match) {
        const shareCode = match[1]
        shareInfo.value = await quarkApi.getShareInfo(shareCode, transferPasscode.value)
        selectedFiles.value = []
      }
    } else {
      throw new Error('不支持的链接类型')
    }
  } catch (error: any) {
    console.error('解析链接失败:', error)
    alert(error.message || '解析链接失败')
  } finally {
    parsingLink.value = false
  }
}

async function showFolderSelectorDialog() {
  showFolderSelector.value = true
  loadingFolders.value = true
  try {
    if (currentCloudType.value === '115') {
      const data = await cloud115Api.getFolderList()
      folders.value = data.folders
    } else if (currentCloudType.value === 'quark') {
      const data = await quarkApi.getFolderList()
      folders.value = data.folders
    }
  } catch (error: any) {
    console.error('获取文件夹列表失败:', error)
  } finally {
    loadingFolders.value = false
  }
}

function selectFolder(folder: FolderItem) {
  selectedFolder.value = folder
  showFolderSelector.value = false
}

async function handleSaveFiles() {
  if (selectedFiles.value.length === 0 || !selectedFolder.value) return
  
  saving.value = true
  try {
    const fileTokens = shareInfo.value?.list
      .filter(f => selectedFiles.value.includes(f.file_id))
      .map(f => f.file_id_token || '')
    
    if (currentCloudType.value === '115') {
      await cloud115Api.saveFile({
        shareCode: transferLink.value.match(/s\/([a-zA-Z0-9]+)/)?.[1] || '',
        passcode: transferPasscode.value,
        folderId: selectedFolder.value.cid,
        fileIds: selectedFiles.value,
      })
    } else if (currentCloudType.value === 'quark') {
      await quarkApi.saveFile({
        folderId: selectedFolder.value.cid,
        fileIds: selectedFiles.value,
        fileTokens: fileTokens || [],
        pwdId: shareInfo.value?.pwd_id || '',
        stoken: shareInfo.value?.stoken || '',
      })
    }
    
    alert('转存成功！')
    shareInfo.value = null
    selectedFiles.value = []
    selectedFolder.value = null
    transferLink.value = ''
    transferPasscode.value = ''
  } catch (error: any) {
    console.error('转存失败:', error)
    alert(error.message || '转存失败')
  } finally {
    saving.value = false
  }
}

function handleTransfer(link: string) {
  transferLink.value = link
  currentTab.value = 'transfer'
}

function getCloudTypeClass(type: string) {
  const classes: Record<string, string> = {
    'pan115': 'bg-blue-100 text-blue-600',
    'quark': 'bg-purple-100 text-purple-600',
    'baidu': 'bg-green-100 text-green-600',
    'aliyun': 'bg-orange-100 text-orange-600',
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

function getCloudTypeName(type: string) {
  const names: Record<string, string> = {
    'pan115': '115网盘',
    'quark': '夸克网盘',
    'baidu': '百度网盘',
    'aliyun': '阿里云盘',
  }
  return names[type] || '未知'
}

function formatDate(date: string) {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString()
}

function openDoubanLink(url: string) {
  window.open(url, '_blank')
}

onMounted(() => {
  // 初始化加载豆瓣榜单
  fetchDoubanList()
})
</script>

<style scoped>
.nav-link {
  @apply flex items-center text-gray-600 hover:text-primary transition-colors px-3 py-2 rounded;
}

.nav-link.active {
  @apply text-primary font-medium bg-primary-50;
}

.btn-sm {
  @apply text-sm px-3 py-1;
}
</style>

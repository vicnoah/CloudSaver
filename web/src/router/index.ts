import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/pages/home/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/pages/settings/index.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 防止循环重定向：如果 from 和 to 是同一个路由，直接放行
  if (from.name === to.name && from.path === to.path) {
    next()
    return
  }
  
  // 如果目标路由是登录页
  if (to.name === 'Login') {
    // 如果已登录且不是从首页跳转过来的，重定向到首页
    if (userStore.isLoggedIn && from.name !== 'Home') {
      next({ name: 'Home', replace: true })
      return
    }
    next()
    return
  }
  
  // 如果目标路由需要认证
  if (to.meta.requiresAuth) {
    if (!userStore.isLoggedIn) {
      // 未登录，重定向到登录页
      next({ name: 'Login', query: { redirect: to.fullPath }, replace: true })
      return
    }
    next()
    return
  }
  
  // 其他情况直接放行
  next()
})

export default router

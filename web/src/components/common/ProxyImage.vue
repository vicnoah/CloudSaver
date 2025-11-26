<template>
  <img 
    :src="proxiedSrc"
    :alt="alt"
    :class="className"
    v-bind="$attrs"
    @error="handleError"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
  src: string
  alt?: string
  className?: string
  useProxy?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  className: '',
  useProxy: true,
})

const hasError = ref(false)

const proxiedSrc = computed(() => {
  if (!props.src || hasError.value) {
    return ''
  }
  
  // 如果是豆瓣图片或者启用了代理，使用代理
  if (props.useProxy && (props.src.includes('doubanio.com') || props.src.includes('douban.com'))) {
    return `/api/image/proxy?url=${encodeURIComponent(props.src)}`
  }
  
  return props.src
})

function handleError() {
  // 如果代理失败，标记错误避免无限重试
  hasError.value = true
  console.error('图片加载失败:', props.src)
}
</script>

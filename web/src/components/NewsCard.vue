<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const props = defineProps({
  news: {
    type: Object,
    required: true,
  },
})

const router = useRouter()

// 格式化发布时间
const formattedTime = computed(() => {
  return dayjs(props.news.published).fromNow()
})

// 提取纯文本摘要
const summary = computed(() => {
  if (!props.news.content) return ''
  // 移除 HTML 标签
  const text = props.news.content.replace(/<[^>]+>/g, '')
  // 截取前 150 个字符
  return text.length > 150 ? text.slice(0, 150) + '...' : text
})

// 新标签页打开详情
function goToDetail() {
  const url = router.resolve({ name: 'detail', params: { id: props.news.id } }).href
  window.open(url, '_blank')
}
</script>

<template>
  <div class="news-card" @click="goToDetail">
    <h3 class="title" :title="news.title">{{ news.title }}</h3>
    <p class="summary">{{ summary }}</p>
    <div class="meta">
      <span class="author" v-if="news.author">
        <span class="meta-label">作者:</span> {{ news.author }}
      </span>
      <span class="divider" v-if="news.author">|</span>
      <span class="time">
        <span class="meta-label">发布:</span> {{ formattedTime }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.news-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #f0f0f0;
}

.news-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(211, 47, 47, 0.15);
  border-color: rgba(211, 47, 47, 0.3);
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: rgba(0, 0, 0, 0.88);
  margin-bottom: 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.news-card:hover .title {
  color: #d32f2f;
}

.summary {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.55);
  line-height: 1.6;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 13px;
  color: rgba(0, 0, 0, 0.45);
}

.meta-label {
  color: rgba(0, 0, 0, 0.35);
}

.divider {
  color: rgba(0, 0, 0, 0.15);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .news-card {
    padding: 16px;
  }

  .title {
    font-size: 16px;
  }

  .summary {
    font-size: 13px;
  }
}
</style>

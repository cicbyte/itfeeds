<script setup>
import { onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Skeleton, Button, Tag } from 'ant-design-vue'
import { ArrowLeftOutlined, LinkOutlined, ClockCircleOutlined, UserOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import { useRssStore } from '@/stores/rss'

const route = useRoute()
const router = useRouter()
const rssStore = useRssStore()

// 当前文章
const article = computed(() => rssStore.currentDetail)

// 格式化时间
const formattedTime = computed(() => {
  if (!article.value?.published) return ''
  return dayjs(article.value.published).format('YYYY-MM-DD HH:mm')
})

// 加载详情
onMounted(async () => {
  const id = route.params.id
  if (id) {
    await rssStore.fetchDetail(Number(id))
  }
})

// 返回列表
function goBack() {
  router.push({ name: 'home' })
}

// 打开原文
function openOriginal() {
  if (article.value?.url) {
    window.open(article.value.url, '_blank')
  }
}
</script>

<template>
  <div class="detail-view">
    <div class="page-container">
      <!-- 加载状态 -->
      <div v-if="rssStore.detailLoading" class="skeleton-wrapper">
        <Skeleton active :paragraph="{ rows: 10 }" />
      </div>

      <!-- 文章内容 -->
      <article v-else-if="article" class="article-wrapper">
        <!-- 文章头部 -->
        <header class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="article-meta">
            <span class="meta-item" v-if="article.author">
              <UserOutlined />
              <span>{{ article.author }}</span>
            </span>
            <span class="meta-item">
              <ClockCircleOutlined />
              <span>{{ formattedTime }}</span>
            </span>
            <Tag color="error" v-if="article.guid">ITFeeds</Tag>
          </div>
        </header>

        <!-- 文章内容 -->
        <div class="article-content" v-html="article.content"></div>

        <!-- 底部操作栏 -->
        <footer class="article-footer">
          <Button @click="goBack">
            <template #icon><ArrowLeftOutlined /></template>
            返回列表
          </Button>
          <Button type="primary" @click="openOriginal" v-if="article.url">
            <template #icon><LinkOutlined /></template>
            查看原文
          </Button>
        </footer>
      </article>

      <!-- 空状态 -->
      <div v-else class="empty-state">
        <p class="empty-text">文章不存在或已删除</p>
        <Button type="primary" @click="goBack">返回首页</Button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.detail-view {
  flex: 1;
  padding: 24px 0;
  background: #fff;
}

.page-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 24px;
}

.skeleton-wrapper {
  padding: 24px 0;
}

.article-wrapper {
  background: #fff;
}

.article-header {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f0f0f0;
}

.article-title {
  font-size: 32px;
  font-weight: 600;
  color: #d32f2f;
  line-height: 1.4;
  margin-bottom: 16px;
}

.article-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  color: rgba(0, 0, 0, 0.45);
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.article-content {
  font-size: 16px;
  line-height: 1.8;
  color: rgba(0, 0, 0, 0.75);
}

.article-content :deep(img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 16px auto;
  border-radius: 4px;
}

.article-content :deep(p) {
  margin-bottom: 16px;
}

.article-content :deep(h2) {
  font-size: 24px;
  font-weight: 600;
  margin: 32px 0 16px;
  color: rgba(0, 0, 0, 0.88);
}

.article-content :deep(h3) {
  font-size: 20px;
  font-weight: 600;
  margin: 24px 0 12px;
  color: rgba(0, 0, 0, 0.88);
}

.article-content :deep(ul),
.article-content :deep(ol) {
  padding-left: 24px;
  margin-bottom: 16px;
}

.article-content :deep(li) {
  margin-bottom: 8px;
}

.article-content :deep(blockquote) {
  border-left: 4px solid #d32f2f;
  padding-left: 16px;
  margin: 16px 0;
  color: rgba(0, 0, 0, 0.55);
}

.article-content :deep(a) {
  color: #d32f2f;
  text-decoration: none;
}

.article-content :deep(a:hover) {
  text-decoration: underline;
}

.article-content :deep(pre) {
  background: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 16px 0;
}

.article-content :deep(code) {
  font-family: 'SF Mono', Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
  font-size: 14px;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 48px;
  padding-top: 24px;
  border-top: 1px solid #f0f0f0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-text {
  font-size: 16px;
  color: rgba(0, 0, 0, 0.45);
  margin-bottom: 24px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .detail-view {
    padding: 16px 0;
  }

  .page-container {
    padding: 0 16px;
  }

  .article-title {
    font-size: 24px;
  }

  .article-meta {
    gap: 12px;
    font-size: 13px;
  }

  .article-content {
    font-size: 15px;
  }

  .article-footer {
    flex-direction: column;
    gap: 12px;
  }

  .article-footer :deep(.ant-btn) {
    width: 100%;
  }
}
</style>

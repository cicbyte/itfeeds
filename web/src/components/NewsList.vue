<script setup>
import { Skeleton } from 'ant-design-vue'
import NewsCard from './NewsCard.vue'

defineProps({
  list: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
})
</script>

<template>
  <div class="news-list">
    <!-- 加载状态 -->
    <template v-if="loading">
      <div v-for="i in 6" :key="i" class="skeleton-card">
        <Skeleton active :paragraph="{ rows: 3 }" />
      </div>
    </template>

    <!-- 新闻列表 -->
    <template v-else-if="list.length > 0">
      <NewsCard
        v-for="news in list"
        :key="news.id"
        :news="news"
      />
    </template>

    <!-- 空状态 -->
    <template v-else>
      <div class="empty-state">
        <div class="empty-icon">:(</div>
        <p class="empty-text">暂无新闻数据</p>
        <p class="empty-hint">请稍后再来查看</p>
      </div>
    </template>
  </div>
</template>

<style scoped>
.news-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.skeleton-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  border: 1px solid #f0f0f0;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 48px;
  color: #d32f2f;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: rgba(0, 0, 0, 0.45);
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.25);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .news-list {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .skeleton-card {
    padding: 16px;
  }
}
</style>

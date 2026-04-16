<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Pagination, Button, message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { useRssStore } from '@/stores/rss'
import { syncRss } from '@/api/rss'
import NewsList from '@/components/NewsList.vue'

const route = useRoute()
const rssStore = useRssStore()
const syncing = ref(false)

async function handleSync() {
  syncing.value = true
  try {
    const res = await syncRss()
    message.success(res.message || '拉取完成')
    rssStore.fetchList()
  } catch {
    message.error('拉取失败')
  } finally {
    syncing.value = false
  }
}

// 初始化加载
onMounted(() => {
  // 从 URL 获取搜索关键词
  const keyword = route.query.keyword || ''
  if (keyword) {
    rssStore.search(keyword)
  } else {
    rssStore.fetchList()
  }
})

// 监听路由变化
watch(
  () => route.query.keyword,
  (newKeyword) => {
    if (newKeyword !== rssStore.searchKeyword) {
      rssStore.search(newKeyword || '')
    }
  }
)

// 分页变化
function handlePageChange(page, pageSize) {
  rssStore.fetchList({ pageNum: page, pageSize })
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<template>
  <div class="home-view">
    <div class="page-container">
      <!-- 页面标题 -->
      <div class="page-header">
        <div class="page-header-top">
          <h1 class="page-title">最新资讯</h1>
          <Button :loading="syncing" @click="handleSync" size="small">
            <template #icon><ReloadOutlined /></template>
            拉取
          </Button>
        </div>
        <p class="page-desc" v-if="rssStore.searchKeyword || (rssStore.dateRange.startDate && rssStore.dateRange.endDate)">
          <template v-if="rssStore.searchKeyword">
            搜索 "<span class="keyword">{{ rssStore.searchKeyword }}</span>"
          </template>
          <template v-if="rssStore.dateRange.startDate && rssStore.dateRange.endDate">
            <span v-if="rssStore.searchKeyword">，</span>
            时间范围：<span class="keyword">{{ rssStore.dateRange.startDate }}</span> 至 <span class="keyword">{{ rssStore.dateRange.endDate }}</span>
          </template>
          <span>，共 {{ rssStore.total }} 条</span>
        </p>
        <p class="page-desc" v-else>
          为您提供最新的科技资讯
        </p>
      </div>

      <!-- 新闻列表 -->
      <NewsList
        :list="rssStore.list"
        :loading="rssStore.loading"
      />

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="rssStore.total > rssStore.pageSize">
        <Pagination
          v-model:current="rssStore.currentPage"
          v-model:page-size="rssStore.pageSize"
          :total="rssStore.total"
          :show-size-changer="false"
          :show-quick-jumper="true"
          :show-total="(total) => `共 ${total} 条`"
          @change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  flex: 1;
  padding: 24px 0;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 28px;
  font-weight: 600;
  color: #d32f2f;
  margin-bottom: 8px;
}

.page-desc {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.45);
}

.page-desc .keyword {
  color: #d32f2f;
  font-weight: 500;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 24px 0;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .home-view {
    padding: 16px 0;
  }

  .page-container {
    padding: 0 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .pagination-wrapper {
    margin-top: 24px;
    padding: 16px 0;
  }
}
</style>

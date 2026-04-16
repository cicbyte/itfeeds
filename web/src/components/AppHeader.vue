<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Input, DatePicker } from 'ant-design-vue'
import { SearchOutlined, CloseOutlined, CalendarOutlined } from '@ant-design/icons-vue'
import { useRssStore } from '@/stores/rss'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const rssStore = useRssStore()

const searchValue = ref('')
const activeTimeFilter = ref('')
const customDateRange = ref(null)
const showDatePicker = ref(false)

// 时间范围快捷选项
const timeFilters = [
  { key: 'today', label: '今天', days: 0 },
  { key: 'yesterday', label: '昨天', days: 1 },
  { key: '3days', label: '最近三天', days: 2 },
  { key: '7days', label: '最近七天', days: 6 },
]

// 从URL同步搜索关键词到搜索框
onMounted(() => {
  if (route.query.keyword) {
    searchValue.value = route.query.keyword
  }
})

// 监听路由变化
watch(() => route.query.keyword, (newKeyword) => {
  searchValue.value = newKeyword || ''
})

// 执行搜索
function handleSearch() {
  if (searchValue.value.trim()) {
    rssStore.search(searchValue.value.trim())
    router.push({ name: 'home', query: { keyword: searchValue.value.trim() } })
  } else {
    // 搜索为空时清除搜索
    clearSearch()
  }
}

// 清除搜索
function clearSearch() {
  searchValue.value = ''
  rssStore.search('')
  router.push({ name: 'home' })
}

// 监听搜索框变化
function handleSearchChange(e) {
  if (!e.target.value) {
    clearSearch()
  }
}

// 清空搜索并返回首页
function goHome() {
  searchValue.value = ''
  activeTimeFilter.value = ''
  customDateRange.value = null
  showDatePicker.value = false
  rssStore.search('')
  rssStore.clearDateRange()
  router.push({ name: 'home' })
}

// 时间范围快捷选择
function handleTimeFilter(filter) {
  activeTimeFilter.value = filter.key
  customDateRange.value = null
  showDatePicker.value = false
  const today = dayjs()
  const endDate = today.format('YYYY-MM-DD')
  let startDate

  if (filter.key === 'today') {
    startDate = today.format('YYYY-MM-DD')
  } else if (filter.key === 'yesterday') {
    startDate = today.subtract(1, 'day').format('YYYY-MM-DD')
  } else {
    startDate = today.subtract(filter.days, 'day').format('YYYY-MM-DD')
  }

  rssStore.setDateRange(startDate, endDate)
  router.push({ name: 'home' })
}

// 自定义日期选择
function handleCustomDateChange(dates) {
  if (dates && dates.length === 2) {
    customDateRange.value = dates
    activeTimeFilter.value = 'custom'
    const startDate = dates[0].format('YYYY-MM-DD')
    const endDate = dates[1].format('YYYY-MM-DD')
    rssStore.setDateRange(startDate, endDate)
    router.push({ name: 'home' })
    showDatePicker.value = false
  }
}

// 切换日期选择器显示
function toggleDatePicker() {
  showDatePicker.value = !showDatePicker.value
  if (showDatePicker.value) {
    activeTimeFilter.value = 'custom'
  }
}

// 清除时间筛选
function clearTimeFilter() {
  activeTimeFilter.value = ''
  customDateRange.value = null
  showDatePicker.value = false
  rssStore.clearDateRange()
}
</script>

<template>
  <header class="app-header">
    <div class="header-container">
      <!-- Logo -->
      <div class="logo" @click="goHome">
        <span class="logo-icon">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 11a9 9 0 0 1 9 9"/>
            <path d="M4 4a16 16 0 0 1 16 16"/>
            <circle cx="5" cy="19" r="1"/>
          </svg>
        </span>
        <span class="logo-text">IT<span class="logo-accent">Feeds</span></span>
      </div>

      <!-- 时间范围筛选 -->
      <div class="time-filter">
        <span
          v-for="filter in timeFilters"
          :key="filter.key"
          class="filter-item"
          :class="{ active: activeTimeFilter === filter.key }"
          @click="handleTimeFilter(filter)"
        >
          {{ filter.label }}
        </span>
        <span
          class="filter-item custom"
          :class="{ active: activeTimeFilter === 'custom' }"
          @click="toggleDatePicker"
        >
          <CalendarOutlined />
          <span class="custom-label">自定义</span>
        </span>
        <span
          v-if="activeTimeFilter"
          class="filter-clear"
          @click="clearTimeFilter"
        >
          <CloseOutlined />
        </span>
      </div>

      <!-- 自定义日期选择器 -->
      <div class="date-picker-wrapper" v-if="showDatePicker">
        <DatePicker.RangePicker
          :value="customDateRange"
          @change="handleCustomDateChange"
          format="YYYY-MM-DD"
          :placeholder="['开始日期', '结束日期']"
          allow-clear
        />
      </div>

      <!-- 搜索框 -->
      <div class="search-box">
        <Input
          v-model:value="searchValue"
          placeholder="搜索新闻..."
          allow-clear
          @press-enter="handleSearch"
          @change="handleSearchChange"
        >
          <template #prefix>
            <SearchOutlined />
          </template>
          <template #suffix>
            <span class="search-btn" @click="handleSearch">搜索</span>
          </template>
        </Input>
      </div>
    </div>

    <!-- 移动端日期选择器 -->
    <div class="mobile-date-picker" v-if="showDatePicker">
      <DatePicker.RangePicker
        :value="customDateRange"
        @change="handleCustomDateChange"
        format="YYYY-MM-DD"
        :placeholder="['开始日期', '结束日期']"
        allow-clear
        size="small"
      />
    </div>
  </header>
</template>

<style scoped>
.app-header {
  position: sticky;
  top: 0;
  z-index: 100;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  flex-shrink: 0;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #d32f2f 0%, #f44336 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: rgba(0, 0, 0, 0.85);
  letter-spacing: -0.5px;
}

.logo-accent {
  color: #d32f2f;
}

/* 时间筛选 */
.time-filter {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.filter-item {
  padding: 4px 12px;
  font-size: 13px;
  color: rgba(0, 0, 0, 0.65);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
  white-space: nowrap;
  display: flex;
  align-items: center;
  gap: 4px;
}

.filter-item:hover {
  color: #d32f2f;
  background: rgba(211, 47, 47, 0.08);
}

.filter-item.active {
  color: #fff;
  background: #d32f2f;
}

.filter-item.custom .custom-label {
  display: inline;
}

.filter-clear {
  padding: 4px 8px;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
}

.filter-clear:hover {
  color: #d32f2f;
  background: rgba(211, 47, 47, 0.08);
}

/* 日期选择器 */
.date-picker-wrapper {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: #fff;
  padding: 16px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 200;
  margin-top: 8px;
}

.mobile-date-picker {
  display: none;
  padding: 12px 24px;
  border-top: 1px solid #f0f0f0;
}

.search-box {
  flex: 1;
  max-width: 300px;
}

.search-box :deep(.ant-input-affix-wrapper) {
  border-radius: 20px;
}

.search-box :deep(.ant-input-affix-wrapper:focus),
.search-box :deep(.ant-input-affix-wrapper-focused) {
  border-color: #d32f2f;
  box-shadow: 0 0 0 2px rgba(211, 47, 47, 0.1);
}

.search-btn {
  color: #d32f2f;
  cursor: pointer;
  font-size: 14px;
}

.search-btn:hover {
  color: #f44336;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .header-container {
    padding: 0 16px;
    gap: 8px;
    flex-wrap: wrap;
    height: auto;
    padding-top: 12px;
    padding-bottom: 0;
  }

  .logo-text {
    display: none;
  }

  .time-filter {
    order: 3;
    width: 100%;
    justify-content: flex-start;
    padding-top: 8px;
    border-top: 1px solid #f0f0f0;
    margin-top: 8px;
    overflow-x: auto;
    padding-bottom: 8px;
  }

  .filter-item {
    padding: 4px 8px;
    font-size: 12px;
    flex-shrink: 0;
  }

  .filter-item.custom .custom-label {
    display: none;
  }

  .date-picker-wrapper {
    display: none;
  }

  .mobile-date-picker {
    display: block;
    padding: 8px 16px 12px;
  }

  .search-box {
    flex: 1;
    max-width: none;
  }
}
</style>

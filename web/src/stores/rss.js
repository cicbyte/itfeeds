import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { getRssList, getRssDetail } from '@/api/rss'

export const useRssStore = defineStore('rss', () => {
  // 状态
  const list = ref([])
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const loading = ref(false)
  const searchKeyword = ref('')
  const dateRange = ref({ startDate: '', endDate: '' })
  const currentDetail = ref(null)
  const detailLoading = ref(false)

  // 获取新闻列表
  async function fetchList(params = {}) {
    loading.value = true
    try {
      const res = await getRssList({
        pageNum: params.pageNum || currentPage.value,
        pageSize: params.pageSize || pageSize.value,
        title: params.title ?? searchKeyword.value,
        startDate: params.startDate ?? dateRange.value.startDate,
        endDate: params.endDate ?? dateRange.value.endDate,
      })
      list.value = res.list || []
      total.value = res.total || 0
      currentPage.value = params.pageNum || currentPage.value
    } catch (error) {
      console.error('获取新闻列表失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 获取新闻详情
  async function fetchDetail(id) {
    detailLoading.value = true
    try {
      const res = await getRssDetail(id)
      currentDetail.value = res
      return res
    } catch (error) {
      console.error('获取新闻详情失败:', error)
      return null
    } finally {
      detailLoading.value = false
    }
  }

  // 搜索
  function search(keyword) {
    searchKeyword.value = keyword
    currentPage.value = 1
    fetchList({ pageNum: 1, title: keyword })
  }

  // 设置时间范围
  function setDateRange(startDate, endDate) {
    dateRange.value = { startDate, endDate }
    currentPage.value = 1
    fetchList({ pageNum: 1, startDate, endDate })
  }

  // 清空时间范围
  function clearDateRange() {
    dateRange.value = { startDate: '', endDate: '' }
    currentPage.value = 1
    fetchList({ pageNum: 1, startDate: '', endDate: '' })
  }

  // 清空详情
  function clearDetail() {
    currentDetail.value = null
  }

  // 是否有数据
  const hasData = computed(() => list.value.length > 0)

  return {
    list,
    total,
    currentPage,
    pageSize,
    loading,
    searchKeyword,
    dateRange,
    currentDetail,
    detailLoading,
    hasData,
    fetchList,
    fetchDetail,
    search,
    setDateRange,
    clearDateRange,
    clearDetail,
  }
})

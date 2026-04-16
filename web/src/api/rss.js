import request from './request'

/**
 * 获取新闻列表
 * @param {Object} params - 查询参数
 * @param {number} params.pageNum - 页码
 * @param {number} params.pageSize - 每页数量
 * @param {string} params.title - 搜索关键词
 * @param {string} params.startDate - 开始日期
 * @param {string} params.endDate - 结束日期
 */
export function getRssList(params) {
  return request({
    url: '/rss_entries/list',
    method: 'get',
    params,
  })
}

/**
 * 获取新闻详情
 * @param {number} id - 新闻 ID
 */
export function getRssDetail(id) {
  return request({
    url: '/rss_entries/detail',
    method: 'get',
    params: { id },
  })
}

// 手动拉取RSS
export function syncRss() {
  return request({
    url: '/rss_entries/sync',
    method: 'post',
  })
}

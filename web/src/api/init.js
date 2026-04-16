import request from './request'

// 查询初始化状态
export function apiInitStatus() {
  return request({
    url: '/init/status',
    method: 'get',
  })
}

// 测试数据库连接
export function apiTestConnection(data) {
  return request({
    url: '/init/test-connection',
    method: 'post',
    data,
  })
}

// 执行初始化
export function apiSetup(data) {
  return request({
    url: '/init/setup',
    method: 'post',
    data,
  })
}

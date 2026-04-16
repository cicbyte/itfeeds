<script setup>
import { RouterView, useRoute } from 'vue-router'
import AppHeader from '@/components/AppHeader.vue'
import AppFooter from '@/components/AppFooter.vue'

const route = useRoute()
</script>

<template>
  <RouterView v-if="route.meta.noLayout" />
  <div v-else class="app-container">
    <AppHeader />
    <main class="main-content">
      <RouterView />
    </main>
    <AppFooter />
  </div>
</template>

<style>
/* 全局重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #f5f5f5;
}

#app {
  height: 100%;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* Ant Design Vue 主题色覆盖 - ITFeeds 红色 */
.ant-btn-primary {
  background: #d32f2f !important;
  border-color: #d32f2f !important;
}

.ant-btn-primary:hover {
  background: #f44336 !important;
  border-color: #f44336 !important;
}

.ant-pagination .ant-pagination-item-active {
  border-color: #d32f2f;
}

.ant-pagination .ant-pagination-item-active a {
  color: #d32f2f;
}

.ant-pagination .ant-pagination-item:hover {
  border-color: #d32f2f;
}

.ant-pagination .ant-pagination-item:hover a {
  color: #d32f2f;
}

.ant-pagination .ant-pagination-prev:hover .ant-pagination-item-link,
.ant-pagination .ant-pagination-next:hover .ant-pagination-item-link {
  color: #d32f2f;
  border-color: #d32f2f;
}

.ant-input:focus,
.ant-input-focused {
  border-color: #d32f2f !important;
  box-shadow: 0 0 0 2px rgba(211, 47, 47, 0.1) !important;
}

.ant-input-affix-wrapper:focus,
.ant-input-affix-wrapper-focused {
  border-color: #d32f2f !important;
  box-shadow: 0 0 0 2px rgba(211, 47, 47, 0.1) !important;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>

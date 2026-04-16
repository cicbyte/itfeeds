import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { apiInitStatus } from '@/api/init'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/init',
      name: 'init',
      component: () => import('../views/InitView.vue'),
      meta: { noLayout: true },
    },
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/detail/:id',
      name: 'detail',
      component: () => import('../views/DetailView.vue'),
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})

let needsInit = null

router.beforeEach(async (to) => {
  if (needsInit === null) {
    try {
      const res = await apiInitStatus()
      needsInit = !res.initialized
    } catch {
      needsInit = true
    }
  }

  if (to.name === 'init' && !needsInit) {
    return { name: 'home' }
  }

  if (needsInit && to.name !== 'init') {
    return { name: 'init' }
  }
  return true
})

export default router

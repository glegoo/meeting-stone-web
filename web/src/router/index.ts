import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/index/HomePage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomePage
    },
    {
      path: '/activity/:gameKey/:id',
      name: 'activity',
      component: () => import('@/pages/activity/ActivityPage.vue')
    },
    {
      path: '/lineup/:gameKey/:activityId',
      name: 'lineup',
      component: () => import('@/pages/lineup/LineupPage.vue')
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('@/pages/settings/SettingsPage.vue')
    },
  ],
})

export default router

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { getActivities } from '@/api/activities'
import { useGameStore } from '@/stores/game'
import type { ActivitySummary } from '@/types/activity'
import { useResponsive } from '@/composables/useResponsive'

const gameStore = useGameStore()
const { isMobile } = useResponsive()

const activeTopTab = ref<'mine' | 'lobby' | 'contacts'>('mine')
const activeMineTab = ref<'created' | 'joined'>('created')

const activities = ref<ActivitySummary[]>([])
const loading = ref(false)
const errorText = ref<string | null>(null)

const setMineTab = (tab: 'created' | 'joined') => {
  activeMineTab.value = tab
  loadMyActivities()
}

const loadMyActivities = () => {
  loading.value = true
  errorText.value = null
  return getActivities(gameStore.currentGameKey, { type: activeMineTab.value, page: 1, pageSize: 10, sort: 'desc' })
    .then((res) => {
      // 兼容文档示例：data 可能直接是数组，也可能是 {data,hasMore}
      const payload = res.data.data as any
      activities.value = Array.isArray(payload) ? payload : (payload.data || [])
      loading.value = false
    })
    .catch((err) => {
      loading.value = false
      errorText.value = err && (err as any).message ? (err as any).message : '加载失败'
    })
}

onMounted(() => {
  gameStore.loadGames().catch(() => {})
  gameStore.loadManifest(gameStore.currentGameKey).catch(() => {})
  loadMyActivities()
})
</script>

<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-4">
    <div class="mb-3 flex items-center justify-between gap-3">
      <div class="text-lg font-semibold">Meeting Stone</div>
      <div class="flex items-center gap-2">
        <Button variant="outline" size="sm" @click="loadMyActivities">刷新</Button>
        <Button size="sm">创建</Button>
      </div>
    </div>

    <!-- 顶部 Tab（移动端横向，桌面端也可用） -->
    <div class="mb-3 flex gap-2">
      <Button
        variant="outline"
        :class="activeTopTab === 'mine' ? 'border-primary' : ''"
        @click="activeTopTab = 'mine'"
      >
        我的活动
      </Button>
      <Button
        variant="outline"
        :class="activeTopTab === 'lobby' ? 'border-primary' : ''"
        @click="activeTopTab = 'lobby'"
      >
        大厅
      </Button>
      <Button
        variant="outline"
        :class="activeTopTab === 'contacts' ? 'border-primary' : ''"
        @click="activeTopTab = 'contacts'"
      >
        联系人
      </Button>
    </div>

    <div v-if="activeTopTab === 'mine'">
      <!-- 二级 Tab -->
      <div class="mb-3 flex gap-2">
        <Button
          variant="outline"
          size="sm"
          :class="activeMineTab === 'created' ? 'border-primary' : ''"
          @click="setMineTab('created')"
        >
          我创建的
        </Button>
        <Button variant="outline" size="sm" :disabled="true">模板（待接入）</Button>
        <Button
          variant="outline"
          size="sm"
          :class="activeMineTab === 'joined' ? 'border-primary' : ''"
          @click="setMineTab('joined')"
        >
          我参与的
        </Button>
      </div>

      <div v-if="loading" class="text-sm text-muted-foreground">加载中...</div>
      <div v-else-if="errorText" class="text-sm text-destructive">{{ errorText }}</div>
      <div v-else class="grid gap-3" :class="isMobile ? '' : 'md:grid-cols-2'">
        <div
          v-for="item in activities"
          :key="item.id"
          class="rounded-lg border bg-card p-4"
        >
          <div class="mb-1 flex items-center justify-between gap-3">
            <div class="truncate font-medium">{{ item.title }}</div>
            <div class="text-xs text-muted-foreground">#{{ item.id }}</div>
          </div>
          <div class="text-sm text-muted-foreground">
            {{ item.size }}人 · 截止 {{ item.deadline }}
          </div>
        </div>
        <div v-if="activities.length === 0" class="text-sm text-muted-foreground">暂无数据</div>
      </div>
    </div>

    <div v-else class="text-sm text-muted-foreground">该模块后续按 UX 指南继续补齐。</div>
  </div>
</template>


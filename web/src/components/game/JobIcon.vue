<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'

interface Props {
  gameKey: string
  jobKey: string
  size?: number
  showLabel?: boolean
}

const props = defineProps<Props>()
const gameStore = useGameStore()

const job = computed(() => {
  const manifest = gameStore.manifests[props.gameKey]
  const jobs = manifest && manifest.jobs ? manifest.jobs : []
  return jobs.find((j) => j.key === props.jobKey) || null
})

const label = computed(() => (job.value ? job.value.label : props.jobKey))
const color = computed(() => (job.value && job.value.color ? job.value.color : undefined))
const sizePx = computed(() => `${props.size || 18}px`)
</script>

<template>
  <span class="inline-flex items-center gap-1">
    <span
      class="inline-block rounded-sm bg-muted"
      :style="{ width: sizePx, height: sizePx, backgroundColor: color || undefined }"
      aria-hidden="true"
    />
    <span v-if="showLabel" class="text-sm" :style="{ color: color || undefined }">{{ label }}</span>
  </span>
</template>


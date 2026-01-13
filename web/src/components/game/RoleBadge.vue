<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  roleKey: string
  label?: string
}

const props = defineProps<Props>()

const roleClass = computed(() => {
  // 轻量映射：后续可从 manifest.roles 注入颜色/图标
  if (props.roleKey === 'tank') return 'bg-blue-500/20 text-blue-200 border-blue-500/30'
  if (props.roleKey === 'healer') return 'bg-green-500/20 text-green-200 border-green-500/30'
  return 'bg-red-500/20 text-red-200 border-red-500/30'
})
</script>

<template>
  <span
    class="inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium"
    :class="roleClass"
  >
    {{ label || roleKey }}
  </span>
</template>


<script setup lang="ts">
import { computed } from 'vue'
import { Input } from '@/components/ui/input'
import { useGameStore } from '@/stores/game'

interface Props {
  gameKey: string
  modelValue: string
  placeholder?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const gameStore = useGameStore()

const enabled = computed(() => {
  const manifest = gameStore.manifests[props.gameKey]
  return !!(manifest && manifest.features && manifest.features.serverRegion)
})
</script>

<template>
  <div v-if="enabled">
    <Input
      :model-value="modelValue"
      :placeholder="placeholder || '输入服务器...'"
      @update:model-value="emit('update:modelValue', String($event))"
    />
  </div>
</template>


<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'
import { Button } from '@/components/ui/button'

interface Props {
  gameKey: string
  jobKey: string
  modelValue: string
}

const props = defineProps<Props>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const gameStore = useGameStore()

const specs = computed(() => {
  const manifest = gameStore.manifests[props.gameKey]
  const all = manifest && manifest.specs ? manifest.specs : []
  return all.filter((s) => s.jobKey === props.jobKey)
})

const onPick = (specKey: string) => {
  emit('update:modelValue', specKey)
}
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <Button
      v-for="spec in specs"
      :key="spec.key"
      type="button"
      variant="outline"
      :class="spec.key === modelValue ? 'border-primary' : ''"
      @click="onPick(spec.key)"
    >
      {{ spec.label }}
    </Button>
  </div>
</template>


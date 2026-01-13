import { computed } from 'vue'
import { useGameStore } from '@/stores/game'

export function useGameManifest(gameKey: string) {
  const gameStore = useGameStore()

  const manifest = computed(() => gameStore.manifests[gameKey] || null)

  const loadManifest = () => {
    if (manifest.value) {
      return Promise.resolve(manifest.value)
    }
    return gameStore.loadManifest(gameKey)
  }

  return { manifest, loadManifest }
}


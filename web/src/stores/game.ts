import { defineStore } from 'pinia'
import type { GameInfo, GameManifest } from '@/types/game'
import { getGameManifest, getGames } from '@/api/games'

interface GameState {
  currentGameKey: string
  games: GameInfo[]
  manifests: Record<string, GameManifest>
  loading: boolean
}

export const useGameStore = defineStore('game', {
  state: (): GameState => ({
    currentGameKey: 'wow-classic',
    games: [],
    manifests: {},
    loading: false
  }),
  getters: {
    currentManifest: (state) => state.manifests[state.currentGameKey] || null
  },
  actions: {
    setCurrentGameKey(gameKey: string) {
      this.currentGameKey = gameKey
    },
    loadGames() {
      this.loading = true
      return getGames()
        .then((res) => {
          this.games = res.data.data || []
          this.loading = false
          return this.games
        })
        .catch((err) => {
          this.loading = false
          return Promise.reject(err)
        })
    },
    loadManifest(gameKey: string) {
      if (this.manifests[gameKey]) {
        return Promise.resolve(this.manifests[gameKey])
      }
      return getGameManifest(gameKey)
        .then((res) => {
          this.manifests[gameKey] = res.data.data
          return this.manifests[gameKey]
        })
        .catch((err) => Promise.reject(err))
    }
  }
})


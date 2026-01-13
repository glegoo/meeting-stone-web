import type { ApiSuccessResponse } from '@/types/api'
import type { GameInfo, GameManifest } from '@/types/game'
import { apiClient } from './client'

export function getGames() {
  return apiClient.get<ApiSuccessResponse<GameInfo[]>>('/games')
}

export function getGameManifest(gameKey: string) {
  return apiClient.get<ApiSuccessResponse<GameManifest>>(`/games/${encodeURIComponent(gameKey)}/manifest`)
}


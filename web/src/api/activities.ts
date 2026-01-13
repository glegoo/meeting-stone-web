import type { ApiSuccessResponse } from '@/types/api'
import type { ActivitySummary, ListResponse, LobbySearchParams } from '@/types/activity'
import { apiClient } from './client'

export interface GetActivitiesParams {
  type: 'created' | 'joined'
  page?: number
  pageSize?: number
  sort?: 'asc' | 'desc'
  includeExpire?: boolean
}

export function getActivities(gameKey: string, params: GetActivitiesParams) {
  return apiClient.get<ApiSuccessResponse<ListResponse<ActivitySummary>>>(
    `/games/${encodeURIComponent(gameKey)}/activities`,
    { params }
  )
}

export function getLobbyActivities(gameKey: string, params: LobbySearchParams) {
  return apiClient.get<ApiSuccessResponse<ListResponse<ActivitySummary>>>(
    `/games/${encodeURIComponent(gameKey)}/lobby/activities`,
    { params }
  )
}


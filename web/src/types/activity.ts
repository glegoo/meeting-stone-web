export interface ActivityOwner {
  id: string
  name: string
}

export interface ActivitySummary {
  id: number
  gameKey: string
  title: string
  deadline: string
  size: number
  isPublic: boolean
  accessMode?: 'public' | 'private' | 'code' | 'invite'
  owner?: ActivityOwner
}

export interface ListResponse<T> {
  data: T[]
  hasMore?: boolean
}

export interface LobbySearchParams {
  server?: string
  region?: string
  faction?: number
  date?: string
  size?: number
  keyword?: string
  page?: number
  pageSize?: number
}


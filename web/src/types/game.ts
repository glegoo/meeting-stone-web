export interface GameInfo {
  gameKey: string
  name: string
  gameFamily?: string
}

export interface GameManifest {
  gameKey: string
  gameFamily?: string
  manifestVersion?: number
  roles?: Array<{ key: string; label: string }>
  jobs?: Array<{ key: string; label: string; color?: string }>
  specs?: Array<{ jobKey: string; key: string; label: string }>
  features?: {
    serverRegion?: boolean
    faction?: boolean
    buffComposition?: boolean
    scoreProvider?: string
  }
  partySizes?: number[]
  rosterShape?: { type: 'grid'; cols: number }
  buffComposition?: unknown
}


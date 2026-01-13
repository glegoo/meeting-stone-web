import axios from 'axios'
import type { AxiosError, AxiosInstance } from 'axios'

interface ApiClientOptions {
  /**
   * 仅用于 dev 环境直连后端（避免 CORS），生产环境推荐同域部署（/api/v1）
   * 示例：VITE_API_ORIGIN=http://localhost:8080
   */
  apiOrigin?: string
}

let accessToken: string | null = null

export function setAccessToken(token: string | null) {
  accessToken = token
}

export function createApiClient(options: ApiClientOptions = {}): AxiosInstance {
  const baseURL = options.apiOrigin ? `${options.apiOrigin}/api/v1` : '/api/v1'

  const client = axios.create({
    baseURL,
    timeout: 20_000
  })

  client.interceptors.request.use((config) => {
    if (accessToken) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${accessToken}`
    }
    return config
  })

  client.interceptors.response.use(
    (res) => res,
    (error: AxiosError) => {
      // 这里先做最小封装：把错误往外抛，页面层/Store 再统一 toast
      // 后续在 auth store 里补 refreshToken 流程（HttpOnly cookie）
      return Promise.reject(error)
    }
  )

  return client
}

export const apiClient = createApiClient({
  apiOrigin: import.meta.env.VITE_API_ORIGIN
})


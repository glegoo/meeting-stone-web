export interface ApiSuccessResponse<T> {
  code: number
  message?: string
  data: T
}

export interface ApiErrorResponse {
  code: number
  message?: string
  error?: string
}


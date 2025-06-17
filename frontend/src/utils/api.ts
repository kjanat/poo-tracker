export const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:3002'

export const createAuthHeaders = (token: string): Record<string, string> => ({
  Authorization: `Bearer ${token}`,
  'Content-Type': 'application/json'
})

export const createFormDataHeaders = (token: string): Record<string, string> => ({
  Authorization: `Bearer ${token}`
})

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export const handleApiResponse = async <T>(response: Response): Promise<T> => {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new ApiError(response.status, errorData.error ?? 'Request failed')
  }
  return (await response.json()) as T
}

export interface ApiConfig {
  baseUrl: string
  timeout?: number
}

export interface AuthProvider {
  getToken: () => string | null
}

export class ApiError extends Error {
  constructor (public status: number, message: string, public data?: unknown) {
    super(message)
    this.name = 'ApiError'
  }
}

export interface ApiResponse<T> {
  data: T
  status: number
  headers: Headers
}

export class ApiClient {
  constructor (
    private readonly config: ApiConfig,
    private readonly authProvider: AuthProvider
  ) {}

  private createHeaders (includeAuth = true): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json'
    }

    if (includeAuth) {
      const token = this.authProvider.getToken()
      if (token != null) {
        headers.Authorization = `Bearer ${token}`
      }
    }

    return headers
  }

  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }))
      throw new ApiError(response.status, errorData.error ?? 'Request failed', errorData)
    }

    const data = await response.json() as T
    return {
      data,
      status: response.status,
      headers: response.headers
    }
  }

  async get<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.config.baseUrl}${endpoint}`, {
      method: 'GET',
      headers: this.createHeaders(),
      ...options
    })

    return await this.handleResponse<T>(response)
  }

  async post<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.config.baseUrl}${endpoint}`, {
      method: 'POST',
      headers: this.createHeaders(),
      body: body != null ? JSON.stringify(body) : undefined,
      ...options
    })

    return await this.handleResponse<T>(response)
  }

  async put<T>(endpoint: string, body?: unknown, options?: RequestInit): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.config.baseUrl}${endpoint}`, {
      method: 'PUT',
      headers: this.createHeaders(),
      body: body != null ? JSON.stringify(body) : undefined,
      ...options
    })

    return await this.handleResponse<T>(response)
  }

  async delete<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
    const response = await fetch(`${this.config.baseUrl}${endpoint}`, {
      method: 'DELETE',
      headers: this.createHeaders(),
      ...options
    })

    return await this.handleResponse<T>(response)
  }

  async uploadFile<T>(endpoint: string, file: File, fieldName = 'image'): Promise<ApiResponse<T>> {
    const formData = new FormData()
    formData.append(fieldName, file)

    const response = await fetch(`${this.config.baseUrl}${endpoint}`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${this.authProvider.getToken() ?? ''}`
      },
      body: formData
    })

    return await this.handleResponse<T>(response)
  }
}

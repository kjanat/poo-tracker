import type { ApiClient } from '../../core/api/ApiClient'
import type {
  Entry,
  CreateEntryRequest,
  UpdateEntryRequest,
  EntryListResponse,
  EntryFilters
} from './types'

export class EntryService {
  constructor(private readonly apiClient: ApiClient) {}

  async getEntries(filters: EntryFilters = {}): Promise<EntryListResponse> {
    const params = new URLSearchParams()

    if (filters.page != null) params.append('page', filters.page.toString())
    if (filters.limit != null) params.append('limit', filters.limit.toString())
    if (filters.sortBy != null) params.append('sortBy', filters.sortBy)
    if (filters.sortOrder != null) params.append('sortOrder', filters.sortOrder)
    if (filters.bristolType != null) params.append('bristolType', filters.bristolType.toString())
    if (filters.dateFrom != null) params.append('dateFrom', filters.dateFrom)
    if (filters.dateTo != null) params.append('dateTo', filters.dateTo)

    const queryString = params.toString()
    const endpoint = queryString !== '' ? `/api/entries?${queryString}` : '/api/entries'

    const response = await this.apiClient.get<EntryListResponse>(endpoint)
    return response.data
  }

  async getEntry(id: string): Promise<Entry> {
    const response = await this.apiClient.get<Entry>(`/api/entries/${id}`)
    return response.data
  }

  async createEntry(entry: CreateEntryRequest): Promise<Entry> {
    const response = await this.apiClient.post<Entry>('/api/entries', entry)
    return response.data
  }

  async updateEntry(id: string, entry: UpdateEntryRequest): Promise<Entry> {
    const response = await this.apiClient.put<Entry>(`/api/entries/${id}`, entry)
    return response.data
  }

  async deleteEntry(id: string): Promise<void> {
    await this.apiClient.delete(`/api/entries/${id}`)
  }

  async uploadEntryPhoto(file: File): Promise<{ url: string }> {
    const response = await this.apiClient.uploadFile<{ url: string }>('/api/uploads', file)
    return response.data
  }
}

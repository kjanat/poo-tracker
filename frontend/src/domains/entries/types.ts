export interface Entry {
  id: string
  bristolType: number
  volume?: string
  color?: string
  consistency?: string
  notes?: string
  createdAt: string
  satisfaction?: number
  pain?: number
  strain?: number
  floaters?: boolean
  smell?: string
  photoUrl?: string
}

export interface CreateEntryRequest {
  bristolType: number
  volume?: string
  color?: string
  consistency?: string
  notes?: string
  satisfaction?: number
  pain?: number
  strain?: number
  floaters?: boolean
  smell?: string
  photoUrl?: string
}

export type UpdateEntryRequest = Partial<CreateEntryRequest>

export interface EntryListResponse {
  entries: Entry[]
  pagination: {
    page: number
    limit: number
    total: number
    pages: number
  }
}

export interface EntryFilters {
  page?: number
  limit?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  bristolType?: number
  dateFrom?: string
  dateTo?: string
}

export interface CreateEntryData {
  bristolType: number
  volume?: string
  color?: string
  notes?: string
  photo?: File
}

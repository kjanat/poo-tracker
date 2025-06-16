export interface Entry {
  id: string
  bristolType: number
  volume?: string
  color?: string
  consistency?: string
  floaters: boolean
  pain?: number
  strain?: number
  satisfaction?: number
  notes?: string
  smell?: string
  photoUrl?: string
  createdAt: Date
  updatedAt: Date
  userId: string
}

export interface CreateEntryRequest {
  bristolType: number
  volume?: string | undefined
  color?: string | undefined
  consistency?: string | undefined
  floaters?: boolean | undefined
  pain?: number | undefined
  strain?: number | undefined
  satisfaction?: number | undefined
  notes?: string | undefined
  smell?: string | undefined
  photoUrl?: string | undefined
}

export interface UpdateEntryRequest extends Partial<CreateEntryRequest> {}

export interface EntryFilters {
  page?: number
  limit?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  bristolType?: number
  dateFrom?: Date
  dateTo?: Date
}

export interface EntryListResponse {
  entries: Entry[]
  pagination: {
    page: number
    limit: number
    total: number
    pages: number
  }
}

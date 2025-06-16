export interface Entry {
  id: string
  bristolType: number
  volume?: string | null
  color?: string | null
  consistency?: string | null
  floaters: boolean
  pain?: number | null
  strain?: number | null
  satisfaction?: number | null
  notes?: string | null
  smell?: string | null
  photoUrl?: string | null
  createdAt: Date
  updatedAt: Date
  userId: string
}

export interface CreateEntryRequest {
  bristolType: number
  volume?: string | null | undefined
  color?: string | null | undefined
  consistency?: string | null | undefined
  floaters?: boolean | undefined
  pain?: number | null | undefined
  strain?: number | null | undefined
  satisfaction?: number | null | undefined
  notes?: string | null | undefined
  smell?: string | null | undefined
  photoUrl?: string | null | undefined
}

export type UpdateEntryRequest = Partial<CreateEntryRequest>

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

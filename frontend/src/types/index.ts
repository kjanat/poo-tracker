export interface Meal {
  id: string
  name: string
  category?: string
  description?: string
  cuisine?: string
  spicyLevel?: number
  fiberRich: boolean
  dairy: boolean
  gluten: boolean
  notes?: string
  photoUrl?: string
  mealTime: string
  createdAt: string
  linkedEntries?: Entry[]
}

export interface Entry {
  id: string
  bristolType: number
  volume?: string
  color?: string
  consistency?: string
  notes?: string
  createdAt: string
}

export interface MealFormData {
  name: string
  category: string
  description: string
  cuisine: string
  spicyLevel: number
  fiberRich: boolean
  dairy: boolean
  gluten: boolean
  notes: string
  photoUrl?: string
}

export interface EntryResponse {
  id: string
  bristolType: number
  volume?: string
  color?: string
  notes?: string
  createdAt: string
  userId: string
}

export interface EntriesApiResponse {
  entries: EntryResponse[]
  pagination: {
    page: number
    limit: number
    total: number
    pages: number
  }
}

export interface AnalyticsSummary {
  totalEntries: number
  bristolDistribution: Array<{ type: number, count: number }>
  recentEntries: Array<{
    id: string
    bristolType: number
    createdAt: string
    satisfaction?: number
  }>
  averageSatisfaction?: number
}

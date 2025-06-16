export interface Meal {
  id: string
  name: string
  description?: string | null
  mealTime: Date
  category?: string | null
  cuisine?: string | null
  spicyLevel?: number | null
  fiberRich: boolean
  dairy: boolean
  gluten: boolean
  notes?: string | null
  photoUrl?: string | null
  createdAt: Date
  updatedAt: Date
  userId: string
}

export interface CreateMealRequest {
  name: string
  description?: string | null | undefined
  mealTime: Date
  category?: string | null | undefined
  cuisine?: string | null | undefined
  spicyLevel?: number | null | undefined
  fiberRich?: boolean | undefined
  dairy?: boolean | undefined
  gluten?: boolean | undefined
  notes?: string | null | undefined
  photoUrl?: string | null | undefined
}

export type UpdateMealRequest = Partial<CreateMealRequest>

export interface MealFilters {
  page?: number
  limit?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  category?: string
  dateFrom?: Date
  dateTo?: Date
  fiberRich?: boolean
  dairy?: boolean
  gluten?: boolean
}

export interface MealListResponse {
  meals: Meal[]
  pagination: {
    page: number
    limit: number
    total: number
    pages: number
  }
}

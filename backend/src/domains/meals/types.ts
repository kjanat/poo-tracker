export interface Meal {
  id: string
  name: string
  description?: string
  mealTime: Date
  category?: string
  cuisine?: string
  spicyLevel?: number
  fiberRich: boolean
  dairy: boolean
  gluten: boolean
  notes?: string
  photoUrl?: string
  createdAt: Date
  updatedAt: Date
  userId: string
}

export interface CreateMealRequest {
  name: string
  description?: string | undefined
  mealTime: Date
  category?: string | undefined
  cuisine?: string | undefined
  spicyLevel?: number | undefined
  fiberRich?: boolean | undefined
  dairy?: boolean | undefined
  gluten?: boolean | undefined
  notes?: string | undefined
  photoUrl?: string | undefined
}

export interface UpdateMealRequest extends Partial<CreateMealRequest> {}

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

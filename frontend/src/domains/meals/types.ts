export interface Meal {
  id: string
  name: string
  category?: 'BREAKFAST' | 'LUNCH' | 'DINNER' | 'SNACK'
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
  linkedEntries?: Array<{ id: string; bristolType: number; createdAt: string }>
}

export interface CreateMealRequest {
  name: string
  category: 'BREAKFAST' | 'LUNCH' | 'DINNER' | 'SNACK'
  description?: string
  cuisine?: string
  spicyLevel: number
  fiberRich: boolean
  dairy: boolean
  gluten: boolean
  notes?: string
  photoUrl?: string
}

export type UpdateMealRequest = Partial<CreateMealRequest>

export interface MealListResponse {
  meals: Meal[]
  pagination?: {
    page: number
    limit: number
    total: number
    pages: number
  }
}

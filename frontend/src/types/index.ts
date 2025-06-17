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
  linkedBowelMovements?: BowelMovement[]
  linkedEntries?: Array<{ id: string; bristolType: number; createdAt: string }>
}

export interface BowelMovement {
  id: string
  bristolType: number
  recordedAt: string
  volume?: 'SMALL' | 'MEDIUM' | 'LARGE' | 'MASSIVE'
  color?: 'BROWN' | 'DARK_BROWN' | 'LIGHT_BROWN' | 'YELLOW' | 'GREEN' | 'RED' | 'BLACK'
  consistency?: 'SOLID' | 'SOFT' | 'LOOSE' | 'WATERY'
  floaters: boolean
  pain: number
  strain: number
  satisfaction: number
  smell?: 'NONE' | 'MILD' | 'MODERATE' | 'STRONG' | 'TOXIC'
  photoUrl?: string
  notes?: string
  createdAt: string
  updatedAt: string
}

// Legacy interface for backward compatibility
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
  category: 'BREAKFAST' | 'LUNCH' | 'DINNER' | 'SNACK' | ''
  description: string
  cuisine: string
  spicyLevel: number
  fiberRich: boolean
  dairy: boolean
  gluten: boolean
  notes: string
  photoUrl: string
}

export interface BowelMovementResponse {
  id: string
  bristolType: number
  recordedAt: string
  volume?: 'SMALL' | 'MEDIUM' | 'LARGE' | 'MASSIVE'
  color?: 'BROWN' | 'DARK_BROWN' | 'LIGHT_BROWN' | 'YELLOW' | 'GREEN' | 'RED' | 'BLACK'
  consistency?: 'SOLID' | 'SOFT' | 'LOOSE' | 'WATERY'
  floaters: boolean
  pain: number
  strain: number
  satisfaction: number
  smell?: 'NONE' | 'MILD' | 'MODERATE' | 'STRONG' | 'TOXIC'
  photoUrl?: string
  notes?: string
  createdAt: string
  updatedAt: string
  userId: string
}

// Legacy for backward compatibility
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
  bristolDistribution: Array<{ type: number; count: number }>
  recentEntries: Array<{
    id: string
    bristolType: number
    createdAt: string
    satisfaction?: number
  }>
  averageSatisfaction?: number
}

export interface Symptom {
  id: string
  userId: string
  bowelMovementId?: string | null
  createdAt: string
  recordedAt: string
  type:
    | 'BLOATING'
    | 'CRAMPS'
    | 'NAUSEA'
    | 'HEARTBURN'
    | 'CONSTIPATION'
    | 'DIARRHEA'
    | 'GAS'
    | 'FATIGUE'
    | 'OTHER'
  severity: number
  notes?: string | null
}

export interface Medication {
  id: string
  userId: string
  createdAt: string
  updatedAt: string
  name: string
  dosage?: string | null
  frequency?: string | null
  startDate: string
  endDate?: string | null
  notes?: string | null
}

export interface UserSettings {
  id: string
  userId: string
  timezone: string
  reminderEnabled: boolean
  reminderTime: string
  dataRetentionDays: number
  createdAt: string
  updatedAt: string
}

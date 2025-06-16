export interface BowelMovement {
  id: string
  userId: string
  createdAt: Date
  updatedAt: Date
  recordedAt: Date
  bristolType: number
  volume?: 'SMALL' | 'MEDIUM' | 'LARGE' | 'MASSIVE' | null
  color?: 'BROWN' | 'DARK_BROWN' | 'LIGHT_BROWN' | 'YELLOW' | 'GREEN' | 'RED' | 'BLACK' | null
  consistency?: 'SOLID' | 'SOFT' | 'LOOSE' | 'WATERY' | null
  floaters: boolean
  pain: number
  strain: number
  satisfaction: number
  photoUrl?: string | null
  smell?: 'NONE' | 'MILD' | 'MODERATE' | 'STRONG' | 'TOXIC' | null
}

export interface CreateBowelMovementRequest {
  bristolType: number
  recordedAt?: Date | undefined
  volume?: 'SMALL' | 'MEDIUM' | 'LARGE' | 'MASSIVE' | undefined
  color?: 'BROWN' | 'DARK_BROWN' | 'LIGHT_BROWN' | 'YELLOW' | 'GREEN' | 'RED' | 'BLACK' | undefined
  consistency?: 'SOLID' | 'SOFT' | 'LOOSE' | 'WATERY' | undefined
  floaters?: boolean | undefined
  pain?: number | undefined
  strain?: number | undefined
  satisfaction?: number | undefined
  photoUrl?: string | undefined
  smell?: 'NONE' | 'MILD' | 'MODERATE' | 'STRONG' | 'TOXIC' | undefined
  notes?: string | undefined
}

export type UpdateBowelMovementRequest = Partial<CreateBowelMovementRequest>

export interface BowelMovementFilters {
  page?: number
  limit?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  bristolType?: number
  dateFrom?: Date
  dateTo?: Date
}

export interface BowelMovementListResponse {
  bowelMovements: BowelMovement[]
  pagination: {
    page: number
    limit: number
    total: number
    pages: number
  }
}

export interface BowelMovementDetails {
  id: string
  bowelMovementId: string
  notes?: string | null
  aiAnalysis?: Record<string, unknown> | null
}

// Additional new models from the schema
export interface Symptom {
  id: string
  userId: string
  bowelMovementId?: string | null
  createdAt: Date
  recordedAt: Date
  type: 'BLOATING' | 'CRAMPS' | 'NAUSEA' | 'HEARTBURN' | 'CONSTIPATION' | 'DIARRHEA' | 'GAS' | 'FATIGUE' | 'OTHER'
  severity: number
  notes?: string | null
}

export interface Medication {
  id: string
  userId: string
  createdAt: Date
  updatedAt: Date
  name: string
  dosage?: string | null
  frequency?: string | null
  startDate: Date
  endDate?: Date | null
  notes?: string | null
}

export interface UserSettings {
  id: string
  userId: string
  timezone: string
  reminderEnabled: boolean
  reminderTime: string
  dataRetentionDays: number
  createdAt: Date
  updatedAt: Date
}

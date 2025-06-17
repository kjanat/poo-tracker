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

export interface TrendData {
  date: string
  bristolType: number
  satisfaction?: number
  count: number
}

export interface HealthMetrics {
  averageBristolType: number
  consistencyScore: number
  healthTrend: 'improving' | 'stable' | 'declining'
  recommendationsNeeded: boolean
}

import { useState, useEffect, useCallback } from 'react'
import { container } from '../../../core/services'
import type { AnalyticsSummary, TrendData, HealthMetrics } from '../types'

export interface UseAnalyticsResult {
  summary: AnalyticsSummary | null
  trendData: TrendData[]
  healthMetrics: HealthMetrics | null
  isLoading: boolean
  error: string | null
  refreshData: () => Promise<void>
}

export function useAnalytics(): UseAnalyticsResult {
  const [summary, setSummary] = useState<AnalyticsSummary | null>(null)
  const [trendData, setTrendData] = useState<TrendData[]>([])
  const [healthMetrics, setHealthMetrics] = useState<HealthMetrics | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const analyticsService = container.get('analyticsService')

  const loadData = useCallback(async (): Promise<void> => {
    try {
      setIsLoading(true)
      setError(null)

      // Load all analytics data in parallel
      const [summaryData, trendsData, healthData] = await Promise.all([
        analyticsService.getAnalyticsSummary(),
        analyticsService.getTrendData(),
        analyticsService.getHealthMetrics()
      ])

      setSummary(summaryData)
      setTrendData(trendsData)
      setHealthMetrics(healthData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load analytics data')
    } finally {
      setIsLoading(false)
    }
  }, [analyticsService])

  useEffect(() => {
    loadData()
  }, [loadData])

  return {
    summary,
    trendData,
    healthMetrics,
    isLoading,
    error,
    refreshData: loadData
  }
}

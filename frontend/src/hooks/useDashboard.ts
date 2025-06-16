import { useState, useEffect } from 'react'
import { useAuthStore } from '../stores/authStore'
import { API_BASE_URL, createAuthHeaders, handleApiResponse } from '../utils/api'
import { getThisWeekCount } from '../utils/date'
import { getAverageBristolType } from '../utils/bristol'
import type { AnalyticsSummary, EntryResponse, EntriesApiResponse } from '../types'

export interface UseDashboardReturn {
  analytics: AnalyticsSummary | null
  recentEntries: EntryResponse[]
  isLoading: boolean
  error: string
  thisWeekCount: number
  averageBristolType: number
  refreshData: () => Promise<void>
}

export function useDashboard (): UseDashboardReturn {
  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null)
  const [recentEntries, setRecentEntries] = useState<EntryResponse[]>([])
  const [isLoading, setIsLoading] = useState<boolean>(true)
  const [error, setError] = useState<string>('')
  
  const { token } = useAuthStore()

  const fetchAnalytics = async (): Promise<void> => {
    if (token == null) {
      setError('Authentication token missing')
      setIsLoading(false)
      return
    }
    
    try {
      const response = await fetch(`${API_BASE_URL}/api/analytics/summary`, {
        headers: createAuthHeaders(token)
      })

      const data = await handleApiResponse<AnalyticsSummary>(response)
      setAnalytics(data)
    } catch (error) {
      console.error('Failed to fetch analytics:', error)
      setError(error instanceof Error ? error.message : 'Failed to fetch analytics')
    }
  }

  const fetchRecentEntries = async (): Promise<void> => {
    if (token == null) return

    try {
      const response = await fetch(`${API_BASE_URL}/api/entries?limit=20`, {
        headers: createAuthHeaders(token)
      })

      const data = await handleApiResponse<EntriesApiResponse>(response)
      setRecentEntries(data.entries)
    } catch (error) {
      console.error('Failed to fetch recent entries:', error)
      setError(error instanceof Error ? error.message : 'Failed to fetch recent entries')
    }
  }

  const refreshData = async (): Promise<void> => {
    setIsLoading(true)
    setError('')
    
    try {
      await Promise.all([
        fetchAnalytics(),
        fetchRecentEntries()
      ])
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    void refreshData()
  }, [token])

  // Computed values
  const thisWeekCount = getThisWeekCount(recentEntries)
  const averageBristolType = analytics != null 
    ? getAverageBristolType(analytics.bristolDistribution)
    : 0

  return {
    analytics,
    recentEntries,
    isLoading,
    error,
    thisWeekCount,
    averageBristolType,
    refreshData
  }
}

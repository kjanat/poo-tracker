import type { ApiClient } from '../../core/api/ApiClient'
import type { AnalyticsSummary, TrendData, HealthMetrics } from './types'
import type { Entry } from '../entries/types'
import { BristolAnalyzer } from '../bristol/BristolAnalyzer'

export class AnalyticsService {
  constructor (private readonly apiClient: ApiClient) {}

  async getSummary(): Promise<AnalyticsSummary> {
    const response = await this.apiClient.get<AnalyticsSummary>('/api/analytics/summary')
    return response.data
  }

  async getAnalyticsSummary(): Promise<AnalyticsSummary> {
    return this.getSummary()
  }

  async getTrendData(): Promise<TrendData[]> {
    // For now, we'll calculate trend data from recent entries
    // In a real app, this might be a separate API endpoint
    const entriesResponse = await this.apiClient.get<{ entries: Entry[] }>('/api/entries?limit=100')
    return this.generateTrendData(entriesResponse.data.entries)
  }

  async getHealthMetrics(): Promise<HealthMetrics> {
    // Calculate health metrics from recent entries
    const entriesResponse = await this.apiClient.get<{ entries: Entry[] }>('/api/entries?limit=50')
    return this.calculateHealthMetrics(entriesResponse.data.entries)
  }

  calculateThisWeekCount(entries: Entry[]): number {
    const oneWeekAgo = new Date()
    oneWeekAgo.setDate(oneWeekAgo.getDate() - 7)

    return entries.filter((entry) => new Date(entry.createdAt) >= oneWeekAgo).length
  }

  generateTrendData(entries: Entry[]): TrendData[] {
    const groupedByDate = entries.reduce((acc, entry) => {
      const date = new Date(entry.createdAt).toISOString().split('T')[0]
      
      if (!date) {
        return acc
      }
      
      if (!(date in acc)) {
        acc[date] = {
          date,
          bristolTypes: [],
          satisfactions: [],
          count: 0
        }
      }
      
      const dateGroup = acc[date]
      if (dateGroup) {
        dateGroup.bristolTypes.push(entry.bristolType)
        if (entry.satisfaction != null) {
          dateGroup.satisfactions.push(entry.satisfaction)
        }
        dateGroup.count++
      }
      
      return acc
    }, {} as Record<string, { date: string, bristolTypes: number[], satisfactions: number[], count: number }>)

    return Object.values(groupedByDate).map(group => ({
      date: group.date,
      bristolType: group.bristolTypes.reduce((sum, type) => sum + type, 0) / group.bristolTypes.length,
      satisfaction: group.satisfactions.length > 0 
        ? group.satisfactions.reduce((sum, sat) => sum + sat, 0) / group.satisfactions.length 
        : 5, // Default satisfaction if none provided
      count: group.count
    })).sort((a, b) => a.date.localeCompare(b.date))
  }

  calculateHealthMetrics(entries: Entry[]): HealthMetrics {
    if (entries.length === 0) {
      return {
        averageBristolType: 0,
        consistencyScore: 0,
        healthTrend: 'stable',
        recommendationsNeeded: false
      }
    }

    const recentEntries = entries.slice(0, 10)
    const olderEntries = entries.slice(10, 20)

    const averageBristolType = recentEntries.reduce((sum, entry) => sum + entry.bristolType, 0) / recentEntries.length
    
    // Calculate consistency score (how close to ideal range 3-4)
    const idealCount = recentEntries.filter(entry => entry.bristolType === 3 || entry.bristolType === 4).length
    const consistencyScore = (idealCount / recentEntries.length) * 100

    // Calculate trend
    let healthTrend: 'improving' | 'stable' | 'declining' = 'stable'
    if (olderEntries.length > 0) {
      const recentAvg = averageBristolType
      const olderAvg = olderEntries.reduce((sum, entry) => sum + entry.bristolType, 0) / olderEntries.length
      
      const difference = Math.abs(recentAvg - 3.5) - Math.abs(olderAvg - 3.5)
      if (difference < -0.3) healthTrend = 'improving'
      else if (difference > 0.3) healthTrend = 'declining'
    }

    // Check if recommendations needed
    const extremeCount = recentEntries.filter(entry => entry.bristolType === 1 || entry.bristolType === 7).length
    const recommendationsNeeded = extremeCount > recentEntries.length * 0.3 || consistencyScore < 50

    return {
      averageBristolType: Number(averageBristolType.toFixed(1)),
      consistencyScore: Number(consistencyScore.toFixed(1)),
      healthTrend,
      recommendationsNeeded
    }
  }

  generateRecommendations(metrics: HealthMetrics, entries: Entry[]): string[] {
    const recommendations: string[] = []

    if (metrics.averageBristolType < 2.5) {
      recommendations.push('Consider increasing fiber intake and hydration')
      recommendations.push('Regular exercise can help with constipation')
    } else if (metrics.averageBristolType > 5.5) {
      recommendations.push('Consider reducing dairy or identifying trigger foods')
      recommendations.push('Stay hydrated and consider probiotics')
    }

    if (metrics.consistencyScore < 30) {
      recommendations.push('Track food intake to identify patterns')
      recommendations.push('Consider consulting a healthcare provider')
    }

    const hasExtremes = entries.some(entry => BristolAnalyzer.needsAttention(entry.bristolType))
    if (hasExtremes) {
      recommendations.push('Monitor extreme types and consider medical advice if persistent')
    }

    return recommendations
  }
}

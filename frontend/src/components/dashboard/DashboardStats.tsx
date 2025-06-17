import { type ReactElement } from 'react'
import type { AnalyticsSummary } from '../../types'

interface DashboardStatsProps {
  analytics: AnalyticsSummary | null
  thisWeekCount: number
  averageBristolType: number
}

export function DashboardStats({
  analytics,
  thisWeekCount,
  averageBristolType
}: DashboardStatsProps): ReactElement {
  if (analytics == null) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="bg-white p-6 rounded-lg shadow-md animate-pulse">
            <div className="h-4 bg-gray-200 rounded mb-2"></div>
            <div className="h-8 bg-gray-200 rounded"></div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h3 className="text-sm font-medium text-gray-500 mb-2">Total Entries</h3>
        <p className="text-3xl font-bold text-amber-600">{analytics.totalEntries}</p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <h3 className="text-sm font-medium text-gray-500 mb-2">This Week</h3>
        <p className="text-3xl font-bold text-green-600">{thisWeekCount}</p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <h3 className="text-sm font-medium text-gray-500 mb-2">Average Type</h3>
        <p className="text-3xl font-bold text-blue-600">{averageBristolType}</p>
      </div>

      <div className="bg-white p-6 rounded-lg shadow-md">
        <h3 className="text-sm font-medium text-gray-500 mb-2">Avg Satisfaction</h3>
        <p className="text-3xl font-bold text-purple-600">
          {analytics.averageSatisfaction != null ? analytics.averageSatisfaction.toFixed(1) : 'N/A'}
        </p>
      </div>
    </div>
  )
}

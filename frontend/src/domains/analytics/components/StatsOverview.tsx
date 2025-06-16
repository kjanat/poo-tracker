import React from 'react'
import type { AnalyticsSummary } from '../types'

interface StatsOverviewProps {
  summary: AnalyticsSummary
}

export function StatsOverview({ summary }: StatsOverviewProps): JSX.Element {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-2xl font-bold text-blue-600">
          {summary.totalEntries}
        </h3>
        <p className="text-gray-600 mt-2">Total Entries</p>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-2xl font-bold text-green-600">
          {summary.averageSatisfaction 
            ? `${summary.averageSatisfaction.toFixed(1)}/10`
            : 'N/A'
          }
        </h3>
        <p className="text-gray-600 mt-2">Average Satisfaction</p>
      </div>

      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h3 className="text-2xl font-bold text-purple-600">
          {summary.recentEntries.length}
        </h3>
        <p className="text-gray-600 mt-2">Recent Entries</p>
      </div>
    </div>
  )
}

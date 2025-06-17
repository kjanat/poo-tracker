import type { HealthMetrics } from '../types'

interface HealthInsightsProps {
  healthMetrics: HealthMetrics
  recommendations: string[]
}

export function HealthInsights({ healthMetrics, recommendations }: HealthInsightsProps) {
  const getTrendIcon = (trend: HealthMetrics['healthTrend']) => {
    switch (trend) {
      case 'improving':
        return '📈'
      case 'declining':
        return '📉'
      default:
        return '➡️'
    }
  }

  const getTrendColor = (trend: HealthMetrics['healthTrend']) => {
    switch (trend) {
      case 'improving':
        return 'text-green-600'
      case 'declining':
        return 'text-red-600'
      default:
        return 'text-blue-600'
    }
  }

  const getConsistencyColor = (score: number) => {
    if (score >= 70) return 'text-green-600'
    if (score >= 40) return 'text-yellow-600'
    return 'text-red-600'
  }

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <h3 className="text-lg font-semibold mb-4">Health Insights</h3>

      <div className="space-y-4">
        <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
          <span className="text-sm font-medium">Average Bristol Type:</span>
          <span className="text-lg font-semibold">
            {healthMetrics.averageBristolType.toFixed(1)}
          </span>
        </div>

        <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
          <span className="text-sm font-medium">Consistency Score:</span>
          <span
            className={`text-lg font-semibold ${getConsistencyColor(healthMetrics.consistencyScore)}`}
          >
            {healthMetrics.consistencyScore.toFixed(1)}%
          </span>
        </div>

        <div className="flex items-center justify-between p-3 bg-gray-50 rounded">
          <span className="text-sm font-medium">Health Trend:</span>
          <span className={`text-lg font-semibold ${getTrendColor(healthMetrics.healthTrend)}`}>
            {getTrendIcon(healthMetrics.healthTrend)} {healthMetrics.healthTrend}
          </span>
        </div>

        {recommendations.length > 0 && (
          <div className="mt-6">
            <h4 className="text-md font-semibold mb-3">Recommendations</h4>
            <ul className="space-y-2">
              {recommendations.map((rec, index) => (
                <li key={index} className="flex items-start">
                  <span className="text-blue-500 mr-2">•</span>
                  <span className="text-sm text-gray-700">{rec}</span>
                </li>
              ))}
            </ul>
          </div>
        )}

        {healthMetrics.recommendationsNeeded && (
          <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded">
            <p className="text-sm text-yellow-800">
              <strong>Note:</strong> Consider tracking patterns more closely or consulting a
              healthcare provider.
            </p>
          </div>
        )}
      </div>
    </div>
  )
}

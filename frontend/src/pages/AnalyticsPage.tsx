
import { useAuthStore } from '../stores/authStore'
import Logo from '../components/Logo'
import { StatsOverview } from '../domains/analytics/components/StatsOverview'
import { BristolDistributionChart } from '../domains/analytics/components/BristolDistributionChart'
import { HealthInsights } from '../domains/analytics/components/HealthInsights'
import { RecentEntriesTable } from '../domains/analytics/components/RecentEntriesTable'
import { useAnalytics } from '../domains/analytics/hooks/useAnalytics'
import { container } from '../core/services'
import type { AnalyticsService } from '../domains/analytics/AnalyticsService'

export function AnalyticsPage() {
  const { token } = useAuthStore()
  const { summary, healthMetrics, isLoading, error } = useAnalytics()

  if (!token) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">Please log in to access this page.</p>
        </div>
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <div className="flex items-center justify-between mb-8">
            <div className="flex items-center space-x-4">
              <Logo size="40" />
              <h1 className="text-3xl font-bold mb-8 flex items-center gap-2">
                📊 Analytics Dashboard
              </h1>
            </div>
          </div>
          <div className="text-center py-8">
            <p className="text-gray-500">Loading analytics data...</p>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <div className="flex items-center justify-between mb-8">
            <div className="flex items-center space-x-4">
              <Logo size="40" />
              <h1 className="text-3xl font-bold mb-8 flex items-center gap-2">
                📊 Analytics Dashboard
              </h1>
            </div>
          </div>
          <div className="bg-red-50 border border-red-200 rounded-md p-4">
            <p className="text-red-600">{error}</p>
          </div>
        </div>
      </div>
    )
  }

  if (!summary) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <div className="flex items-center justify-between mb-8">
            <div className="flex items-center space-x-4">
              <Logo size="40" />
              <h1 className="text-3xl font-bold mb-8 flex items-center gap-2">
                📊 Analytics Dashboard
              </h1>
            </div>
          </div>
          <div className="text-center py-8">
            <p className="text-gray-500">No analytics data available.</p>
          </div>
        </div>
      </div>
    )
  }

  // Generate recommendations if we have health metrics
  const recommendations = healthMetrics 
    ? container.get<AnalyticsService>('analyticsService').generateRecommendations(healthMetrics, summary.recentEntries)
    : []

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center space-x-4">
            <Logo size="40" />
            <h1 className="text-3xl font-bold flex items-center gap-2">
              📊 Analytics Dashboard
            </h1>
          </div>
        </div>

        {/* Stats Overview */}
        <StatsOverview summary={summary} />

        {/* Charts and Insights */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          <BristolDistributionChart summary={summary} />
          
          {healthMetrics && (
            <HealthInsights 
              healthMetrics={healthMetrics} 
              recommendations={recommendations}
            />
          )}
        </div>

        {/* Recent Entries Table */}
        <RecentEntriesTable summary={summary} />
      </div>
    </div>
  )
}

import React, { type ReactElement } from 'react'
import Logo from '../components/Logo'
import { DashboardStats } from '../components/dashboard/DashboardStats'
import { BristolChart } from '../components/dashboard/BristolChart'
import { RecentEntries } from '../components/dashboard/RecentEntries'
import { useDashboard } from '../hooks/useDashboard'

export function DashboardPage (): ReactElement {
  const {
    analytics,
    recentEntries,
    isLoading,
    error,
    thisWeekCount,
    averageBristolType,
    refreshData
  } = useDashboard()

  return (
    <div className="min-h-screen bg-amber-50">
      <nav className="bg-white shadow-sm border-b border-amber-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Logo />
            <div className="flex items-center gap-4">
              <h1 className="text-2xl font-bold text-amber-800">Dashboard</h1>
              <button
                onClick={() => void refreshData()}
                disabled={isLoading}
                className="text-amber-600 hover:text-amber-700 disabled:opacity-50"
              >
                {isLoading ? 'Loading...' : 'Refresh'}
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {error !== '' && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {error}
          </div>
        )}

        <DashboardStats 
          analytics={analytics}
          thisWeekCount={thisWeekCount}
          averageBristolType={averageBristolType}
        />

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <BristolChart 
            bristolDistribution={analytics?.bristolDistribution ?? []}
          />
          
          <RecentEntries 
            entries={recentEntries}
            loading={isLoading}
          />
        </div>

        {analytics?.totalEntries === 0 && !isLoading && (
          <div className="text-center py-12">
            <div className="text-gray-500 mb-4">
              <svg className="mx-auto h-24 w-24 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">No data yet</h3>
            <p className="text-gray-500 mb-6">Start tracking your bowel movements to see insights here.</p>
            <a
              href="/new-entry"
              className="bg-amber-600 text-white px-6 py-3 rounded-md hover:bg-amber-700 transition-colors"
            >
              Add Your First Entry
            </a>
          </div>
        )}
      </main>
    </div>
  )
}

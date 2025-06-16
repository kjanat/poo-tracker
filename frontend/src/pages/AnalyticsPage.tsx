import { useState, useEffect } from 'react'
import { useAuthStore } from '../stores/authStore'
import Logo from '../components/Logo'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
  Legend
} from 'recharts'

interface AnalyticsData {
  totalEntries: number
  bristolDistribution: Array<{
    type: number
    count: number
  }>
  recentEntries: Array<{
    id: string
    bristolType: number
    createdAt: string
    satisfaction?: number
  }>
  averageSatisfaction: number | null
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3002'

const bristolTypeDescriptions = {
  1: 'Hard lumps',
  2: 'Lumpy sausage',
  3: 'Cracked sausage',
  4: 'Smooth sausage',
  5: 'Soft blobs',
  6: 'Fluffy pieces',
  7: 'Watery'
}

const bristolColors = [
  '#8B5CF6',
  '#EC4899',
  '#F59E0B',
  '#10B981',
  '#3B82F6',
  '#EF4444',
  '#6B7280'
]

export function AnalyticsPage () {
  const { token } = useAuthStore()
  const [analyticsData, setAnalyticsData] = useState<AnalyticsData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>('')

  useEffect(() => {
    const fetchAnalytics = async () => {
      if (!token) return

      try {
        const response = await fetch(`${API_BASE_URL}/api/analytics/summary`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })

        if (!response.ok) {
          throw new Error('Failed to fetch analytics data')
        }

        const data = await response.json()
        setAnalyticsData(data)
      } catch (err) {
        console.error('Error fetching analytics:', err)
        setError(
          err instanceof Error ? err.message : 'Failed to load analytics'
        )
      } finally {
        setLoading(false)
      }
    }

    fetchAnalytics()
  }, [token])

  // Prepare chart data
  const bristolChartData =
    ((analyticsData?.bristolDistribution.map((item) => ({
      type: `Type ${item.type}`,
      description:
        bristolTypeDescriptions[
          item.type as keyof typeof bristolTypeDescriptions
        ],
      count: item.count,
      percentage: Math.round((item.count / analyticsData.totalEntries) * 100)
    }))) != null) || []

  // Get frequency data from recent entries (group by day)
  const frequencyData =
    ((analyticsData?.recentEntries.reduce(
      (acc, entry) => {
        const date = new Date(entry.createdAt).toLocaleDateString()
        const existing = acc.find((item) => item.date === date)
        if (existing != null) {
          existing.count += 1
        } else {
          acc.push({ date, count: 1 })
        }
        return acc
      },
      [] as Array<{ date: string, count: number }>
    )) != null) || []

  if (loading) {
    return (
      <div className='max-w-6xl mx-auto'>
        <h1 className='text-3xl font-bold mb-8 flex items-center gap-2'>
          <Logo size={32} /> Analytics
        </h1>
        <div className='card text-center py-8'>
          <p className='text-gray-600'>Loading analytics data...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className='max-w-6xl mx-auto'>
        <h1 className='text-3xl font-bold mb-8 flex items-center gap-2'>
          <Logo size={32} /> Analytics
        </h1>
        <div className='card bg-red-50 border-red-200'>
          <p className='text-red-600'>{error}</p>
        </div>
      </div>
    )
  }

  if ((analyticsData == null) || analyticsData.totalEntries === 0) {
    return (
      <div className='max-w-6xl mx-auto'>
        <h1 className='text-3xl font-bold mb-8 flex items-center gap-2'>
          <Logo size={32} /> Analytics
        </h1>
        <div className='card text-center py-8'>
          <p className='text-gray-600 flex items-center justify-center gap-2'>
            No data to analyze yet. Start logging your entries!{' '}
            <Logo size={24} />
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className='max-w-6xl mx-auto'>
      <h1 className='text-3xl font-bold mb-8 flex items-center gap-2'>
        <Logo size={32} /> Analytics
      </h1>

      {/* Summary Stats */}
      <div className='grid grid-cols-1 md:grid-cols-3 gap-6 mb-8'>
        <div className='card text-center'>
          <h3 className='text-2xl font-bold text-blue-600'>
            {analyticsData.totalEntries}
          </h3>
          <p className='text-gray-600'>Total Entries</p>
        </div>

        <div className='card text-center'>
          <h3 className='text-2xl font-bold text-green-600'>
            {analyticsData.bristolDistribution.find((d) => d.type === 4)
              ?.count || 0}
          </h3>
          <p className='text-gray-600'>Ideal Entries (Type 4)</p>
        </div>

        <div className='card text-center'>
          <h3 className='text-2xl font-bold text-purple-600'>
            {analyticsData.averageSatisfaction
              ? `${analyticsData.averageSatisfaction.toFixed(1)}/10`
              : 'N/A'}
          </h3>
          <p className='text-gray-600'>Avg Satisfaction</p>
        </div>
      </div>

      <div className='grid lg:grid-cols-2 gap-6 mb-8'>
        {/* Bristol Chart Distribution */}
        <div className='card'>
          <h3 className='text-lg font-semibold mb-4'>
            Bristol Chart Distribution
          </h3>
          {bristolChartData.length > 0
            ? (
              <ResponsiveContainer width='100%' height={300}>
                <BarChart data={bristolChartData}>
                  <CartesianGrid strokeDasharray='3 3' />
                  <XAxis dataKey='type' />
                  <YAxis />
                  <Tooltip
                    formatter={(value) => [value, 'Count']}
                    labelFormatter={(label) => {
                      const item = bristolChartData.find((d) => d.type === label)
                      return (item != null) ? `${label}: ${item.description}` : label
                    }}
                  />
                  <Bar dataKey='count' fill='#3B82F6' />
                </BarChart>
              </ResponsiveContainer>
              )
            : (
              <p className='text-gray-600 text-center py-8'>
                No distribution data available
              </p>
              )}
        </div>

        {/* Bristol Type Pie Chart */}
        <div className='card'>
          <h3 className='text-lg font-semibold mb-4'>Type Distribution</h3>
          {bristolChartData.length > 0
            ? (
              <ResponsiveContainer width='100%' height={300}>
                <PieChart>
                  <Pie
                    data={bristolChartData}
                    cx='50%'
                    cy='50%'
                    labelLine={false}
                    label={({ percentage }) => `${percentage}%`}
                    outerRadius={80}
                    fill='#8884d8'
                    dataKey='count'
                  >
                    {bristolChartData.map((_, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={bristolColors[index % bristolColors.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => [value, 'Count']} />
                  <Legend
                    formatter={(value) => {
                      const item = bristolChartData.find((d) => d.type === value)
                      return (item != null) ? `${value}: ${item.description}` : value
                    }}
                  />
                </PieChart>
              </ResponsiveContainer>
              )
            : (
              <p className='text-gray-600 text-center py-8'>
                No distribution data available
              </p>
              )}
        </div>
      </div>

      <div className='grid lg:grid-cols-2 gap-6 mb-8'>
        {/* Recent Activity */}
        <div className='card'>
          <h3 className='text-lg font-semibold mb-4'>Recent Activity</h3>
          {frequencyData.length > 0
            ? (
              <ResponsiveContainer width='100%' height={250}>
                <LineChart data={frequencyData}>
                  <CartesianGrid strokeDasharray='3 3' />
                  <XAxis dataKey='date' />
                  <YAxis />
                  <Tooltip />
                  <Line
                    type='monotone'
                    dataKey='count'
                    stroke='#10B981'
                    strokeWidth={2}
                  />
                </LineChart>
              </ResponsiveContainer>
              )
            : (
              <p className='text-gray-600 text-center py-8'>
                No recent activity data
              </p>
              )}
        </div>

        {/* Health Insights */}
        <div className='card'>
          <h3 className='text-lg font-semibold mb-4'>Health Insights</h3>
          <div className='space-y-4'>
            {/* Bristol Type 4 Percentage */}
            <div>
              <div className='flex justify-between mb-1'>
                <span className='text-sm font-medium'>
                  Ideal Movements (Type 4)
                </span>
                <span className='text-sm text-gray-600'>
                  {bristolChartData.find((d) => d.type === 'Type 4')
                    ?.percentage || 0}
                  %
                </span>
              </div>
              <div className='w-full bg-gray-200 rounded-full h-2'>
                <div
                  className='bg-green-600 h-2 rounded-full'
                  style={{
                    width: `${bristolChartData.find((d) => d.type === 'Type 4')?.percentage || 0}%`
                  }}
                />
              </div>
            </div>

            {/* Constipation Indicator */}
            <div>
              <div className='flex justify-between mb-1'>
                <span className='text-sm font-medium'>
                  Constipation Risk (Types 1-2)
                </span>
                <span className='text-sm text-gray-600'>
                  {(bristolChartData.find((d) => d.type === 'Type 1')
                    ?.percentage || 0) +
                    (bristolChartData.find((d) => d.type === 'Type 2')
                      ?.percentage || 0)}
                  %
                </span>
              </div>
              <div className='w-full bg-gray-200 rounded-full h-2'>
                <div
                  className='bg-red-600 h-2 rounded-full'
                  style={{
                    width: `${
                      (bristolChartData.find((d) => d.type === 'Type 1')
                        ?.percentage || 0) +
                      (bristolChartData.find((d) => d.type === 'Type 2')
                        ?.percentage || 0)
                    }%`
                  }}
                />
              </div>
            </div>

            {/* Diarrhea Indicator */}
            <div>
              <div className='flex justify-between mb-1'>
                <span className='text-sm font-medium'>
                  Diarrhea Risk (Types 6-7)
                </span>
                <span className='text-sm text-gray-600'>
                  {(bristolChartData.find((d) => d.type === 'Type 6')
                    ?.percentage || 0) +
                    (bristolChartData.find((d) => d.type === 'Type 7')
                      ?.percentage || 0)}
                  %
                </span>
              </div>
              <div className='w-full bg-gray-200 rounded-full h-2'>
                <div
                  className='bg-orange-600 h-2 rounded-full'
                  style={{
                    width: `${
                      (bristolChartData.find((d) => d.type === 'Type 6')
                        ?.percentage || 0) +
                      (bristolChartData.find((d) => d.type === 'Type 7')
                        ?.percentage || 0)
                    }%`
                  }}
                />
              </div>
            </div>

            <div className='mt-4 p-3 bg-blue-50 rounded'>
              <h4 className='font-medium text-blue-800 mb-2'>Quick Tips</h4>
              <ul className='text-sm text-blue-700 space-y-1'>
                <li>• Aim for Types 3-4 for optimal digestive health</li>
                <li>• Stay hydrated and eat fiber-rich foods</li>
                <li>• Track patterns to identify triggers</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Entries Table */}
      <div className='card'>
        <h3 className='text-lg font-semibold mb-4'>Recent Entries</h3>
        {analyticsData.recentEntries.length > 0
          ? (
            <div className='overflow-x-auto'>
              <table className='min-w-full divide-y divide-gray-200'>
                <thead className='bg-gray-50'>
                  <tr>
                    <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                      Date
                    </th>
                    <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                      Bristol Type
                    </th>
                    <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                      Description
                    </th>
                    <th className='px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider'>
                      Satisfaction
                    </th>
                  </tr>
                </thead>
                <tbody className='bg-white divide-y divide-gray-200'>
                  {analyticsData.recentEntries.map((entry) => (
                    <tr key={entry.id}>
                      <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-900'>
                        {new Date(entry.createdAt).toLocaleDateString()}
                      </td>
                      <td className='px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900'>
                        Type {entry.bristolType}
                      </td>
                      <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                        {
                        bristolTypeDescriptions[
                          entry.bristolType as keyof typeof bristolTypeDescriptions
                        ]
                      }
                      </td>
                      <td className='px-6 py-4 whitespace-nowrap text-sm text-gray-500'>
                        {entry.satisfaction ? `${entry.satisfaction}/10` : 'N/A'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            )
          : (
            <p className='text-gray-600 text-center py-8'>
              No recent entries found
            </p>
            )}
      </div>
    </div>
  )
}

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'
import { BristolAnalyzer } from '../../bristol/BristolAnalyzer'
import type { AnalyticsSummary } from '../types'

interface BristolDistributionChartProps {
  summary: AnalyticsSummary
}

interface TooltipProps {
  active?: boolean
  payload?: Array<{
    payload: {
      type: number
      description: string
      count: number
    }
  }>
}

const CustomTooltip = ({ active, payload }: TooltipProps) => {
  if (active && payload && payload.length) {
    const data = payload[0]?.payload
    if (!data) return null

    return (
      <div className="bg-white p-3 border border-gray-200 rounded shadow">
        <p className="font-medium">{`Type ${data.type}`}</p>
        <p className="text-sm text-gray-600">{data.description}</p>
        <p className="text-blue-600">{`Count: ${data.count}`}</p>
      </div>
    )
  }
  return null
}

/**
 * Render a bar chart showing the distribution of Bristol stool types from the provided summary.
 *
 * @param summary - AnalyticsSummary whose `bristolDistribution` array is used to build the chart data
 * @returns The JSX element containing the rendered distribution chart
 */
export function BristolDistributionChart({ summary }: BristolDistributionChartProps) {
  const chartData = summary.bristolDistribution.map((item) => ({
    type: item.type,
    count: item.count,
    description: BristolAnalyzer.getDescription(item.type),
    color: BristolAnalyzer.getColor(item.type)
  }))

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
      <h3 className="text-lg font-semibold mb-4">Bristol Stool Distribution</h3>
      <div className="h-64">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="type"
              label={{ value: 'Bristol Type', position: 'insideBottom', offset: -5 }}
            />
            <YAxis label={{ value: 'Count', angle: -90, position: 'insideLeft' }} />
            <Tooltip content={<CustomTooltip />} />
            <Bar dataKey="count" fill="#8B5CF6" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}
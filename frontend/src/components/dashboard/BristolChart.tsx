import { type ReactElement } from 'react'
import { getBristolTypeDescription, getBristolTypeCategory } from '../../utils/bristol'

interface BristolChartProps {
  bristolDistribution: Array<{ type: number, count: number }>
}

export function BristolChart ({ bristolDistribution }: BristolChartProps): ReactElement {
  const maxCount = Math.max(...bristolDistribution.map(item => item.count), 1)
  
  const getBarColor = (type: number): string => {
    const category = getBristolTypeCategory(type)
    switch (category) {
      case 'Constipated': return 'bg-red-500'
      case 'Normal': return 'bg-green-500'
      case 'Loose': return 'bg-yellow-500'
      default: return 'bg-gray-500'
    }
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">Bristol Stool Chart Distribution</h3>
      
      {bristolDistribution.length === 0 ? (
        <div className="text-gray-500 text-center py-8">
          No data available yet
        </div>
      ) : (
        <div className="space-y-3">
          {bristolDistribution.map((item) => {
            const percentage = (item.count / maxCount) * 100
            return (
              <div key={item.type} className="flex items-center">
                <div className="w-20 text-sm font-medium text-gray-700">
                  Type {item.type}
                </div>
                <div className="flex-1 mx-3">
                  <div className="bg-gray-200 rounded-full h-4 relative">
                    <div
                      className={`${getBarColor(item.type)} h-4 rounded-full transition-all duration-300`}
                      style={{ width: `${percentage}%` }}
                    />
                  </div>
                </div>
                <div className="w-12 text-sm text-gray-600 text-right">
                  {item.count}
                </div>
                <div className="w-32 text-xs text-gray-500 ml-2">
                  {getBristolTypeDescription(item.type)}
                </div>
              </div>
            )
          })}
        </div>
      )}
      
      <div className="mt-6 text-xs text-gray-500">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-1">
            <div className="w-3 h-3 bg-red-500 rounded"></div>
            <span>Constipated (1-2)</span>
          </div>
          <div className="flex items-center gap-1">
            <div className="w-3 h-3 bg-green-500 rounded"></div>
            <span>Normal (3-4)</span>
          </div>
          <div className="flex items-center gap-1">
            <div className="w-3 h-3 bg-yellow-500 rounded"></div>
            <span>Loose (5-7)</span>
          </div>
        </div>
      </div>
    </div>
  )
}

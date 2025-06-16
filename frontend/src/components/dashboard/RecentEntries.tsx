import React, { type ReactElement } from 'react'
import { formatDate } from '../../utils/date'
import { getBristolTypeDescription } from '../../utils/bristol'
import type { EntryResponse } from '../../types'

interface RecentEntriesProps {
  entries: EntryResponse[]
  loading: boolean
}

export function RecentEntries ({ entries, loading }: RecentEntriesProps): ReactElement {
  if (loading) {
    return (
      <div className="bg-white p-6 rounded-lg shadow-md">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">Recent Entries</h3>
        <div className="space-y-3">
          {[...Array(5)].map((_, i) => (
            <div key={i} className="border-b border-gray-100 pb-3 animate-pulse">
              <div className="h-4 bg-gray-200 rounded mb-2"></div>
              <div className="h-3 bg-gray-200 rounded w-1/2"></div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">Recent Entries</h3>
      
      {entries.length === 0 ? (
        <div className="text-gray-500 text-center py-8">
          No entries yet. Add your first entry!
        </div>
      ) : (
        <div className="space-y-3">
          {entries.slice(0, 10).map((entry) => (
            <div key={entry.id} className="border-b border-gray-100 pb-3 last:border-b-0">
              <div className="flex justify-between items-start">
                <div>
                  <span className="font-medium text-gray-800">
                    Type {entry.bristolType}: {getBristolTypeDescription(entry.bristolType)}
                  </span>
                  {entry.volume != null && (
                    <span className="text-sm text-gray-600 ml-2">
                      Volume: {entry.volume}
                    </span>
                  )}
                  {entry.color != null && (
                    <span className="text-sm text-gray-600 ml-2">
                      Color: {entry.color}
                    </span>
                  )}
                </div>
                <span className="text-xs text-gray-500">
                  {formatDate(entry.createdAt)}
                </span>
              </div>
              
              {entry.notes != null && entry.notes !== '' && (
                <div className="text-sm text-gray-600 mt-1">
                  Notes: {entry.notes}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
      
      <div className="mt-4 text-center">
        <a
          href="/entries"
          className="text-amber-600 hover:text-amber-700 text-sm font-medium"
        >
          View all entries →
        </a>
      </div>
    </div>
  )
}

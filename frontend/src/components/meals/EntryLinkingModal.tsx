import React, { type ReactElement } from 'react'
import { formatDate } from '../../utils/date'
import { getBristolTypeDescription } from '../../utils/bristol'
import type { Meal, Entry } from '../../types'

interface EntryLinkingModalProps {
  show: boolean
  onClose: () => void
  meal: Meal | null
  availableEntries: Entry[]
  linkedEntries: Entry[]
  onLinkEntry: (entryId: string) => Promise<void>
  onUnlinkEntry: (entryId: string) => Promise<void>
}

export function EntryLinkingModal ({
  show,
  onClose,
  meal,
  availableEntries,
  linkedEntries,
  onLinkEntry,
  onUnlinkEntry
}: EntryLinkingModalProps): ReactElement | null {
  if (!show || meal == null) return null

  const linkedEntryIds = new Set(linkedEntries.map(entry => entry.id))
  const unlinkedEntries = availableEntries.filter(entry => !linkedEntryIds.has(entry.id))

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-semibold text-gray-800">
            Link Entries to {meal.name}
          </h3>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        {linkedEntries.length > 0 && (
          <div className="mb-6">
            <h4 className="font-medium text-gray-700 mb-2">Currently Linked Entries:</h4>
            <div className="space-y-2 max-h-40 overflow-y-auto">
              {linkedEntries.map((entry) => (
                <div
                  key={entry.id}
                  className="flex justify-between items-center p-2 bg-green-50 rounded border"
                >
                  <div className="text-sm">
                    <span className="font-medium">
                      Type {entry.bristolType}: {getBristolTypeDescription(entry.bristolType)}
                    </span>
                    <div className="text-gray-600">
                      {formatDate(entry.createdAt)}
                    </div>
                    {entry.notes != null && (
                      <div className="text-gray-600 text-xs">
                        Notes: {entry.notes}
                      </div>
                    )}
                  </div>
                  <button
                    onClick={async () => await onUnlinkEntry(entry.id)}
                    className="text-red-600 hover:text-red-800 text-sm px-2 py-1 border border-red-300 rounded hover:bg-red-50"
                  >
                    Unlink
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {unlinkedEntries.length > 0 && (
          <div>
            <h4 className="font-medium text-gray-700 mb-2">Available Entries to Link:</h4>
            <div className="space-y-2 max-h-60 overflow-y-auto">
              {unlinkedEntries.map((entry) => (
                <div
                  key={entry.id}
                  className="flex justify-between items-center p-2 bg-gray-50 rounded border"
                >
                  <div className="text-sm">
                    <span className="font-medium">
                      Type {entry.bristolType}: {getBristolTypeDescription(entry.bristolType)}
                    </span>
                    <div className="text-gray-600">
                      {formatDate(entry.createdAt)}
                    </div>
                    {entry.notes != null && (
                      <div className="text-gray-600 text-xs">
                        Notes: {entry.notes}
                      </div>
                    )}
                  </div>
                  <button
                    onClick={async () => await onLinkEntry(entry.id)}
                    className="text-green-600 hover:text-green-800 text-sm px-2 py-1 border border-green-300 rounded hover:bg-green-50"
                  >
                    Link
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {unlinkedEntries.length === 0 && linkedEntries.length === 0 && (
          <div className="text-gray-500 text-center py-4">
            No entries available to link.
          </div>
        )}

        <div className="mt-6 flex justify-end">
          <button
            onClick={onClose}
            className="bg-gray-500 text-white px-4 py-2 rounded-md hover:bg-gray-600"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}

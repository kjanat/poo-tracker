import React from 'react'
import { BristolAnalyzer } from '../../bristol/BristolAnalyzer'
import type { Entry } from '../types'

interface EntryListProps {
  entries: Entry[]
  onEdit: (entry: Entry) => void
  onDelete: (id: string) => void
  isSubmitting: boolean
}

const bristolAnalyzer = new BristolAnalyzer()

export function EntryList({ entries, onEdit, onDelete, isSubmitting }: EntryListProps): JSX.Element {
  if (entries.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-500">No entries yet. Create your first entry above!</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {entries.map((entry) => (
        <EntryCard
          key={entry.id}
          entry={entry}
          onEdit={onEdit}
          onDelete={onDelete}
          isSubmitting={isSubmitting}
        />
      ))}
    </div>
  )
}

interface EntryCardProps {
  entry: Entry
  onEdit: (entry: Entry) => void
  onDelete: (id: string) => void
  isSubmitting: boolean
}

function EntryCard({ entry, onEdit, onDelete, isSubmitting }: EntryCardProps): JSX.Element {
  const handleEdit = (): void => {
    onEdit(entry)
  }

  const handleDelete = (): void => {
    onDelete(entry.id)
  }

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
      <div className="flex justify-between items-start mb-3">
        <div className="flex-1">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium text-gray-900">
              Type {entry.bristolType}
            </h3>
            <span className="text-sm text-gray-500">
              {new Date(entry.createdAt).toLocaleDateString('en-US', {
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
              })}
            </span>
          </div>
          <p className="text-sm text-gray-600 mt-1">
            {bristolAnalyzer.getDescription(entry.bristolType)}
          </p>
        </div>
        <div className="flex space-x-2 ml-4">
          <button
            onClick={handleEdit}
            className="text-blue-600 hover:text-blue-800 text-sm font-medium"
            disabled={isSubmitting}
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className="text-red-600 hover:text-red-800 text-sm font-medium"
            disabled={isSubmitting}
          >
            Delete
          </button>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-3 text-sm">
        {entry.volume && (
          <div>
            <span className="font-medium">Volume:</span> {entry.volume}
          </div>
        )}
        {entry.color && (
          <div>
            <span className="font-medium">Color:</span> {entry.color}
          </div>
        )}
      </div>

      {entry.photoUrl && (
        <div className="mb-3">
          <img
            src={entry.photoUrl}
            alt="Entry photo"
            className="max-w-xs max-h-48 object-cover rounded border"
          />
        </div>
      )}

      {entry.notes && (
        <div className="mt-3 p-3 bg-gray-50 rounded">
          <p className="text-sm text-gray-700">{entry.notes}</p>
        </div>
      )}
    </div>
  )
}

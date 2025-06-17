import React from 'react'
import { BristolSelector } from './BristolSelector'
import { EntryDetailsForm } from './EntryDetailsForm'
import type { CreateEntryData } from '../types'

interface EntryFormProps {
  formData: CreateEntryData
  onUpdate: (updates: Partial<CreateEntryData>) => void
  onSubmit: () => void
  onCancel: (() => void) | null
  isSubmitting: boolean
  isEditing?: boolean
}

export function EntryForm({
  formData,
  onUpdate,
  onSubmit,
  onCancel,
  isSubmitting,
  isEditing = false
}: EntryFormProps) {
  const handleSubmit = (e: React.FormEvent): void => {
    e.preventDefault()
    onSubmit()
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <h2 className="text-xl font-semibold mb-6">{isEditing ? 'Edit Entry' : 'New Entry'}</h2>

        <BristolSelector
          selectedType={formData.bristolType}
          onTypeSelect={(type) => onUpdate({ bristolType: type })}
          disabled={isSubmitting}
        />

        <EntryDetailsForm
          volume={formData.volume}
          color={formData.color}
          notes={formData.notes}
          photo={formData.photo}
          onVolumeChange={(volume) => onUpdate({ volume })}
          onColorChange={(color) => onUpdate({ color })}
          onNotesChange={(notes) => onUpdate({ notes })}
          onPhotoChange={(photo) => onUpdate({ photo })}
          disabled={isSubmitting}
        />

        <div className="flex justify-end space-x-3 mt-6">
          {isEditing && onCancel && (
            <button
              type="button"
              onClick={onCancel}
              className="px-4 py-2 text-gray-600 hover:text-gray-800 font-medium"
              disabled={isSubmitting}
            >
              Cancel
            </button>
          )}
          <button
            type="submit"
            className="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed font-medium"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Saving...' : isEditing ? 'Update Entry' : 'Save Entry'}
          </button>
        </div>
      </div>
    </form>
  )
}

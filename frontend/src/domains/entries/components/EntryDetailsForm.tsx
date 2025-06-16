import React, { type ReactElement } from 'react'

interface EntryDetailsFormProps {
  volume: string
  color: string
  notes: string
  photo?: File
  onVolumeChange: (volume: string) => void
  onColorChange: (color: string) => void
  onNotesChange: (notes: string) => void
  onPhotoChange: (photo: File | undefined) => void
  disabled?: boolean
}

export function EntryDetailsForm({
  volume,
  color,
  notes,
  photo,
  onVolumeChange,
  onColorChange,
  onNotesChange,
  onPhotoChange,
  disabled = false
}: EntryDetailsFormProps): ReactElement {
  const handlePhotoChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const file = e.target.files?.[0]
    onPhotoChange(file)
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Volume</label>
          <select
            value={volume}
            onChange={(e) => onVolumeChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={disabled}
          >
            <option value="">Select volume</option>
            <option value="small">Small</option>
            <option value="medium">Medium</option>
            <option value="large">Large</option>
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Color</label>
          <select
            value={color}
            onChange={(e) => onColorChange(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={disabled}
          >
            <option value="">Select color</option>
            <option value="brown">Brown</option>
            <option value="light-brown">Light Brown</option>
            <option value="dark-brown">Dark Brown</option>
            <option value="yellow">Yellow</option>
            <option value="green">Green</option>
            <option value="black">Black</option>
            <option value="red">Red</option>
          </select>
        </div>
      </div>

      {/* Photo Upload */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Photo (optional)</label>
        <input
          type="file"
          accept="image/*"
          onChange={handlePhotoChange}
          className="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          disabled={disabled}
        />
        {photo && (
          <div className="mt-2">
            <p className="text-sm text-gray-600">Selected: {photo.name}</p>
          </div>
        )}
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
        <textarea
          value={notes}
          onChange={(e) => onNotesChange(e.target.value)}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="Any additional observations..."
          disabled={disabled}
        />
      </div>
    </div>
  )
}

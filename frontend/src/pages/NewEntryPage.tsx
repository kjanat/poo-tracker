import React, { useState, useEffect } from 'react'
import { useAuthStore } from '../stores/authStore'
import Logo from '../components/Logo'

interface StoolEntry {
    bristolType: number;
    volume?: string;
    color?: string;
    notes?: string;
    photoUrl?: string;
}

interface EntryResponse {
    id: string;
    bristolType: number;
    volume?: string;
    color?: string;
    notes?: string;
    photoUrl?: string;
    createdAt: string;
    userId: string;
}

interface EntriesApiResponse {
    entries: EntryResponse[];
  pagination: {
      page: number;
      limit: number;
      total: number;
      pages: number;
  };
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3002'

const getBristolTypeDescription = (type: number): string => {
  const descriptions = {
    1: 'Hard lumps (Severe constipation)',
    2: 'Lumpy sausage (Mild constipation)',
    3: 'Cracked sausage (Normal)',
    4: 'Smooth sausage (Ideal)',
    5: 'Soft blobs (Lacking fiber)',
    6: 'Fluffy pieces (Mild diarrhea)',
    7: 'Watery (Severe diarrhea)'
  }
  return descriptions[type as keyof typeof descriptions] || 'Unknown'
}

export function NewEntryPage(): JSX.Element {
  const { token } = useAuthStore()

  const [formData, setFormData] = useState<StoolEntry>({
    bristolType: 0,
    volume: '',
    color: '',
    notes: ''
  })

  const [selectedImage, setSelectedImage] = useState<File | null>(null)
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [submitStatus, setSubmitStatus] = useState<
      'idle' | 'success' | 'error'
  >('idle')
  const [errorMessage, setErrorMessage] = useState<string>('')
  const [entries, setEntries] = useState<EntryResponse[]>([])
  const [isLoadingEntries, setIsLoadingEntries] = useState(true)
  const [editingEntry, setEditingEntry] = useState<EntryResponse | null>(null)

  const fetchEntries = async () => {
    try {
      const response = await fetch(
        `${API_BASE_URL}/api/entries?limit=10&sortOrder=desc`,
        {
          headers: {
            Authorization: `Bearer ${token}`
          }
        }
      )

      if (!response.ok) {
        throw new Error('Failed to fetch entries')
      }

      const data: EntriesApiResponse = await response.json()
      setEntries(data.entries)
    } catch (error) {
      console.error('Error fetching entries:', error)
    } finally {
      setIsLoadingEntries(false)
    }
  }

  useEffect(() => {
    if (token) {
      fetchEntries()
    }
  }, [token])

  const handleInputChange = (
    field: keyof StoolEntry,
    value: string | number
  ) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value
    }))
  }

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
      if (file) {
      // Validate file type
      if (!file.type.startsWith('image/')) {
        setErrorMessage('Please select a valid image file')
        setSubmitStatus('error')
        return
      }

      // Validate file size (5MB limit)
      if (file.size > 5 * 1024 * 1024) {
        setErrorMessage('Image size must be less than 5MB')
        setSubmitStatus('error')
        return
      }

      setSelectedImage(file)

      // Create preview
      const reader = new FileReader()
      reader.onload = (e) => {
        setImagePreview(e.target?.result as string)
      }
      reader.readAsDataURL(file)

      // Clear any previous errors
      setSubmitStatus('idle')
      setErrorMessage('')
    }
  }

  const removeImage = () => {
    setSelectedImage(null)
    setImagePreview(null)
  }

  const uploadImage = async (): Promise<string | null> => {
      if (!selectedImage) return null

    const uploadFormData = new FormData()
    uploadFormData.append('image', selectedImage)

    try {
      const response = await fetch(`${API_BASE_URL}/api/uploads`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`
        },
        body: uploadFormData
      })

      if (!response.ok) {
        throw new Error('Failed to upload image')
      }

      const data = await response.json()
      return data.imageUrl
    } catch (error) {
      console.error('Error uploading image:', error)
      throw error
    }
  }

  const startEdit = (entry: EntryResponse) => {
    setEditingEntry(entry)
    setFormData({
      bristolType: entry.bristolType,
      volume: entry.volume || '',
      color: entry.color || '',
      notes: entry.notes || '',
      photoUrl: entry.photoUrl
    })
    // Reset image states when editing (user can upload new image if desired)
    setSelectedImage(null)
    setImagePreview(null)
    // Scroll to form
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const cancelEdit = () => {
    setEditingEntry(null)
    setFormData({
      bristolType: 0,
      volume: '',
      color: '',
      notes: ''
    })
    setSelectedImage(null)
    setImagePreview(null)
    setSubmitStatus('idle')
    setErrorMessage('')
  }

  const deleteEntry = async (entryId: string) => {
    if (!confirm('Are you sure you want to delete this entry?')) return

    try {
      const response = await fetch(`${API_BASE_URL}/api/entries/${entryId}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`
        }
      })

      if (!response.ok) {
        throw new Error('Failed to delete entry')
      }

      // Remove from local state
      setEntries((prev) => prev.filter((entry) => entry.id !== entryId))
    } catch (error) {
      console.error('Error deleting entry:', error)
      setErrorMessage(
        error instanceof Error ? error.message : 'Failed to delete entry'
      )
      setSubmitStatus('error')
    }
  }

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault()

    if (formData.bristolType === 0) {
      setErrorMessage('Please select a Bristol Type')
      setSubmitStatus('error')
      return
    }

    setIsSubmitting(true)
    setSubmitStatus('idle')
    setErrorMessage('')

    try {
      // Upload image if selected
      let photoUrl = formData.photoUrl
        if (selectedImage) {
        const uploadedUrl = await uploadImage()
        if (uploadedUrl) {
          photoUrl = uploadedUrl
        }
      }

      const submitData = {
        bristolType: formData.bristolType,
        volume: formData.volume || undefined,
        color: formData.color || undefined,
        notes: formData.notes || undefined,
        photoUrl: photoUrl || undefined
      }

      const isEditing = editingEntry !== null
      const url = isEditing
        ? `${API_BASE_URL}/api/entries/${editingEntry.id}`
        : `${API_BASE_URL}/api/entries`
      const method = isEditing ? 'PUT' : 'POST'

      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify(submitData)
      })

      if (!response.ok) {
        const error = await response.json()
        throw new Error(
          error.error || `Failed to ${isEditing ? 'update' : 'save'} entry`
        )
      }

      const savedEntry: EntryResponse = await response.json()
      console.log(
        `Entry ${isEditing ? 'updated' : 'saved'} successfully:`,
        savedEntry
      )

      setSubmitStatus('success')

      if (isEditing) {
        // Update the entry in the list
        setEntries((prev) =>
          prev.map((entry) =>
            entry.id === editingEntry.id ? savedEntry : entry
          )
        )
        setEditingEntry(null)
      } else {
        // Add new entry to the list
        setEntries((prev) => [savedEntry, ...prev])
      }

      // Reset form
      setFormData({
        bristolType: 0,
        volume: '',
        color: '',
        notes: ''
      })
      setSelectedImage(null)
      setImagePreview(null)
    } catch (error) {
      console.error(
          `Error ${editingEntry ? 'updating' : 'saving'} entry:`,
        error
      )
      setErrorMessage(
        error instanceof Error ? error.message : 'An error occurred'
      )
      setSubmitStatus('error')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
      <div className="max-w-2xl mx-auto">
          <h1 className="text-3xl font-bold mb-8">
              {editingEntry ? 'Edit Entry' : 'Log New Entry'}
      </h1>

          <div className="card">
              <p className="text-center text-gray-600 mb-8 flex items-center justify-center gap-2">
                  {editingEntry
            ? 'Update your masterpiece!'
            : 'Time to document another masterpiece!'}{' '}
          <Logo size={24} />
        </p>

        {submitStatus === 'success' && (
                  <div className="mb-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
                      Entry {editingEntry ? 'updated' : 'saved'} successfully! Keep
            tracking your progress.
          </div>
        )}

        {submitStatus === 'error' && (
                  <div className="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {errorMessage}
          </div>
        )}

              <form onSubmit={handleSubmit} className="space-y-6">
          <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
              Bristol Stool Chart Type (1-7) *
            </label>
            <select
                          className="input-field"
              value={formData.bristolType}
              onChange={(e) =>
                  handleInputChange('bristolType', parseInt(e.target.value) || 0)
              }
              required
            >
                          <option value="">Select Bristol Type</option>
                          <option value="1">
                Type 1 - Hard lumps (Severe constipation)
              </option>
                          <option value="2">
                Type 2 - Lumpy sausage (Mild constipation)
              </option>
                          <option value="3">Type 3 - Cracked sausage (Normal)</option>
                          <option value="4">Type 4 - Smooth sausage (Ideal)</option>
                          <option value="5">Type 5 - Soft blobs (Lacking fiber)</option>
                          <option value="6">Type 6 - Fluffy pieces (Mild diarrhea)</option>
                          <option value="7">Type 7 - Watery (Severe diarrhea)</option>
            </select>
          </div>

                  <div className="grid md:grid-cols-2 gap-4">
            <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                Volume
              </label>
              <select
                              className="input-field"
                value={formData.volume}
                onChange={(e) => handleInputChange('volume', e.target.value)}
              >
                              <option value="">Select volume</option>
                              <option value="Small">Small</option>
                              <option value="Medium">Medium</option>
                              <option value="Large">Large</option>
                              <option value="Massive">Massive</option>
              </select>
            </div>

            <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                Color
              </label>
              <select
                              className="input-field"
                value={formData.color}
                onChange={(e) => handleInputChange('color', e.target.value)}
              >
                              <option value="">Select color</option>
                              <option value="Brown">Brown</option>
                              <option value="Dark Brown">Dark Brown</option>
                              <option value="Light Brown">Light Brown</option>
                              <option value="Yellow">Yellow</option>
                              <option value="Green">Green</option>
                              <option value="Red">Red</option>
                              <option value="Black">Black</option>
              </select>
            </div>
          </div>

          <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
              Notes (Optional)
            </label>
            <textarea
                          className="input-field"
              rows={3}
                          placeholder="Any additional observations..."
              value={formData.notes}
              onChange={(e) => handleInputChange('notes', e.target.value)}
            />
          </div>

          {/* Image Upload Section */}
          <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
              Photo (Optional)
            </label>
                      <div className="space-y-3">
              <input
                              type="file"
                              accept="image/*"
                onChange={handleImageChange}
                              className="input-field"
              />

              {imagePreview && (
                              <div className="relative">
                  <img
                    src={imagePreview}
                                      alt="Preview"
                                      className="max-w-xs max-h-48 object-cover rounded border"
                  />
                  <button
                                      type="button"
                    onClick={removeImage}
                                      className="absolute top-2 right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-sm hover:bg-red-600"
                  >
                    ×
                  </button>
                </div>
              )}

              {editingEntry?.photoUrl && !imagePreview && (
                              <div className="text-sm text-gray-600">
                  Current photo:
                  <img
                    src={`${API_BASE_URL}${editingEntry.photoUrl}`}
                                      alt="Current entry photo"
                                      className="mt-2 max-w-xs max-h-48 object-cover rounded border"
                  />
                                  <p className="mt-1 text-xs">
                    Upload a new image to replace it
                  </p>
                </div>
              )}

                          <p className="text-xs text-gray-500">
                Max file size: 5MB. Accepted formats: JPG, PNG, GIF, WebP
              </p>
            </div>
          </div>

                  <div className="flex gap-3">
            <button
                          type="submit"
                          className="btn-primary flex-1"
              disabled={isSubmitting}
            >
              {isSubmitting
                              ? editingEntry
                                  ? 'Updating...'
                                  : 'Saving...'
                              : editingEntry
                                  ? 'Update Entry'
                                  : 'Save Entry'}
            </button>

                      {editingEntry && (
              <button
                              type="button"
                onClick={cancelEdit}
                              className="btn-secondary"
                disabled={isSubmitting}
              >
                Cancel
              </button>
            )}
          </div>
        </form>
      </div>

      {/* Recent Entries List */}
          <div className="mt-8">
              <h2 className="text-2xl font-bold mb-4">Recent Entries</h2>
              {isLoadingEntries ? (
                  <div className="card text-center py-8">
                      <p className="text-gray-600">Loading entries...</p>
                  </div>
              ) : entries.length === 0 ? (
                  <div className="card text-center py-8">
                      <p className="text-gray-600 flex items-center justify-center gap-2">
                          No entries yet. Create your first one above! <Logo size={24} />
                      </p>
                  </div>
                  ) : (
                      <div className="space-y-4">
                          {entries.map((entry) => (
                <div key={entry.id} className="card">
                    <div className="flex justify-between items-start mb-2">
                        <div>
                            <h3 className="font-semibold text-lg">
                                Bristol Type {entry.bristolType}
                            </h3>
                            <p className="text-sm text-gray-600">
                                {getBristolTypeDescription(entry.bristolType)}
                            </p>
                        </div>
                        <div className="flex items-center gap-2">
                            <span className="text-sm text-gray-500">
                                {new Date(entry.createdAt).toLocaleDateString()} at{' '}
                                {new Date(entry.createdAt).toLocaleTimeString([], {
                                    hour: '2-digit',
                                    minute: '2-digit'
                                })}
                            </span>
                            <div className="flex gap-1">
                                <button
                                    onClick={() => startEdit(entry)}
                                    className="text-blue-600 hover:text-blue-800 text-sm font-medium"
                                    disabled={isSubmitting}
                                >
                                    Edit
                                </button>
                                <button
                                    onClick={() => deleteEntry(entry.id)}
                                    className="text-red-600 hover:text-red-800 text-sm font-medium"
                                    disabled={isSubmitting}
                                >
                                    Delete
                                </button>
                            </div>
                        </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4 mb-3 text-sm">
                        {entry.volume && (
                            <div>
                                <span className="font-medium">Volume:</span>{' '}
                                {entry.volume}
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
                                src={`${API_BASE_URL}${entry.photoUrl}`}
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
            ))}
                  </div>
              )}
      </div>
    </div>
  )
}

import { useAuthStore } from '../stores/authStore'
import Logo from '../components/Logo'
import { EntryForm } from '../domains/entries/components/EntryForm'
import { EntryList } from '../domains/entries/components/EntryList'
import { useEntries } from '../domains/entries/hooks/useEntries'
import { useEntryForm } from '../domains/entries/hooks/useEntryForm'
import type { Entry } from '../domains/entries/types'

export function NewEntryPage() {
  const { token } = useAuthStore()
  const { entries, isLoading, error, isSubmitting, submitEntry, deleteEntry } = useEntries()

  const { formData, isEditing, updateFormData, resetForm, startEditing, cancelEditing } =
    useEntryForm()

  const handleSubmit = async (): Promise<void> => {
    try {
      await submitEntry(formData)
      resetForm()
    } catch {
      // Error is handled by the hook
    }
  }

  const handleEdit = (entry: Entry): void => {
    startEditing(entry)
  }

  const handleDelete = async (id: string): Promise<void> => {
    if (confirm('Are you sure you want to delete this entry?')) {
      await deleteEntry(id)
    }
  }

  if (!token) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">Please log in to access this page.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center space-x-4">
            <Logo size="40" />
            <h1 className="text-2xl font-bold text-gray-900">Track Entry</h1>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
            <p className="text-red-600 text-sm">{error}</p>
          </div>
        )}

        {/* Entry Form */}
        <EntryForm
          formData={formData}
          onUpdate={updateFormData}
          onSubmit={handleSubmit}
          onCancel={isEditing ? cancelEditing : undefined}
          isSubmitting={isSubmitting}
          isEditing={isEditing}
        />

        {/* Recent Entries */}
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Recent Entries</h2>

          {isLoading ? (
            <div className="text-center py-8">
              <p className="text-gray-500">Loading entries...</p>
            </div>
          ) : (
            <EntryList
              entries={entries}
              onEdit={handleEdit}
              onDelete={handleDelete}
              isSubmitting={isSubmitting}
            />
          )}
        </div>
      </div>
    </div>
  )
}

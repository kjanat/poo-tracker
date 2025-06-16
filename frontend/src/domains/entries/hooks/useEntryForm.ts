import { useState } from 'react'
import { EntryFactory } from '../EntryFactory'
import type { Entry, CreateEntryData } from '../types'

export interface UseEntryFormResult {
  formData: CreateEntryData
  isEditing: boolean
  editingEntry: Entry | null
  updateFormData: (updates: Partial<CreateEntryData>) => void
  resetForm: () => void
  startEditing: (entry: Entry) => void
  cancelEditing: () => void
}

export function useEntryForm(): UseEntryFormResult {
  const [formData, setFormData] = useState<CreateEntryData>(EntryFactory.createEmptyEntryData())
  const [editingEntry, setEditingEntry] = useState<Entry | null>(null)

  const updateFormData = (updates: Partial<CreateEntryData>): void => {
    setFormData(prev => ({ ...prev, ...updates }))
  }

  const resetForm = (): void => {
    setFormData(EntryFactory.createEmptyEntryData())
    setEditingEntry(null)
  }

  const startEditing = (entry: Entry): void => {
    setEditingEntry(entry)
    setFormData({
      bristolType: entry.bristolType,
      volume: entry.volume || '',
      color: entry.color || '',
      notes: entry.notes || '',
      photo: undefined // Photo editing not supported in this version
    })
  }

  const cancelEditing = (): void => {
    setEditingEntry(null)
    setFormData(EntryFactory.createEmptyEntryData())
  }

  return {
    formData,
    isEditing: editingEntry !== null,
    editingEntry,
    updateFormData,
    resetForm,
    startEditing,
    cancelEditing
  }
}

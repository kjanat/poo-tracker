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
    setFormData((prev) => ({ ...prev, ...updates }))
  }

  const resetForm = (): void => {
    setFormData(EntryFactory.createEmptyEntryData())
    setEditingEntry(null)
  }

  const startEditing = (entry: Entry): void => {
    setEditingEntry(entry)
    const newFormData: CreateEntryData = {
      bristolType: entry.bristolType,
      volume: entry.volume || '',
      color: entry.color || '',
      notes: entry.notes || '',
      photo: null, // Photo editing not supported in this version
      floaters: entry.floaters ?? false
    }

    // Only add optional properties if they have values
    if (entry.satisfaction !== undefined) newFormData.satisfaction = entry.satisfaction
    if (entry.pain !== undefined) newFormData.pain = entry.pain
    if (entry.strain !== undefined) newFormData.strain = entry.strain
    if (entry.smell) newFormData.smell = entry.smell

    setFormData(newFormData)
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

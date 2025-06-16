import { useState, useEffect, useCallback } from 'react'
import { container } from '../../../core/services'
import type { Entry, CreateEntryData, CreateEntryRequest } from '../types'
import type { EntryService } from '../EntryService'

export interface UseEntriesResult {
  entries: Entry[]
  isLoading: boolean
  error: string | null
  isSubmitting: boolean
  submitEntry: (data: CreateEntryData) => Promise<void>
  deleteEntry: (id: string) => Promise<void>
  refreshEntries: () => Promise<void>
}

export function useEntries(): UseEntriesResult {
  const [entries, setEntries] = useState<Entry[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const entryService = container.get<EntryService>('entryService')

  const loadEntries = useCallback(async (): Promise<void> => {
    try {
      setIsLoading(true)
      setError(null)
      const result = await entryService.getEntries()
      setEntries(result.entries)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load entries')
    } finally {
      setIsLoading(false)
    }
  }, [entryService])

  const submitEntry = async (data: CreateEntryData): Promise<void> => {
    try {
      setIsSubmitting(true)
      setError(null)
      
      // Convert CreateEntryData to CreateEntryRequest
      const requestData: CreateEntryRequest = {
        bristolType: data.bristolType,
        notes: data.notes,
        floaters: data.floaters ?? false,
      }

      // Only add optional properties if they have values
      if (data.volume && data.volume.trim()) requestData.volume = data.volume
      if (data.color && data.color.trim()) requestData.color = data.color
      if (data.satisfaction !== undefined) requestData.satisfaction = data.satisfaction
      if (data.pain !== undefined) requestData.pain = data.pain
      if (data.strain !== undefined) requestData.strain = data.strain
      if (data.smell) requestData.smell = data.smell
      
      await entryService.createEntry(requestData)
      await loadEntries() // Refresh the list
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create entry')
      throw err
    } finally {
      setIsSubmitting(false)
    }
  }

  const deleteEntry = async (id: string): Promise<void> => {
    try {
      setIsSubmitting(true)
      setError(null)
      await entryService.deleteEntry(id)
      await loadEntries() // Refresh the list
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete entry')
      throw err
    } finally {
      setIsSubmitting(false)
    }
  }

  useEffect(() => {
    loadEntries()
  }, [loadEntries])

  return {
    entries,
    isLoading,
    error,
    isSubmitting,
    submitEntry,
    deleteEntry,
    refreshEntries: loadEntries
  }
}

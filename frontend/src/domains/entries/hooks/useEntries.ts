import { useState, useEffect, useCallback } from 'react'
import { container } from '../../../core/services'
import type { Entry, CreateEntryData } from '../types'

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

  const entryService = container.get('entryService')

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
      await entryService.createEntry(data)
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

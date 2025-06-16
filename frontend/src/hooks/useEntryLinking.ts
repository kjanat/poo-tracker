import { useState, useEffect } from 'react'
import { useAuthStore } from '../stores/authStore'
import { API_BASE_URL, createAuthHeaders, handleApiResponse } from '../utils/api'
import type { Entry } from '../types'

export interface UseEntryLinkingReturn {
  availableEntries: Entry[]
  linkedEntries: Entry[]
  loading: boolean
  error: string
  setError: (error: string) => void
  fetchLinkedEntries: (mealId: string) => Promise<void>
  linkEntry: (mealId: string, entryId: string) => Promise<void>
  unlinkEntry: (mealId: string, entryId: string) => Promise<void>
}

export function useEntryLinking (): UseEntryLinkingReturn {
  const [availableEntries, setAvailableEntries] = useState<Entry[]>([])
  const [linkedEntries, setLinkedEntries] = useState<Entry[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string>('')
  
  const token = useAuthStore((state) => state.token)

  // Fetch all available entries
  useEffect(() => {
    const fetchAvailableEntries = async (): Promise<void> => {
      if (token == null) return

      try {
        const response = await fetch(`${API_BASE_URL}/api/entries`, {
          headers: createAuthHeaders(token)
        })

        const data = await handleApiResponse<{ entries: Entry[] }>(response)
        setAvailableEntries(data.entries)
      } catch (error) {
        console.error('Failed to fetch available entries:', error)
        setError(error instanceof Error ? error.message : 'Failed to fetch entries')
      }
    }

    void fetchAvailableEntries()
  }, [token])

  const fetchLinkedEntries = async (mealId: string): Promise<void> => {
    if (token == null) return

    try {
      setLoading(true)
      const response = await fetch(`${API_BASE_URL}/api/meals/${mealId}/entries`, {
        headers: createAuthHeaders(token)
      })

      const data = await handleApiResponse<Entry[]>(response)
      setLinkedEntries(data)
    } catch (error) {
      console.error('Failed to fetch linked entries:', error)
      setError(error instanceof Error ? error.message : 'Failed to fetch linked entries')
    } finally {
      setLoading(false)
    }
  }

  const linkEntry = async (mealId: string, entryId: string): Promise<void> => {
    if (token == null) return

    try {
      const response = await fetch(`${API_BASE_URL}/api/meals/${mealId}/entries/${entryId}`, {
        method: 'POST',
        headers: createAuthHeaders(token)
      })

      await handleApiResponse(response)
      
      // Update local state
      const entryToLink = availableEntries.find(entry => entry.id === entryId)
      if (entryToLink != null) {
        setLinkedEntries(prev => [...prev, entryToLink])
      }
    } catch (error) {
      console.error('Failed to link entry:', error)
      setError(error instanceof Error ? error.message : 'Failed to link entry')
    }
  }

  const unlinkEntry = async (mealId: string, entryId: string): Promise<void> => {
    if (token == null) return

    try {
      const response = await fetch(`${API_BASE_URL}/api/meals/${mealId}/entries/${entryId}`, {
        method: 'DELETE',
        headers: createAuthHeaders(token)
      })

      await handleApiResponse(response)
      
      // Update local state
      setLinkedEntries(prev => prev.filter(entry => entry.id !== entryId))
    } catch (error) {
      console.error('Failed to unlink entry:', error)
      setError(error instanceof Error ? error.message : 'Failed to unlink entry')
    }
  }

  return {
    availableEntries,
    linkedEntries,
    loading,
    error,
    setError,
    fetchLinkedEntries,
    linkEntry,
    unlinkEntry
  }
}

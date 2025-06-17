import { useState, useEffect, useCallback } from 'react'
import { useAuthStore } from '../stores/authStore'
import { API_BASE_URL, createAuthHeaders, handleApiResponse } from '../utils/api'
import type { Meal, MealFormData } from '../types'

export interface UseMealsReturn {
  meals: Meal[]
  loading: boolean
  error: string
  success: string
  setError: (error: string) => void
  setSuccess: (success: string) => void
  refreshMeals: () => Promise<void>
  createMeal: (formData: MealFormData) => Promise<void>
  updateMeal: (id: string, formData: MealFormData) => Promise<void>
  deleteMeal: (id: string) => Promise<void>
}

export function useMeals(): UseMealsReturn {
  const [meals, setMeals] = useState<Meal[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string>('')
  const [success, setSuccess] = useState<string>('')

  const token = useAuthStore((state) => state.token)

  const refreshMeals = useCallback(async (): Promise<void> => {
    if (token == null) return

    try {
      setLoading(true)
      const response = await fetch(`${API_BASE_URL}/api/meals`, {
        headers: createAuthHeaders(token)
      })

      const mealsData = await handleApiResponse<Meal[]>(response)

      // Fetch linked entries for each meal
      const mealsWithLinkedEntries = await Promise.all(
        mealsData.map(async (meal: Meal) => {
          try {
            const linkedResponse = await fetch(`${API_BASE_URL}/api/meals/${meal.id}/entries`, {
              headers: createAuthHeaders(token)
            })

            if (linkedResponse.ok) {
              const linkedData = await linkedResponse.json()
              return { ...meal, linkedEntries: linkedData }
            }
            return meal
          } catch (error) {
            console.error(`Failed to fetch linked entries for meal ${meal.id}:`, error)
            return meal
          }
        })
      )

      setMeals(mealsWithLinkedEntries)
    } catch (error) {
      console.error('Failed to fetch meals:', error)
      setError(error instanceof Error ? error.message : 'Failed to fetch meals')
    } finally {
      setLoading(false)
    }
  }, [token])

  const createMeal = async (formData: MealFormData): Promise<void> => {
    if (token == null) throw new Error('No authentication token')

    setLoading(true)
    try {
      const response = await fetch(`${API_BASE_URL}/api/meals`, {
        method: 'POST',
        headers: createAuthHeaders(token),
        body: JSON.stringify(formData)
      })

      await handleApiResponse(response)
      setSuccess('Meal created successfully!')
      await refreshMeals()
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create meal'
      setError(errorMessage)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const updateMeal = async (id: string, formData: MealFormData): Promise<void> => {
    if (token == null) throw new Error('No authentication token')

    setLoading(true)
    try {
      const response = await fetch(`${API_BASE_URL}/api/meals/${id}`, {
        method: 'PUT',
        headers: createAuthHeaders(token),
        body: JSON.stringify(formData)
      })

      await handleApiResponse(response)
      setSuccess('Meal updated successfully!')
      await refreshMeals()
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update meal'
      setError(errorMessage)
      throw error
    } finally {
      setLoading(false)
    }
  }

  const deleteMeal = async (id: string): Promise<void> => {
    if (token == null) throw new Error('No authentication token')

    if (!window.confirm('Are you sure you want to delete this meal?')) {
      return
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/meals/${id}`, {
        method: 'DELETE',
        headers: createAuthHeaders(token)
      })

      await handleApiResponse(response)
      setSuccess('Meal deleted successfully!')
      await refreshMeals()
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete meal'
      setError(errorMessage)
    }
  }

  useEffect(() => {
    void refreshMeals()
  }, [refreshMeals])

  return {
    meals,
    loading,
    error,
    success,
    setError,
    setSuccess,
    refreshMeals,
    createMeal,
    updateMeal,
    deleteMeal
  }
}

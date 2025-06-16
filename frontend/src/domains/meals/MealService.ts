import type { ApiClient } from '../../core/api/ApiClient'
import type { Meal, CreateMealRequest, UpdateMealRequest } from './types'
import type { Entry } from '../entries/types'

export class MealService {
  constructor (private readonly apiClient: ApiClient) {}

  async getMeals(): Promise<Meal[]> {
    const response = await this.apiClient.get<Meal[]>('/api/meals')
    return response.data
  }

  async getMeal(id: string): Promise<Meal> {
    const response = await this.apiClient.get<Meal>(`/api/meals/${id}`)
    return response.data
  }

  async createMeal(meal: CreateMealRequest): Promise<Meal> {
    const response = await this.apiClient.post<Meal>('/api/meals', meal)
    return response.data
  }

  async updateMeal(id: string, meal: UpdateMealRequest): Promise<Meal> {
    const response = await this.apiClient.put<Meal>(`/api/meals/${id}`, meal)
    return response.data
  }

  async deleteMeal(id: string): Promise<void> {
    await this.apiClient.delete(`/api/meals/${id}`)
  }

  async getLinkedEntries(mealId: string): Promise<Entry[]> {
    const response = await this.apiClient.get<Entry[]>(`/api/meals/${mealId}/entries`)
    return response.data
  }

  async linkEntry(mealId: string, entryId: string): Promise<void> {
    await this.apiClient.post(`/api/meals/${mealId}/entries/${entryId}`)
  }

  async unlinkEntry(mealId: string, entryId: string): Promise<void> {
    await this.apiClient.delete(`/api/meals/${mealId}/entries/${entryId}`)
  }

  async uploadMealPhoto(file: File): Promise<{ url: string }> {
    const response = await this.apiClient.uploadFile<{ url: string }>('/api/uploads', file)
    return response.data
  }
}

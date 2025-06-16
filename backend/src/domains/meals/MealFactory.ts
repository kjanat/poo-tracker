import type { CreateMealRequest, Meal } from './types'

export class MealFactory {
  static createFromRequest(request: CreateMealRequest, userId: string): Omit<Meal, 'id' | 'createdAt' | 'updatedAt'> {
    return {
      name: request.name.trim(),
      description: request.description?.trim(),
      mealTime: request.mealTime,
      category: request.category,
      cuisine: request.cuisine?.trim(),
      spicyLevel: request.spicyLevel,
      fiberRich: request.fiberRich ?? false,
      dairy: request.dairy ?? false,
      gluten: request.gluten ?? false,
      notes: request.notes?.trim(),
      photoUrl: request.photoUrl,
      userId
    }
  }

  static validateMealName(name: string): boolean {
    return typeof name === 'string' && name.trim().length > 0 && name.trim().length <= 200
  }

  static validateSpicyLevel(level: number | undefined): boolean {
    if (level === undefined) return true
    return Number.isInteger(level) && level >= 1 && level <= 10
  }

  static validateMealTime(mealTime: Date): boolean {
    return mealTime instanceof Date && !isNaN(mealTime.getTime())
  }

  static validateCategory(category: string | undefined): category is 'Breakfast' | 'Lunch' | 'Dinner' | 'Snack' | undefined {
    if (category === undefined) return true
    return ['Breakfast', 'Lunch', 'Dinner', 'Snack'].includes(category)
  }

  static getDefaultValues(): Partial<CreateMealRequest> {
    return {
      fiberRich: false,
      dairy: false,
      gluten: false,
      spicyLevel: 1
    }
  }

  static sanitizeNotes(notes: string | undefined): string | undefined {
    if (!notes) return undefined
    return notes.trim().substring(0, 1000) || undefined
  }

  static sanitizeName(name: string): string {
    return name.trim().substring(0, 200)
  }

  static sanitizeDescription(description: string | undefined): string | undefined {
    if (!description) return undefined
    return description.trim().substring(0, 500) || undefined
  }
}

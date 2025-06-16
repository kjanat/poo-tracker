import type { CreateMealRequest, Meal } from './types'

export function createFromRequest (
  request: CreateMealRequest,
  userId: string
): Omit<Meal, 'id' | 'createdAt' | 'updatedAt'> {
  return {
    name: request.name.trim(),
    description: request.description?.trim() ?? null,
    mealTime: request.mealTime,
    category: request.category ?? null,
    cuisine: request.cuisine?.trim() ?? null,
    spicyLevel: request.spicyLevel ?? null,
    fiberRich: request.fiberRich ?? false,
    dairy: request.dairy ?? false,
    gluten: request.gluten ?? false,
    notes: request.notes?.trim() ?? null,
    photoUrl: request.photoUrl ?? null,
    userId
  }
}

export function validateMealName (name: string): boolean {
  return typeof name === 'string' && name.trim().length > 0 && name.trim().length <= 200
}

export function validateSpicyLevel (level: number | undefined): boolean {
  if (level === undefined) return true
  return Number.isInteger(level) && level >= 1 && level <= 10
}

export function validateMealTime (mealTime: Date): boolean {
  return mealTime instanceof Date && !isNaN(mealTime.getTime())
}

export function validateCategory (
  category: string | undefined
): category is 'Breakfast' | 'Lunch' | 'Dinner' | 'Snack' | undefined {
  if (category === undefined) return true
  return ['Breakfast', 'Lunch', 'Dinner', 'Snack'].includes(category)
}

export function getDefaultValues (): Partial<CreateMealRequest> {
  return {
    fiberRich: false,
    dairy: false,
    gluten: false,
    spicyLevel: 1
  }
}

export function sanitizeNotes (notes: string | null | undefined): string | undefined {
  if (!notes) return undefined
  const sanitized = notes.trim().substring(0, 1000)
  return sanitized.length > 0 ? sanitized : undefined
}

export function sanitizeName (name: string): string {
  return name.trim().substring(0, 200)
}

export function sanitizeDescription (
  description: string | null | undefined
): string | undefined {
  if (!description) return undefined
  const sanitized = description.trim().substring(0, 500)
  return sanitized.length > 0 ? sanitized : undefined
}

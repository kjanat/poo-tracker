import type { CreateMealRequest, Meal } from './types'

export class MealFactory {
  static createEmpty(): CreateMealRequest {
    return {
      name: '',
      category: undefined,
      description: undefined,
      cuisine: undefined,
      spicyLevel: 1,
      fiberRich: false,
      dairy: false,
      gluten: false,
      notes: undefined,
      photoUrl: undefined
    }
  }

  static createFromFormData(formData: Record<string, any>): CreateMealRequest {
    return {
      name: formData.name || '',
      category: formData.category || undefined,
      description: formData.description || undefined,
      cuisine: formData.cuisine || undefined,
      spicyLevel: parseInt(formData.spicyLevel) || 1,
      fiberRich: Boolean(formData.fiberRich),
      dairy: Boolean(formData.dairy),
      gluten: Boolean(formData.gluten),
      notes: formData.notes || undefined,
      photoUrl: formData.photoUrl || undefined
    }
  }

  static createUpdatePayload(current: Meal, updates: Partial<CreateMealRequest>): CreateMealRequest {
    return {
      name: updates.name ?? current.name,
      category: updates.category ?? current.category,
      description: updates.description ?? current.description,
      cuisine: updates.cuisine ?? current.cuisine,
      spicyLevel: updates.spicyLevel ?? current.spicyLevel ?? 1,
      fiberRich: updates.fiberRich ?? current.fiberRich,
      dairy: updates.dairy ?? current.dairy,
      gluten: updates.gluten ?? current.gluten,
      notes: updates.notes ?? current.notes,
      photoUrl: updates.photoUrl ?? current.photoUrl
    }
  }

  static validateSpicyLevel(level: number): boolean {
    return Number.isInteger(level) && level >= 1 && level <= 5
  }

  static validateName(name: string): boolean {
    return typeof name === 'string' && name.trim().length > 0
  }

  static getDietaryTags(meal: Meal): string[] {
    const tags: string[] = []
    if (meal.fiberRich) tags.push('High Fiber')
    if (meal.dairy) tags.push('Dairy')
    if (meal.gluten) tags.push('Gluten')
    return tags
  }
}

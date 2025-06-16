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

  static createFromFormData(formData: Record<string, unknown>): CreateMealRequest {
    const getValue = (key: string): string | undefined => {
      const value = formData[key]
      return typeof value === 'string' && value.trim() !== '' ? value : undefined
    }

    const getNumberValue = (key: string): number | undefined => {
      const value = formData[key]
      if (typeof value === 'string' && value.trim() !== '') {
        const parsed = parseInt(value, 10)
        return isNaN(parsed) ? undefined : parsed
      }
      return undefined
    }

    const nameValue = getValue('name')
    if (!nameValue) {
      throw new Error('Meal name is required')
    }

    return {
      name: nameValue,
      category: getValue('category'),
      description: getValue('description'),
      cuisine: getValue('cuisine'),
      spicyLevel: getNumberValue('spicyLevel') ?? 1,
      fiberRich: Boolean(formData.fiberRich),
      dairy: Boolean(formData.dairy),
      gluten: Boolean(formData.gluten),
      notes: getValue('notes'),
      photoUrl: getValue('photoUrl')
    }
  }

  static createUpdatePayload(
    current: Meal,
    updates: Partial<CreateMealRequest>
  ): CreateMealRequest {
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

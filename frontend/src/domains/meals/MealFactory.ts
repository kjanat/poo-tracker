import type { CreateMealRequest, Meal } from './types'

export class MealFactory {
  static createEmpty(): CreateMealRequest {
    return {
      name: '',
      category: 'BREAKFAST',
      spicyLevel: 1,
      fiberRich: false,
      dairy: false,
      gluten: false,
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

    const categoryValue = getValue('category') as 'BREAKFAST' | 'LUNCH' | 'DINNER' | 'SNACK' | undefined
    if (!categoryValue || !['BREAKFAST', 'LUNCH', 'DINNER', 'SNACK'].includes(categoryValue)) {
      throw new Error('Valid meal category is required')
    }

    const result: CreateMealRequest = {
      name: nameValue,
      category: categoryValue,
      spicyLevel: getNumberValue('spicyLevel') ?? 1,
      fiberRich: Boolean(formData.fiberRich),
      dairy: Boolean(formData.dairy),
      gluten: Boolean(formData.gluten),
    }

    // Only add optional properties if they have values
    const description = getValue('description')
    if (description) result.description = description

    const cuisine = getValue('cuisine')
    if (cuisine) result.cuisine = cuisine

    const notes = getValue('notes')
    if (notes) result.notes = notes

    const photoUrl = getValue('photoUrl')
    if (photoUrl) result.photoUrl = photoUrl

    return result
  }

  static createUpdatePayload(
    current: Meal,
    updates: Partial<CreateMealRequest>
  ): CreateMealRequest {
    const result: CreateMealRequest = {
      name: updates.name ?? current.name,
      category: updates.category ?? current.category ?? 'BREAKFAST',
      spicyLevel: updates.spicyLevel ?? current.spicyLevel ?? 1,
      fiberRich: updates.fiberRich ?? current.fiberRich,
      dairy: updates.dairy ?? current.dairy,
      gluten: updates.gluten ?? current.gluten,
    }

    // Only add optional properties if they have values
    const description = updates.description ?? current.description
    if (description) result.description = description

    const cuisine = updates.cuisine ?? current.cuisine
    if (cuisine) result.cuisine = cuisine

    const notes = updates.notes ?? current.notes
    if (notes) result.notes = notes

    const photoUrl = updates.photoUrl ?? current.photoUrl
    if (photoUrl) result.photoUrl = photoUrl

    return result
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

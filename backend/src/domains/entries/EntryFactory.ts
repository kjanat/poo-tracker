import type { CreateEntryRequest, Entry } from './types'

export class EntryFactory {
  static createFromRequest (request: CreateEntryRequest, userId: string): Omit<Entry, 'id' | 'createdAt' | 'updatedAt'> {
    return {
      bristolType: request.bristolType,
      volume: request.volume,
      color: request.color,
      consistency: request.consistency,
      floaters: request.floaters ?? false,
      pain: request.pain,
      strain: request.strain,
      satisfaction: request.satisfaction,
      notes: request.notes,
      smell: request.smell,
      photoUrl: request.photoUrl,
      userId
    }
  }

  static validateBristolType (type: number): boolean {
    return Number.isInteger(type) && type >= 1 && type <= 7
  }

  static validateRating (rating: number | undefined): boolean {
    if (rating === undefined) return true
    return Number.isInteger(rating) && rating >= 1 && rating <= 10
  }

  static validateEnumValue<T extends string>(value: string | undefined, allowedValues: readonly T[]): value is T | undefined {
    if (value === undefined) return true
    return allowedValues.includes(value as T)
  }

  static getDefaultValues (): Partial<CreateEntryRequest> {
    return {
      bristolType: 4,
      floaters: false,
      notes: ''
    }
  }

  static sanitizeNotes (notes: string | undefined): string | undefined {
    if (!notes) return undefined
    return notes.trim().substring(0, 1000) || undefined
  }
}

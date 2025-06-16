import type { CreateEntryRequest, Entry } from './types'

export class EntryFactory {

  static createFromRequest (
    request: CreateEntryRequest,
    userId: string
  ): Omit<Entry, 'id' | 'createdAt' | 'updatedAt'> {
    return {
      bristolType: request.bristolType,
      volume: request.volume ?? null,
      color: request.color ?? null,
      consistency: request.consistency ?? null,
      floaters: request.floaters ?? false,
      pain: request.pain ?? null,
      strain: request.strain ?? null,
      satisfaction: request.satisfaction ?? null,
      notes: request.notes ?? null,
      smell: request.smell ?? null,
      photoUrl: request.photoUrl ?? null,
      userId
    }
  }

  static validateBristolType (type: number): boolean {
    return Number.isInteger(type) && type >= 1 && type <= 7
  }

  static validateRating (rating: number | null | undefined): boolean {
    if (rating == null) return true
    return Number.isInteger(rating) && rating >= 1 && rating <= 10
  }

  static validateEnumValue<T extends string> (
    value: string | undefined,
    allowedValues: readonly T[]
  ): value is T | undefined {
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

  static sanitizeNotes (notes: string | null | undefined): string | null {
    if (!notes) return null
    const sanitized = notes.trim().substring(0, 1000)
    return sanitized.length > 0 ? sanitized : null
  }
}

export default EntryFactory

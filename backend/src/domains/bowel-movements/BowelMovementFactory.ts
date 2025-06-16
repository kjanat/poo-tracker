import type { CreateBowelMovementRequest, BowelMovement } from './types'

export class BowelMovementFactory {

  static createFromRequest (
    request: CreateBowelMovementRequest,
    userId: string
  ): Omit<BowelMovement, 'id' | 'createdAt' | 'updatedAt'> {
    return {
      bristolType: request.bristolType,
      recordedAt: request.recordedAt ?? new Date(),
      volume: request.volume ?? null,
      color: request.color ?? null,
      consistency: request.consistency ?? null,
      floaters: request.floaters ?? false,
      pain: request.pain ?? 1,
      strain: request.strain ?? 1,
      satisfaction: request.satisfaction ?? 5,
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

  static getDefaultValues (): Partial<CreateBowelMovementRequest> {
    return {
      bristolType: 4,
      floaters: false,
      pain: 1,
      strain: 1,
      satisfaction: 5
    }
  }

  static sanitizeNotes (notes: string | null | undefined): string | null {
    if (!notes || typeof notes !== 'string') {
      return null
    }

    // Basic sanitization: trim whitespace and limit length
    const trimmed = notes.trim()
    if (trimmed.length === 0) {
      return null
    }

    // Limit length to 1000 characters
    return trimmed.length > 1000 ? trimmed.substring(0, 1000) : trimmed
  }
}

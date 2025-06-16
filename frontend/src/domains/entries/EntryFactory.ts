import type { CreateEntryRequest, Entry, CreateEntryData } from './types'

export class EntryFactory {
  static createEmpty(): CreateEntryRequest {
    return {
      bristolType: 4,
      volume: undefined,
      color: undefined,
      consistency: undefined,
      notes: '',
      satisfaction: undefined,
      pain: undefined,
      strain: undefined,
      floaters: false,
      smell: undefined,
      photoUrl: undefined
    }
  }

  static createEmptyEntryData(): CreateEntryData {
    return {
      bristolType: 4,
      volume: '',
      color: '',
      notes: '',
      photo: undefined
    }
  }

  static createFromFormData(formData: Record<string, unknown>): CreateEntryRequest {
    return {
      bristolType: parseInt(formData.bristolType) || 4,
      volume: formData.volume || undefined,
      color: formData.color || undefined,
      consistency: formData.consistency || undefined,
      notes: formData.notes || '',
      satisfaction: formData.satisfaction ? parseInt(formData.satisfaction) : undefined,
      pain: formData.pain ? parseInt(formData.pain) : undefined,
      strain: formData.strain ? parseInt(formData.strain) : undefined,
      floaters: Boolean(formData.floaters),
      smell: formData.smell || undefined,
      photoUrl: formData.photoUrl || undefined
    }
  }

  static createUpdatePayload(current: Entry, updates: Partial<CreateEntryRequest>): CreateEntryRequest {
    return {
      bristolType: updates.bristolType ?? current.bristolType,
      volume: updates.volume ?? current.volume,
      color: updates.color ?? current.color,
      consistency: updates.consistency ?? current.consistency,
      notes: updates.notes ?? current.notes,
      satisfaction: updates.satisfaction ?? current.satisfaction,
      pain: updates.pain ?? current.pain,
      strain: updates.strain ?? current.strain,
      floaters: updates.floaters ?? current.floaters,
      smell: updates.smell ?? current.smell,
      photoUrl: updates.photoUrl ?? current.photoUrl
    }
  }

  static validateBristolType(type: number): boolean {
    return Number.isInteger(type) && type >= 1 && type <= 7
  }

  static validateRating(rating: number | undefined): boolean {
    if (rating === undefined) return true
    return Number.isInteger(rating) && rating >= 1 && rating <= 10
  }
}

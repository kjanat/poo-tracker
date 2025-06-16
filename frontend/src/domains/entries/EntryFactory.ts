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

    return {
      bristolType: getNumberValue('bristolType') ?? 4,
      volume: getValue('volume'),
      color: getValue('color'),
      consistency: getValue('consistency'),
      notes: getValue('notes') ?? '',
      satisfaction: getNumberValue('satisfaction'),
      pain: getNumberValue('pain'),
      strain: getNumberValue('strain'),
      floaters: Boolean(formData.floaters),
      smell: getValue('smell'),
      photoUrl: getValue('photoUrl')
    }
  }

  static createUpdatePayload(
    current: Entry,
    updates: Partial<CreateEntryRequest>
  ): CreateEntryRequest {
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

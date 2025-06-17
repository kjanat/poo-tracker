import type { CreateEntryRequest, Entry, CreateEntryData } from './types'

export class EntryFactory {
  static createEmpty(): CreateEntryRequest {
    return {
      bristolType: 4,
      notes: '',
      floaters: false
    }
  }

  static createEmptyEntryData(): CreateEntryData {
    return {
      bristolType: 4,
      volume: '',
      color: '',
      notes: '',
      photo: null,
      floaters: false
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

    const result: CreateEntryRequest = {
      bristolType: getNumberValue('bristolType') ?? 4,
      notes: getValue('notes') ?? '',
      floaters: Boolean(formData.floaters)
    }

    // Only add optional properties if they have values
    const volume = getValue('volume')
    if (volume) result.volume = volume

    const color = getValue('color')
    if (color) result.color = color

    const consistency = getValue('consistency')
    if (consistency) result.consistency = consistency

    const satisfaction = getNumberValue('satisfaction')
    if (satisfaction !== undefined) result.satisfaction = satisfaction

    const pain = getNumberValue('pain')
    if (pain !== undefined) result.pain = pain

    const strain = getNumberValue('strain')
    if (strain !== undefined) result.strain = strain

    const smell = getValue('smell')
    if (smell) result.smell = smell

    const photoUrl = getValue('photoUrl')
    if (photoUrl) result.photoUrl = photoUrl

    return result
  }

  static createUpdatePayload(
    current: Entry,
    updates: Partial<CreateEntryRequest>
  ): CreateEntryRequest {
    const result: CreateEntryRequest = {
      bristolType: updates.bristolType ?? current.bristolType,
      notes: updates.notes ?? current.notes ?? '',
      floaters: updates.floaters ?? current.floaters ?? false
    }

    // Only add optional properties if they have values
    const volume = updates.volume ?? current.volume
    if (volume) result.volume = volume

    const color = updates.color ?? current.color
    if (color) result.color = color

    const consistency = updates.consistency ?? current.consistency
    if (consistency) result.consistency = consistency

    const satisfaction = updates.satisfaction ?? current.satisfaction
    if (satisfaction !== undefined) result.satisfaction = satisfaction

    const pain = updates.pain ?? current.pain
    if (pain !== undefined) result.pain = pain

    const strain = updates.strain ?? current.strain
    if (strain !== undefined) result.strain = strain

    const smell = updates.smell ?? current.smell
    if (smell) result.smell = smell

    const photoUrl = updates.photoUrl ?? current.photoUrl
    if (photoUrl) result.photoUrl = photoUrl

    return result
  }

  static validateBristolType(type: number): boolean {
    return Number.isInteger(type) && type >= 1 && type <= 7
  }

  static validateRating(rating: number | undefined): boolean {
    if (rating === undefined) return true
    return Number.isInteger(rating) && rating >= 1 && rating <= 10
  }
}

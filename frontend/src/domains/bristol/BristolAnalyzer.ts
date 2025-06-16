export enum BristolType {
  HARD_LUMPS = 1,
  LUMPY_SAUSAGE = 2,
  CRACKED_SAUSAGE = 3,
  SMOOTH_SAUSAGE = 4,
  SOFT_BLOBS = 5,
  FLUFFY_PIECES = 6,
  WATERY = 7
}

export enum BristolCategory {
  CONSTIPATED = 'Constipated',
  NORMAL = 'Normal',
  LOOSE = 'Loose'
}

export interface BristolTypeInfo {
  type: BristolType
  description: string
  category: BristolCategory
  severity: string
  color: string
}

export class BristolAnalyzer {
  private static readonly typeInfo: Record<BristolType, BristolTypeInfo> = {
    [BristolType.HARD_LUMPS]: {
      type: BristolType.HARD_LUMPS,
      description: 'Hard lumps',
      category: BristolCategory.CONSTIPATED,
      severity: 'Severe constipation',
      color: '#dc2626'
    },
    [BristolType.LUMPY_SAUSAGE]: {
      type: BristolType.LUMPY_SAUSAGE,
      description: 'Lumpy sausage',
      category: BristolCategory.CONSTIPATED,
      severity: 'Mild constipation',
      color: '#ea580c'
    },
    [BristolType.CRACKED_SAUSAGE]: {
      type: BristolType.CRACKED_SAUSAGE,
      description: 'Cracked sausage',
      category: BristolCategory.NORMAL,
      severity: 'Normal',
      color: '#16a34a'
    },
    [BristolType.SMOOTH_SAUSAGE]: {
      type: BristolType.SMOOTH_SAUSAGE,
      description: 'Smooth sausage',
      category: BristolCategory.NORMAL,
      severity: 'Ideal',
      color: '#059669'
    },
    [BristolType.SOFT_BLOBS]: {
      type: BristolType.SOFT_BLOBS,
      description: 'Soft blobs',
      category: BristolCategory.LOOSE,
      severity: 'Lacking fiber',
      color: '#ca8a04'
    },
    [BristolType.FLUFFY_PIECES]: {
      type: BristolType.FLUFFY_PIECES,
      description: 'Fluffy pieces',
      category: BristolCategory.LOOSE,
      severity: 'Mild diarrhea',
      color: '#d97706'
    },
    [BristolType.WATERY]: {
      type: BristolType.WATERY,
      description: 'Watery',
      category: BristolCategory.LOOSE,
      severity: 'Severe diarrhea',
      color: '#dc2626'
    }
  }

  static getTypeInfo(type: BristolType): BristolTypeInfo {
    return this.typeInfo[type]
  }

  static getDescription(type: BristolType): string {
    return this.typeInfo[type]?.description ?? 'Unknown'
  }

  static getCategory(type: BristolType): BristolCategory {
    return this.typeInfo[type]?.category ?? BristolCategory.NORMAL
  }

  static getSeverity(type: BristolType): string {
    return this.typeInfo[type]?.severity ?? 'Unknown'
  }

  static getColor(type: BristolType): string {
    return this.typeInfo[type]?.color ?? '#6b7280'
  }

  static calculateAverage(distribution: Array<{ type: number; count: number }>): number {
    if (distribution.length === 0) return 0

    const totalCount = distribution.reduce((sum, item) => sum + item.count, 0)
    const weightedSum = distribution.reduce((sum, item) => sum + item.type * item.count, 0)

    return totalCount > 0 ? Number((weightedSum / totalCount).toFixed(1)) : 0
  }

  static isHealthy(type: BristolType): boolean {
    return type === BristolType.CRACKED_SAUSAGE || type === BristolType.SMOOTH_SAUSAGE
  }

  static needsAttention(type: BristolType): boolean {
    return type === BristolType.HARD_LUMPS || type === BristolType.WATERY
  }

  static getAllTypes(): BristolTypeInfo[] {
    return Object.values(this.typeInfo)
  }
}

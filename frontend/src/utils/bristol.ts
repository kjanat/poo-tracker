export const getBristolTypeDescription = (type: number): string => {
  const descriptions: Record<number, string> = {
    1: 'Hard lumps',
    2: 'Lumpy sausage',
    3: 'Cracked sausage',
    4: 'Smooth sausage',
    5: 'Soft blobs',
    6: 'Fluffy pieces',
    7: 'Watery'
  }
  return descriptions[type] ?? 'Unknown'
}

export const getBristolTypeCategory = (type: number): string => {
  if (type <= 2) return 'Constipated'
  if (type <= 4) return 'Normal'
  return 'Loose'
}

export const getAverageBristolType = (
  bristolDistribution: Array<{ type: number, count: number }>
): number => {
  if (bristolDistribution.length === 0) return 0

  const totalCount = bristolDistribution.reduce(
    (sum, item) => sum + item.count,
    0
  )
  const weightedSum = bristolDistribution.reduce(
    (sum, item) => sum + item.type * item.count,
    0
  )

  return totalCount > 0 ? Number((weightedSum / totalCount).toFixed(1)) : 0
}

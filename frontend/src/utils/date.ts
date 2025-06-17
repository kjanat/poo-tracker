import type { EntryResponse } from '../types'

export const getThisWeekCount = (entries: EntryResponse[]): number => {
  const oneWeekAgo = new Date()
  oneWeekAgo.setDate(oneWeekAgo.getDate() - 7)

  return entries.filter((entry) => new Date(entry.createdAt) >= oneWeekAgo).length
}

export const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

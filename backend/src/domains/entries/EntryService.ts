import { PrismaClient } from '@prisma/client'
import type { Entry, CreateEntryRequest, UpdateEntryRequest, EntryFilters, EntryListResponse } from './types'
import { EntryFactory } from './EntryFactory'

export class EntryService {
  constructor (private readonly prisma: PrismaClient) {}

  async findByUserId (userId: string, filters: EntryFilters = {}): Promise<EntryListResponse> {
    const {
      page = 1,
      limit = 20,
      sortBy = 'createdAt',
      sortOrder = 'desc',
      bristolType,
      dateFrom,
      dateTo
    } = filters

    const skip = (page - 1) * limit

    // Build where clause
    const where: any = { userId }

    if (bristolType) {
      where.bristolType = bristolType
    }

    if ((dateFrom != null) || (dateTo != null)) {
      where.createdAt = {}
      if (dateFrom != null) where.createdAt.gte = dateFrom
      if (dateTo != null) where.createdAt.lte = dateTo
    }

    // Execute queries in parallel
    const [entries, total] = await Promise.all([
      this.prisma.entry.findMany({
        where,
        orderBy: { [sortBy]: sortOrder },
        skip,
        take: limit
      }),
      this.prisma.entry.count({ where })
    ])

    return {
      entries,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    }
  }

  async findById (id: string, userId: string): Promise<Entry | null> {
    return await this.prisma.entry.findFirst({
      where: { id, userId }
    })
  }

  async create (request: CreateEntryRequest, userId: string): Promise<Entry> {
    const entryData = EntryFactory.createFromRequest(request, userId)

    return await this.prisma.entry.create({
      data: {
        ...entryData,
        notes: EntryFactory.sanitizeNotes(entryData.notes) || null
      }
    })
  }

  async update (id: string, request: UpdateEntryRequest, userId: string): Promise<Entry | null> {
    const existingEntry = await this.findById(id, userId)
    if (existingEntry == null) {
      return null
    }

    const updateData: any = {}

    // Only update provided fields
    if (request.bristolType !== undefined) updateData.bristolType = request.bristolType
    if (request.volume !== undefined) updateData.volume = request.volume
    if (request.color !== undefined) updateData.color = request.color
    if (request.consistency !== undefined) updateData.consistency = request.consistency
    if (request.floaters !== undefined) updateData.floaters = request.floaters
    if (request.pain !== undefined) updateData.pain = request.pain
    if (request.strain !== undefined) updateData.strain = request.strain
    if (request.satisfaction !== undefined) updateData.satisfaction = request.satisfaction
    if (request.notes !== undefined) updateData.notes = EntryFactory.sanitizeNotes(request.notes)
    if (request.smell !== undefined) updateData.smell = request.smell
    if (request.photoUrl !== undefined) updateData.photoUrl = request.photoUrl

    return await this.prisma.entry.update({
      where: { id },
      data: updateData
    })
  }

  async delete (id: string, userId: string): Promise<boolean> {
    const existingEntry = await this.findById(id, userId)
    if (existingEntry == null) {
      return false
    }

    await this.prisma.entry.delete({
      where: { id }
    })

    return true
  }

  async getAnalytics (userId: string): Promise<{
    totalEntries: number
    bristolDistribution: Array<{ type: number, count: number }>
    averageSatisfaction: number | null
    recentEntries: Entry[]
  }> {
    const [totalEntries, bristolStats, satisfactionAvg, recentEntries] = await Promise.all([
      this.prisma.entry.count({ where: { userId } }),

      this.prisma.entry.groupBy({
        by: ['bristolType'],
        where: { userId },
        _count: { bristolType: true }
      }),

      this.prisma.entry.aggregate({
        where: { userId, satisfaction: { not: null } },
        _avg: { satisfaction: true }
      }),

      this.prisma.entry.findMany({
        where: { userId },
        orderBy: { createdAt: 'desc' },
        take: 10
      })
    ])

    return {
      totalEntries,
      bristolDistribution: bristolStats.map(stat => ({
        type: stat.bristolType,
        count: stat._count.bristolType
      })),
      averageSatisfaction: satisfactionAvg._avg.satisfaction,
      recentEntries
    }
  }
}

import { PrismaClient } from '@prisma/client'
import type {
  BowelMovement,
  CreateBowelMovementRequest,
  UpdateBowelMovementRequest,
  BowelMovementFilters,
  BowelMovementListResponse
} from './types'
import { BowelMovementFactory } from './BowelMovementFactory'

export class BowelMovementService {
  constructor(private readonly prisma: PrismaClient) {}

  async findByUserId(
    userId: string,
    filters: BowelMovementFilters = {}
  ): Promise<BowelMovementListResponse> {
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
    const where: { userId: string; bristolType?: number; createdAt?: { gte?: Date; lte?: Date } } =
      { userId }

    if (bristolType) {
      where.bristolType = bristolType
    }

    if (dateFrom != null || dateTo != null) {
      where.createdAt = {}
      if (dateFrom != null) where.createdAt.gte = dateFrom
      if (dateTo != null) where.createdAt.lte = dateTo
    }

    // Execute queries in parallel
    const [bowelMovements, total] = await Promise.all([
      this.prisma.bowelMovement.findMany({
        where,
        orderBy: { [sortBy]: sortOrder },
        skip,
        take: limit,
        include: {
          details: true,
          symptoms: true
        }
      }),
      this.prisma.bowelMovement.count({ where })
    ])

    return {
      bowelMovements,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    }
  }

  async findById(id: string, userId: string): Promise<BowelMovement | null> {
    return await this.prisma.bowelMovement.findFirst({
      where: { id, userId },
      include: {
        details: true,
        symptoms: true,
        meals: {
          include: {
            meal: true
          }
        }
      }
    })
  }

  async create(request: CreateBowelMovementRequest, userId: string): Promise<BowelMovement> {
    const bowelMovementData = BowelMovementFactory.createFromRequest(request, userId)

    const result = await this.prisma.bowelMovement.create({
      data: bowelMovementData,
      include: {
        details: true,
        symptoms: true
      }
    })

    // Create details separately if notes exist
    if (request.notes) {
      await this.prisma.bowelMovementDetails.create({
        data: {
          bowelMovementId: result.id,
          notes: BowelMovementFactory.sanitizeNotes(request.notes)
        }
      })
    }

    return result
  }

  async update(
    id: string,
    request: UpdateBowelMovementRequest,
    userId: string
  ): Promise<BowelMovement | null> {
    const existingBowelMovement = await this.findById(id, userId)
    if (existingBowelMovement == null) {
      return null
    }

    const updateData: Partial<Omit<BowelMovement, 'id' | 'createdAt' | 'updatedAt' | 'userId'>> = {}

    // Only update provided fields
    if (request.bristolType !== undefined) updateData.bristolType = request.bristolType
    if (request.recordedAt !== undefined) updateData.recordedAt = request.recordedAt
    if (request.volume !== undefined) updateData.volume = request.volume
    if (request.color !== undefined) updateData.color = request.color
    if (request.consistency !== undefined) updateData.consistency = request.consistency
    if (request.floaters !== undefined) updateData.floaters = request.floaters
    if (request.pain !== undefined) updateData.pain = request.pain
    if (request.strain !== undefined) updateData.strain = request.strain
    if (request.satisfaction !== undefined) updateData.satisfaction = request.satisfaction
    if (request.smell !== undefined) updateData.smell = request.smell
    if (request.photoUrl !== undefined) updateData.photoUrl = request.photoUrl

    const result = await this.prisma.bowelMovement.update({
      where: { id },
      data: updateData,
      include: {
        details: true,
        symptoms: true
      }
    })

    // Handle notes separately in BowelMovementDetails
    if (request.notes !== undefined) {
      const sanitizedNotes = BowelMovementFactory.sanitizeNotes(request.notes)
      await this.prisma.bowelMovementDetails.upsert({
        where: { bowelMovementId: id },
        update: { notes: sanitizedNotes },
        create: {
          bowelMovementId: id,
          notes: sanitizedNotes
        }
      })
    }

    return result
  }

  async delete(id: string, userId: string): Promise<boolean> {
    const existingBowelMovement = await this.findById(id, userId)
    if (existingBowelMovement == null) {
      return false
    }

    await this.prisma.bowelMovement.delete({
      where: { id }
    })

    return true
  }

  async getAnalytics(userId: string): Promise<{
    totalBowelMovements: number
    bristolDistribution: Array<{ type: number; count: number }>
    averageSatisfaction: number | null
    averagePain: number | null
    recentBowelMovements: BowelMovement[]
  }> {
    const [totalBowelMovements, bristolStats, satisfactionAvg, painAvg, recentBowelMovements] =
      await Promise.all([
        this.prisma.bowelMovement.count({ where: { userId } }),

        this.prisma.bowelMovement.groupBy({
          by: ['bristolType'],
          where: { userId },
          _count: { bristolType: true }
        }),

        this.prisma.bowelMovement.aggregate({
          where: { userId },
          _avg: { satisfaction: true }
        }),

        this.prisma.bowelMovement.aggregate({
          where: { userId },
          _avg: { pain: true }
        }),

        this.prisma.bowelMovement.findMany({
          where: { userId },
          orderBy: { createdAt: 'desc' },
          take: 10,
          include: {
            details: true
          }
        })
      ])

    return {
      totalBowelMovements,
      bristolDistribution: bristolStats.map(
        (stat: { bristolType: number; _count: { bristolType: number } }) => ({
          type: stat.bristolType,
          count: stat._count.bristolType
        })
      ),
      averageSatisfaction: satisfactionAvg._avg.satisfaction,
      averagePain: painAvg._avg.pain,
      recentBowelMovements
    }
  }
}

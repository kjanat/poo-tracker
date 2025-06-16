import { PrismaClient } from '@prisma/client'
import type { Meal, CreateMealRequest, UpdateMealRequest, MealFilters, MealListResponse } from './types'
import { MealFactory } from './MealFactory'

export class MealService {
  constructor(private readonly prisma: PrismaClient) {}

  async findByUserId(userId: string, filters: MealFilters = {}): Promise<MealListResponse> {
    const {
      page = 1,
      limit = 20,
      sortBy = 'mealTime',
      sortOrder = 'desc',
      category,
      dateFrom,
      dateTo,
      fiberRich,
      dairy,
      gluten
    } = filters

    const skip = (page - 1) * limit

    // Build where clause
    const where: any = { userId }
    
    if (category) {
      where.category = category
    }

    if (dateFrom || dateTo) {
      where.mealTime = {}
      if (dateFrom) where.mealTime.gte = dateFrom
      if (dateTo) where.mealTime.lte = dateTo
    }

    if (fiberRich !== undefined) {
      where.fiberRich = fiberRich
    }

    if (dairy !== undefined) {
      where.dairy = dairy
    }

    if (gluten !== undefined) {
      where.gluten = gluten
    }

    // Execute queries in parallel
    const [meals, total] = await Promise.all([
      this.prisma.meal.findMany({
        where,
        orderBy: { [sortBy]: sortOrder },
        skip,
        take: limit
      }),
      this.prisma.meal.count({ where })
    ])

    return {
      meals,
      pagination: {
        page,
        limit,
        total,
        pages: Math.ceil(total / limit)
      }
    }
  }

  async findById(id: string, userId: string): Promise<Meal | null> {
    return this.prisma.meal.findFirst({
      where: { id, userId }
    })
  }

  async create(request: CreateMealRequest, userId: string): Promise<Meal> {
    const mealData = MealFactory.createFromRequest(request, userId)
    
    return this.prisma.meal.create({
      data: {
        ...mealData,
        name: MealFactory.sanitizeName(mealData.name),
        description: MealFactory.sanitizeDescription(mealData.description),
        notes: MealFactory.sanitizeNotes(mealData.notes)
      }
    })
  }

  async update(id: string, request: UpdateMealRequest, userId: string): Promise<Meal | null> {
    const existingMeal = await this.findById(id, userId)
    if (!existingMeal) {
      return null
    }

    const updateData: any = {}
    
    // Only update provided fields
    if (request.name !== undefined) updateData.name = MealFactory.sanitizeName(request.name)
    if (request.description !== undefined) updateData.description = MealFactory.sanitizeDescription(request.description)
    if (request.mealTime !== undefined) updateData.mealTime = request.mealTime
    if (request.category !== undefined) updateData.category = request.category
    if (request.cuisine !== undefined) updateData.cuisine = request.cuisine?.trim()
    if (request.spicyLevel !== undefined) updateData.spicyLevel = request.spicyLevel
    if (request.fiberRich !== undefined) updateData.fiberRich = request.fiberRich
    if (request.dairy !== undefined) updateData.dairy = request.dairy
    if (request.gluten !== undefined) updateData.gluten = request.gluten
    if (request.notes !== undefined) updateData.notes = MealFactory.sanitizeNotes(request.notes)
    if (request.photoUrl !== undefined) updateData.photoUrl = request.photoUrl

    return this.prisma.meal.update({
      where: { id },
      data: updateData
    })
  }

  async delete(id: string, userId: string): Promise<boolean> {
    const existingMeal = await this.findById(id, userId)
    if (!existingMeal) {
      return false
    }

    await this.prisma.meal.delete({
      where: { id }
    })

    return true
  }

  async getAnalytics(userId: string): Promise<{
    totalMeals: number
    mealsByCategory: Array<{ category: string, count: number }>
    averageSpicyLevel: number | null
    dietaryDistribution: {
      fiberRich: number
      dairy: number
      gluten: number
    }
    recentMeals: Meal[]
  }> {
    const [totalMeals, categoryStats, spicyAvg, dietaryStats, recentMeals] = await Promise.all([
      this.prisma.meal.count({ where: { userId } }),
      
      this.prisma.meal.groupBy({
        by: ['category'],
        where: { userId },
        _count: { category: true }
      }),
      
      this.prisma.meal.aggregate({
        where: { userId, spicyLevel: { not: null } },
        _avg: { spicyLevel: true }
      }),

      Promise.all([
        this.prisma.meal.count({ where: { userId, fiberRich: true } }),
        this.prisma.meal.count({ where: { userId, dairy: true } }),
        this.prisma.meal.count({ where: { userId, gluten: true } })
      ]),
      
      this.prisma.meal.findMany({
        where: { userId },
        orderBy: { mealTime: 'desc' },
        take: 10
      })
    ])

    return {
      totalMeals,
      mealsByCategory: categoryStats.map(stat => ({
        category: stat.category || 'Uncategorized',
        count: stat._count.category
      })),
      averageSpicyLevel: spicyAvg._avg.spicyLevel,
      dietaryDistribution: {
        fiberRich: dietaryStats[0],
        dairy: dietaryStats[1],
        gluten: dietaryStats[2]
      },
      recentMeals
    }
  }
}

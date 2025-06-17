import { Router, Response, NextFunction } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router: Router = Router()
const prisma = new PrismaClient()

// Apply authentication to all routes
router.use(authenticateToken)

// Validation schemas
const createEntrySchema = z.object({
  bristolType: z.number().int().min(1).max(7),
  volume: z.enum(['Small', 'Medium', 'Large', 'Massive']).optional(),
  color: z
    .enum(['Brown', 'Dark Brown', 'Light Brown', 'Yellow', 'Green', 'Red', 'Black'])
    .optional(),
  consistency: z.enum(['Solid', 'Soft', 'Loose', 'Watery']).optional(),
  floaters: z.boolean().default(false),
  pain: z.number().int().min(1).max(10).optional(),
  strain: z.number().int().min(1).max(10).optional(),
  satisfaction: z.number().int().min(1).max(10).optional(),
  notes: z.string().optional(),
  smell: z.enum(['None', 'Mild', 'Moderate', 'Strong', 'Toxic']).optional(),
  photoUrl: z.string().url().optional()
})

const updateEntrySchema = createEntrySchema.partial()

// GET /api/entries - Get all user entries
router.get(
  '/',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { page = '1', limit = '20', sortBy = 'createdAt', sortOrder = 'desc' } = req.query

      const pageNum = parseInt(page as string)
      const limitNum = parseInt(limit as string)
      const skip = (pageNum - 1) * limitNum

      // Ensure userId is defined
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const entries = await prisma.entry.findMany({
        where: { userId: req.userId },
        orderBy: { [sortBy as string]: sortOrder },
        skip,
        take: limitNum,
        include: {
          user: {
            select: { name: true, email: true }
          }
        }
      })

      const total = await prisma.entry.count({
        where: { userId: req.userId }
      })

      res.json({
        entries,
        pagination: {
          page: pageNum,
          limit: limitNum,
          total,
          pages: Math.ceil(total / limitNum)
        }
      })
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/entries/:id - Get specific entry
router.get(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const entry = await prisma.entry.findFirst({
        where: {
          id,
          userId: req.userId
        },
        include: {
          user: {
            select: { name: true, email: true }
          }
        }
      })

      if (!entry) {
        res.status(404).json({ error: 'Entry not found' })
        return
      }

      res.json(entry)
    } catch (error) {
      next(error)
    }
  }
)

// POST /api/entries - Create new entry
router.post(
  '/',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      console.log('Creating entry for userId:', req.userId)

      const validatedData = createEntrySchema.parse(req.body)

      // Convert undefined values to null for Prisma
      const entryData = {
        userId: req.userId,
        bristolType: validatedData.bristolType,
        floaters: validatedData.floaters,
        volume: validatedData.volume ?? null,
        color: validatedData.color ?? null,
        consistency: validatedData.consistency ?? null,
        pain: validatedData.pain ?? null,
        strain: validatedData.strain ?? null,
        satisfaction: validatedData.satisfaction ?? null,
        notes: validatedData.notes ?? null,
        smell: validatedData.smell ?? null,
        photoUrl: validatedData.photoUrl ?? null
      }

      console.log('Entry data to be created:', entryData)

      const entry = await prisma.entry.create({
        data: entryData
      })

      res.status(201).json(entry)
    } catch (error) {
      if (error instanceof z.ZodError) {
        res.status(400).json({ error: error.errors })
        return
      }
      next(error)
    }
  }
)

// PUT /api/entries/:id - Update entry
router.put(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const validatedData = updateEntrySchema.parse(req.body)

      const existingEntry = await prisma.entry.findFirst({
        where: { id, userId: req.userId }
      })

      if (!existingEntry) {
        res.status(404).json({ error: 'Entry not found' })
        return
      }

      // Convert undefined values to null for Prisma
      const updateData: any = {}
      if (validatedData.bristolType !== undefined)
        updateData.bristolType = validatedData.bristolType
      if (validatedData.volume !== undefined) updateData.volume = validatedData.volume ?? null
      if (validatedData.color !== undefined) updateData.color = validatedData.color ?? null
      if (validatedData.consistency !== undefined)
        updateData.consistency = validatedData.consistency ?? null
      if (validatedData.floaters !== undefined) updateData.floaters = validatedData.floaters
      if (validatedData.pain !== undefined) updateData.pain = validatedData.pain ?? null
      if (validatedData.strain !== undefined) updateData.strain = validatedData.strain ?? null
      if (validatedData.satisfaction !== undefined)
        updateData.satisfaction = validatedData.satisfaction ?? null
      if (validatedData.notes !== undefined) updateData.notes = validatedData.notes ?? null
      if (validatedData.smell !== undefined) updateData.smell = validatedData.smell ?? null
      if (validatedData.photoUrl !== undefined) updateData.photoUrl = validatedData.photoUrl ?? null

      const entry = await prisma.entry.update({
        where: { id },
        data: updateData
      })

      res.json(entry)
    } catch (error) {
      if (error instanceof z.ZodError) {
        res.status(400).json({ error: error.errors })
        return
      }
      next(error)
    }
  }
)

// DELETE /api/entries/:id - Delete entry
router.delete(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const existingEntry = await prisma.entry.findFirst({
        where: { id, userId: req.userId }
      })

      if (!existingEntry) {
        res.status(404).json({ error: 'Entry not found' })
        return
      }

      await prisma.entry.delete({
        where: { id }
      })

      res.status(204).send()
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/entries/:id/meals - Get all meals linked to an entry
router.get(
  '/:id/meals',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      // Verify the entry belongs to the user
      const entry = await prisma.entry.findFirst({
        where: { id, userId: req.userId }
      })

      if (!entry) {
        res.status(404).json({ error: 'Entry not found' })
        return
      }

      // Get all meals linked to this entry
      const linkedMeals = await prisma.meal.findMany({
        where: {
          userId: req.userId,
          entries: {
            some: { entryId: id }
          }
        },
        orderBy: { mealTime: 'desc' }
      })

      res.json(linkedMeals)
    } catch (error) {
      next(error)
    }
  }
)

export { router as entryRoutes }

import { Router, Response, NextFunction } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'
import { EntryService } from '../domains/entries/EntryService'
import { EntryFactory } from '../domains/entries/EntryFactory'
import type { CreateEntryRequest, UpdateEntryRequest, EntryFilters } from '../domains/entries/types'

const router: Router = Router()
const prisma = new PrismaClient()
const entryService = new EntryService(prisma)

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
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const filters: EntryFilters = {
        page: req.query.page ? parseInt(req.query.page as string) : 1,
        limit: req.query.limit ? parseInt(req.query.limit as string) : 20,
        sortBy: (req.query.sortBy as string) || 'createdAt',
        sortOrder: (req.query.sortOrder as 'asc' | 'desc') || 'desc'
      }

      // Add optional filters
      if (req.query.bristolType) {
        filters.bristolType = parseInt(req.query.bristolType as string)
      }
      if (req.query.dateFrom) {
        filters.dateFrom = new Date(req.query.dateFrom as string)
      }
      if (req.query.dateTo) {
        filters.dateTo = new Date(req.query.dateTo as string)
      }

      const result = await entryService.findByUserId(req.userId, filters)
      res.json(result)
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/entries/:id - Get single entry
router.get(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const entry = await entryService.findById(req.params.id!, req.userId)
      
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

      // Validate request body
      const validationResult = createEntrySchema.safeParse(req.body)
      if (!validationResult.success) {
        res.status(400).json({ 
          error: 'Validation failed', 
          details: validationResult.error.errors 
        })
        return
      }

      const createRequest: CreateEntryRequest = validationResult.data
      
      // Additional validation using factory
      if (!EntryFactory.validateBristolType(createRequest.bristolType)) {
        res.status(400).json({ error: 'Invalid Bristol stool type' })
        return
      }

      if (!EntryFactory.validateRating(createRequest.satisfaction)) {
        res.status(400).json({ error: 'Invalid satisfaction rating' })
        return
      }

      const entry = await entryService.create(createRequest, req.userId)
      res.status(201).json(entry)
    } catch (error) {
      next(error)
    }
  }
)

// PUT /api/entries/:id - Update entry
router.put(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      // Validate request body
      const validationResult = updateEntrySchema.safeParse(req.body)
      if (!validationResult.success) {
        res.status(400).json({ 
          error: 'Validation failed', 
          details: validationResult.error.errors 
        })
        return
      }

      const updateRequest: UpdateEntryRequest = validationResult.data

      // Additional validation using factory
      if (updateRequest.bristolType && !EntryFactory.validateBristolType(updateRequest.bristolType)) {
        res.status(400).json({ error: 'Invalid Bristol stool type' })
        return
      }

      if (!EntryFactory.validateRating(updateRequest.satisfaction)) {
        res.status(400).json({ error: 'Invalid satisfaction rating' })
        return
      }

      const entry = await entryService.update(req.params.id!, updateRequest, req.userId)
      
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

// DELETE /api/entries/:id - Delete entry
router.delete(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const deleted = await entryService.delete(req.params.id!, req.userId)
      
      if (!deleted) {
        res.status(404).json({ error: 'Entry not found' })
        return
      }

      res.status(204).send()
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/entries/analytics - Get analytics data
router.get(
  '/analytics',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const analytics = await entryService.getAnalytics(req.userId)
      res.json(analytics)
    } catch (error) {
      next(error)
    }
  }
)

export default router

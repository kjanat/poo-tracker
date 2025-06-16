import { Router, Response, NextFunction } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'
import { BowelMovementService } from '../domains/bowel-movements/BowelMovementService'
import type { CreateBowelMovementRequest, UpdateBowelMovementRequest, BowelMovementFilters } from '../domains/bowel-movements/types'

const router: Router = Router()
const prisma = new PrismaClient()
const bowelMovementService = new BowelMovementService(prisma)

// Apply authentication to all routes
router.use(authenticateToken)

// Validation schemas
const createBowelMovementSchema = z.object({
  bristolType: z.number().int().min(1).max(7),
  recordedAt: z.string().datetime().optional().transform(val => val ? new Date(val) : undefined),
  volume: z.enum(['SMALL', 'MEDIUM', 'LARGE', 'MASSIVE']).optional(),
  color: z
    .enum(['BROWN', 'DARK_BROWN', 'LIGHT_BROWN', 'YELLOW', 'GREEN', 'RED', 'BLACK'])
    .optional(),
  consistency: z.enum(['SOLID', 'SOFT', 'LOOSE', 'WATERY']).optional(),
  floaters: z.boolean().default(false),
  pain: z.number().int().min(1).max(10).default(1),
  strain: z.number().int().min(1).max(10).default(1),
  satisfaction: z.number().int().min(1).max(10).default(5),
  notes: z.string().optional(),
  smell: z.enum(['NONE', 'MILD', 'MODERATE', 'STRONG', 'TOXIC']).optional(),
  photoUrl: z.string().url().optional()
})

const updateBowelMovementSchema = createBowelMovementSchema.partial()

// GET /api/entries - Get all user bowel movements (keep old endpoint for backward compatibility)
router.get(
  '/',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const filters: BowelMovementFilters = {
        page: req.query.page ? parseInt(req.query.page as string) : 1,
        limit: req.query.limit ? parseInt(req.query.limit as string) : 20,
        sortBy: (req.query.sortBy as string) || 'createdAt',
        sortOrder: (req.query.sortOrder as 'asc' | 'desc') || 'desc'
      }

      // Optional filters
      if (req.query.bristolType) {
        filters.bristolType = parseInt(req.query.bristolType as string)
      }
      if (req.query.dateFrom) {
        filters.dateFrom = new Date(req.query.dateFrom as string)
      }
      if (req.query.dateTo) {
        filters.dateTo = new Date(req.query.dateTo as string)
      }

      const result = await bowelMovementService.findByUserId(req.userId, filters)
      res.json(result)
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/entries/:id - Get specific bowel movement
router.get(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId || !req.params.id) {
        res.status(401).json({ error: 'User not authenticated or invalid request' })
        return
      }

      const bowelMovement = await bowelMovementService.findById(req.params.id, req.userId)
      if (!bowelMovement) {
        res.status(404).json({ error: 'Bowel movement not found' })
        return
      }

      res.json(bowelMovement)
    } catch (error) {
      next(error)
    }
  }
)

// POST /api/entries - Create new bowel movement
router.post(
  '/',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const validationResult = createBowelMovementSchema.safeParse(req.body)
      if (!validationResult.success) {
        res.status(400).json({
          error: 'Invalid request data',
          details: validationResult.error.errors
        })
        return
      }

      const request: CreateBowelMovementRequest = {
        ...validationResult.data,
        recordedAt: validationResult.data.recordedAt ?? new Date()
      }
      const bowelMovement = await bowelMovementService.create(request, req.userId)

      res.status(201).json(bowelMovement)
    } catch (error) {
      next(error)
    }
  }
)

// PUT /api/entries/:id - Update bowel movement
router.put(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId || !req.params.id) {
        res.status(401).json({ error: 'User not authenticated or invalid request' })
        return
      }

      const validationResult = updateBowelMovementSchema.safeParse(req.body)
      if (!validationResult.success) {
        res.status(400).json({
          error: 'Invalid request data',
          details: validationResult.error.errors
        })
        return
      }

      const request = validationResult.data as UpdateBowelMovementRequest
      const bowelMovement = await bowelMovementService.update(req.params.id, request, req.userId)

      if (!bowelMovement) {
        res.status(404).json({ error: 'Bowel movement not found' })
        return
      }

      res.json(bowelMovement)
    } catch (error) {
      next(error)
    }
  }
)

// DELETE /api/entries/:id - Delete bowel movement
router.delete(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId || !req.params.id) {
        res.status(401).json({ error: 'User not authenticated or invalid request' })
        return
      }

      const success = await bowelMovementService.delete(req.params.id, req.userId)
      if (!success) {
        res.status(404).json({ error: 'Bowel movement not found' })
        return
      }

      res.status(204).send()
    } catch (error) {
      next(error)
    }
  }
)

export default router

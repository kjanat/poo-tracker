import { Router, Response, NextFunction } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'
import { MealService } from '../domains/meals/MealService'
import { MealFactory } from '../domains/meals/MealFactory'
import type { CreateMealRequest, UpdateMealRequest } from '../domains/meals/types'

const router: Router = Router()
const prisma = new PrismaClient()
const mealService = new MealService(prisma)

// Apply authentication to all routes
router.use(authenticateToken)

// Validation schemas
const createMealSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  description: z.string().optional(),
  mealTime: z.string().datetime(),
  category: z.enum(['Breakfast', 'Lunch', 'Dinner', 'Snack']).optional(),
  cuisine: z.string().optional(),
  spicyLevel: z.number().int().min(1).max(10).optional(),
  fiberRich: z.boolean().default(false),
  dairy: z.boolean().default(false),
  gluten: z.boolean().default(false),
  notes: z.string().optional(),
  photoUrl: z.string().url().optional()
})

const updateMealSchema = createMealSchema.partial()

// GET /api/meals - Get all user meals
router.get(
  '/',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const page = parseInt(req.query.page as string) || 1
      const limit = parseInt(req.query.limit as string) || 20

      const result = await mealService.findByUserId(req.userId, { page, limit })
      res.json(result)
    } catch (error) {
      next(error)
    }
  }
)

// POST /api/meals - Create new meal
router.post(
  '/',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const validationResult = createMealSchema.safeParse(req.body)
      if (!validationResult.success) {
        res.status(400).json({
          error: 'Validation failed',
          details: validationResult.error.errors
        })
        return
      }

      const createRequest: CreateMealRequest = {
        ...validationResult.data,
        mealTime: new Date(validationResult.data.mealTime)
      }

      // Additional business validation using factory
      if (!MealFactory.validateMealName(createRequest.name)) {
        res.status(400).json({ error: 'Invalid meal name' })
        return
      }

      if (createRequest.spicyLevel !== undefined && createRequest.spicyLevel !== null && !MealFactory.validateSpicyLevel(createRequest.spicyLevel)) {
        res.status(400).json({ error: 'Invalid spicy level' })
        return
      }

      const meal = await mealService.create(createRequest, req.userId)
      res.status(201).json(meal)
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/meals/:id - Get specific meal
router.get(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const meal = await mealService.findById(id, req.userId)

      if (meal == null) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      res.json(meal)
    } catch (error) {
      next(error)
    }
  }
)

// PUT /api/meals/:id - Update meal
router.put(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const validationResult = updateMealSchema.safeParse(req.body)
      if (!validationResult.success) {
        res.status(400).json({
          error: 'Validation failed',
          details: validationResult.error.errors
        })
        return
      }

      const updateRequest = {
        ...validationResult.data,
        mealTime: validationResult.data.mealTime ? new Date(validationResult.data.mealTime) : undefined
      } as UpdateMealRequest

      // Additional business validation using factory
      if (updateRequest.name !== undefined && !MealFactory.validateMealName(updateRequest.name)) {
        res.status(400).json({ error: 'Invalid meal name' })
        return
      }

      if (updateRequest.spicyLevel !== undefined && updateRequest.spicyLevel !== null && !MealFactory.validateSpicyLevel(updateRequest.spicyLevel)) {
        res.status(400).json({ error: 'Invalid spicy level' })
        return
      }

      const meal = await mealService.update(id, updateRequest, req.userId)

      if (meal == null) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      res.json(meal)
    } catch (error) {
      next(error)
    }
  }
)

// DELETE /api/meals/:id - Delete meal
router.delete(
  '/:id',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const success = await mealService.delete(id, req.userId)

      if (!success) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      res.json({ message: 'Meal deleted successfully' })
    } catch (error) {
      next(error)
    }
  }
)

// POST /api/meals/:id/link-entry - Link an entry to a meal
router.post(
  '/:id/link-entry',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params
      const { entryId } = req.body

      if (!req.userId || !id || !entryId) {
        res.status(400).json({ error: 'Invalid request - meal ID and entry ID required' })
        return
      }

      const success = await mealService.linkEntry(id, entryId, req.userId)

      if (!success) {
        res.status(400).json({ error: 'Unable to link entry to meal - meal/entry not found or already linked' })
        return
      }

      res.status(201).json({ message: 'Entry linked to meal successfully' })
    } catch (error) {
      next(error)
    }
  }
)

// DELETE /api/meals/:id/unlink-entry - Unlink an entry from a meal
router.delete(
  '/:id/unlink-entry',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params
      const { entryId } = req.body

      if (!req.userId || !id || !entryId) {
        res.status(400).json({ error: 'Invalid request - meal ID and entry ID required' })
        return
      }

      const success = await mealService.unlinkEntry(id, entryId, req.userId)

      if (!success) {
        res.status(400).json({ error: 'Unable to unlink entry from meal - meal not found or entry not linked' })
        return
      }

      res.json({ message: 'Entry unlinked from meal successfully' })
    } catch (error) {
      next(error)
    }
  }
)

// GET /api/meals/:id/entries - Get all entries linked to a meal
router.get(
  '/:id/entries',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const { id } = req.params

      if (!req.userId || !id) {
        res.status(400).json({ error: 'Invalid request' })
        return
      }

      const linkedEntries = await mealService.getLinkedEntries(id, req.userId)

      res.json(linkedEntries)
    } catch (error) {
      next(error)
    }
  }
)

export { router as mealRoutes }

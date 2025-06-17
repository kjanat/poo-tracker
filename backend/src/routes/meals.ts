import { Router, Response, NextFunction } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router: Router = Router()
const prisma = new PrismaClient()

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

      const meals = await prisma.meal.findMany({
        where: { userId: req.userId },
        orderBy: { mealTime: 'desc' }
      })

      res.json(meals)
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

      const validatedData = createMealSchema.parse(req.body)

      // Convert undefined values to null for Prisma
      const mealData = {
        name: validatedData.name,
        description: validatedData.description ?? null,
        mealTime: new Date(validatedData.mealTime),
        category: validatedData.category ?? null,
        cuisine: validatedData.cuisine ?? null,
        spicyLevel: validatedData.spicyLevel ?? null,
        fiberRich: validatedData.fiberRich,
        dairy: validatedData.dairy,
        gluten: validatedData.gluten,
        notes: validatedData.notes ?? null,
        photoUrl: validatedData.photoUrl ?? null,
        userId: req.userId
      }

      const meal = await prisma.meal.create({
        data: mealData
      })

      res.status(201).json(meal)
    } catch (error) {
      if (error instanceof z.ZodError) {
        res.status(400).json({ error: error.errors })
        return
      }
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

      const meal = await prisma.meal.findFirst({
        where: {
          id,
          userId: req.userId
        }
      })

      if (!meal) {
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

      const validatedData = updateMealSchema.parse(req.body)

      const existingMeal = await prisma.meal.findFirst({
        where: { id, userId: req.userId }
      })

      if (!existingMeal) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      // Convert undefined values to null for Prisma
      const updateData: any = {}
      if (validatedData.name !== undefined) updateData.name = validatedData.name
      if (validatedData.description !== undefined)
        updateData.description = validatedData.description ?? null
      if (validatedData.mealTime !== undefined)
        updateData.mealTime = new Date(validatedData.mealTime)
      if (validatedData.category !== undefined) updateData.category = validatedData.category ?? null
      if (validatedData.cuisine !== undefined) updateData.cuisine = validatedData.cuisine ?? null
      if (validatedData.spicyLevel !== undefined)
        updateData.spicyLevel = validatedData.spicyLevel ?? null
      if (validatedData.fiberRich !== undefined) updateData.fiberRich = validatedData.fiberRich
      if (validatedData.dairy !== undefined) updateData.dairy = validatedData.dairy
      if (validatedData.gluten !== undefined) updateData.gluten = validatedData.gluten
      if (validatedData.notes !== undefined) updateData.notes = validatedData.notes ?? null

      const meal = await prisma.meal.update({
        where: { id },
        data: updateData
      })

      res.json(meal)
    } catch (error) {
      if (error instanceof z.ZodError) {
        res.status(400).json({ error: error.errors })
        return
      }
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

      const existingMeal = await prisma.meal.findFirst({
        where: { id, userId: req.userId }
      })

      if (!existingMeal) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      await prisma.meal.delete({
        where: { id }
      })

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

      // Verify the meal belongs to the user
      const meal = await prisma.meal.findFirst({
        where: { id, userId: req.userId }
      })

      if (!meal) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      // Verify the entry belongs to the user
      const entry = await prisma.entry.findFirst({
        where: { id: entryId, userId: req.userId }
      })

      if (!entry) {
        res.status(404).json({ error: 'Entry not found' })
        return
      }

      // Check if the relation already exists
      const existingRelation = await prisma.mealEntryRelation.findFirst({
        where: { mealId: id, entryId }
      })

      if (existingRelation) {
        res.status(409).json({ error: 'Entry is already linked to this meal' })
        return
      }

      // Create the relation
      const relation = await prisma.mealEntryRelation.create({
        data: { mealId: id, entryId }
      })

      res.status(201).json({ message: 'Entry linked to meal successfully', relation })
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

      // Verify the meal belongs to the user
      const meal = await prisma.meal.findFirst({
        where: { id, userId: req.userId }
      })

      if (!meal) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      // Find and delete the relation
      const relation = await prisma.mealEntryRelation.findFirst({
        where: { mealId: id, entryId }
      })

      if (!relation) {
        res.status(404).json({ error: 'Entry is not linked to this meal' })
        return
      }

      await prisma.mealEntryRelation.delete({
        where: { id: relation.id }
      })

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

      // Verify the meal belongs to the user
      const meal = await prisma.meal.findFirst({
        where: { id, userId: req.userId }
      })

      if (!meal) {
        res.status(404).json({ error: 'Meal not found' })
        return
      }

      // Get all entries linked to this meal
      const linkedEntries = await prisma.entry.findMany({
        where: {
          userId: req.userId,
          meals: {
            some: { mealId: id }
          }
        },
        orderBy: { createdAt: 'desc' }
      })

      res.json(linkedEntries)
    } catch (error) {
      next(error)
    }
  }
)

export { router as mealRoutes }

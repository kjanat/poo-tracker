import { Router } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router = Router()
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
  notes: z.string().optional()
})

// GET /api/meals - Get all user meals
router.get('/', async (req: AuthenticatedRequest, res, next) => {
  try {
    const meals = await prisma.meal.findMany({
      where: { userId: req.userId },
      orderBy: { mealTime: 'desc' }
    })
    
    res.json(meals)
  } catch (error) {
    next(error)
  }
})

// POST /api/meals - Create new meal
router.post('/', async (req: AuthenticatedRequest, res, next) => {
  try {
    const validatedData = createMealSchema.parse(req.body)
    
    const meal = await prisma.meal.create({
      data: {
        ...validatedData,
        mealTime: new Date(validatedData.mealTime),
        userId: req.userId!
      }
    })
    
    res.status(201).json(meal)
  } catch (error) {
    if (error instanceof z.ZodError) {
      return res.status(400).json({ error: error.errors })
    }
    next(error)
  }
})

export { router as mealRoutes }

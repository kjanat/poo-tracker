import { Router } from 'express'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router = Router()
const prisma = new PrismaClient()

// Apply authentication to all routes
router.use(authenticateToken)

// Validation schemas
const createEntrySchema = z.object({
  bristolType: z.number().int().min(1).max(7),
  volume: z.enum(['Small', 'Medium', 'Large', 'Massive']).optional(),
  color: z.enum(['Brown', 'Dark Brown', 'Light Brown', 'Yellow', 'Green', 'Red', 'Black']).optional(),
  consistency: z.enum(['Solid', 'Soft', 'Loose', 'Watery']).optional(),
  floaters: z.boolean().default(false),
  pain: z.number().int().min(1).max(10).optional(),
  strain: z.number().int().min(1).max(10).optional(),
  satisfaction: z.number().int().min(1).max(10).optional(),
  notes: z.string().optional(),
  smell: z.enum(['None', 'Mild', 'Moderate', 'Strong', 'Toxic']).optional()
})

const updateEntrySchema = createEntrySchema.partial()

// GET /api/entries - Get all user entries
router.get('/', async (req: AuthenticatedRequest, res, next) => {
  try {
    const { page = '1', limit = '20', sortBy = 'createdAt', sortOrder = 'desc' } = req.query
    
    const pageNum = parseInt(page as string)
    const limitNum = parseInt(limit as string)
    const skip = (pageNum - 1) * limitNum
    
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
})

// GET /api/entries/:id - Get specific entry
router.get('/:id', async (req: AuthenticatedRequest, res, next) => {
  try {
    const { id } = req.params
    
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
      return res.status(404).json({ error: 'Entry not found' })
    }
    
    res.json(entry)
  } catch (error) {
    next(error)
  }
})

// POST /api/entries - Create new entry
router.post('/', async (req: AuthenticatedRequest, res, next) => {
  try {
    const validatedData = createEntrySchema.parse(req.body)
    
    const entry = await prisma.entry.create({
      data: {
        ...validatedData,
        userId: req.userId!
      }
    })
    
    res.status(201).json(entry)
  } catch (error) {
    if (error instanceof z.ZodError) {
      return res.status(400).json({ error: error.errors })
    }
    next(error)
  }
})

// PUT /api/entries/:id - Update entry
router.put('/:id', async (req: AuthenticatedRequest, res, next) => {
  try {
    const { id } = req.params
    const validatedData = updateEntrySchema.parse(req.body)
    
    const existingEntry = await prisma.entry.findFirst({
      where: { id, userId: req.userId }
    })
    
    if (!existingEntry) {
      return res.status(404).json({ error: 'Entry not found' })
    }
    
    const entry = await prisma.entry.update({
      where: { id },
      data: validatedData
    })
    
    res.json(entry)
  } catch (error) {
    if (error instanceof z.ZodError) {
      return res.status(400).json({ error: error.errors })
    }
    next(error)
  }
})

// DELETE /api/entries/:id - Delete entry
router.delete('/:id', async (req: AuthenticatedRequest, res, next) => {
  try {
    const { id } = req.params
    
    const existingEntry = await prisma.entry.findFirst({
      where: { id, userId: req.userId }
    })
    
    if (!existingEntry) {
      return res.status(404).json({ error: 'Entry not found' })
    }
    
    await prisma.entry.delete({
      where: { id }
    })
    
    res.status(204).send()
  } catch (error) {
    next(error)
  }
})

export { router as entryRoutes }

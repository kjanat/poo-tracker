import { Router, Response, NextFunction } from 'express'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router: Router = Router()
const prisma = new PrismaClient()

// Apply authentication to all routes
router.use(authenticateToken)

// GET /api/analytics/summary - Get user's poop summary
router.get(
  '/summary',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      const userId = req.userId

      // Ensure userId is defined
      if (userId === undefined) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      // Get basic stats
      const totalEntries = await prisma.entry.count({
        where: { userId }
      })

      // Get bristol type distribution
      const bristolDistribution = await prisma.entry.groupBy({
        by: ['bristolType'],
        where: { userId },
        _count: true,
        orderBy: { bristolType: 'asc' }
      })

      // Get recent entries
      const recentEntries = await prisma.entry.findMany({
        where: { userId },
        orderBy: { createdAt: 'desc' },
        take: 5,
        select: {
          id: true,
          bristolType: true,
          createdAt: true,
          satisfaction: true
        }
      })

      // Calculate average satisfaction if available
      const avgSatisfaction = await prisma.entry.aggregate({
        where: {
          userId,
          satisfaction: { not: null }
        },
        _avg: {
          satisfaction: true
        }
      })

      res.json({
        totalEntries,
        bristolDistribution: bristolDistribution.map((item) => ({
          type: item.bristolType,
          count: item._count
        })),
        recentEntries,
        averageSatisfaction: avgSatisfaction._avg?.satisfaction ?? null
      })
    } catch (error) {
      next(error)
    }
  }
)

export { router as analyticsRoutes }

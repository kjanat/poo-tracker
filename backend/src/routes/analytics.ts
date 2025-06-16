import { Router, Response, NextFunction } from 'express'
import { PrismaClient } from '@prisma/client'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'
import { BowelMovementService } from '../domains/bowel-movements/BowelMovementService'

const router: Router = Router()
const prisma = new PrismaClient()
const bowelMovementService = new BowelMovementService(prisma)

// Apply authentication to all routes
router.use(authenticateToken)

// GET /api/analytics/summary - Get user's bowel movement summary
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

      const analytics = await bowelMovementService.getAnalytics(userId)
      res.json(analytics)
    } catch (error) {
      next(error)
    }
  }
)

export { router as analyticsRoutes }

import { Router } from 'express'
import bowelMovementsRouter from './bowel-movements'

const router: Router = Router()

// Redirect all entry routes to bowel-movements for backward compatibility
router.use('/', bowelMovementsRouter)

export default router

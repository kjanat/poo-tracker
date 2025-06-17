import { Request, Response, NextFunction } from 'express'
import jwt from 'jsonwebtoken'
import { config } from '../config'
import { logger } from '../utils/logger'

export interface AuthenticatedRequest extends Request {
  userId?: string
}

export const authenticateToken = (
  req: AuthenticatedRequest,
  res: Response,
  next: NextFunction
): void => {
  const authHeader = req.headers.authorization
  const token = authHeader?.split(' ')[1] // Bearer TOKEN

  logger.debug(
    `🔐 Auth middleware - Headers: ${typeof authHeader === 'string' ? 'Present' : 'Missing'}`
  )

  if (token === undefined || token === '') {
    logger.warn('❌ Auth middleware - No token provided')
    res.status(401).json({ error: 'Access token required' })
    return
  }

  logger.debug(`🎫 Auth middleware - Token received: ${token.substring(0, 20)}...`)

  try {
    const decoded = jwt.verify(token, config.jwt.secret) as { userId: string }
    logger.info(`✅ Auth middleware - Token valid, userId: ${decoded.userId}`)
    req.userId = decoded.userId
    next()
  } catch (error) {
    logger.error(
      `❌ Auth middleware - Token invalid: ${
        error instanceof Error ? error.message : 'Unknown error'
      }`
    )
    res.status(403).json({ error: 'Invalid or expired token' })
  }
}

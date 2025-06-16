import { Request, Response, NextFunction } from 'express'
import jwt from 'jsonwebtoken'
import { config } from '../config'

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

  console.log(
    '🔐 Auth middleware - Headers:',
    typeof authHeader === 'string' ? 'Present' : 'Missing'
  )

  if (token === undefined || token === '') {
    console.log('❌ Auth middleware - No token provided')
    res.status(401).json({ error: 'Access token required' })
    return
  }

  console.log('🎫 Auth middleware - Token received:', token.substring(0, 20) + '...')

  try {
    const decoded = jwt.verify(token, config.jwt.secret) as { userId: string }
    console.log('✅ Auth middleware - Token valid, userId:', decoded.userId)
    req.userId = decoded.userId
    next()
  } catch (error) {
    console.log(
      '❌ Auth middleware - Token invalid:',
      error instanceof Error ? error.message : 'Unknown error'
    )
    res.status(403).json({ error: 'Invalid or expired token' })
  }
}

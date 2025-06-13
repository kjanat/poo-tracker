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

  if (!token) {
    res.status(401).json({ error: 'Access token required' })
    return
  }

  try {
    const decoded = jwt.verify(token, config.jwt.secret) as { userId: string }
    req.userId = decoded.userId
    next()
  } catch (error) {
    res.status(403).json({ error: 'Invalid or expired token' })
  }
}

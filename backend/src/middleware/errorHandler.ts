import { Request, Response, NextFunction } from 'express'

export interface ApiError extends Error {
  statusCode?: number
}

export const errorHandler = (
  err: ApiError,
  _req: Request,
  res: Response,
  // Next function is required by Express but not used
  _next: NextFunction
): void => {
  const statusCode = err.statusCode ?? 500
  const message = (err.message ?? '') !== '' ? err.message : 'Internal Server Error'

  // Log error details (but not in production)
  if (process.env.NODE_ENV !== 'production') {
    console.error('Error:', {
      message: err.message,
      stack: err.stack,
      statusCode
    })
  }

  res.status(statusCode).json({
    error: message,
    ...(process.env.NODE_ENV === 'development' && { stack: err.stack })
  })
}

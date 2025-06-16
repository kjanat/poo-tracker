import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { Request, Response, NextFunction } from 'express'
import { errorHandler, ApiError } from '../errorHandler'

describe('errorHandler middleware', () => {
  let mockReq: Partial<Request>
  let mockRes: Partial<Response>
  let mockNext: NextFunction
  let originalNodeEnv: string | undefined

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks()

    // Setup mock objects
    mockReq = {}
    mockRes = {
      status: vi.fn().mockReturnThis(),
      json: vi.fn().mockReturnThis()
    }
    mockNext = vi.fn()

    // Mock console.error to avoid noise in test output

    vi.spyOn(console, 'error').mockImplementation(() => {})

    // Store original NODE_ENV
    originalNodeEnv = process.env.NODE_ENV
  })

  afterEach(() => {
    vi.restoreAllMocks()
    // Restore original NODE_ENV
    process.env.NODE_ENV = originalNodeEnv
  })

  describe('basic error handling', () => {
    it('should handle error with custom status code and message', () => {
      // Setup
      const error: ApiError = new Error('Custom error message')
      error.statusCode = 400

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(400)
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Custom error message'
      })
    })

    it('should default to 500 status code when none provided', () => {
      // Setup
      const error: ApiError = new Error('Server error')

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(500)
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Server error'
      })
    })

    it('should default to "Internal Server Error" message when none provided', () => {
      // Setup
      const error: ApiError = new Error('')
      error.statusCode = 500

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(500)
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Internal Server Error'
      })
    })

    it('should handle error without message property', () => {
      // Setup
      const error: ApiError = { name: 'Error', message: '' } as ApiError
      error.statusCode = 422

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(422)
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Internal Server Error'
      })
    })
  })

  describe('environment-specific behavior', () => {
    it('should include stack trace in development environment', () => {
      // Setup
      process.env.NODE_ENV = 'development'
      const error: ApiError = new Error('Development error')
      error.statusCode = 400
      error.stack = 'Error stack trace...'

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Development error',
        stack: 'Error stack trace...'
      })
    })

    it('should not include stack trace in production environment', () => {
      // Setup
      process.env.NODE_ENV = 'production'
      const error: ApiError = new Error('Production error')
      error.statusCode = 400
      error.stack = 'Error stack trace...'

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Production error'
      })
    })

    it('should not include stack trace in test environment', () => {
      // Setup
      process.env.NODE_ENV = 'test'
      const error: ApiError = new Error('Test error')
      error.statusCode = 400
      error.stack = 'Error stack trace...'

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Test error'
      })
    })
  })

  describe('logging behavior', () => {
    it('should log error details in non-production environment', () => {
      // Setup
      process.env.NODE_ENV = 'development'
      const consoleSpy = vi.spyOn(console, 'error')
      const error: ApiError = new Error('Logged error')
      error.statusCode = 404
      error.stack = 'Stack trace...'

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(consoleSpy).toHaveBeenCalledWith('Error:', {
        message: 'Logged error',
        stack: 'Stack trace...',
        statusCode: 404
      })

      consoleSpy.mockRestore()
    })

    it('should not log error details in production environment', () => {
      // Setup
      process.env.NODE_ENV = 'production'
      const consoleSpy = vi.spyOn(console, 'error')
      const error: ApiError = new Error('Production error')
      error.statusCode = 500

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(consoleSpy).not.toHaveBeenCalled()

      consoleSpy.mockRestore()
    })

    it('should log in test environment', () => {
      // Setup
      process.env.NODE_ENV = 'test'
      const consoleSpy = vi.spyOn(console, 'error')
      const error: ApiError = new Error('Test error')

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(consoleSpy).toHaveBeenCalled()

      consoleSpy.mockRestore()
    })
  })

  describe('edge cases', () => {
    it('should handle undefined error message', () => {
      // Setup
      const error: any = { statusCode: 400 }

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Internal Server Error'
      })
    })

    it('should handle null error message', () => {
      // Setup
      const error: any = { message: null, statusCode: 400 }

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.json).toHaveBeenCalledWith({
        error: 'Internal Server Error'
      })
    })

    it('should handle error with statusCode 0 as-is', () => {
      // Setup
      const error: ApiError = new Error('Zero status')
      error.statusCode = 0

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify - actually uses the provided status code
      expect(mockRes.status).toHaveBeenCalledWith(0)
    })

    it('should handle negative status code as-is', () => {
      // Setup
      const error: ApiError = new Error('Negative status')
      error.statusCode = -1

      // Execute
      errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

      // Verify - actually uses the provided status code
      expect(mockRes.status).toHaveBeenCalledWith(-1)
    })
  })

  describe('common HTTP status codes', () => {
    const testCases = [
      { code: 400, message: 'Bad Request' },
      { code: 401, message: 'Unauthorized' },
      { code: 403, message: 'Forbidden' },
      { code: 404, message: 'Not Found' },
      { code: 422, message: 'Unprocessable Entity' },
      { code: 500, message: 'Internal Server Error' }
    ]

    testCases.forEach(({ code, message }) => {
      it(`should handle ${code} status code correctly`, () => {
        // Setup
        const error: ApiError = new Error(message)
        error.statusCode = code

        // Execute
        errorHandler(error, mockReq as Request, mockRes as Response, mockNext)

        // Verify
        expect(mockRes.status).toHaveBeenCalledWith(code)
        expect(mockRes.json).toHaveBeenCalledWith({
          error: message
        })
      })
    })
  })
})

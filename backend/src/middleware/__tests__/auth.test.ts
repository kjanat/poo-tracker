import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { Response, NextFunction } from 'express'
import jwt from 'jsonwebtoken'
import { authenticateToken, AuthenticatedRequest } from '../auth'
import { config } from '../../config'

// Mock dependencies
vi.mock('jsonwebtoken')
vi.mock('../../config', () => ({
  config: {
    jwt: {
      secret: 'test-secret-key'
    }
  }
}))

describe('authenticateToken middleware', () => {
  let mockReq: Partial<AuthenticatedRequest>
  let mockRes: Partial<Response>
  let mockNext: NextFunction
  const mockJwt = vi.mocked(jwt)

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks()

    // Setup mock request
    mockReq = {
      headers: {}
    }

    // Setup mock response with chainable methods
    mockRes = {
      status: vi.fn().mockReturnThis(),
      json: vi.fn().mockReturnThis()
    }

    // Setup mock next function
    mockNext = vi.fn()

    // Mock console methods to avoid noise in test output
    vi.spyOn(console, 'log').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('when no authorization header is provided', () => {
    it('should return 401 with error message', () => {
      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(401)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Access token required' })
      expect(mockNext).not.toHaveBeenCalled()
      expect(mockReq.userId).toBeUndefined()
    })
  })

  describe('when authorization header exists but no token', () => {
    it('should return 401 when only "Bearer" without token', () => {
      // Setup
      mockReq.headers = { authorization: 'Bearer' }

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(401)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Access token required' })
      expect(mockNext).not.toHaveBeenCalled()
    })

    it('should return 403 when authorization header has wrong format', () => {
      // Setup - "InvalidFormat token123" will split to ["InvalidFormat", "token123"]
      // So token will be "token123" which is invalid when verified
      mockReq.headers = { authorization: 'InvalidFormat token123' }

      const tokenError = new Error('invalid token')
      mockJwt.verify.mockImplementation(() => {
        throw tokenError
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify - this gets treated as invalid token since jwt.verify is called
      expect(mockRes.status).toHaveBeenCalledWith(403)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Invalid or expired token' })
      expect(mockNext).not.toHaveBeenCalled()
    })
  })

  describe('when valid token is provided', () => {
    it('should authenticate successfully and call next()', () => {
      // Setup
      const token = 'valid-jwt-token'
      const userId = 'user-123'
      mockReq.headers = { authorization: `Bearer ${token}` }

      mockJwt.verify.mockReturnValue({ userId } as any)

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockJwt.verify).toHaveBeenCalledWith(token, config.jwt.secret)
      expect(mockReq.userId).toBe(userId)
      expect(mockNext).toHaveBeenCalled()
      expect(mockRes.status).not.toHaveBeenCalled()
      expect(mockRes.json).not.toHaveBeenCalled()
    })

    it('should not handle token with extra spaces in authorization header', () => {
      // Setup - The current implementation doesn't trim the authorization header
      // So "  Bearer   token  " splits to ["", "", "Bearer", "", "", "token", "", ""]
      // The token becomes "" which is falsy
      const token = 'valid-jwt-token'
      const userId = 'user-456'
      mockReq.headers = { authorization: `  Bearer   ${token}  ` }

      mockJwt.verify.mockReturnValue({ userId } as any)

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify - should fail because of space handling
      expect(mockRes.status).toHaveBeenCalledWith(401)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Access token required' })
      expect(mockNext).not.toHaveBeenCalled()
    })
  })

  describe('when invalid token is provided', () => {
    it('should return 403 for expired token', () => {
      // Setup
      const token = 'expired-token'
      mockReq.headers = { authorization: `Bearer ${token}` }

      const tokenError = new Error('jwt expired')
      mockJwt.verify.mockImplementation(() => {
        throw tokenError
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockJwt.verify).toHaveBeenCalledWith(token, config.jwt.secret)
      expect(mockRes.status).toHaveBeenCalledWith(403)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Invalid or expired token' })
      expect(mockNext).not.toHaveBeenCalled()
      expect(mockReq.userId).toBeUndefined()
    })

    it('should return 403 for malformed token', () => {
      // Setup
      const token = 'malformed.token'
      mockReq.headers = { authorization: `Bearer ${token}` }

      const tokenError = new Error('invalid token')
      mockJwt.verify.mockImplementation(() => {
        throw tokenError
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(403)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Invalid or expired token' })
      expect(mockNext).not.toHaveBeenCalled()
    })

    it('should return 403 for token with wrong signature', () => {
      // Setup
      const token = 'wrong-signature-token'
      mockReq.headers = { authorization: `Bearer ${token}` }

      const tokenError = new Error('invalid signature')
      mockJwt.verify.mockImplementation(() => {
        throw tokenError
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(403)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Invalid or expired token' })
      expect(mockNext).not.toHaveBeenCalled()
    })

    it('should handle non-Error exceptions gracefully', () => {
      // Setup
      const token = 'problematic-token'
      mockReq.headers = { authorization: `Bearer ${token}` }

      // Throw a non-Error object
      mockJwt.verify.mockImplementation(() => {
        throw 'String error'
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(403)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Invalid or expired token' })
      expect(mockNext).not.toHaveBeenCalled()
    })
  })

  describe('edge cases', () => {
    it('should handle empty authorization header', () => {
      // Setup
      mockReq.headers = { authorization: '' }

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(401)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Access token required' })
    })

    it('should handle authorization header with only spaces', () => {
      // Setup
      mockReq.headers = { authorization: '   ' }

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify
      expect(mockRes.status).toHaveBeenCalledWith(401)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Access token required' })
    })

    it('should handle case-insensitive Bearer token as invalid', () => {
      // Setup - "bearer token" splits to ["bearer", "token"]
      // The token "token" will be passed to jwt.verify and fail
      const token = 'valid-token'
      mockReq.headers = { authorization: `bearer ${token}` }

      const tokenError = new Error('invalid token')
      mockJwt.verify.mockImplementation(() => {
        throw tokenError
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify - this should fail with 403 since jwt.verify is called
      expect(mockRes.status).toHaveBeenCalledWith(403)
      expect(mockRes.json).toHaveBeenCalledWith({ error: 'Invalid or expired token' })
    })
  })

  describe('logging behavior', () => {
    it('should log appropriate messages for successful authentication', () => {
      // Setup
      const consoleSpy = vi.spyOn(console, 'log')
      const token = 'valid-jwt-token'
      const userId = 'user-123'
      mockReq.headers = { authorization: `Bearer ${token}` }

      mockJwt.verify.mockReturnValue({ userId } as any)

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify logging
      expect(consoleSpy).toHaveBeenCalledWith('🔐 Auth middleware - Headers:', 'Present')
      expect(consoleSpy).toHaveBeenCalledWith(
        '🎫 Auth middleware - Token received:',
        'valid-jwt-token...'
      )
      expect(consoleSpy).toHaveBeenCalledWith('✅ Auth middleware - Token valid, userId:', userId)

      consoleSpy.mockRestore()
    })

    it('should log appropriate messages for missing token', () => {
      // Setup
      const consoleSpy = vi.spyOn(console, 'log')

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify logging
      expect(consoleSpy).toHaveBeenCalledWith('🔐 Auth middleware - Headers:', 'Missing')
      expect(consoleSpy).toHaveBeenCalledWith('❌ Auth middleware - No token provided')

      consoleSpy.mockRestore()
    })

    it('should log appropriate messages for invalid token', () => {
      // Setup
      const consoleSpy = vi.spyOn(console, 'log')
      const token = 'invalid-token'
      mockReq.headers = { authorization: `Bearer ${token}` }

      const tokenError = new Error('jwt malformed')
      mockJwt.verify.mockImplementation(() => {
        throw tokenError
      })

      // Execute
      authenticateToken(mockReq as AuthenticatedRequest, mockRes as Response, mockNext)

      // Verify logging
      expect(consoleSpy).toHaveBeenCalledWith(
        '❌ Auth middleware - Token invalid:',
        'jwt malformed'
      )

      consoleSpy.mockRestore()
    })
  })
})

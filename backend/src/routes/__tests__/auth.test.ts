import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import request from 'supertest'
import express from 'express'
import bcrypt from 'bcryptjs'
import jwt from 'jsonwebtoken'
import { authRoutes } from '../auth'
import { config } from '../../config'
import { errorHandler } from '../../middleware/errorHandler'

// Mock Prisma Client using vi.hoisted for proper scoping
const mockPrismaClient = vi.hoisted(() => ({
  user: {
    findUnique: vi.fn(),
    create: vi.fn()
  },
  userAuth: {
    update: vi.fn()
  }
}))

vi.mock('@prisma/client', () => ({
  PrismaClient: vi.fn().mockImplementation(() => mockPrismaClient)
}))

vi.mock('../../config', () => ({
  config: {
    jwt: {
      secret: 'test-secret-key',
      expiresIn: '7d'
    }
  }
}))

describe('Auth Routes Integration Tests', () => {
  let app: express.Application

  beforeEach(() => {
    // Create Express app for testing
    app = express()
    app.use(express.json())
    app.use('/auth', authRoutes)
    app.use(errorHandler) // Add error handler middleware

    // Clear all mocks
    vi.clearAllMocks()
    mockPrismaClient.user.findUnique.mockClear()
    mockPrismaClient.user.create.mockClear()
    mockPrismaClient.userAuth.update.mockClear()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('POST /auth/register', () => {
    it('should register a new user successfully', async () => {
      // Setup
      const userData = {
        email: 'test@example.com',
        name: 'Test User',
        password: 'password123'
      }

      const createdUser = {
        id: 'user-123',
        email: userData.email,
        name: userData.name
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(null) // User doesn't exist
      mockPrismaClient.user.create.mockResolvedValue(createdUser)

      // Mock jwt.sign
      const mockToken = 'mock-jwt-token'
      vi.spyOn(jwt, 'sign').mockReturnValue(mockToken as any)

      // Execute
      const response = await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(201)

      // Verify
      expect(response.body).toEqual({
        message: 'User created successfully',
        token: mockToken,
        user: {
          id: createdUser.id,
          email: createdUser.email,
          name: createdUser.name
        }
      })

      expect(mockPrismaClient.user.findUnique).toHaveBeenCalledWith({
        where: { email: userData.email }
      })
      expect(mockPrismaClient.user.create).toHaveBeenCalledWith({
        data: {
          email: userData.email,
          name: userData.name,
          auth: {
            create: {
              password: expect.any(String) // Hashed password
            }
          }
        }
      })

      expect(jwt.sign).toHaveBeenCalledWith(
        { userId: createdUser.id },
        config.jwt.secret,
        { expiresIn: config.jwt.expiresIn }
      )
    })

    it('should register user without name', async () => {
      // Setup
      const userData = {
        email: 'test@example.com',
        password: 'password123'
      }

      const createdUser = {
        id: 'user-123',
        email: userData.email,
        name: null
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(null)
      mockPrismaClient.user.create.mockResolvedValue(createdUser)
      vi.spyOn(jwt, 'sign').mockReturnValue('mock-token' as any)

      // Execute
      const response = await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(201)

      // Verify
      expect(mockPrismaClient.user.create).toHaveBeenCalledWith({
        data: {
          email: userData.email,
          name: null, // When not provided, it should be null, not undefined
          auth: {
            create: {
              password: expect.any(String)
            }
          }
        }
      })

      expect(response.body.user.name).toBeNull()
    })

    it('should return 400 if user already exists', async () => {
      // Setup
      const userData = {
        email: 'existing@example.com',
        password: 'password123'
      }

      mockPrismaClient.user.findUnique.mockResolvedValue({ id: 'existing-user' })

      // Execute
      const response = await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(400)

      // Verify
      expect(response.body).toEqual({
        error: 'User already exists with this email'
      })

      expect(mockPrismaClient.user.create).not.toHaveBeenCalled()
    })

    it('should return 400 for invalid email', async () => {
      // Setup
      const userData = {
        email: 'invalid-email',
        password: 'password123'
      }

      // Execute
      const response = await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(400)

      // Verify
      expect(response.body.error).toBeDefined()
      expect(response.body.error[0].message).toBe('Invalid email format')
      expect(mockPrismaClient.user.findUnique).not.toHaveBeenCalled()
    })

    it('should return 400 for short password', async () => {
      // Setup
      const userData = {
        email: 'test@example.com',
        password: '123' // Too short
      }

      // Execute
      const response = await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(400)

      // Verify
      expect(response.body.error).toBeDefined()
      expect(response.body.error[0].message).toBe('Password must be at least 6 characters')
    })

    it('should return 400 for missing fields', async () => {
      // Execute
      const response = await request(app)
        .post('/auth/register')
        .send({}) // No fields
        .expect(400)

      // Verify
      expect(response.body.error).toBeDefined()
      expect(response.body.error).toHaveLength(2) // email, password required
      expect(response.body.error.some((err: any) => err.path[0] === 'email' && err.message === 'Required')).toBe(true)
      expect(response.body.error.some((err: any) => err.path[0] === 'password' && err.message === 'Required')).toBe(true)
    })

    it('should hash password before storing', async () => {
      // Setup
      const userData = {
        email: 'test@example.com',
        password: 'plainpassword'
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(null)
      mockPrismaClient.user.create.mockResolvedValue({ id: 'user-123', email: userData.email, name: null })
      vi.spyOn(jwt, 'sign').mockReturnValue('mock-token' as any)

      // Execute
      await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(201)

      // Verify password was hashed
      const createCall = mockPrismaClient.user.create.mock.calls[0][0]
      const hashedPassword = createCall.data.auth.create.password
      expect(hashedPassword).toBeDefined()
      expect(hashedPassword).not.toBe(userData.password)
      expect(typeof hashedPassword).toBe('string')
    })
  })

  describe('POST /auth/login', () => {
    it('should login successfully with valid credentials', async () => {
      // Setup
      const loginData = {
        email: 'test@example.com',
        password: 'password123'
      }

      const user = {
        id: 'user-123',
        email: loginData.email,
        name: 'Test User',
        auth: {
          password: 'hashed-password'
        }
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(user)
      mockPrismaClient.userAuth.update.mockResolvedValue({}) // Mock the update call
      vi.spyOn(bcrypt, 'compare').mockResolvedValue(true as any)
      vi.spyOn(jwt, 'sign').mockReturnValue('login-token' as any)

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(200)

      // Verify
      expect(response.body).toEqual({
        message: 'Login successful',
        token: 'login-token',
        user: {
          id: user.id,
          email: user.email,
          name: user.name
        }
      })

      expect(mockPrismaClient.user.findUnique).toHaveBeenCalledWith({
        where: { email: loginData.email },
        include: { auth: true }
      })

      expect(bcrypt.compare).toHaveBeenCalledWith(loginData.password, user.auth.password)
      expect(jwt.sign).toHaveBeenCalledWith(
        { userId: user.id },
        config.jwt.secret,
        { expiresIn: config.jwt.expiresIn }
      )
    })

    it('should return 401 for non-existent user', async () => {
      // Setup
      const loginData = {
        email: 'nonexistent@example.com',
        password: 'password123'
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(null)

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(401)

      // Verify
      expect(response.body).toEqual({
        error: 'Invalid credentials'
      })
    })

    it('should return 401 for user without auth record', async () => {
      // Setup
      const loginData = {
        email: 'test@example.com',
        password: 'password123'
      }

      const userWithoutAuth = {
        id: 'user-123',
        email: loginData.email,
        name: 'Test User',
        auth: null
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(userWithoutAuth)

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(401)

      // Verify
      expect(response.body).toEqual({
        error: 'Invalid credentials'
      })
    })

    it('should return 401 for invalid password', async () => {
      // Setup
      const loginData = {
        email: 'test@example.com',
        password: 'wrongpassword'
      }

      const user = {
        id: 'user-123',
        email: loginData.email,
        name: 'Test User',
        auth: {
          password: 'hashed-password'
        }
      }

      mockPrismaClient.user.findUnique.mockResolvedValue(user)
      vi.spyOn(bcrypt, 'compare').mockResolvedValue(false as any)

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(401)

      // Verify
      expect(response.body).toEqual({
        error: 'Invalid credentials'
      })

      expect(bcrypt.compare).toHaveBeenCalledWith(loginData.password, user.auth.password)
    })

    it('should return 400 for invalid email format', async () => {
      // Setup
      const loginData = {
        email: 'invalid-email',
        password: 'password123'
      }

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(400)

      // Verify
      expect(response.body.error).toBeDefined()
      expect(response.body.error[0].message).toBe('Invalid email format')
    })

    it('should return 400 for missing password', async () => {
      // Setup
      const loginData = {
        email: 'test@example.com'
        // password missing
      }

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(400)

      // Verify
      expect(response.body.error[0].message).toBe('Required')
      expect(response.body.error[0].path[0]).toBe('password')
    })
  })

  describe('Error handling', () => {
    it('should handle database errors during registration', async () => {
      // Setup
      const userData = {
        email: 'test@example.com',
        password: 'password123'
      }

      mockPrismaClient.user.findUnique.mockRejectedValue(new Error('Database error'))

      // Execute - should return 500 due to error handler middleware
      const response = await request(app)
        .post('/auth/register')
        .send(userData)
        .expect(500)

      // Verify error response structure
      expect(response.body.error).toBe('Database error')
    })

    it('should handle database errors during login', async () => {
      // Setup
      const loginData = {
        email: 'test@example.com',
        password: 'password123'
      }

      mockPrismaClient.user.findUnique.mockRejectedValue(new Error('Database connection failed'))

      // Execute
      const response = await request(app)
        .post('/auth/login')
        .send(loginData)
        .expect(500)

      // Verify
      expect(response.body.error).toBe('Database connection failed')
    })
  })
})

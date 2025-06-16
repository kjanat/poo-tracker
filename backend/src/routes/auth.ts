import { Router, Request, Response, NextFunction } from 'express'
import bcrypt from 'bcryptjs'
import jwt from 'jsonwebtoken'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { config } from '../config'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router: Router = Router()
const prisma = new PrismaClient()

// Validation schemas
const registerSchema = z.object({
  email: z.string().email('Invalid email format'),
  name: z.string().min(1, 'Name is required').optional(),
  password: z.string().min(6, 'Password must be at least 6 characters')
})

const loginSchema = z.object({
  email: z.string().email('Invalid email format'),
  password: z.string().min(1, 'Password is required')
})

// Register
router.post('/register', async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  try {
    const { email, name, password } = registerSchema.parse(req.body)

    // Check if user already exists
    const existingUser = await prisma.user.findUnique({
      where: { email }
    })

    if (existingUser != null) {
      res.status(400).json({ error: 'User already exists with this email' })
      return
    }

    // Hash password
    const hashedPassword = await bcrypt.hash(password, 12)

    // Create user with separate auth record
    const user = await prisma.user.create({
      data: {
        email,
        name: name ?? null,
        auth: {
          create: {
            password: hashedPassword
          }
        }
      }
    })

    // Generate JWT
    const token = jwt.sign({ userId: user.id }, config.jwt.secret, {
      expiresIn: config.jwt.expiresIn
    } as jwt.SignOptions)

    res.status(201).json({
      message: 'User created successfully',
      token,
      user: {
        id: user.id,
        email: user.email,
        name: user.name
      }
    })
  } catch (error) {
    if (error instanceof z.ZodError) {
      res.status(400).json({ error: error.errors })
      return
    }
    next(error)
  }
})

// Login
router.post('/login', async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  try {
    const { email, password } = loginSchema.parse(req.body)

    // Find user with auth data
    const user = await prisma.user.findUnique({
      where: { email },
      include: { auth: true }
    })

    if ((user == null) || (user.auth == null)) {
      res.status(401).json({ error: 'Invalid credentials' })
      return
    }

    // Verify password
    const isPasswordValid = await bcrypt.compare(password, user.auth.password)
    if (!isPasswordValid) {
      res.status(401).json({ error: 'Invalid credentials' })
      return
    }

    // Update last login
    await prisma.userAuth.update({
      where: { userId: user.id },
      data: { lastLogin: new Date() }
    })

    // Generate JWT
    const token = jwt.sign({ userId: user.id }, config.jwt.secret, {
      expiresIn: config.jwt.expiresIn
    } as jwt.SignOptions)

    res.json({
      message: 'Login successful',
      token,
      user: {
        id: user.id,
        email: user.email,
        name: user.name
      }
    })
  } catch (error) {
    if (error instanceof z.ZodError) {
      res.status(400).json({ error: error.errors })
      return
    }
    next(error)
  }
})

// GET /api/auth/profile - Get user profile with auth info
router.get(
  '/profile',
  authenticateToken,
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const user = await prisma.user.findUnique({
        where: {
          id: req.userId
        }
      })

      const userAuth = await prisma.userAuth.findUnique({
        where: {
          userId: req.userId
        },
        select: {
          lastLogin: true,
          createdAt: true,
          updatedAt: true
        }
      })

      if (user == null) {
        res.status(404).json({ error: 'User not found' })
        return
      }

      res.json({
        user,
        auth: userAuth
      })
    } catch (error) {
      next(error)
    }
  }
)

// PUT /api/auth/profile - Update user profile
router.put(
  '/profile',
  authenticateToken,
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const updateProfileSchema = z.object({
        name: z.string().optional(),
        email: z.string().email().optional(),
        currentPassword: z.string().optional(),
        newPassword: z.string().min(6).optional()
      })

      const { name, email, currentPassword, newPassword } = updateProfileSchema.parse(req.body)

      // If changing password, verify current password
      if (newPassword) {
        if (!currentPassword) {
          res.status(400).json({ error: 'Current password is required to change password' })
          return
        }

        const userAuth = await prisma.userAuth.findUnique({
          where: { userId: req.userId }
        })

        if (userAuth == null) {
          res.status(404).json({ error: 'User auth not found' })
          return
        }

        const isCurrentPasswordValid = await bcrypt.compare(currentPassword, userAuth.password)
        if (!isCurrentPasswordValid) {
          res.status(400).json({ error: 'Current password is incorrect' })
          return
        }

        // Hash new password and update
        const hashedNewPassword = await bcrypt.hash(newPassword, 12)
        await prisma.userAuth.update({
          where: { userId: req.userId },
          data: { password: hashedNewPassword }
        })
      }

      // Update user profile
      const updateData: any = {}
      if (name !== undefined) updateData.name = name
      if (email !== undefined) updateData.email = email

      const user = await prisma.user.update({
        where: { id: req.userId },
        data: updateData
      })

      const userAuth = await prisma.userAuth.findUnique({
        where: { userId: req.userId },
        select: {
          lastLogin: true,
          createdAt: true,
          updatedAt: true
        }
      })

      res.json({
        user,
        auth: userAuth,
        message: 'Profile updated successfully'
      })
    } catch (error) {
      if (error instanceof z.ZodError) {
        res.status(400).json({ error: error.errors })
        return
      }
      next(error)
    }
  }
)

export { router as authRoutes }

import { Router, Request, Response, NextFunction } from 'express'
import bcrypt from 'bcryptjs'
import jwt from 'jsonwebtoken'
import { z } from 'zod'
import { PrismaClient } from '@prisma/client'
import { config } from '../config'

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

    if (existingUser) {
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

    if (!user || !user.auth) {
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

// Get user profile (no auth data exposed)
router.get('/profile', async (req: Request, res: Response, next: NextFunction): Promise<void> => {
  try {
    const authHeader = req.headers.authorization
    const token = authHeader?.split(' ')[1]

    if (!token) {
      res.status(401).json({ error: 'Access token required' })
      return
    }

    const decoded = jwt.verify(token, config.jwt.secret) as { userId: string }

    // Fetch only user profile data - NO AUTH DATA
    const user = await prisma.user.findUnique({
      where: { id: decoded.userId },
      select: {
        id: true,
        email: true,
        name: true,
        createdAt: true
        // Explicitly exclude auth relation - clean separation!
        // auth: false // Not needed, default behavior
      }
    })

    if (!user) {
      res.status(404).json({ error: 'User not found' })
      return
    }

    res.json({
      user,
      message: '✨ Clean profile data - no auth secrets leaked!'
    })
  } catch (error) {
    if (error instanceof jwt.JsonWebTokenError) {
      res.status(403).json({ error: 'Invalid token' })
      return
    }
    next(error)
  }
})

export { router as authRoutes }

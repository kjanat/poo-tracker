import request from 'supertest'
import express from 'express'

process.env.DATABASE_URL = 'postgresql://user:pass@localhost:5432/test'
process.env.JWT_SECRET = 'test-secret'
process.env.JWT_EXPIRES_IN = '1d'

const prismaMock = {
  user: {
    findUnique: jest.fn(),
    create: jest.fn(),
    update: jest.fn()
  },
  userAuth: {
    findUnique: jest.fn(),
    update: jest.fn()
  }
}

jest.mock('@prisma/client', () => ({
  PrismaClient: jest.fn(() => prismaMock)
}))
jest.mock('bcryptjs', () => ({ compare: jest.fn(), hash: jest.fn() }))
jest.mock('jsonwebtoken', () => ({ sign: jest.fn() }))

import { authRoutes } from '../routes/auth'
import { errorHandler } from '../middleware/errorHandler'
import bcrypt from 'bcryptjs'
import jwt from 'jsonwebtoken'

afterEach(() => {
  jest.clearAllMocks()
})

const createApp = () => {
  const app = express()
  app.use(express.json())
  const healthHandler: express.RequestHandler = (_req, res) => {
    res.json({ status: 'OK' })
  }
  app.get('/health', healthHandler)
  app.use('/api/auth', authRoutes)
  app.use(errorHandler)
  return app
}

describe('health endpoint', () => {
  it('responds with status OK', async () => {
    const app = createApp()
    const res = await request(app).get('/health')
    expect(res.status).toBe(200)
    expect(res.body.status).toBe('OK')
  })
})

describe('POST /api/auth/login', () => {
  it('returns 401 for invalid credentials', async () => {
    prismaMock.user.findUnique.mockResolvedValue(null)
    const app = createApp()
    const res = await request(app)
      .post('/api/auth/login')
      .send({ email: 'a@b.com', password: 'secret' })
    expect(res.status).toBe(401)
    expect(res.body.error).toBe('Invalid credentials')
  })

  it('logs in successfully', async () => {
    prismaMock.user.findUnique.mockResolvedValue({
      id: '1',
      email: 'a@b.com',
      name: 'Test',
      auth: { password: 'hashed' }
    })
    ;(bcrypt.compare as jest.Mock).mockResolvedValue(true)
    prismaMock.userAuth.update.mockResolvedValue({})
    ;(jwt.sign as jest.Mock).mockReturnValue('token')

    const app = createApp()
    const res = await request(app)
      .post('/api/auth/login')
      .send({ email: 'a@b.com', password: 'secret' })

    expect(res.status).toBe(200)
    expect(res.body.token).toBe('token')
    expect(prismaMock.userAuth.update).toHaveBeenCalled()
  })
})

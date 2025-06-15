process.env.DATABASE_URL = 'postgresql://user:pass@localhost:5432/test'
process.env.JWT_SECRET = 'testsecret'
process.env.NODE_ENV = 'test'
import request from 'supertest'
import app from '../index'
import bcrypt from 'bcryptjs'
import jwt from 'jsonwebtoken'

jest.mock('@prisma/client', () => {
  const mPrisma = {
    user: {
      findUnique: jest.fn(),
      create: jest.fn(),
      update: jest.fn()
    },
    userAuth: {
      findUnique: jest.fn(),
      update: jest.fn()
    },
    entry: {
      count: jest.fn(),
      groupBy: jest.fn(),
      findMany: jest.fn(),
      aggregate: jest.fn()
    }
  }
  return { PrismaClient: jest.fn(() => mPrisma), __esModule: true, default: mPrisma }
})

const { PrismaClient } = jest.requireMock('@prisma/client') as any
const prisma = new PrismaClient()

describe('API Endpoints', () => {
  afterEach(() => {
    jest.clearAllMocks()
  })

  it('GET /health should return ok', async () => {
    const res = await request(app).get('/health')
    expect(res.statusCode).toBe(200)
    expect(res.body.status).toBeDefined()
  })

  it('POST /api/auth/login success', async () => {
    prisma.user.findUnique.mockResolvedValue({
      id: '1',
      email: 'test@example.com',
      name: 'Tester',
      auth: { password: 'hashed' }
    })
    prisma.userAuth.update.mockResolvedValue({})
    jest.spyOn(bcrypt as any, 'compare').mockResolvedValue(true as any)
    jest.spyOn(jwt, 'sign').mockReturnValue('token123' as any)

    const res = await request(app)
      .post('/api/auth/login')
      .send({ email: 'test@example.com', password: 'password' })

    expect(res.statusCode).toBe(200)
    expect(res.body.token).toBe('token123')
    expect(prisma.user.findUnique).toHaveBeenCalled()
  })

  it('GET /api/entries requires auth', async () => {
    const res = await request(app).get('/api/entries')
    expect(res.statusCode).toBe(401)
  })
})

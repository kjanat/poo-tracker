process.env.NODE_ENV = 'test'
process.env.DATABASE_URL = process.env.DATABASE_URL || 'postgresql://user:pass@localhost:5432/test'

import request from 'supertest'

jest.mock('@prisma/client', () => ({ PrismaClient: jest.fn(() => ({})) }))

import app from '../index'

describe('GET /health', () => {
  it('should return ok status', async () => {
    const res = await request(app).get('/health')
    expect(res.status).toBe(200)
    expect(res.body).toHaveProperty('status', 'OK')
    expect(res.body).toHaveProperty('timestamp')
  })
})

import dotenv from 'dotenv'
import { z } from 'zod'

dotenv.config()

const envSchema = z.object({
  API_PORT: z.coerce.number().int().min(1).default(3002),
  NODE_ENV: z.enum(['development', 'test', 'production']).default('development'),
  CORS_ORIGIN: z.string().url().default('http://localhost:5173'),

  DATABASE_URL: z.string().url(),

  MINIO_ENDPOINT: z.string().nonempty(),
  MINIO_ACCESS_KEY: z.string().nonempty(),
  MINIO_SECRET_KEY: z.string().nonempty(),
  MINIO_BUCKET_NAME: z.string().nonempty(),
  MINIO_USE_SSL: z.coerce.boolean().default(false),

  AI_SERVICE_URL: z.string().url(),

  REDIS_URL: z.string().url(),

  JWT_SECRET: z.string().nonempty(),
  JWT_EXPIRES_IN: z.string().nonempty(),

  UPLOAD_DIR: z.string().default('./uploads'),
  UPLOAD_BASE_URL: z.string().url().default('http://localhost:3002'),
  MAX_FILE_SIZE: z.coerce
    .number()
    .int()
    .default(5 * 1024 * 1024)
})

const parsed = envSchema.safeParse(process.env)

if (!parsed.success) {
  const formatted = JSON.stringify(parsed.error.flatten().fieldErrors, null, 2)
  throw new Error(`Invalid environment variables:\n${formatted}`)
}

const env = parsed.data

export const config = {
  port: env.API_PORT,
  nodeEnv: env.NODE_ENV,
  corsOrigin: env.CORS_ORIGIN,

  // Database
  databaseUrl: env.DATABASE_URL,

  // Storage
  minio: {
    endpoint: env.MINIO_ENDPOINT,
    accessKey: env.MINIO_ACCESS_KEY,
    secretKey: env.MINIO_SECRET_KEY,
    bucketName: env.MINIO_BUCKET_NAME,
    useSSL: env.MINIO_USE_SSL
  },

  // AI Service
  aiServiceUrl: env.AI_SERVICE_URL,

  // Redis
  redisUrl: env.REDIS_URL,

  // JWT
  jwt: {
    secret: env.JWT_SECRET,
    expiresIn: env.JWT_EXPIRES_IN
  },

  // Image uploads
  uploads: {
    directory: env.UPLOAD_DIR,
    baseUrl: env.UPLOAD_BASE_URL,
    maxFileSize: env.MAX_FILE_SIZE,
    allowedTypes: ['image/jpeg', 'image/png', 'image/webp', 'image/gif'] as const
  }
} as const

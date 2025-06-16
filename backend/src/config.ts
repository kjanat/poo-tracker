import dotenv from 'dotenv'

dotenv.config()

export const config = {
  port: parseInt(process.env.API_PORT ?? '3002'),
  nodeEnv: process.env.NODE_ENV ?? 'development',
  corsOrigin: process.env.CORS_ORIGIN ?? 'http://localhost:5173',

  // Database
  databaseUrl: process.env.DATABASE_URL ?? 'postgresql://postgres:postgres@localhost:5432/poo_tracker',

  // Storage
  minio: {
    endpoint: process.env.MINIO_ENDPOINT ?? 'localhost:9000',
    accessKey: process.env.MINIO_ACCESS_KEY ?? 'minioadmin',
    secretKey: process.env.MINIO_SECRET_KEY ?? 'minioadmin123',
    bucketName: process.env.MINIO_BUCKET_NAME ?? 'poo-photos',
    useSSL: process.env.MINIO_USE_SSL === 'true'
  },

  // AI Service
  aiServiceUrl: process.env.AI_SERVICE_URL ?? 'http://localhost:8001',

  // Redis
  redisUrl: process.env.REDIS_URL ?? 'redis://localhost:6379',

  // JWT
  jwt: {
    secret: process.env.JWT_SECRET ?? 'your-super-secret-jwt-key-change-in-production',
    expiresIn: process.env.JWT_EXPIRES_IN ?? '7d'
  },

  // Image uploads
  uploads: {
    directory: process.env.UPLOAD_DIR ?? './uploads',
    baseUrl: process.env.UPLOAD_BASE_URL ?? 'http://localhost:3002',
    maxFileSize: parseInt(process.env.MAX_FILE_SIZE ?? '5242880'), // 5MB
    allowedTypes: ['image/jpeg', 'image/png', 'image/webp', 'image/gif'] as const
  }
} as const

// Validate required environment variables
const requiredEnvVars = ['DATABASE_URL']

for (const envVar of requiredEnvVars) {
  if ((process.env[envVar] ?? '') === '') {
    throw new Error(`Missing required environment variable: ${envVar}`)
  }
}

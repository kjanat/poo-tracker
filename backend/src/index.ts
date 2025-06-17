import express from 'express'
import cors from 'cors'
import helmet from 'helmet'
import compression from 'compression'
import morgan from 'morgan'
import rateLimit from 'express-rate-limit'
import path from 'path'
import { config } from './config'
import { errorHandler } from './middleware/errorHandler'
import { authRoutes } from './routes/auth'
import entryRoutes from './routes/entries'
import bowelMovementsRoutes from './routes/bowel-movements'
import { mealRoutes } from './routes/meals'
import { uploadRoutes } from './routes/uploads'
import { analyticsRoutes } from './routes/analytics'
import { ImageProcessingFactory } from './services/ImageProcessingService'

// Initialize image processing factory
ImageProcessingFactory.configure(config.uploads.directory, config.uploads.baseUrl)

const app: express.Application = express()

// Security middleware
app.use(helmet())
app.use(
  cors({
    origin: config.corsOrigin,
    credentials: true
  })
)

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // limit each IP to 100 requests per windowMs
  message: 'Too many requests from this IP, please try again later.'
})
app.use(limiter)

// General middleware
app.use(compression())
app.use(morgan('combined'))
app.use(express.json({ limit: '10mb' }))
app.use(express.urlencoded({ extended: true, limit: '10mb' }))

// Serve static assets (logos, etc.)
app.use('/assets', express.static(path.join(__dirname, '../public')))

// Health check
app.get('/health', (_req, res) => {
  res.json({ status: 'OK', timestamp: new Date().toISOString() })
})

// API routes
app.use('/api/auth', authRoutes)
app.use('/api/entries', entryRoutes) // Legacy route (redirects to bowel-movements)
app.use('/api/bowel-movements', bowelMovementsRoutes)
app.use('/api/meals', mealRoutes)
app.use('/api/uploads', uploadRoutes)
app.use('/api/analytics', analyticsRoutes)

// Serve uploaded images
app.use('/uploads', express.static(config.uploads.directory))

// 404 handler
app.use((_req, res) => {
  res.status(404).json({ error: 'Not found' })
})

// Error handling
app.use(errorHandler)

const server = app.listen(config.port, () => {
  console.log(`🚽 Poo Tracker API running on port ${config.port}`)
  console.log(`📊 Health check available at http://localhost:${config.port}/health`)
})

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully')
  server.close(() => {
    console.log('Process terminated')
  })
})

export default app

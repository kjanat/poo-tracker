import { Router, Response, NextFunction } from 'express'
import multer from 'multer'
import { Client as MinioClient } from 'minio'
import { v4 as uuidv4 } from 'uuid'
import { config } from '../config'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'

const router: Router = Router()

// Initialize MinIO client
const minioClient = new MinioClient({
  endPoint: config.minio.endpoint.split(':')[0]!,
  port: parseInt(config.minio.endpoint.split(':')[1]!),
  useSSL: config.minio.useSSL,
  accessKey: config.minio.accessKey,
  secretKey: config.minio.secretKey
})

// Configure multer for memory storage
const upload = multer({
  storage: multer.memoryStorage(),
  limits: {
    fileSize: 10 * 1024 * 1024 // 10MB limit
  },
  fileFilter: (_req, file, cb) => {
    // Only allow images
    if (file.mimetype.startsWith('image/')) {
      cb(null, true)
    } else {
      cb(new Error('Only image files are allowed'))
    }
  }
})

// Apply authentication to all routes
router.use(authenticateToken)

// Ensure bucket exists
const ensureBucket = async (): Promise<void> => {
  const bucketExists = await minioClient.bucketExists(config.minio.bucketName)
  if (!bucketExists) {
    await minioClient.makeBucket(config.minio.bucketName)
  }
}

// POST /api/uploads/photo - Upload photo
const uploadPhotoHandler = async (req: any, res: Response, next: NextFunction): Promise<void> => {
  try {
    const authenticatedReq = req as AuthenticatedRequest

    if (!req.file) {
      res.status(400).json({ error: 'No file provided' })
      return
    }

    if (!authenticatedReq.userId) {
      res.status(401).json({ error: 'User not authenticated' })
      return
    }

    await ensureBucket()

    const fileExtension = req.file.originalname.split('.').pop()
    const filename = `${authenticatedReq.userId}/${uuidv4()}.${fileExtension}`

    // Upload to MinIO - fix the putObject signature
    await minioClient.putObject(config.minio.bucketName, filename, req.file.buffer, req.file.size, {
      'Content-Type': req.file.mimetype
    })

    // Generate URL (in production you'd want signed URLs)
    const photoUrl = `http://${config.minio.endpoint}/${config.minio.bucketName}/${filename}`

    res.json({
      photoUrl,
      filename,
      size: req.file.size,
      mimetype: req.file.mimetype
    })
  } catch (error) {
    next(error)
  }
}

router.post('/photo', upload.single('photo') as any, uploadPhotoHandler)

export { router as uploadRoutes }

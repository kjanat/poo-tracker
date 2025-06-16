import { Router, Response, NextFunction } from 'express'
import multer from 'multer'
import { config } from '../config'
import { authenticateToken, AuthenticatedRequest } from '../middleware/auth'
import { ImageProcessingFactory } from '../services/ImageProcessingService'
import { sanitizeFilename } from '../utils/filename'

const router: Router = Router()

// Configure multer for memory storage
const upload = multer({
  storage: multer.memoryStorage(),
  limits: {
    fileSize: config.uploads.maxFileSize
  },
  fileFilter: (_req, file, cb) => {
    // Only allow specified image types
    if (
      config.uploads.allowedTypes.includes(
        file.mimetype as (typeof config.uploads.allowedTypes)[number]
      )
    ) {
      cb(null, true)
    } else {
      cb(new Error(`Only these image types are allowed: ${config.uploads.allowedTypes.join(', ')}`))
    }
  }
})

// Apply authentication to all routes
router.use(authenticateToken)

// POST /api/uploads/photo - Upload and process photo
router.post(
  '/photo',
  upload.single('photo'),
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (req.file == null) {
        res.status(400).json({ error: 'No file uploaded' })
        return
      }

      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const imageProcessor = ImageProcessingFactory.getInstance()

      // Process the image (resize, convert to WebP)
      const processedImage = await imageProcessor.processImage(
        req.file.buffer,
        req.file.originalname,
        {
          maxWidth: 1024,
          maxHeight: 1024,
          quality: 85,
          format: 'webp'
        }
      )

      res.status(201).json({
        message: 'Photo uploaded successfully',
        photo: {
          filename: processedImage.filename,
          url: processedImage.url,
          size: processedImage.size,
          width: processedImage.width,
          height: processedImage.height
        }
      })
    } catch (error) {
      if (error instanceof multer.MulterError) {
        if (error.code === 'LIMIT_FILE_SIZE') {
          res.status(400).json({
            error: `File too large. Maximum size is ${Math.round(config.uploads.maxFileSize / 1024 / 1024)}MB`
          })
          return
        }
      }
      next(error)
    }
  }
)

// DELETE /api/uploads/photo/:filename - Delete photo
router.delete(
  '/photo/:filename',
  async (req: AuthenticatedRequest, res: Response, next: NextFunction): Promise<void> => {
    try {
      if (!req.userId) {
        res.status(401).json({ error: 'User not authenticated' })
        return
      }

      const { filename } = req.params

      if (!filename) {
        res.status(400).json({ error: 'Filename is required' })
        return
      }

      const safeFilename = sanitizeFilename(filename)
      if (!safeFilename) {
        res.status(400).json({ error: 'Invalid filename' })
        return
      }

      const imageProcessor = ImageProcessingFactory.getInstance()
      await imageProcessor.deleteImage(safeFilename)

      res.json({ message: 'Photo deleted successfully' })
    } catch (error) {
      next(error)
    }
  }
)

export { router as uploadRoutes }

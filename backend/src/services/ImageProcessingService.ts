import sharp from 'sharp'
import { v4 as uuidv4 } from 'uuid'
import path from 'path'
import fs from 'fs/promises'
import { sanitizeFilename } from '../utils/filename'

export interface ImageProcessingConfig {
  maxWidth: number
  maxHeight: number
  quality: number
  format: 'webp' | 'jpeg' | 'png'
}

export interface ProcessedImage {
  filename: string
  path: string
  url: string
  size: number
  width: number
  height: number
}

export interface ImageProcessingService {
  processImage: (
    buffer: Buffer,
    originalName: string,
    config?: Partial<ImageProcessingConfig>
  ) => Promise<ProcessedImage>
  deleteImage: (filename: string) => Promise<void>
}

export class SharpImageProcessingService implements ImageProcessingService {
  private readonly uploadDir: string
  private readonly baseUrl: string
  private readonly defaultConfig: ImageProcessingConfig = {
    maxWidth: 1024,
    maxHeight: 1024,
    quality: 80,
    format: 'webp'
  }

  constructor (uploadDir: string, baseUrl: string) {
    this.uploadDir = uploadDir
    this.baseUrl = baseUrl
  }

  async processImage (
    buffer: Buffer,
    _originalName: string,
    config: Partial<ImageProcessingConfig> = {}
  ): Promise<ProcessedImage> {
    const finalConfig = { ...this.defaultConfig, ...config }

    // Generate unique filename
    const fileId = uuidv4()
    const extension = finalConfig.format
    const filename = `${fileId}.${extension}`
    const filePath = path.join(this.uploadDir, filename)

    // Ensure upload directory exists
    await fs.mkdir(this.uploadDir, { recursive: true })

    // Process image
    const processedBuffer = await sharp(buffer)
      .resize(finalConfig.maxWidth, finalConfig.maxHeight, {
        fit: 'inside',
        withoutEnlargement: true
      })
      .toFormat(finalConfig.format, { quality: finalConfig.quality })
      .toBuffer()

    // Get image metadata
    const metadata = await sharp(processedBuffer).metadata()

    // Save processed image
    await fs.writeFile(filePath, processedBuffer)

    return {
      filename,
      path: filePath,
      url: `${this.baseUrl}/uploads/${filename}`,
      size: processedBuffer.length,
      width: metadata.width || 0,
      height: metadata.height || 0
    }
  }

  async deleteImage (filename: string): Promise<void> {
    const safeFilename = sanitizeFilename(filename)
    if (!safeFilename) {
      throw new Error(`Invalid filename: ${filename}`)
    }
    const filePath = path.join(this.uploadDir, safeFilename)
    try {
      await fs.unlink(filePath)
    } catch (error) {
      // File might not exist, which is okay
      console.warn(`Failed to delete image ${filename}:`, error)
    }
  }
}

// Factory for dependency injection
export class ImageProcessingFactory {
  private static instance: ImageProcessingService

  static configure (uploadDir: string, baseUrl: string): void {
    this.instance = new SharpImageProcessingService(uploadDir, baseUrl)
  }

  static getInstance (): ImageProcessingService {
    if (!this.instance) {
      throw new Error('ImageProcessingFactory not configured. Call configure() first.')
    }
    return this.instance
  }

  // For testing - allows injection of mock service
  static setInstance (service: ImageProcessingService): void {
    this.instance = service
  }
}

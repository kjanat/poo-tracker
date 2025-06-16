import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import fs from 'fs/promises'
import path from 'path'
import sharp from 'sharp'
import { SharpImageProcessingService, ImageProcessingFactory } from '../ImageProcessingService'
import { sanitizeFilename } from '../../utils/filename'

// Mock dependencies
vi.mock('fs/promises')
vi.mock('sharp')
vi.mock('../../utils/filename')
vi.mock('uuid', () => ({
  v4: vi.fn(() => 'test-uuid-123')
}))

describe('SharpImageProcessingService', () => {
  let service: SharpImageProcessingService
  const mockUploadDir = '/mock/uploads'
  const mockBaseUrl = 'http://localhost:3000'

  // Mock implementations
  const mockFs = vi.mocked(fs)
  const mockSharp = vi.mocked(sharp)
  const mockSanitizeFilename = vi.mocked(sanitizeFilename)

  beforeEach(() => {
    // Reset all mocks
    vi.clearAllMocks()

    service = new SharpImageProcessingService(mockUploadDir, mockBaseUrl)

    // Default mock implementations
    mockFs.mkdir.mockResolvedValue(undefined)
    mockFs.writeFile.mockResolvedValue(undefined)
    mockFs.unlink.mockResolvedValue(undefined)
    mockSanitizeFilename.mockReturnValue('valid-filename.webp')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('processImage', () => {
    it('should process image with default configuration', async () => {
      // Setup
      const mockBuffer = Buffer.from('mock-image-data')
      const originalName = 'test-image.jpg'
      const processedBuffer = Buffer.from('processed-image-data')

      // Mock Sharp chain
      const mockSharpInstance = {
        resize: vi.fn().mockReturnThis(),
        toFormat: vi.fn().mockReturnThis(),
        toBuffer: vi.fn().mockResolvedValue(processedBuffer),
        metadata: vi.fn().mockResolvedValue({ width: 800, height: 600 })
      }

      mockSharp.mockReturnValue(mockSharpInstance as any)

      // Execute
      const result = await service.processImage(mockBuffer, originalName)

      // Verify
      expect(mockFs.mkdir).toHaveBeenCalledWith(mockUploadDir, { recursive: true })
      expect(mockSharp).toHaveBeenCalledWith(mockBuffer)
      expect(mockSharpInstance.resize).toHaveBeenCalledWith(1024, 1024, {
        fit: 'inside',
        withoutEnlargement: true
      })
      expect(mockSharpInstance.toFormat).toHaveBeenCalledWith('webp', { quality: 80 })
      expect(mockFs.writeFile).toHaveBeenCalledWith(
        path.join(mockUploadDir, 'test-uuid-123.webp'),
        processedBuffer
      )

      expect(result).toEqual({
        filename: 'test-uuid-123.webp',
        path: path.join(mockUploadDir, 'test-uuid-123.webp'),
        url: `${mockBaseUrl}/uploads/test-uuid-123.webp`,
        size: processedBuffer.length,
        width: 800,
        height: 600
      })
    })

    it('should process image with custom configuration', async () => {
      // Setup
      const mockBuffer = Buffer.from('mock-image-data')
      const originalName = 'test-image.jpg'
      const processedBuffer = Buffer.from('processed-image-data')
      const customConfig = {
        maxWidth: 512,
        maxHeight: 512,
        quality: 90,
        format: 'jpeg' as const
      }

      const mockSharpInstance = {
        resize: vi.fn().mockReturnThis(),
        toFormat: vi.fn().mockReturnThis(),
        toBuffer: vi.fn().mockResolvedValue(processedBuffer),
        metadata: vi.fn().mockResolvedValue({ width: 400, height: 300 })
      }

      mockSharp.mockReturnValue(mockSharpInstance as any)

      // Execute
      const result = await service.processImage(mockBuffer, originalName, customConfig)

      // Verify custom config was used
      expect(mockSharpInstance.resize).toHaveBeenCalledWith(512, 512, {
        fit: 'inside',
        withoutEnlargement: true
      })
      expect(mockSharpInstance.toFormat).toHaveBeenCalledWith('jpeg', { quality: 90 })
      expect(result.filename).toBe('test-uuid-123.jpeg')
    })

    it('should handle missing metadata gracefully', async () => {
      // Setup
      const mockBuffer = Buffer.from('mock-image-data')
      const processedBuffer = Buffer.from('processed-image-data')

      const mockSharpInstance = {
        resize: vi.fn().mockReturnThis(),
        toFormat: vi.fn().mockReturnThis(),
        toBuffer: vi.fn().mockResolvedValue(processedBuffer),
        metadata: vi.fn().mockResolvedValue({}) // No width/height
      }

      mockSharp.mockReturnValue(mockSharpInstance as any)

      // Execute
      const result = await service.processImage(mockBuffer, 'test.jpg')

      // Verify defaults are used for missing metadata
      expect(result.width).toBe(0)
      expect(result.height).toBe(0)
    })

    it('should throw error when Sharp fails', async () => {
      // Setup
      const mockBuffer = Buffer.from('mock-image-data')
      const mockSharpInstance = {
        resize: vi.fn().mockReturnThis(),
        toFormat: vi.fn().mockReturnThis(),
        toBuffer: vi.fn().mockRejectedValue(new Error('Invalid image format'))
      }

      mockSharp.mockReturnValue(mockSharpInstance as any)

      // Execute & Verify
      await expect(service.processImage(mockBuffer, 'invalid.txt')).rejects.toThrow(
        'Invalid image format'
      )
    })

    it('should throw error when file system fails', async () => {
      // Setup
      const mockBuffer = Buffer.from('mock-image-data')
      const processedBuffer = Buffer.from('processed-image-data')

      const mockSharpInstance = {
        resize: vi.fn().mockReturnThis(),
        toFormat: vi.fn().mockReturnThis(),
        toBuffer: vi.fn().mockResolvedValue(processedBuffer),
        metadata: vi.fn().mockResolvedValue({ width: 800, height: 600 })
      }

      mockSharp.mockReturnValue(mockSharpInstance as any)
      mockFs.writeFile.mockRejectedValue(new Error('Permission denied'))

      // Execute & Verify
      await expect(service.processImage(mockBuffer, 'test.jpg')).rejects.toThrow(
        'Permission denied'
      )
    })
  })

  describe('deleteImage', () => {
    it('should delete image with valid filename', async () => {
      // Setup
      const filename = 'test-image.webp'
      const safeFilename = 'safe-test-image.webp'
      mockSanitizeFilename.mockReturnValue(safeFilename)

      // Execute
      await service.deleteImage(filename)

      // Verify
      expect(mockSanitizeFilename).toHaveBeenCalledWith(filename)
      expect(mockFs.unlink).toHaveBeenCalledWith(path.join(mockUploadDir, safeFilename))
    })

    it('should throw error for invalid filename', async () => {
      // Setup
      const filename = '../../../etc/passwd'
      mockSanitizeFilename.mockReturnValue(null)

      // Execute & Verify
      await expect(service.deleteImage(filename)).rejects.toThrow(
        'Invalid filename: ../../../etc/passwd'
      )

      expect(mockFs.unlink).not.toHaveBeenCalled()
    })

    it('should handle file not found gracefully', async () => {
      // Setup
      const filename = 'nonexistent.webp'
      const safeFilename = 'nonexistent.webp'
      mockSanitizeFilename.mockReturnValue(safeFilename)
      mockFs.unlink.mockRejectedValue(new Error('ENOENT: no such file or directory'))

      // Spy on console.warn

      const consoleSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})

      // Execute (should not throw)
      await expect(service.deleteImage(filename)).resolves.toBeUndefined()

      // Verify warning was logged
      expect(consoleSpy).toHaveBeenCalledWith(
        `Failed to delete image ${filename}:`,
        expect.any(Error)
      )

      consoleSpy.mockRestore()
    })
  })
})

describe('ImageProcessingFactory', () => {
  afterEach(() => {
    // Reset factory instance
    ImageProcessingFactory.setInstance(null as any)
  })

  it('should configure and return instance', () => {
    // Execute
    ImageProcessingFactory.configure('/uploads', 'http://localhost')
    const instance = ImageProcessingFactory.getInstance()

    // Verify
    expect(instance).toBeInstanceOf(SharpImageProcessingService)
  })

  it('should throw error when not configured', () => {
    // Execute & Verify
    expect(() => ImageProcessingFactory.getInstance()).toThrow(
      'ImageProcessingFactory not configured. Call configure() first.'
    )
  })

  it('should allow setting custom instance for testing', () => {
    // Setup
    const mockService = {} as any

    // Execute
    ImageProcessingFactory.setInstance(mockService)
    const instance = ImageProcessingFactory.getInstance()

    // Verify
    expect(instance).toBe(mockService)
  })
})

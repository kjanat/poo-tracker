import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import fs from 'fs/promises'
import sharp from 'sharp'
import { SharpImageProcessingService } from '../ImageProcessingService'

describe('SharpImageProcessingService file system errors', () => {
  const mockFs = vi.mocked(fs)
  const mockSharp = vi.mocked(sharp)
  let service: SharpImageProcessingService

  beforeEach(() => {
    vi.mock('fs/promises')
    vi.mock('sharp')
    vi.clearAllMocks()
    service = new SharpImageProcessingService('/uploads', 'http://localhost')
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('should throw when directory creation fails', async () => {
    const buffer = Buffer.from('data')
    const error = new Error('mkdir failed')
    const mockSharpInstance = {
      resize: vi.fn().mockReturnThis(),
      toFormat: vi.fn().mockReturnThis(),
      toBuffer: vi.fn().mockResolvedValue(buffer),
      metadata: vi.fn().mockResolvedValue({ width: 1, height: 1 })
    }
    mockSharp.mockReturnValue(mockSharpInstance as any)
    mockFs.mkdir.mockRejectedValue(error)

    await expect(service.processImage(buffer, 'a.png')).rejects.toThrow('mkdir failed')
  })

  it('deleteImage should log non-ENOENT errors', async () => {
    const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})
    vi.mock('../../utils/filename', () => ({ sanitizeFilename: () => 'file.webp' }))
    mockFs.unlink.mockRejectedValue(new Error('EACCES'))

    await service.deleteImage('file.webp')

    expect(warnSpy).toHaveBeenCalledWith('Failed to delete image file.webp:', expect.any(Error))
    warnSpy.mockRestore()
  })
})

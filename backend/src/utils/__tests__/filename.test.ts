import { sanitizeFilename } from '../filename'

describe('sanitizeFilename', () => {
  it('returns filename when valid', () => {
    expect(sanitizeFilename('test.png')).toBe('test.png')
  })

  it('rejects traversal paths', () => {
    expect(sanitizeFilename('../etc/passwd')).toBeNull()
  })

  it('rejects filenames with slashes', () => {
    expect(sanitizeFilename('a/b.png')).toBeNull()
  })
})

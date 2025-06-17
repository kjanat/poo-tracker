import { describe, it, expect } from 'vitest'
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

  it('rejects filenames with spaces', () => {
    expect(sanitizeFilename('file name.png')).toBeNull()
  })

  it('rejects filenames with underscores', () => {
    expect(sanitizeFilename('file_name.png')).toBeNull()
  })

  it('rejects filenames with special characters', () => {
    expect(sanitizeFilename('file@name.png')).toBeNull()
    expect(sanitizeFilename('file#name.png')).toBeNull()
    expect(sanitizeFilename('file$name.png')).toBeNull()
    expect(sanitizeFilename('file%name.png')).toBeNull()
    expect(sanitizeFilename('file&name.png')).toBeNull()
    expect(sanitizeFilename('file*name.png')).toBeNull()
    expect(sanitizeFilename('file+name.png')).toBeNull()
    expect(sanitizeFilename('file=name.png')).toBeNull()
    expect(sanitizeFilename('file?name.png')).toBeNull()
    expect(sanitizeFilename('file[name].png')).toBeNull()
    expect(sanitizeFilename('file{name}.png')).toBeNull()
    expect(sanitizeFilename('file|name.png')).toBeNull()
    expect(sanitizeFilename('file\\name.png')).toBeNull()
    expect(sanitizeFilename('file:name.png')).toBeNull()
    expect(sanitizeFilename('file;name.png')).toBeNull()
    expect(sanitizeFilename('file"name.png')).toBeNull()
    expect(sanitizeFilename("file'name.png")).toBeNull()
    expect(sanitizeFilename('file<name>.png')).toBeNull()
    expect(sanitizeFilename('file>name.png')).toBeNull()
    expect(sanitizeFilename('file,name.png')).toBeNull()
    expect(sanitizeFilename('file~name.png')).toBeNull()
    expect(sanitizeFilename('file`name.png')).toBeNull()
  })

  it('rejects filenames with directory traversal sequences', () => {
    expect(sanitizeFilename('file..name.png')).toBeNull()
    expect(sanitizeFilename('..filename.png')).toBeNull()
    expect(sanitizeFilename('filename..png')).toBeNull()
  })

  it('rejects empty filenames', () => {
    expect(sanitizeFilename('')).toBeNull()
  })

  it('rejects filenames with backslashes', () => {
    expect(sanitizeFilename('path\\to\\file.png')).toBeNull()
  })

  it('accepts valid filenames with allowed characters', () => {
    expect(sanitizeFilename('valid-filename.png')).toBe('valid-filename.png')
    expect(sanitizeFilename('file.name.png')).toBe('file.name.png')
    expect(sanitizeFilename('123.png')).toBe('123.png')
    expect(sanitizeFilename('a.b.c.d.png')).toBe('a.b.c.d.png')
    expect(sanitizeFilename('file-name.png')).toBe('file-name.png')
    expect(sanitizeFilename('UPPERCASE.PNG')).toBe('UPPERCASE.PNG')
    expect(sanitizeFilename('MiXeDcAsE.JpG')).toBe('MiXeDcAsE.JpG')
  })
})

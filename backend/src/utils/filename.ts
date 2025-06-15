import path from 'path'

export function sanitizeFilename(filename: string): string | null {
  const base = path.basename(filename)
  // reject if path has directory components or invalid characters
  if (base !== filename) return null
  if (!/^[a-zA-Z0-9.-]+$/.test(base)) return null
  // prevent path traversal using '..'
  if (base.includes('..')) return null
  return base
}

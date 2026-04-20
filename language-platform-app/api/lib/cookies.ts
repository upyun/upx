export function parseCookies(cookieHeader?: string): Record<string, string> {
  if (!cookieHeader) return {}

  return cookieHeader
    .split(';')
    .map((p) => p.trim())
    .filter(Boolean)
    .reduce<Record<string, string>>((acc, part) => {
      const idx = part.indexOf('=')
      if (idx <= 0) return acc
      const key = decodeURIComponent(part.slice(0, idx).trim())
      const value = decodeURIComponent(part.slice(idx + 1).trim())
      acc[key] = value
      return acc
    }, {})
}

export function buildCookie(
  name: string,
  value: string,
  options?: {
    httpOnly?: boolean
    secure?: boolean
    sameSite?: 'Lax' | 'Strict' | 'None'
    path?: string
    maxAgeSeconds?: number
  },
): string {
  const parts: string[] = []
  parts.push(`${encodeURIComponent(name)}=${encodeURIComponent(value)}`)
  parts.push(`Path=${options?.path ?? '/'}`)
  parts.push(`SameSite=${options?.sameSite ?? 'Lax'}`)
  if (options?.httpOnly ?? true) parts.push('HttpOnly')
  if (options?.secure) parts.push('Secure')
  if (typeof options?.maxAgeSeconds === 'number') {
    parts.push(`Max-Age=${Math.max(0, Math.floor(options.maxAgeSeconds))}`)
  }
  return parts.join('; ')
}


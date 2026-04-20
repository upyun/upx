/**
 * This is a user authentication API route demo.
 * Handle user registration, login, token management, etc.
 */
import { Router, type Request, type Response } from 'express'
import crypto from 'crypto'
import { buildCookie, parseCookies } from '../lib/cookies.js'
import { hashPassword, randomToken, verifyPassword } from '../lib/password.js'
import { loadDb, saveDb } from '../storage/db.js'

const router = Router()

router.post('/register', async (req: Request, res: Response): Promise<void> => {
  const email = String(req.body?.email ?? '').trim().toLowerCase()
  const password = String(req.body?.password ?? '')
  const uiLocale = String(req.body?.uiLocale ?? 'en').trim() || 'en'

  if (!email || !email.includes('@') || password.length < 8) {
    res.status(400).json({ success: false, error: 'Invalid registration payload' })
    return
  }

  const db = await loadDb()
  const existing = db.users.find((u) => u.email === email)
  if (existing) {
    res.status(409).json({ success: false, error: 'Email already registered' })
    return
  }

  const { salt, hash } = hashPassword(password)
  const user = {
    id: crypto.randomUUID(),
    email,
    passwordHash: hash,
    passwordSalt: salt,
    uiLocale,
    createdAt: new Date().toISOString(),
  }
  db.users.push(user)

  const sessionId = randomToken(32)
  const now = Date.now()
  const expiresAt = new Date(now + 1000 * 60 * 60 * 24 * 14).toISOString()
  db.sessions.push({
    id: sessionId,
    userId: user.id,
    createdAt: new Date(now).toISOString(),
    expiresAt,
  })

  await saveDb(db)

  const secure = process.env.NODE_ENV === 'production'
  res.setHeader(
    'Set-Cookie',
    buildCookie('sid', sessionId, { httpOnly: true, secure, sameSite: 'Lax', path: '/' }),
  )
  res.status(201).json({
    success: true,
    user: { id: user.id, email: user.email, uiLocale: user.uiLocale },
  })
})

router.post('/login', async (req: Request, res: Response): Promise<void> => {
  const email = String(req.body?.email ?? '').trim().toLowerCase()
  const password = String(req.body?.password ?? '')

  if (!email || !email.includes('@') || !password) {
    res.status(400).json({ success: false, error: 'Invalid login payload' })
    return
  }

  const db = await loadDb()
  const user = db.users.find((u) => u.email === email)
  if (!user || !verifyPassword(password, user.passwordSalt, user.passwordHash)) {
    res.status(401).json({ success: false, error: 'Invalid credentials' })
    return
  }

  const sessionId = randomToken(32)
  const now = Date.now()
  const expiresAt = new Date(now + 1000 * 60 * 60 * 24 * 14).toISOString()
  db.sessions.push({
    id: sessionId,
    userId: user.id,
    createdAt: new Date(now).toISOString(),
    expiresAt,
  })
  await saveDb(db)

  const secure = process.env.NODE_ENV === 'production'
  res.setHeader(
    'Set-Cookie',
    buildCookie('sid', sessionId, { httpOnly: true, secure, sameSite: 'Lax', path: '/' }),
  )
  res.status(200).json({
    success: true,
    user: { id: user.id, email: user.email, uiLocale: user.uiLocale },
  })
})

router.post('/logout', async (req: Request, res: Response): Promise<void> => {
  const cookies = parseCookies(req.headers.cookie)
  const sid = cookies.sid

  if (sid) {
    const db = await loadDb()
    db.sessions = db.sessions.filter((s) => s.id !== sid)
    await saveDb(db)
  }

  const secure = process.env.NODE_ENV === 'production'
  res.setHeader(
    'Set-Cookie',
    buildCookie('sid', '', {
      httpOnly: true,
      secure,
      sameSite: 'Lax',
      path: '/',
      maxAgeSeconds: 0,
    }),
  )
  res.status(200).json({ success: true })
})

router.get('/me', async (req: Request, res: Response): Promise<void> => {
  const cookies = parseCookies(req.headers.cookie)
  const sid = cookies.sid
  if (!sid) {
    res.status(401).json({ success: false, error: 'Not authenticated' })
    return
  }

  const db = await loadDb()
  const session = db.sessions.find((s) => s.id === sid)
  if (!session) {
    res.status(401).json({ success: false, error: 'Not authenticated' })
    return
  }
  if (new Date(session.expiresAt).getTime() < Date.now()) {
    db.sessions = db.sessions.filter((s) => s.id !== sid)
    await saveDb(db)
    res.status(401).json({ success: false, error: 'Not authenticated' })
    return
  }

  const user = db.users.find((u) => u.id === session.userId)
  if (!user) {
    res.status(401).json({ success: false, error: 'Not authenticated' })
    return
  }

  res.status(200).json({ success: true, user: { id: user.id, email: user.email, uiLocale: user.uiLocale } })
})

export default router

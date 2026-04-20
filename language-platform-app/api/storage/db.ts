import fs from 'fs/promises'
import path from 'path'
import { fileURLToPath } from 'url'

export type DbUser = {
  id: string
  email: string
  passwordHash: string
  passwordSalt: string
  uiLocale: string
  createdAt: string
}

export type DbSession = {
  id: string
  userId: string
  createdAt: string
  expiresAt: string
}

export type Db = {
  users: DbUser[]
  sessions: DbSession[]
}

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const dataDir = path.join(__dirname, '..', '.data')
const dbFilePath = path.join(dataDir, 'db.json')

const emptyDb: Db = { users: [], sessions: [] }

async function ensureDataDir() {
  await fs.mkdir(dataDir, { recursive: true })
}

export async function loadDb(): Promise<Db> {
  await ensureDataDir()
  try {
    const raw = await fs.readFile(dbFilePath, 'utf8')
    const parsed = JSON.parse(raw) as Partial<Db>
    return {
      users: parsed.users ?? [],
      sessions: parsed.sessions ?? [],
    }
  } catch {
    await fs.writeFile(dbFilePath, JSON.stringify(emptyDb, null, 2), 'utf8')
    return { ...emptyDb }
  }
}

export async function saveDb(db: Db): Promise<void> {
  await ensureDataDir()
  await fs.writeFile(dbFilePath, JSON.stringify(db, null, 2), 'utf8')
}


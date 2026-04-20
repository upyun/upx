import crypto from 'crypto'

const ITERATIONS = 150_000
const KEYLEN = 32
const DIGEST = 'sha256'

export function hashPassword(
  password: string,
  salt = crypto.randomBytes(16).toString('hex'),
): { salt: string; hash: string } {
  const hash = crypto
    .pbkdf2Sync(password, salt, ITERATIONS, KEYLEN, DIGEST)
    .toString('hex')
  return { salt, hash }
}

export function verifyPassword(password: string, salt: string, hash: string): boolean {
  const attempt = crypto
    .pbkdf2Sync(password, salt, ITERATIONS, KEYLEN, DIGEST)
    .toString('hex')
  return crypto.timingSafeEqual(Buffer.from(attempt, 'hex'), Buffer.from(hash, 'hex'))
}

export function randomToken(bytes = 32): string {
  return crypto.randomBytes(bytes).toString('hex')
}


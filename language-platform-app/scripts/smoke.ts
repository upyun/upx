import assert from 'node:assert/strict'
import { once } from 'node:events'
import app from '../api/app.js'

async function main() {
  const server = app.listen(0)
  await once(server, 'listening')
  const address = server.address()
  assert.ok(address && typeof address === 'object')
  const base = `http://127.0.0.1:${address.port}`

  const health = await fetch(`${base}/api/health`)
  assert.equal(health.status, 200)

  const register = await fetch(`${base}/api/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: `smoke_${Date.now()}@example.com`,
      password: 'password123',
      uiLocale: 'en',
    }),
  })
  assert.equal(register.status, 201)
  const setCookie = register.headers.get('set-cookie')
  assert.ok(setCookie)

  const me = await fetch(`${base}/api/auth/me`, {
    headers: { cookie: setCookie },
  })
  assert.equal(me.status, 200)
  const meJson = (await me.json()) as any
  assert.equal(meJson.success, true)
  assert.ok(meJson.user?.id)

  server.close()
}

main().catch((e) => {
  console.error(e)
  process.exit(1)
})


export async function apiRequest<T>(
  input: string,
  init?: RequestInit & { json?: unknown },
): Promise<T> {
  const headers = new Headers(init?.headers)
  if (init?.json !== undefined) headers.set('Content-Type', 'application/json')

  const res = await fetch(input, {
    ...init,
    headers,
    credentials: 'include',
    body: init?.json !== undefined ? JSON.stringify(init.json) : init?.body,
  })

  const text = await res.text()
  const data = text ? (JSON.parse(text) as unknown) : null

  if (!res.ok) {
    const message =
      typeof data === 'object' && data && 'error' in data ? String((data as any).error) : 'Request failed'
    throw new Error(message)
  }

  return data as T
}


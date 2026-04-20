import { useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useI18n } from '@/i18n/useI18n'
import { useAuthStore } from '@/stores/useAuthStore'

export default function Login() {
  const { t } = useI18n()
  const login = useAuthStore((s) => s.login)
  const loading = useAuthStore((s) => s.loading)
  const error = useAuthStore((s) => s.error)
  const clearError = useAuthStore((s) => s.clearError)
  const navigate = useNavigate()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const canSubmit = useMemo(() => email.includes('@') && password.length >= 1 && !loading, [email, password, loading])

  return (
    <div className="mx-auto max-w-md space-y-4">
      <div className="text-xl font-semibold">{t('auth.login')}</div>

      {error ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          <div className="flex items-center justify-between gap-2">
            <div className="min-w-0">{error}</div>
            <button type="button" className="text-xs underline" onClick={clearError}>
              Close
            </button>
          </div>
        </div>
      ) : null}

      <div className="space-y-2">
        <label className="block text-xs text-zinc-500 dark:text-zinc-400">{t('auth.email')}</label>
        <input
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="h-10 w-full rounded-xl border border-zinc-200 bg-white px-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-700 dark:bg-zinc-900"
          placeholder="you@example.com"
        />
      </div>

      <div className="space-y-2">
        <label className="block text-xs text-zinc-500 dark:text-zinc-400">{t('auth.password')}</label>
        <input
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="h-10 w-full rounded-xl border border-zinc-200 bg-white px-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-700 dark:bg-zinc-900"
          type="password"
          placeholder="••••••••"
        />
      </div>

      <button
        type="button"
        disabled={!canSubmit}
        onClick={async () => {
          await login({ email, password })
          if (!useAuthStore.getState().error) navigate('/catalog')
        }}
        className="h-10 w-full rounded-xl bg-zinc-900 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
      >
        {t('auth.login')}
      </button>

      <div className="text-xs text-zinc-500 dark:text-zinc-400">
        {t('auth.createAccount')}?{' '}
        <Link to="/register" className="text-zinc-900 underline dark:text-white">
          {t('auth.register')}
        </Link>
      </div>
    </div>
  )
}


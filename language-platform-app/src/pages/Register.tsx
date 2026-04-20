import { useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useI18n } from '@/i18n/useI18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { useSettingsStore } from '@/stores/useSettingsStore'

export default function Register() {
  const { t, uiLocale, setUiLocale } = useI18n()
  const register = useAuthStore((s) => s.register)
  const loading = useAuthStore((s) => s.loading)
  const error = useAuthStore((s) => s.error)
  const clearError = useAuthStore((s) => s.clearError)
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const navigate = useNavigate()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const canSubmit = useMemo(
    () => email.includes('@') && password.length >= 8 && !loading,
    [email, password, loading],
  )

  return (
    <div className="mx-auto max-w-md space-y-4">
      <div className="text-xl font-semibold">{t('auth.register')}</div>

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

      <div className="grid grid-cols-2 gap-2">
        <label className="text-xs text-zinc-500 dark:text-zinc-400">{t('landing.uiLanguage')}</label>
        <select
          value={uiLocale}
          onChange={(e) => setUiLocale(e.target.value as any)}
          className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-xs dark:border-zinc-700 dark:bg-zinc-900"
        >
          <option value="en">EN</option>
          <option value="ja">日本語</option>
          <option value="ko">한국어</option>
        </select>
      </div>

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
          placeholder="8+ characters"
        />
      </div>

      <button
        type="button"
        disabled={!canSubmit}
        onClick={async () => {
          await register({ email, password, uiLocale })
          if (!useAuthStore.getState().error) navigate(`/onboarding?target=${targetLanguage}`)
        }}
        className="h-10 w-full rounded-xl bg-zinc-900 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
      >
        {t('auth.register')}
      </button>

      <div className="text-xs text-zinc-500 dark:text-zinc-400">
        {t('auth.haveAccount')}?{' '}
        <Link to="/login" className="text-zinc-900 underline dark:text-white">
          {t('auth.login')}
        </Link>
      </div>
    </div>
  )
}


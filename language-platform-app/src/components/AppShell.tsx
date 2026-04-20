import type { ReactNode } from 'react'
import { Link, NavLink, useNavigate } from 'react-router-dom'
import { cn } from '@/lib/utils'
import { useI18n } from '@/i18n/useI18n'
import { useSettingsStore } from '@/stores/useSettingsStore'
import { useAuthStore } from '@/stores/useAuthStore'

type Props = {
  children: ReactNode
}

const navItemBase =
  'flex items-center gap-2 rounded-lg px-3 py-2 text-sm transition-colors hover:bg-zinc-900/5 dark:hover:bg-white/10'

const navItemActive = 'bg-zinc-900/5 text-zinc-900 dark:bg-white/10 dark:text-white'

export default function AppShell({ children }: Props) {
  const { t, uiLocale, setUiLocale } = useI18n()
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const setTargetLanguage = useSettingsStore((s) => s.setTargetLanguage)
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)
  const navigate = useNavigate()

  return (
    <div className="min-h-dvh bg-zinc-50 text-zinc-900 dark:bg-zinc-950 dark:text-zinc-50">
      <div className="mx-auto grid max-w-6xl grid-cols-1 gap-6 px-4 py-6 md:grid-cols-[240px_1fr]">
        <aside className="rounded-2xl border border-zinc-200 bg-white p-4 shadow-sm dark:border-zinc-800 dark:bg-zinc-900/40">
          <Link to="/" className="block rounded-xl px-2 py-2">
            <div className="text-xs font-medium tracking-wide text-zinc-500 dark:text-zinc-400">
              {t('app.name')}
            </div>
          </Link>

          <div className="mt-3 space-y-1">
            <NavLink to="/" end className={({ isActive }) => cn(navItemBase, isActive && navItemActive)}>
              {t('nav.home')}
            </NavLink>
            <NavLink to="/catalog" className={({ isActive }) => cn(navItemBase, isActive && navItemActive)}>
              {t('nav.catalog')}
            </NavLink>
            <NavLink to="/path" className={({ isActive }) => cn(navItemBase, isActive && navItemActive)}>
              {t('nav.path')}
            </NavLink>
            <NavLink to="/progress" className={({ isActive }) => cn(navItemBase, isActive && navItemActive)}>
              {t('nav.progress')}
            </NavLink>
            <NavLink to="/community" className={({ isActive }) => cn(navItemBase, isActive && navItemActive)}>
              {t('nav.community')}
            </NavLink>
            <NavLink to="/achievements" className={({ isActive }) => cn(navItemBase, isActive && navItemActive)}>
              {t('nav.achievements')}
            </NavLink>
          </div>

          <div className="mt-4 border-t border-zinc-200 pt-4 dark:border-zinc-800">
            <div className="grid grid-cols-2 gap-2">
              <label className="text-xs text-zinc-500 dark:text-zinc-400">{t('landing.uiLanguage')}</label>
              <select
                value={uiLocale}
                onChange={(e) => setUiLocale(e.target.value as any)}
                className="h-8 rounded-lg border border-zinc-200 bg-white px-2 text-xs dark:border-zinc-700 dark:bg-zinc-900"
              >
                <option value="en">EN</option>
                <option value="ja">日本語</option>
                <option value="ko">한국어</option>
              </select>
              <label className="text-xs text-zinc-500 dark:text-zinc-400">{t('landing.targetLanguage')}</label>
              <select
                value={targetLanguage}
                onChange={(e) => setTargetLanguage(e.target.value as any)}
                className="h-8 rounded-lg border border-zinc-200 bg-white px-2 text-xs dark:border-zinc-700 dark:bg-zinc-900"
              >
                <option value="ja">日本語</option>
                <option value="ko">한국어</option>
                <option value="en">English</option>
              </select>
            </div>

            <div className="mt-4 flex items-center justify-between gap-2">
              <div className="min-w-0">
                <div className="truncate text-xs text-zinc-500 dark:text-zinc-400">
                  {user ? user.email : 'Guest'}
                </div>
              </div>
              {user ? (
                <button
                  type="button"
                  onClick={async () => {
                    await logout()
                    navigate('/')
                  }}
                  className="h-8 rounded-lg border border-zinc-200 bg-white px-2 text-xs hover:bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-900 dark:hover:bg-zinc-900/60"
                >
                  {t('nav.logout')}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={() => navigate('/login')}
                  className="h-8 rounded-lg border border-zinc-200 bg-white px-2 text-xs hover:bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-900 dark:hover:bg-zinc-900/60"
                >
                  {t('nav.login')}
                </button>
              )}
            </div>
          </div>
        </aside>

        <main className="rounded-2xl border border-zinc-200 bg-white p-5 shadow-sm dark:border-zinc-800 dark:bg-zinc-900/40">
          {children}
        </main>
      </div>
    </div>
  )
}

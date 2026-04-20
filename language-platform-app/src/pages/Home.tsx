import { useNavigate } from 'react-router-dom'
import { useI18n } from '@/i18n/useI18n'
import { useSettingsStore } from '@/stores/useSettingsStore'
import { useAuthStore } from '@/stores/useAuthStore'

export default function Home() {
  const { t } = useI18n()
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const user = useAuthStore((s) => s.user)
  const navigate = useNavigate()

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-gradient-to-br from-zinc-900 to-zinc-800 p-6 text-white">
        <div className="text-xs uppercase tracking-widest text-white/70">Immersive learning</div>
        <div className="mt-2 text-3xl font-semibold">{t('app.name')}</div>
        <div className="mt-2 max-w-2xl text-sm text-white/80">{t('landing.tagline')}</div>
        <div className="mt-5 flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => navigate(user ? '/catalog' : '/register')}
            className="h-10 rounded-xl bg-white px-4 text-sm font-medium text-zinc-900 hover:bg-white/90"
          >
            {t('landing.start')}
          </button>
          <button
            type="button"
            onClick={() => navigate('/catalog')}
            className="h-10 rounded-xl border border-white/25 bg-white/0 px-4 text-sm text-white hover:bg-white/10"
          >
            {t('nav.catalog')}
          </button>
        </div>
      </div>

      <div className="grid gap-3 md:grid-cols-3">
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
          <div className="text-sm font-medium">Leveled system</div>
          <div className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">A1 → C2 progression per target language.</div>
        </div>
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
          <div className="text-sm font-medium">Interactive modules</div>
          <div className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">
            Vocabulary, grammar, listening, and shadowing in one flow.
          </div>
        </div>
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
          <div className="text-sm font-medium">Momentum</div>
          <div className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">
            XP, achievements, and community challenges to keep you going.
          </div>
        </div>
      </div>

      <div className="rounded-2xl border border-zinc-200 p-4 dark:border-zinc-800">
        <div className="flex items-center justify-between gap-3">
          <div>
            <div className="text-sm font-medium">Today’s focus</div>
            <div className="text-xs text-zinc-500 dark:text-zinc-400">
              Recommended starter course for {targetLanguage.toUpperCase()}.
            </div>
          </div>
          <button
            type="button"
            onClick={() => navigate('/path')}
            className="h-9 rounded-xl bg-zinc-900 px-3 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
          >
            {t('nav.path')}
          </button>
        </div>
      </div>
    </div>
  )
}

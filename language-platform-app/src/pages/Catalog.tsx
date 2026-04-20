import { useMemo } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useI18n } from '@/i18n/useI18n'
import { COURSES } from '@/data/courses'
import { useSettingsStore } from '@/stores/useSettingsStore'
import { useProfileStore } from '@/stores/useProfileStore'
import { useLearningStore } from '@/stores/useLearningStore'

export default function Catalog() {
  const { t } = useI18n()
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const startingLevelByLanguage = useProfileStore((s) => s.startingLevelByLanguage)
  const enroll = useLearningStore((s) => s.enroll)
  const enrolled = useLearningStore((s) => s.enrolledCourseIds)
  const navigate = useNavigate()

  const recommendedLevel = startingLevelByLanguage[targetLanguage] ?? 'A1'

  const courses = useMemo(() => {
    return COURSES.filter((c) => c.language === targetLanguage).sort((a, b) => a.level.localeCompare(b.level))
  }, [targetLanguage])

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between gap-3">
        <div>
          <div className="text-xl font-semibold">{t('catalog.title')}</div>
          <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Target: <span className="font-medium">{targetLanguage.toUpperCase()}</span>, suggested start:{' '}
            <span className="font-medium">{recommendedLevel}</span>
          </div>
        </div>
        <Link
          to="/onboarding"
          className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
        >
          Adjust
        </Link>
      </div>

      <div className="grid gap-3 md:grid-cols-2">
        {courses.map((c) => {
          const isEnrolled = enrolled.includes(c.id)
          return (
            <div
              key={c.id}
              className="rounded-2xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900/20"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0">
                  <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">
                    {c.language.toUpperCase()} · {c.level}
                  </div>
                  <div className="mt-1 truncate text-lg font-semibold">{c.title}</div>
                  <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{c.description}</div>
                </div>
                <button
                  type="button"
                  onClick={() => {
                    enroll(c.id)
                    navigate(`/courses/${c.id}`)
                  }}
                  className={[
                    'h-9 shrink-0 rounded-xl px-3 text-sm font-medium',
                    isEnrolled
                      ? 'border border-zinc-200 bg-white hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-950/40 dark:hover:bg-zinc-950/70'
                      : 'bg-zinc-900 text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90',
                  ].join(' ')}
                >
                  {isEnrolled ? 'Open' : t('catalog.enroll')}
                </button>
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}


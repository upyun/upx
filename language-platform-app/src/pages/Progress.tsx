import { useMemo } from 'react'
import { useI18n } from '@/i18n/useI18n'
import { useLearningStore } from '@/stores/useLearningStore'

function isoDay(d: Date) {
  return d.toISOString().slice(0, 10)
}

function computeStreak(days: string[]) {
  const set = new Set(days)
  let streak = 0
  const cursor = new Date()
  for (;;) {
    const key = isoDay(cursor)
    if (!set.has(key)) break
    streak += 1
    cursor.setDate(cursor.getDate() - 1)
  }
  return streak
}

export default function Progress() {
  const { t } = useI18n()
  const xp = useLearningStore((s) => s.xp)
  const completedLessons = useLearningStore((s) => s.completedLessons)
  const activityDays = useLearningStore((s) => s.activityDays)

  const completedCount = Object.keys(completedLessons).length
  const streak = useMemo(() => computeStreak(activityDays), [activityDays])

  return (
    <div className="space-y-4">
      <div className="text-xl font-semibold">{t('progress.title')}</div>

      <div className="grid gap-3 md:grid-cols-3">
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
          <div className="text-xs text-zinc-500 dark:text-zinc-400">XP</div>
          <div className="mt-1 text-2xl font-semibold">{xp}</div>
        </div>
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
          <div className="text-xs text-zinc-500 dark:text-zinc-400">Streak</div>
          <div className="mt-1 text-2xl font-semibold">{streak} days</div>
        </div>
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
          <div className="text-xs text-zinc-500 dark:text-zinc-400">Lessons completed</div>
          <div className="mt-1 text-2xl font-semibold">{completedCount}</div>
        </div>
      </div>

      <div className="rounded-2xl border border-zinc-200 p-4 dark:border-zinc-800">
        <div className="text-sm font-medium">Activity (last 14 days)</div>
        <div className="mt-3 grid grid-cols-14 gap-1">
          {Array.from({ length: 14 }).map((_, i) => {
            const d = new Date()
            d.setDate(d.getDate() - (13 - i))
            const key = isoDay(d)
            const active = activityDays.includes(key)
            return (
              <div
                key={key}
                title={key}
                className={[
                  'h-6 rounded-md border',
                  active
                    ? 'border-zinc-900 bg-zinc-900 dark:border-white dark:bg-white'
                    : 'border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20',
                ].join(' ')}
              />
            )
          })}
        </div>
      </div>
    </div>
  )
}


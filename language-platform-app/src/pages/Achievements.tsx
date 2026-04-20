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

type Badge = {
  code: string
  title: string
  desc: string
  earned: boolean
}

export default function Achievements() {
  const { t } = useI18n()
  const xp = useLearningStore((s) => s.xp)
  const completedLessons = useLearningStore((s) => s.completedLessons)
  const completedCount = Object.keys(completedLessons).length
  const activityDays = useLearningStore((s) => s.activityDays)

  const streak = useMemo(() => computeStreak(activityDays), [activityDays])

  const badges = useMemo<Badge[]>(() => {
    return [
      {
        code: 'first_lesson',
        title: 'First Lesson',
        desc: 'Complete your first lesson.',
        earned: completedCount >= 1,
      },
      {
        code: 'streak_3',
        title: '3-Day Streak',
        desc: 'Learn 3 days in a row.',
        earned: streak >= 3,
      },
      {
        code: 'xp_100',
        title: 'XP 100',
        desc: 'Earn 100 XP.',
        earned: xp >= 100,
      },
      {
        code: 'lessons_5',
        title: 'Finisher',
        desc: 'Complete 5 lessons.',
        earned: completedCount >= 5,
      },
    ]
  }, [completedCount, streak, xp])

  return (
    <div className="space-y-4">
      <div className="text-xl font-semibold">{t('achievements.title')}</div>

      <div className="grid gap-3 md:grid-cols-2">
        {badges.map((b) => (
          <div
            key={b.code}
            className={[
              'rounded-2xl border p-4',
              b.earned
                ? 'border-emerald-200 bg-emerald-50 text-emerald-900 dark:border-emerald-900/40 dark:bg-emerald-900/10 dark:text-emerald-200'
                : 'border-zinc-200 bg-zinc-50 text-zinc-700 dark:border-zinc-800 dark:bg-zinc-900/20 dark:text-zinc-200',
            ].join(' ')}
          >
            <div className="text-sm font-semibold">{b.title}</div>
            <div className="mt-1 text-xs opacity-80">{b.desc}</div>
            <div className="mt-3 text-xs">{b.earned ? 'Earned' : 'Locked'}</div>
          </div>
        ))}
      </div>

      <div className="rounded-2xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
        <div className="text-sm font-medium">Leaderboard (local)</div>
        <div className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
          You: <span className="font-semibold text-zinc-900 dark:text-white">{xp}</span> XP
        </div>
      </div>
    </div>
  )
}

import { useMemo, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { useProfileStore, type LearningGoal } from '@/stores/useProfileStore'
import { useSettingsStore } from '@/stores/useSettingsStore'
import type { CourseLevel } from '@/data/courses'

const GOALS: Array<{ id: LearningGoal; title: string; desc: string }> = [
  { id: 'conversation', title: 'Conversation', desc: 'Speak with confidence in daily life.' },
  { id: 'travel', title: 'Travel', desc: 'Handle common situations abroad.' },
  { id: 'work', title: 'Work', desc: 'Communicate in professional contexts.' },
  { id: 'exam', title: 'Exam', desc: 'Build structured knowledge for tests.' },
]

const LEVELS: CourseLevel[] = ['A1', 'A2', 'B1', 'B2', 'C1', 'C2']

export default function Onboarding() {
  const [params] = useSearchParams()
  const targetFromQuery = params.get('target')
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const effectiveTarget = (targetFromQuery as any) || targetLanguage

  const goal = useProfileStore((s) => s.goal)
  const setGoal = useProfileStore((s) => s.setGoal)
  const startingLevelByLanguage = useProfileStore((s) => s.startingLevelByLanguage)
  const setStartingLevel = useProfileStore((s) => s.setStartingLevel)

  const currentLevel = startingLevelByLanguage[effectiveTarget] ?? 'A1'
  const [quizScore, setQuizScore] = useState<number | null>(null)
  const navigate = useNavigate()

  const suggested = useMemo<CourseLevel>(() => {
    if (quizScore === null) return currentLevel
    if (quizScore <= 2) return 'A1'
    if (quizScore <= 4) return 'A2'
    return 'B1'
  }, [quizScore, currentLevel])

  return (
    <div className="space-y-5">
      <div>
        <div className="text-xl font-semibold">Onboarding</div>
        <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
          Pick a goal and a starting level. This can be changed later.
        </div>
      </div>

      <div className="space-y-2">
        <div className="text-sm font-medium">Goal</div>
        <div className="grid gap-2 md:grid-cols-2">
          {GOALS.map((g) => (
            <button
              key={g.id}
              type="button"
              onClick={() => setGoal(g.id)}
              className={[
                'rounded-2xl border px-4 py-3 text-left transition-colors',
                goal === g.id
                  ? 'border-zinc-900 bg-zinc-900 text-white dark:border-white dark:bg-white dark:text-zinc-900'
                  : 'border-zinc-200 bg-white hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40',
              ].join(' ')}
            >
              <div className="text-sm font-medium">{g.title}</div>
              <div className="mt-1 text-xs opacity-80">{g.desc}</div>
            </button>
          ))}
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between gap-3">
          <div>
            <div className="text-sm font-medium">Placement check (quick)</div>
            <div className="text-xs text-zinc-500 dark:text-zinc-400">A tiny quiz to suggest a starting level.</div>
          </div>
          <button
            type="button"
            onClick={() => setQuizScore(Math.floor(Math.random() * 7))}
            className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
          >
            Run
          </button>
        </div>
        <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 text-sm dark:border-zinc-800 dark:bg-zinc-900/20">
          Suggested level: <span className="font-semibold">{suggested}</span>
          {quizScore !== null ? (
            <span className="ml-2 text-xs text-zinc-500 dark:text-zinc-400">(score: {quizScore}/6)</span>
          ) : null}
        </div>
      </div>

      <div className="space-y-2">
        <div className="text-sm font-medium">Starting level</div>
        <div className="flex flex-wrap gap-2">
          {LEVELS.map((lvl) => (
            <button
              key={lvl}
              type="button"
              onClick={() => setStartingLevel(effectiveTarget, lvl)}
              className={[
                'h-9 rounded-xl border px-3 text-sm',
                (startingLevelByLanguage[effectiveTarget] ?? 'A1') === lvl
                  ? 'border-zinc-900 bg-zinc-900 text-white dark:border-white dark:bg-white dark:text-zinc-900'
                  : 'border-zinc-200 bg-white hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40',
              ].join(' ')}
            >
              {lvl}
            </button>
          ))}
        </div>
      </div>

      <div className="flex justify-end">
        <button
          type="button"
          onClick={() => {
            setStartingLevel(effectiveTarget, suggested)
            navigate('/catalog')
          }}
          className="h-10 rounded-xl bg-zinc-900 px-4 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
        >
          Continue
        </button>
      </div>
    </div>
  )
}


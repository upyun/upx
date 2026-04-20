import { useState } from 'react'
import type { LessonStep } from '@/data/courses'

export default function VocabStep({
  step,
  onXp,
}: {
  step: Extract<LessonStep, { kind: 'vocab' }>
  onXp: (xp: number) => void
}) {
  const [idx, setIdx] = useState(0)
  const [revealed, setRevealed] = useState(false)
  const current = step.items[idx]

  return (
    <div className="space-y-4">
      <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">Vocabulary</div>
      <div className="rounded-2xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900/20">
        <div className="text-2xl font-semibold">{current.term}</div>
        {current.reading ? <div className="mt-1 text-sm text-zinc-500">{current.reading}</div> : null}
        <div className="mt-3">
          {revealed ? (
            <div className="rounded-xl bg-zinc-50 p-3 text-sm text-zinc-800 dark:bg-zinc-950/40 dark:text-zinc-100">
              {current.meaning}
            </div>
          ) : (
            <button
              type="button"
              onClick={() => setRevealed(true)}
              className="h-10 rounded-xl bg-zinc-900 px-4 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
            >
              Reveal meaning
            </button>
          )}
        </div>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="text-xs text-zinc-500 dark:text-zinc-400">
          Card {idx + 1}/{step.items.length}
        </div>
        <div className="flex gap-2">
          <button
            type="button"
            onClick={() => {
              onXp(revealed ? 6 : 3)
              setRevealed(false)
              setIdx((v) => Math.min(step.items.length - 1, v + 1))
            }}
            className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
          >
            I know it
          </button>
          <button
            type="button"
            onClick={() => {
              onXp(2)
              setRevealed(false)
              setIdx((v) => Math.min(step.items.length - 1, v + 1))
            }}
            className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
          >
            Review later
          </button>
        </div>
      </div>
    </div>
  )
}


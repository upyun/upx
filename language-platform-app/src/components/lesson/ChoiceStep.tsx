import { useState } from 'react'

export default function ChoiceStep({
  title,
  prompt,
  choices,
  answerIndex,
  explanation,
  onXp,
}: {
  title: string
  prompt: string
  choices: string[]
  answerIndex: number
  explanation?: string
  onXp: (xp: number) => void
}) {
  const [selected, setSelected] = useState<number | null>(null)
  const [checked, setChecked] = useState(false)
  const correct = checked && selected === answerIndex

  return (
    <div className="space-y-4">
      <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">{title}</div>
      <div className="rounded-2xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900/20">
        <div className="text-sm font-medium">{prompt}</div>
        <div className="mt-3 grid gap-2">
          {choices.map((c, i) => (
            <button
              key={`${c}_${i}`}
              type="button"
              onClick={() => {
                if (checked) return
                setSelected(i)
              }}
              className={[
                'rounded-xl border px-4 py-3 text-left text-sm transition-colors',
                selected === i
                  ? 'border-zinc-900 bg-zinc-900 text-white dark:border-white dark:bg-white dark:text-zinc-900'
                  : 'border-zinc-200 bg-white hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/10 dark:hover:bg-zinc-900/30',
              ].join(' ')}
            >
              {c}
            </button>
          ))}
        </div>

        <div className="mt-4 flex flex-wrap items-center gap-2">
          <button
            type="button"
            disabled={selected === null || checked}
            onClick={() => {
              setChecked(true)
              if (selected === answerIndex) onXp(12)
              else onXp(4)
            }}
            className="h-10 rounded-xl bg-zinc-900 px-4 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
          >
            Check
          </button>
          {checked ? (
            <div className={correct ? 'text-sm text-emerald-600' : 'text-sm text-rose-600'}>
              {correct ? 'Correct' : `Answer: ${choices[answerIndex]}`}
            </div>
          ) : null}
        </div>

        {checked && explanation ? (
          <div className="mt-3 rounded-xl bg-zinc-50 p-3 text-sm text-zinc-700 dark:bg-zinc-950/40 dark:text-zinc-200">
            {explanation}
          </div>
        ) : null}
      </div>
    </div>
  )
}


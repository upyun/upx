import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { getLesson } from '@/data/courses'
import ChoiceStep from '@/components/lesson/ChoiceStep'
import ShadowingStep from '@/components/lesson/ShadowingStep'
import VocabStep from '@/components/lesson/VocabStep'
import { speak, ttsLang } from '@/components/lesson/tts'
import { useLearningStore } from '@/stores/useLearningStore'
import { useSettingsStore } from '@/stores/useSettingsStore'
import { useI18n } from '@/i18n/useI18n'

export default function Learn() {
  const { t } = useI18n()
  const { lessonId } = useParams()
  const lesson = useMemo(() => (lessonId ? getLesson(lessonId) : null), [lessonId])
  const completedLessons = useLearningStore((s) => s.completedLessons)
  const completeLesson = useLearningStore((s) => s.completeLesson)
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const navigate = useNavigate()

  const [stepIndex, setStepIndex] = useState(0)
  const [stepXp, setStepXp] = useState(0)

  useEffect(() => {
    setStepIndex(0)
    setStepXp(0)
  }, [lessonId])

  if (!lesson) {
    return <div className="text-sm text-zinc-500 dark:text-zinc-400">Lesson not found.</div>
  }

  const total = lesson.steps.length
  const step = lesson.steps[Math.min(stepIndex, total - 1)]
  const doneAt = completedLessons[lesson.id]
  const lang = ttsLang(targetLanguage)

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">
            Lesson · {lesson.estimatedMinutes} min
          </div>
          <div className="mt-1 text-2xl font-semibold">{lesson.title}</div>
          <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Step {stepIndex + 1}/{total} · XP this lesson: {stepXp} {doneAt ? `· Completed` : ''}
          </div>
        </div>
        <Link
          to={`/courses/${lesson.courseId}`}
          className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
        >
          Back
        </Link>
      </div>

      <div className="h-2 overflow-hidden rounded-full bg-zinc-100 dark:bg-zinc-800">
        <div
          className="h-full bg-zinc-900 dark:bg-white"
          style={{ width: `${((stepIndex + 1) / total) * 100}%` }}
        />
      </div>

      {step.kind === 'vocab' ? <VocabStep step={step} onXp={(xp) => setStepXp((v) => v + xp)} /> : null}
      {step.kind === 'grammar' ? (
        <ChoiceStep
          title="Grammar"
          prompt={step.prompt}
          choices={step.choices}
          answerIndex={step.answerIndex}
          explanation={step.explanation}
          onXp={(xp) => setStepXp((v) => v + xp)}
        />
      ) : null}
      {step.kind === 'listening' ? (
        <div className="space-y-4">
          <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">Listening</div>
          <div className="rounded-2xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900/20">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div className="text-sm font-medium">Audio</div>
              <button
                type="button"
                onClick={() => speak(step.text, lang)}
                className="h-9 rounded-xl bg-zinc-900 px-3 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
              >
                Play
              </button>
            </div>
            <div className="mt-3 rounded-xl bg-zinc-50 p-3 text-sm text-zinc-700 dark:bg-zinc-950/40 dark:text-zinc-200">
              {step.text}
            </div>
          </div>
          <ChoiceStep
            title="Comprehension"
            prompt={step.question}
            choices={step.choices}
            answerIndex={step.answerIndex}
            onXp={(xp) => setStepXp((v) => v + xp)}
          />
        </div>
      ) : null}
      {step.kind === 'shadowing' ? (
        <ShadowingStep step={step} lang={lang} onXp={(xp) => setStepXp((v) => v + xp)} />
      ) : null}

      <div className="flex items-center justify-between gap-2">
        <button
          type="button"
          disabled={stepIndex === 0}
          onClick={() => setStepIndex((v) => Math.max(0, v - 1))}
          className="h-10 rounded-xl border border-zinc-200 bg-white px-4 text-sm hover:bg-zinc-50 disabled:cursor-not-allowed disabled:opacity-40 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
        >
          Back
        </button>

        {stepIndex < total - 1 ? (
          <button
            type="button"
            onClick={() => setStepIndex((v) => Math.min(total - 1, v + 1))}
            className="h-10 rounded-xl bg-zinc-900 px-4 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
          >
            {t('learn.next')}
          </button>
        ) : (
          <button
            type="button"
            onClick={() => {
              completeLesson(lesson.id, stepXp + 20)
              navigate('/progress')
            }}
            className="h-10 rounded-xl bg-zinc-900 px-4 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
          >
            Finish lesson
          </button>
        )}
      </div>
    </div>
  )
}

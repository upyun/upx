import { useMemo } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getCourse, getLessonsByCourse } from '@/data/courses'
import { useLearningStore } from '@/stores/useLearningStore'
import { useI18n } from '@/i18n/useI18n'

export default function Course() {
  const { t } = useI18n()
  const { courseId } = useParams()
  const course = useMemo(() => (courseId ? getCourse(courseId) : null), [courseId])
  const lessons = useMemo(() => (courseId ? getLessonsByCourse(courseId) : []), [courseId])
  const completedLessons = useLearningStore((s) => s.completedLessons)

  if (!course) {
    return <div className="text-sm text-zinc-500 dark:text-zinc-400">Course not found.</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between gap-3">
        <div>
          <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">
            {course.language.toUpperCase()} · {course.level}
          </div>
          <div className="mt-1 text-2xl font-semibold">{course.title}</div>
          <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{course.description}</div>
        </div>
        <Link
          to="/catalog"
          className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
        >
          Back
        </Link>
      </div>

      <div className="space-y-2">
        <div className="text-sm font-medium">{t('course.curriculum')}</div>
        <div className="grid gap-2">
          {lessons.map((l) => {
            const doneAt = completedLessons[l.id]
            return (
              <div
                key={l.id}
                className="flex items-center justify-between gap-3 rounded-2xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900/20"
              >
                <div className="min-w-0">
                  <div className="truncate text-sm font-medium">{l.title}</div>
                  <div className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">
                    {l.estimatedMinutes} min · {doneAt ? `Completed` : 'Not started'}
                  </div>
                </div>
                <Link
                  to={`/learn/${l.id}`}
                  className="h-9 rounded-xl bg-zinc-900 px-3 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
                >
                  {t('course.startLesson')}
                </Link>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}


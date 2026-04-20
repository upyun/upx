import { Link, useNavigate } from 'react-router-dom'
import { useI18n } from '@/i18n/useI18n'
import { COURSES, LESSONS } from '@/data/courses'
import { useLearningStore } from '@/stores/useLearningStore'
import { useSettingsStore } from '@/stores/useSettingsStore'
import { useProfileStore } from '@/stores/useProfileStore'

function recommendNext(params: {
  targetLanguage: 'en' | 'ja' | 'ko'
  enrolledCourseIds: string[]
  completedLessons: Record<string, string>
  suggestedLevel: string
}) {
  const enrolledLessons = LESSONS.filter((l) => params.enrolledCourseIds.includes(l.courseId))
  const incomplete = enrolledLessons.find((l) => !params.completedLessons[l.id])
  if (incomplete) {
    return {
      kind: 'lesson' as const,
      lessonId: incomplete.id,
      reason: 'Continue your enrolled course from the next unfinished lesson.',
    }
  }

  const candidateCourse =
    COURSES.find((c) => c.language === params.targetLanguage && c.level === params.suggestedLevel) ??
    COURSES.find((c) => c.language === params.targetLanguage) ??
    null
  if (!candidateCourse) return null

  const firstLesson = LESSONS.find((l) => l.courseId === candidateCourse.id) ?? null
  if (!firstLesson) return null

  return {
    kind: 'course' as const,
    courseId: candidateCourse.id,
    lessonId: firstLesson.id,
    reason: 'Start the most suitable leveled course for your target language.',
  }
}

export default function Path() {
  const { t } = useI18n()
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const suggestedLevel = useProfileStore((s) => s.startingLevelByLanguage[targetLanguage] ?? 'A1')
  const enrolledCourseIds = useLearningStore((s) => s.enrolledCourseIds)
  const completedLessons = useLearningStore((s) => s.completedLessons)
  const enroll = useLearningStore((s) => s.enroll)
  const navigate = useNavigate()

  const rec = recommendNext({ targetLanguage, enrolledCourseIds, completedLessons, suggestedLevel })

  if (!rec) {
    return <div className="text-sm text-zinc-500 dark:text-zinc-400">No recommendation available.</div>
  }

  return (
    <div className="space-y-4">
      <div className="flex items-end justify-between gap-3">
        <div>
          <div className="text-xl font-semibold">{t('path.title')}</div>
          <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Target: {targetLanguage.toUpperCase()} · Suggested start: {suggestedLevel}
          </div>
        </div>
        <Link
          to="/catalog"
          className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
        >
          Browse courses
        </Link>
      </div>

      <div className="rounded-2xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900/20">
        <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">Next up</div>
        <div className="mt-2 text-lg font-semibold">
          {rec.kind === 'lesson' ? `Lesson: ${rec.lessonId}` : `Course: ${rec.courseId}`}
        </div>
        <div className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">{rec.reason}</div>

        <div className="mt-4 flex flex-wrap gap-2">
          {rec.kind === 'course' ? (
            <button
              type="button"
              onClick={() => {
                enroll(rec.courseId)
                navigate(`/courses/${rec.courseId}`)
              }}
              className="h-10 rounded-xl bg-zinc-900 px-4 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
            >
              Enroll and open
            </button>
          ) : null}
          <button
            type="button"
            onClick={() => navigate(`/learn/${rec.lessonId}`)}
            className="h-10 rounded-xl border border-zinc-200 bg-white px-4 text-sm hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
          >
            Start now
          </button>
        </div>
      </div>
    </div>
  )
}


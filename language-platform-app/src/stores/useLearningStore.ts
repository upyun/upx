import { create } from 'zustand'

export type LearningState = {
  userId: string | null
  enrolledCourseIds: string[]
  completedLessons: Record<string, string>
  xp: number
  activityDays: string[]
  loadForUser: (userId: string | null) => void
  enroll: (courseId: string) => void
  completeLesson: (lessonId: string, xpEarned: number) => void
}

type Stored = Omit<LearningState, 'loadForUser' | 'enroll' | 'completeLesson'>

function storageKey(userId: string | null) {
  return `lf_learning_v1_${userId ?? 'anon'}`
}

function todayIsoDay() {
  return new Date().toISOString().slice(0, 10)
}

function readStored(key: string): Stored | null {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return null
    const parsed = JSON.parse(raw) as Partial<Stored>
    return {
      userId: (parsed.userId as string) ?? null,
      enrolledCourseIds: parsed.enrolledCourseIds ?? [],
      completedLessons: parsed.completedLessons ?? {},
      xp: typeof parsed.xp === 'number' ? parsed.xp : 0,
      activityDays: parsed.activityDays ?? [],
    }
  } catch {
    return null
  }
}

function writeStored(key: string, value: Stored) {
  localStorage.setItem(key, JSON.stringify(value))
}

export const useLearningStore = create<LearningState>((set, get) => ({
  userId: null,
  enrolledCourseIds: [],
  completedLessons: {},
  xp: 0,
  activityDays: [],
  loadForUser: (userId) => {
    const key = storageKey(userId)
    const stored = readStored(key)
    const next: Stored = stored ?? {
      userId,
      enrolledCourseIds: [],
      completedLessons: {},
      xp: 0,
      activityDays: [],
    }
    next.userId = userId
    set(next)
    writeStored(key, next)
  },
  enroll: (courseId) => {
    const state = get()
    if (state.enrolledCourseIds.includes(courseId)) return
    const next: Stored = {
      userId: state.userId,
      enrolledCourseIds: [...state.enrolledCourseIds, courseId],
      completedLessons: state.completedLessons,
      xp: state.xp,
      activityDays: state.activityDays,
    }
    set(next)
    writeStored(storageKey(state.userId), next)
  },
  completeLesson: (lessonId, xpEarned) => {
    const state = get()
    const day = todayIsoDay()
    const nextActivity = state.activityDays.includes(day)
      ? state.activityDays
      : [...state.activityDays, day].slice(-365)
    const next: Stored = {
      userId: state.userId,
      enrolledCourseIds: state.enrolledCourseIds,
      completedLessons: { ...state.completedLessons, [lessonId]: new Date().toISOString() },
      xp: state.xp + Math.max(0, Math.floor(xpEarned)),
      activityDays: nextActivity,
    }
    set(next)
    writeStored(storageKey(state.userId), next)
  },
}))


import { create } from 'zustand'
import type { TargetLanguage } from '@/stores/useSettingsStore'
import type { CourseLevel } from '@/data/courses'

export type LearningGoal = 'travel' | 'work' | 'exam' | 'conversation'

type ProfileState = {
  goal: LearningGoal | null
  startingLevelByLanguage: Partial<Record<TargetLanguage, CourseLevel>>
  setGoal: (goal: LearningGoal) => void
  setStartingLevel: (language: TargetLanguage, level: CourseLevel) => void
}

const STORAGE_KEY = 'lf_profile_v1'

function readStored(): Pick<ProfileState, 'goal' | 'startingLevelByLanguage'> | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as Partial<ProfileState>
    return {
      goal: (parsed.goal as any) ?? null,
      startingLevelByLanguage: parsed.startingLevelByLanguage ?? {},
    }
  } catch {
    return null
  }
}

function writeStored(value: Pick<ProfileState, 'goal' | 'startingLevelByLanguage'>) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(value))
}

const stored = typeof window !== 'undefined' ? readStored() : null

export const useProfileStore = create<ProfileState>((set, get) => ({
  goal: stored?.goal ?? null,
  startingLevelByLanguage: stored?.startingLevelByLanguage ?? {},
  setGoal: (goal) => {
    set({ goal })
    writeStored({ goal, startingLevelByLanguage: get().startingLevelByLanguage })
  },
  setStartingLevel: (language, level) => {
    const next = { ...get().startingLevelByLanguage, [language]: level }
    set({ startingLevelByLanguage: next })
    writeStored({ goal: get().goal, startingLevelByLanguage: next })
  },
}))


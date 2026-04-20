import { create } from 'zustand'
import type { Locale } from '@/i18n/strings'

export type TargetLanguage = 'en' | 'ja' | 'ko'

type SettingsState = {
  uiLocale: Locale
  targetLanguage: TargetLanguage
  setUiLocale: (locale: Locale) => void
  setTargetLanguage: (lang: TargetLanguage) => void
}

const STORAGE_KEY = 'lf_settings_v1'

function readStored(): Pick<SettingsState, 'uiLocale' | 'targetLanguage'> | null {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as Partial<SettingsState>
    const uiLocale = (parsed.uiLocale as Locale) ?? 'en'
    const targetLanguage = (parsed.targetLanguage as TargetLanguage) ?? 'ja'
    return { uiLocale, targetLanguage }
  } catch {
    return null
  }
}

function writeStored(next: Pick<SettingsState, 'uiLocale' | 'targetLanguage'>) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(next))
}

const stored = typeof window !== 'undefined' ? readStored() : null

export const useSettingsStore = create<SettingsState>((set, get) => ({
  uiLocale: stored?.uiLocale ?? 'en',
  targetLanguage: stored?.targetLanguage ?? 'ja',
  setUiLocale: (uiLocale) => {
    set({ uiLocale })
    writeStored({ uiLocale, targetLanguage: get().targetLanguage })
  },
  setTargetLanguage: (targetLanguage) => {
    set({ targetLanguage })
    writeStored({ uiLocale: get().uiLocale, targetLanguage })
  },
}))


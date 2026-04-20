import { STRINGS, type I18nKey, type Locale } from '@/i18n/strings'
import { useSettingsStore } from '@/stores/useSettingsStore'

export function useI18n() {
  const uiLocale = useSettingsStore((s) => s.uiLocale)
  const setUiLocale = useSettingsStore((s) => s.setUiLocale)

  const t = (key: I18nKey, localeOverride?: Locale) => {
    const locale = localeOverride ?? uiLocale
    const dict = STRINGS[locale] ?? STRINGS.en
    return dict[key] ?? STRINGS.en[key] ?? key
  }

  return { t, uiLocale, setUiLocale }
}


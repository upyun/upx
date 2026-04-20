export function ttsLang(target: 'en' | 'ja' | 'ko') {
  if (target === 'ja') return 'ja-JP'
  if (target === 'ko') return 'ko-KR'
  return 'en-US'
}

export function speak(text: string, lang: string) {
  if (!('speechSynthesis' in window)) return
  const u = new SpeechSynthesisUtterance(text)
  u.lang = lang
  window.speechSynthesis.cancel()
  window.speechSynthesis.speak(u)
}


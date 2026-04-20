import { useEffect, useRef, useState } from 'react'
import type { LessonStep } from '@/data/courses'
import { speak } from '@/components/lesson/tts'

export default function ShadowingStep({
  step,
  onXp,
  lang,
}: {
  step: Extract<LessonStep, { kind: 'shadowing' }>
  onXp: (xp: number) => void
  lang: string
}) {
  const [recording, setRecording] = useState(false)
  const [audioUrl, setAudioUrl] = useState<string | null>(null)
  const recorderRef = useRef<MediaRecorder | null>(null)
  const chunksRef = useRef<Blob[]>([])

  useEffect(() => {
    return () => {
      if (audioUrl) URL.revokeObjectURL(audioUrl)
    }
  }, [audioUrl])

  const canRecord = typeof window !== 'undefined' && 'MediaRecorder' in window

  return (
    <div className="space-y-4">
      <div className="text-xs uppercase tracking-widest text-zinc-500 dark:text-zinc-400">Shadowing</div>
      <div className="rounded-2xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900/20">
        <div className="text-sm font-medium">Reference</div>
        <div className="mt-2 rounded-xl bg-zinc-50 p-3 text-sm dark:bg-zinc-950/40">{step.text}</div>
        <div className="mt-3 flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => speak(step.text, lang)}
            className="h-9 rounded-xl bg-zinc-900 px-3 text-sm font-medium text-white hover:bg-zinc-800 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
          >
            Play
          </button>
          <div className="text-xs text-zinc-500 dark:text-zinc-400">{step.tip}</div>
        </div>

        <div className="mt-4 border-t border-zinc-200 pt-4 dark:border-zinc-800">
          <div className="text-sm font-medium">Record</div>
          {!canRecord ? (
            <div className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">Recording is not supported in this browser.</div>
          ) : (
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <button
                type="button"
                disabled={recording}
                onClick={async () => {
                  const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
                  const recorder = new MediaRecorder(stream)
                  recorderRef.current = recorder
                  chunksRef.current = []
                  recorder.ondataavailable = (e) => {
                    if (e.data.size) chunksRef.current.push(e.data)
                  }
                  recorder.onstop = () => {
                    stream.getTracks().forEach((t) => t.stop())
                    const blob = new Blob(chunksRef.current, { type: 'audio/webm' })
                    const url = URL.createObjectURL(blob)
                    setAudioUrl((prev) => {
                      if (prev) URL.revokeObjectURL(prev)
                      return url
                    })
                    onXp(10)
                  }
                  recorder.start()
                  setRecording(true)
                }}
                className="h-9 rounded-xl border border-zinc-200 bg-white px-3 text-sm hover:bg-zinc-50 disabled:cursor-not-allowed disabled:opacity-40 dark:border-zinc-800 dark:bg-zinc-900/20 dark:hover:bg-zinc-900/40"
              >
                Start
              </button>
              <button
                type="button"
                disabled={!recording}
                onClick={() => {
                  recorderRef.current?.stop()
                  setRecording(false)
                }}
                className="h-9 rounded-xl bg-zinc-900 px-3 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
              >
                Stop
              </button>
              {recording ? <div className="text-xs text-rose-600">Recording…</div> : null}
            </div>
          )}

          {audioUrl ? (
            <div className="mt-3">
              <audio controls src={audioUrl} className="w-full" />
            </div>
          ) : null}
        </div>
      </div>
    </div>
  )
}


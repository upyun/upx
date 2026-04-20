import { useEffect, useMemo, useState } from 'react'
import { useI18n } from '@/i18n/useI18n'
import { useAuthStore } from '@/stores/useAuthStore'
import { useCommunityStore } from '@/stores/useCommunityStore'
import { useSettingsStore } from '@/stores/useSettingsStore'

export default function Community() {
  const { t } = useI18n()
  const user = useAuthStore((s) => s.user)
  const targetLanguage = useSettingsStore((s) => s.targetLanguage)
  const loadForUser = useCommunityStore((s) => s.loadForUser)
  const posts = useCommunityStore((s) => s.posts)
  const createPost = useCommunityStore((s) => s.createPost)
  const toggleLike = useCommunityStore((s) => s.toggleLike)

  const [body, setBody] = useState('')

  useEffect(() => {
    loadForUser(user?.id ?? null)
  }, [loadForUser, user?.id])

  const displayPosts = useMemo(() => posts.slice(0, 30), [posts])

  return (
    <div className="space-y-4">
      <div className="text-xl font-semibold">{t('community.title')}</div>

      <div className="rounded-2xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900/20">
        <div className="text-sm font-medium">Post</div>
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          className="mt-2 h-24 w-full resize-none rounded-xl border border-zinc-200 bg-white p-3 text-sm outline-none focus:border-zinc-400 dark:border-zinc-800 dark:bg-zinc-950/40"
          placeholder="Share a quick win, ask a question, or start a challenge."
        />
        <div className="mt-2 flex justify-end">
          <button
            type="button"
            disabled={body.trim().length < 2}
            onClick={() => {
              createPost({
                author: user?.email ?? 'Guest',
                language: targetLanguage,
                body: body.trim(),
              })
              setBody('')
            }}
            className="h-9 rounded-xl bg-zinc-900 px-3 text-sm font-medium text-white hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-40 dark:bg-white dark:text-zinc-900 dark:hover:bg-white/90"
          >
            Publish
          </button>
        </div>
      </div>

      <div className="space-y-2">
        {displayPosts.length === 0 ? (
          <div className="rounded-2xl border border-zinc-200 bg-zinc-50 p-4 text-sm text-zinc-500 dark:border-zinc-800 dark:bg-zinc-900/20 dark:text-zinc-400">
            No posts yet. Start one above.
          </div>
        ) : null}

        {displayPosts.map((p) => (
          <div
            key={p.id}
            className="rounded-2xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900/20"
          >
            <div className="flex items-start justify-between gap-3">
              <div className="min-w-0">
                <div className="text-xs text-zinc-500 dark:text-zinc-400">
                  {p.author} · {p.language.toUpperCase()} · {new Date(p.createdAt).toLocaleString()}
                </div>
                <div className="mt-2 whitespace-pre-wrap text-sm">{p.body}</div>
              </div>
              <button
                type="button"
                onClick={() => toggleLike(p.id)}
                className="h-8 shrink-0 rounded-xl border border-zinc-200 bg-white px-3 text-xs hover:bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-950/40 dark:hover:bg-zinc-950/70"
              >
                +1 · {p.likes}
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}


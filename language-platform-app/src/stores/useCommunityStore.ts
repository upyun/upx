import { create } from 'zustand'
import { randomUUID } from '@/stores/uuid'

export type CommunityPost = {
  id: string
  author: string
  language: string
  body: string
  likes: number
  createdAt: string
}

type CommunityState = {
  userId: string | null
  posts: CommunityPost[]
  loadForUser: (userId: string | null) => void
  createPost: (payload: { author: string; language: string; body: string }) => void
  toggleLike: (postId: string) => void
}

type Stored = Pick<CommunityState, 'userId' | 'posts'>

function storageKey(userId: string | null) {
  return `lf_community_v1_${userId ?? 'anon'}`
}

function readStored(key: string): Stored | null {
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return null
    const parsed = JSON.parse(raw) as Partial<Stored>
    return {
      userId: (parsed.userId as string) ?? null,
      posts: parsed.posts ?? [],
    }
  } catch {
    return null
  }
}

function writeStored(key: string, value: Stored) {
  localStorage.setItem(key, JSON.stringify(value))
}

export const useCommunityStore = create<CommunityState>((set, get) => ({
  userId: null,
  posts: [],
  loadForUser: (userId) => {
    const key = storageKey(userId)
    const stored = readStored(key)
    const next: Stored = stored ?? { userId, posts: [] }
    next.userId = userId
    set(next)
    writeStored(key, next)
  },
  createPost: ({ author, language, body }) => {
    const state = get()
    const post: CommunityPost = {
      id: randomUUID(),
      author: author || 'Guest',
      language,
      body,
      likes: 0,
      createdAt: new Date().toISOString(),
    }
    const next: Stored = { userId: state.userId, posts: [post, ...state.posts].slice(0, 200) }
    set(next)
    writeStored(storageKey(state.userId), next)
  },
  toggleLike: (postId) => {
    const state = get()
    const nextPosts = state.posts.map((p) => (p.id === postId ? { ...p, likes: p.likes + 1 } : p))
    const next: Stored = { userId: state.userId, posts: nextPosts }
    set(next)
    writeStored(storageKey(state.userId), next)
  },
}))


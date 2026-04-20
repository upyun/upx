import type { TargetLanguage } from '@/stores/useSettingsStore'

export type CourseLevel = 'A1' | 'A2' | 'B1' | 'B2' | 'C1' | 'C2'

export type Course = {
  id: string
  language: TargetLanguage
  level: CourseLevel
  title: string
  description: string
}

export type LessonStep =
  | {
      id: string
      kind: 'vocab'
      items: Array<{ term: string; reading?: string; meaning: string }>
    }
  | {
      id: string
      kind: 'grammar'
      prompt: string
      choices: string[]
      answerIndex: number
      explanation: string
    }
  | {
      id: string
      kind: 'listening'
      text: string
      question: string
      choices: string[]
      answerIndex: number
    }
  | {
      id: string
      kind: 'shadowing'
      text: string
      tip: string
    }

export type Lesson = {
  id: string
  courseId: string
  title: string
  estimatedMinutes: number
  steps: LessonStep[]
}

export const COURSES: Course[] = [
  {
    id: 'ja-a1',
    language: 'ja',
    level: 'A1',
    title: 'Japanese A1: First Steps',
    description: 'Greetings, self-introduction, and survival phrases with pronunciation practice.',
  },
  {
    id: 'ko-a1',
    language: 'ko',
    level: 'A1',
    title: 'Korean A1: Daily Basics',
    description: 'Core patterns for everyday situations with listening and shadowing.',
  },
  {
    id: 'en-a1',
    language: 'en',
    level: 'A1',
    title: 'English A1: Everyday Start',
    description: 'Simple sentences, vocabulary, and listening drills for daily life.',
  },
]

export const LESSONS: Lesson[] = [
  {
    id: 'ja-a1-1',
    courseId: 'ja-a1',
    title: 'Hello & Nice to meet you',
    estimatedMinutes: 8,
    steps: [
      {
        id: 'ja-a1-1-vocab',
        kind: 'vocab',
        items: [
          { term: 'こんにちは', reading: 'konnichiwa', meaning: 'Hello' },
          { term: 'はじめまして', reading: 'hajimemashite', meaning: 'Nice to meet you' },
          { term: 'ありがとう', reading: 'arigatou', meaning: 'Thank you' },
        ],
      },
      {
        id: 'ja-a1-1-grammar',
        kind: 'grammar',
        prompt: 'Choose the best translation for: 「はじめまして」',
        choices: ['Goodbye', 'Nice to meet you', 'Thank you'],
        answerIndex: 1,
        explanation: '「はじめまして」 is used when meeting someone for the first time.',
      },
      {
        id: 'ja-a1-1-listening',
        kind: 'listening',
        text: 'はじめまして。こんにちは。',
        question: 'What does the speaker say?',
        choices: ['Nice to meet you. Hello.', 'Thank you. Hello.', 'Good morning. Goodbye.'],
        answerIndex: 0,
      },
      {
        id: 'ja-a1-1-shadowing',
        kind: 'shadowing',
        text: 'はじめまして。よろしくおねがいします。',
        tip: 'Try to match rhythm and vowel length. Record yourself and compare.',
      },
    ],
  },
  {
    id: 'ko-a1-1',
    courseId: 'ko-a1',
    title: 'Greetings',
    estimatedMinutes: 7,
    steps: [
      {
        id: 'ko-a1-1-vocab',
        kind: 'vocab',
        items: [
          { term: '안녕하세요', reading: 'annyeonghaseyo', meaning: 'Hello' },
          { term: '감사합니다', reading: 'gamsahamnida', meaning: 'Thank you' },
          { term: '만나서 반가워요', reading: 'mannaseo bangawoyo', meaning: 'Nice to meet you' },
        ],
      },
      {
        id: 'ko-a1-1-grammar',
        kind: 'grammar',
        prompt: 'Which phrase means “Thank you” in Korean?',
        choices: ['안녕하세요', '감사합니다', '미안해요'],
        answerIndex: 1,
        explanation: '감사합니다 is a polite way to say “Thank you”.',
      },
      {
        id: 'ko-a1-1-listening',
        kind: 'listening',
        text: '안녕하세요. 감사합니다.',
        question: 'What does the speaker say?',
        choices: ['Hello. Thank you.', 'Goodbye. Sorry.', 'Nice to meet you.'],
        answerIndex: 0,
      },
      {
        id: 'ko-a1-1-shadowing',
        kind: 'shadowing',
        text: '안녕하세요. 만나서 반가워요.',
        tip: 'Focus on final consonants and intonation. Record and replay.',
      },
    ],
  },
]

export function getCourse(courseId: string) {
  return COURSES.find((c) => c.id === courseId) ?? null
}

export function getLessonsByCourse(courseId: string) {
  return LESSONS.filter((l) => l.courseId === courseId)
}

export function getLesson(lessonId: string) {
  return LESSONS.find((l) => l.id === lessonId) ?? null
}


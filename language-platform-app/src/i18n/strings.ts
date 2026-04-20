export type Locale = 'en' | 'ja' | 'ko'

export type I18nKey =
  | 'app.name'
  | 'nav.home'
  | 'nav.catalog'
  | 'nav.path'
  | 'nav.progress'
  | 'nav.community'
  | 'nav.achievements'
  | 'nav.login'
  | 'nav.logout'
  | 'auth.email'
  | 'auth.password'
  | 'auth.login'
  | 'auth.register'
  | 'auth.createAccount'
  | 'auth.haveAccount'
  | 'landing.tagline'
  | 'landing.start'
  | 'landing.uiLanguage'
  | 'landing.targetLanguage'
  | 'catalog.title'
  | 'catalog.enroll'
  | 'course.curriculum'
  | 'course.startLesson'
  | 'learn.check'
  | 'learn.next'
  | 'progress.title'
  | 'path.title'
  | 'community.title'
  | 'achievements.title'

type Dict = Record<I18nKey, string>

export const STRINGS: Record<Locale, Dict> = {
  en: {
    'app.name': 'Linguaforge',
    'nav.home': 'Home',
    'nav.catalog': 'Courses',
    'nav.path': 'Path',
    'nav.progress': 'Progress',
    'nav.community': 'Community',
    'nav.achievements': 'Achievements',
    'nav.login': 'Login',
    'nav.logout': 'Logout',
    'auth.email': 'Email',
    'auth.password': 'Password',
    'auth.login': 'Login',
    'auth.register': 'Register',
    'auth.createAccount': 'Create an account',
    'auth.haveAccount': 'I already have an account',
    'landing.tagline': 'Immersive language learning with levels, practice, and community momentum.',
    'landing.start': 'Start learning',
    'landing.uiLanguage': 'UI language',
    'landing.targetLanguage': 'Target language',
    'catalog.title': 'Leveled Courses',
    'catalog.enroll': 'Enroll',
    'course.curriculum': 'Curriculum',
    'course.startLesson': 'Start lesson',
    'learn.check': 'Check',
    'learn.next': 'Next',
    'progress.title': 'Progress',
    'path.title': 'Your Path',
    'community.title': 'Community',
    'achievements.title': 'Achievements',
  },
  ja: {
    'app.name': 'Linguaforge',
    'nav.home': 'ホーム',
    'nav.catalog': 'コース',
    'nav.path': '学習ルート',
    'nav.progress': '進捗',
    'nav.community': 'コミュニティ',
    'nav.achievements': '実績',
    'nav.login': 'ログイン',
    'nav.logout': 'ログアウト',
    'auth.email': 'メール',
    'auth.password': 'パスワード',
    'auth.login': 'ログイン',
    'auth.register': '登録',
    'auth.createAccount': '新規アカウント作成',
    'auth.haveAccount': 'アカウントを持っています',
    'landing.tagline': 'レベル別コースと反復練習、コミュニティで続く学習体験。',
    'landing.start': '学習を始める',
    'landing.uiLanguage': '表示言語',
    'landing.targetLanguage': '学習言語',
    'catalog.title': 'レベル別コース',
    'catalog.enroll': '受講する',
    'course.curriculum': 'カリキュラム',
    'course.startLesson': 'レッスン開始',
    'learn.check': '確認',
    'learn.next': '次へ',
    'progress.title': '進捗',
    'path.title': 'おすすめルート',
    'community.title': 'コミュニティ',
    'achievements.title': '実績',
  },
  ko: {
    'app.name': 'Linguaforge',
    'nav.home': '홈',
    'nav.catalog': '코스',
    'nav.path': '학습 경로',
    'nav.progress': '진도',
    'nav.community': '커뮤니티',
    'nav.achievements': '업적',
    'nav.login': '로그인',
    'nav.logout': '로그아웃',
    'auth.email': '이메일',
    'auth.password': '비밀번호',
    'auth.login': '로그인',
    'auth.register': '가입',
    'auth.createAccount': '계정 만들기',
    'auth.haveAccount': '이미 계정이 있어요',
    'landing.tagline': '레벨 코스, 인터랙티브 훈련, 커뮤니티 보상으로 몰입 학습.',
    'landing.start': '학습 시작',
    'landing.uiLanguage': 'UI 언어',
    'landing.targetLanguage': '학습 언어',
    'catalog.title': '레벨 코스',
    'catalog.enroll': '등록',
    'course.curriculum': '커리큘럼',
    'course.startLesson': '레슨 시작',
    'learn.check': '확인',
    'learn.next': '다음',
    'progress.title': '진도',
    'path.title': '추천 경로',
    'community.title': '커뮤니티',
    'achievements.title': '업적',
  },
}


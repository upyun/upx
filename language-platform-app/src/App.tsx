import { useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import AppShell from '@/components/AppShell'
import Home from '@/pages/Home'
import Catalog from '@/pages/Catalog'
import Course from '@/pages/Course'
import Learn from '@/pages/Learn'
import Login from '@/pages/Login'
import Onboarding from '@/pages/Onboarding'
import Achievements from '@/pages/Achievements'
import Community from '@/pages/Community'
import Path from '@/pages/Path'
import Progress from '@/pages/Progress'
import Register from '@/pages/Register'
import { useAuthStore } from '@/stores/useAuthStore'
import { useLearningStore } from '@/stores/useLearningStore'

export default function App() {
  const loadMe = useAuthStore((s) => s.loadMe)
  const user = useAuthStore((s) => s.user)
  const loadForUser = useLearningStore((s) => s.loadForUser)

  useEffect(() => {
    loadMe()
  }, [loadMe])

  useEffect(() => {
    loadForUser(user?.id ?? null)
  }, [loadForUser, user?.id])

  return (
    <Router>
      <Routes>
        <Route
          path="/"
          element={
            <AppShell>
              <Home />
            </AppShell>
          }
        />
        <Route
          path="/catalog"
          element={
            <AppShell>
              <Catalog />
            </AppShell>
          }
        />
        <Route
          path="/courses/:courseId"
          element={
            <AppShell>
              <Course />
            </AppShell>
          }
        />
        <Route
          path="/learn/:lessonId"
          element={
            <AppShell>
              <Learn />
            </AppShell>
          }
        />
        <Route
          path="/path"
          element={
            <AppShell>
              <Path />
            </AppShell>
          }
        />
        <Route
          path="/community"
          element={
            <AppShell>
              <Community />
            </AppShell>
          }
        />
        <Route
          path="/achievements"
          element={
            <AppShell>
              <Achievements />
            </AppShell>
          }
        />
        <Route
          path="/progress"
          element={
            <AppShell>
              <Progress />
            </AppShell>
          }
        />
        <Route
          path="/login"
          element={
            <AppShell>
              <Login />
            </AppShell>
          }
        />
        <Route
          path="/register"
          element={
            <AppShell>
              <Register />
            </AppShell>
          }
        />
        <Route
          path="/onboarding"
          element={
            <AppShell>
              <Onboarding />
            </AppShell>
          }
        />
      </Routes>
    </Router>
  )
}

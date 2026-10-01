import { Navigate, Route, Routes, useLocation } from 'react-router-dom'
import TopBar from './components/TopBar'
import AccountPage from './pages/AccountPage'
import AdminHomePage from './pages/admin/AdminHomePage'
import AdminLoginPage from './pages/admin/AdminLoginPage'
import CreateNewProblemPage from './pages/admin/CreateNewProblemPage'
import HomePage from './pages/HomePage'
import LoginPage from './pages/LoginPage'
import ProblemDetailPage from './pages/ProblemDetailPage'
import ProblemsPage from './pages/ProblemsPage'
import SignupPage from './pages/SignupPage'
import { useAuthStore } from './store/authStore'

const publicRoutes = ['/login', '/signup']
const publicAdminRoutes = ['/admin/login']

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const location = useLocation()
  const isLoggedIn = useAuthStore((state) => state.isLoggedIn)

  if (!isLoggedIn && !publicRoutes.includes(location.pathname)) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return children
}

function AdminProtectedRoute({ children }: { children: React.ReactNode }) {
  const location = useLocation()
  const isAdminLoggedIn = useAuthStore((state) => state.isAdminLoggedIn)

  if (!isAdminLoggedIn && !publicAdminRoutes.includes(location.pathname)) {
    return <Navigate to="/admin/login" replace state={{ from: location }} />
  }

  return children
}

function AppShell() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-100">
      <TopBar />

      <main className="mx-auto max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/admin/login" element={<AdminLoginPage />} />

          <Route
            path="/"
            element={
              <ProtectedRoute>
                <HomePage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/problems"
            element={
              <ProtectedRoute>
                <ProblemsPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/problem/:problemId"
            element={
              <ProtectedRoute>
                <ProblemDetailPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/account"
            element={
              <ProtectedRoute>
                <AccountPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/admin/home"
            element={
              <AdminProtectedRoute>
                <AdminHomePage />
              </AdminProtectedRoute>
            }
          />
          <Route
            path="/admin/createnew"
            element={
              <AdminProtectedRoute>
                <CreateNewProblemPage />
              </AdminProtectedRoute>
            }
          />

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  )
}

export default function App() {
  return <AppShell />
}

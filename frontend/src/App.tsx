import { Routes, Route } from 'react-router-dom'
import { Navbar } from './components/Navbar'
import { HomePage } from './pages/HomePage'
import { LoginPage } from './pages/LoginPage'
import { DashboardPage } from './pages/DashboardPage'
import { NewEntryPage } from './pages/NewEntryPage'
import { AnalyticsPage } from './pages/AnalyticsPage'
import { MealsPage } from './pages/MealsPage'
import ProfilePage from './pages/ProfilePage'
import { ProtectedRoute } from './components/ProtectedRoute'

function App() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-orange-50 to-yellow-50">
      <Navbar />

      <main className="container mx-auto px-4 py-8">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />

          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <DashboardPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/new-entry"
            element={
              <ProtectedRoute>
                <NewEntryPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/analytics"
            element={
              <ProtectedRoute>
                <AnalyticsPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/meals"
            element={
              <ProtectedRoute>
                <MealsPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            }
          />
        </Routes>
      </main>
    </div>
  )
}

export default App

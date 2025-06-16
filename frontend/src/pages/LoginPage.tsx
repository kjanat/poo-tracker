import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Logo from '../components/Logo'
import { getLogoProps } from '../utils/branding'
import { useAuthStore } from '../stores/authStore'

export function LoginPage () {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError('')

    try {
      await login(email, password)
      navigate('/dashboard')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className='max-w-md mx-auto'>
      <div className='card'>
        <div className='text-center mb-6'>
          <Logo {...getLogoProps('login')} className='mx-auto mb-4' />
          <h1 className='text-2xl font-bold'>Welcome Back</h1>
        </div>
        <p className='text-center text-gray-600 mb-8'>
          Ready to log some legendary logs?
        </p>

        {error && (
          <div className='bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4'>
            {error}
          </div>
        )}

        <form className='space-y-4' onSubmit={handleSubmit}>
          <div>
            <label className='block text-sm font-medium text-gray-700 mb-1'>
              Email
            </label>
            <input
              type='email'
              className='input-field'
              placeholder='your@email.com'
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete='email'
              required
            />
          </div>

          <div>
            <label className='block text-sm font-medium text-gray-700 mb-1'>
              Password
            </label>
            <input
              type='password'
              className='input-field'
              placeholder='Password'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete='current-password'
              required
            />
          </div>

          <button
            type='submit'
            className='btn-primary w-full'
            disabled={loading}
          >
            {loading ? 'Signing In...' : 'Sign In'}
          </button>
        </form>

        <div className='mt-6 p-4 bg-blue-50 border border-blue-200 rounded'>
          <p className='text-sm text-blue-800 font-medium mb-2'>
            🔐 Test Credentials:
          </p>
          <div className='text-sm text-blue-700 space-y-1'>
            <div>
              <strong>Email:</strong> test@example.com
            </div>
            <div>
              <strong>Password:</strong> password123
            </div>
          </div>
        </div>

        <p className='text-center text-sm text-gray-600 mt-6'>
          Don&apos;t have an account? We&apos;ll create one for you automatically on first
          login.
        </p>
      </div>
    </div>
  )
}

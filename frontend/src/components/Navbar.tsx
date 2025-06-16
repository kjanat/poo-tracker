import { Link } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import {
  LogOut,
  Home,
  BarChart3,
  Plus,
  UtensilsCrossed,
  User
} from 'lucide-react'
import Logo from './Logo'
import { getLogoProps } from '../utils/branding'

export function Navbar () {
  const { isAuthenticated, user, logout } = useAuthStore()

  return (
    <nav className='bg-white shadow-lg border-b-4 border-poo-brown-500'>
      <div className='container mx-auto px-4'>
        <div className='flex justify-between items-center h-16'>
          <Link to='/' className='flex items-center space-x-2'>
            <Logo {...getLogoProps('navbar')} />
            <span className='text-xl font-bold text-poo-brown-700'>
              Poo Tracker
            </span>
          </Link>

          <div className='flex items-center space-x-4'>
            {isAuthenticated
              ? (
                <>
                  <Link
                    to='/dashboard'
                    className='flex items-center space-x-1 px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 transition-colors'
                  >
                    <Home size={18} />
                    <span>Dashboard</span>
                  </Link>

                  <Link
                    to='/new-entry'
                    className='flex items-center space-x-1 px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 transition-colors'
                  >
                    <Plus size={18} />
                    <span>New Entry</span>
                  </Link>

                  <Link
                    to='/meals'
                    className='flex items-center space-x-1 px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 transition-colors'
                  >
                    <UtensilsCrossed size={18} />
                    <span>Meals</span>
                  </Link>

                  <Link
                    to='/analytics'
                    className='flex items-center space-x-1 px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 transition-colors'
                  >
                    <BarChart3 size={18} />
                    <span>Analytics</span>
                  </Link>

                  <Link
                    to='/profile'
                    className='flex items-center space-x-1 px-3 py-2 rounded-md text-gray-700 hover:bg-gray-100 transition-colors'
                  >
                    <User size={18} />
                    <span>Profile</span>
                  </Link>

                  <div className='flex items-center space-x-2'>
                    <span className='text-sm text-gray-600'>
                      Welcome, {user?.name || user?.email}!
                    </span>
                    <button
                      onClick={logout}
                      className='flex items-center space-x-1 px-3 py-2 rounded-md text-red-600 hover:bg-red-50 transition-colors'
                    >
                      <LogOut size={18} />
                      <span>Logout</span>
                    </button>
                  </div>
                </>
                )
              : (
                <Link to='/login' className='btn-primary'>
                  Login
                </Link>
                )}
          </div>
        </div>
      </div>
    </nav>
  )
}

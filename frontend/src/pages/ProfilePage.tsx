import React, { useState, useEffect, useCallback } from 'react'
import { useAuthStore } from '../stores/authStore'
import Logo from '../components/Logo'

interface UserProfile {
  id: string
  email: string
  name?: string
  createdAt: string
  updatedAt: string
}

interface UserAuth {
  lastLogin?: string
  createdAt: string
  updatedAt: string
}

const ProfilePage: React.FC = () => {
  const { token, logout } = useAuthStore()
  const [profile, setProfile] = useState<UserProfile | null>(null)
  const [userAuth, setUserAuth] = useState<UserAuth | null>(null)
  const [isEditing, setIsEditing] = useState(false)
  const [editForm, setEditForm] = useState({
    name: '',
    email: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const fetchUserProfile = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)

      const response = await fetch('/api/auth/profile', {
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        throw new Error('Failed to fetch profile')
      }

      const data = await response.json()
      setProfile(data.user)
      setUserAuth(data.auth)
      setEditForm({
        name: data.user.name || '',
        email: data.user.email,
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load profile')
    } finally {
      setLoading(false)
    }
  }, [token])

  useEffect(() => {
    fetchUserProfile()
  }, [fetchUserProfile])

  const handleEditToggle = () => {
    setIsEditing(!isEditing)
    setError(null)
    setSuccess(null)
    if (!isEditing && profile != null) {
      setEditForm({
        name: profile.name || '',
        email: profile.email,
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      })
    }
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setEditForm((prev) => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSuccess(null)

    // Validate password change if attempting to change password
    if (editForm.newPassword || editForm.confirmPassword) {
      if (!editForm.currentPassword) {
        setError('Current password is required to change password')
        return
      }
      if (editForm.newPassword !== editForm.confirmPassword) {
        setError('New passwords do not match')
        return
      }
      if (editForm.newPassword.length < 6) {
        setError('New password must be at least 6 characters long')
        return
      }
    }

    try {
      setLoading(true)

      const updateData: Record<string, string> = {
        name: editForm.name,
        email: editForm.email
      }

      if (editForm.newPassword) {
        updateData.currentPassword = editForm.currentPassword
        updateData.newPassword = editForm.newPassword
      }

      const response = await fetch('/api/auth/profile', {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updateData)
      })

      if (!response.ok) {
        const errorData = await response.json()
        throw new Error(errorData.error || 'Failed to update profile')
      }

      const data = await response.json()
      setProfile(data.user)
      setUserAuth(data.auth)
      setIsEditing(false)
      setSuccess('Profile updated successfully!')

      // Clear password fields
      setEditForm((prev) => ({
        ...prev,
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      }))
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update profile')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  if (loading && profile == null) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-amber-50 to-orange-100 py-8">
        <div className="max-w-4xl mx-auto px-4">
          <div className="bg-white rounded-2xl shadow-xl p-8">
            <div className="animate-pulse">
              <div className="flex items-center space-x-4 mb-8">
                <div className="w-12 h-12 bg-gray-300 rounded-lg" />
                <div className="h-8 w-48 bg-gray-300 rounded" />
              </div>
              <div className="space-y-4">
                <div className="h-4 w-full bg-gray-300 rounded" />
                <div className="h-4 w-3/4 bg-gray-300 rounded" />
                <div className="h-4 w-1/2 bg-gray-300 rounded" />
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-amber-50 to-orange-100 py-8">
      <div className="max-w-4xl mx-auto px-4">
        <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
          {/* Header */}
          <div className="bg-gradient-to-r from-amber-500 to-orange-500 px-8 py-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4">
                <Logo className="w-12 h-12 text-white" />
                <div>
                  <h1 className="text-3xl font-bold text-white">User Profile</h1>
                  <p className="text-amber-100">Manage your account settings</p>
                </div>
              </div>
              <button
                onClick={() => logout()}
                className="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg transition-colors"
              >
                Logout
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="p-8">
            {error && (
              <div className="mb-6 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
                {error}
              </div>
            )}

            {success && (
              <div className="mb-6 bg-green-50 border border-green-200 text-green-600 px-4 py-3 rounded-lg">
                {success}
              </div>
            )}

            <div className="grid md:grid-cols-2 gap-8">
              {/* Profile Information */}
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <h2 className="text-2xl font-semibold text-gray-800">Profile Information</h2>
                  <button
                    onClick={handleEditToggle}
                    className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                      isEditing
                        ? 'bg-gray-500 hover:bg-gray-600 text-white'
                        : 'bg-amber-500 hover:bg-amber-600 text-white'
                    }`}
                  >
                    {isEditing ? 'Cancel' : 'Edit Profile'}
                  </button>
                </div>

                {!isEditing ? (
                  <div className="space-y-4">
                    <div className="bg-gray-50 p-4 rounded-lg">
                      <label className="block text-sm font-medium text-gray-600 mb-1">Name</label>
                      <p className="text-gray-900">{profile?.name || 'Not set'}</p>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                      <label className="block text-sm font-medium text-gray-600 mb-1">Email</label>
                      <p className="text-gray-900">{profile?.email}</p>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg">
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Member Since
                      </label>
                      <p className="text-gray-900">
                        {profile != null && formatDate(profile.createdAt)}
                      </p>
                    </div>
                  </div>
                ) : (
                  <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Name</label>
                      <input
                        type="text"
                        name="name"
                        value={editForm.name}
                        onChange={handleInputChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-transparent"
                        placeholder="Your name"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Email</label>
                      <input
                        type="email"
                        name="email"
                        value={editForm.email}
                        onChange={handleInputChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-transparent"
                      />
                    </div>

                    <div className="border-t pt-4 mt-6">
                      <h3 className="text-lg font-medium text-gray-800 mb-4">
                        Change Password (Optional)
                      </h3>
                      <div className="space-y-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                            Current Password
                          </label>
                          <input
                            type="password"
                            name="currentPassword"
                            value={editForm.currentPassword}
                            onChange={handleInputChange}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-transparent"
                            placeholder="Enter current password"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                            New Password
                          </label>
                          <input
                            type="password"
                            name="newPassword"
                            value={editForm.newPassword}
                            onChange={handleInputChange}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-transparent"
                            placeholder="Enter new password"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">
                            Confirm New Password
                          </label>
                          <input
                            type="password"
                            name="confirmPassword"
                            value={editForm.confirmPassword}
                            onChange={handleInputChange}
                            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-transparent"
                            placeholder="Confirm new password"
                          />
                        </div>
                      </div>
                    </div>

                    <div className="flex space-x-4 pt-4">
                      <button
                        type="submit"
                        disabled={loading}
                        className="flex-1 bg-amber-500 hover:bg-amber-600 text-white py-2 px-4 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {loading ? 'Saving...' : 'Save Changes'}
                      </button>
                    </div>
                  </form>
                )}
              </div>

              {/* Account Statistics */}
              <div className="space-y-6">
                <h2 className="text-2xl font-semibold text-gray-800">Account Information</h2>

                <div className="space-y-4">
                  <div className="bg-gradient-to-r from-blue-50 to-blue-100 p-4 rounded-lg">
                    <h3 className="font-medium text-blue-800 mb-2">Account Status</h3>
                    <p className="text-blue-600">Active</p>
                  </div>

                  {userAuth?.lastLogin && (
                    <div className="bg-gradient-to-r from-green-50 to-green-100 p-4 rounded-lg">
                      <h3 className="font-medium text-green-800 mb-2">Last Login</h3>
                      <p className="text-green-600">{formatDate(userAuth.lastLogin)}</p>
                    </div>
                  )}

                  <div className="bg-gradient-to-r from-purple-50 to-purple-100 p-4 rounded-lg">
                    <h3 className="font-medium text-purple-800 mb-2">Account Created</h3>
                    <p className="text-purple-600">
                      {profile != null && formatDate(profile.createdAt)}
                    </p>
                  </div>

                  {profile?.updatedAt !== profile?.createdAt && (
                    <div className="bg-gradient-to-r from-orange-50 to-orange-100 p-4 rounded-lg">
                      <h3 className="font-medium text-orange-800 mb-2">Last Updated</h3>
                      <p className="text-orange-600">
                        {profile != null && formatDate(profile.updatedAt)}
                      </p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ProfilePage

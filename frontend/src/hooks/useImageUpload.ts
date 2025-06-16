import { useState } from 'react'
import { useAuthStore } from '../stores/authStore'
import { API_BASE_URL, createFormDataHeaders, handleApiResponse } from '../utils/api'

export interface UseImageUploadReturn {
  selectedImage: File | null
  imagePreview: string | null
  handleImageChange: (e: React.ChangeEvent<HTMLInputElement>) => void
  removeImage: () => void
  uploadImage: () => Promise<string | null>
}

export function useImageUpload (): UseImageUploadReturn {
  const [selectedImage, setSelectedImage] = useState<File | null>(null)
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  
  const token = useAuthStore((state) => state.token)

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
    const file = e.target.files?.[0]
    if (file == null) return

    // Validate file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert('File size must be less than 5MB')
      return
    }

    // Validate file type
    if (!file.type.startsWith('image/')) {
      alert('Please select an image file')
      return
    }

    setSelectedImage(file)

    // Create preview
    const reader = new FileReader()
    reader.onload = (e) => {
      setImagePreview(e.target?.result as string)
    }
    reader.readAsDataURL(file)
  }

  const removeImage = (): void => {
    setSelectedImage(null)
    setImagePreview(null)
  }

  const uploadImage = async (): Promise<string | null> => {
    if (selectedImage == null || token == null) return null

    const uploadFormData = new FormData()
    uploadFormData.append('image', selectedImage)

    try {
      const response = await fetch(`${API_BASE_URL}/api/uploads`, {
        method: 'POST',
        headers: createFormDataHeaders(token),
        body: uploadFormData
      })

      const data = await handleApiResponse<{ url: string }>(response)
      return data.url
    } catch (error) {
      console.error('Failed to upload image:', error)
      throw error
    }
  }

  return {
    selectedImage,
    imagePreview,
    handleImageChange,
    removeImage,
    uploadImage
  }
}

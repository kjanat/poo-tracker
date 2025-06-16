import React, { useState, type ReactElement } from 'react'
import Logo from '../components/Logo'
import { MealForm } from '../components/meals/MealForm'
import { MealCard } from '../components/meals/MealCard'
import { EntryLinkingModal } from '../components/meals/EntryLinkingModal'
import { useMeals } from '../hooks/useMeals'
import { useImageUpload } from '../hooks/useImageUpload'
import { useEntryLinking } from '../hooks/useEntryLinking'
import type { Meal, MealFormData } from '../types'

export function MealsPage (): ReactElement {
  const [formData, setFormData] = useState<MealFormData>({
    name: '',
    category: '',
    description: '',
    cuisine: '',
    spicyLevel: 1,
    fiberRich: false,
    dairy: false,
    gluten: false,
    notes: ''
  })
  const [editingMeal, setEditingMeal] = useState<Meal | null>(null)
  const [linkingMeal, setLinkingMeal] = useState<Meal | null>(null)
  const [showLinkingModal, setShowLinkingModal] = useState<boolean>(false)

  const {
    meals,
    loading,
    error,
    success,
    setError,
    setSuccess,
    createMeal,
    updateMeal,
    deleteMeal
  } = useMeals()

  const {
    selectedImage,
    imagePreview,
    handleImageChange,
    removeImage,
    uploadImage
  } = useImageUpload()

  const {
    availableEntries,
    linkedEntries,
    fetchLinkedEntries,
    linkEntry,
    unlinkEntry,
    setError: setLinkingError
  } = useEntryLinking()

  const resetForm = (): void => {
    setFormData({
      name: '',
      category: '',
      description: '',
      cuisine: '',
      spicyLevel: 1,
      fiberRich: false,
      dairy: false,
      gluten: false,
      notes: ''
    })
    removeImage()
    setEditingMeal(null)
  }

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault()
    setError('')
    setSuccess('')

    try {
      // Upload image if selected
      const photoUrl = await uploadImage()
      const finalFormData = { ...formData, photoUrl: photoUrl ?? undefined }

      if (editingMeal != null) {
        await updateMeal(editingMeal.id, finalFormData)
      } else {
        await createMeal(finalFormData)
      }

      resetForm()
    } catch (error) {
      // Error is already handled in the hooks
      console.error('Failed to submit meal:', error)
    }
  }

  const startEdit = (meal: Meal): void => {
    setFormData({
      name: meal.name,
      category: meal.category ?? '',
      description: meal.description ?? '',
      cuisine: meal.cuisine ?? '',
      spicyLevel: meal.spicyLevel ?? 1,
      fiberRich: meal.fiberRich,
      dairy: meal.dairy,
      gluten: meal.gluten,
      notes: meal.notes ?? '',
      photoUrl: meal.photoUrl
    })
    setEditingMeal(meal)
    
    // Set image preview if meal has photo
    if (meal.photoUrl != null) {
      // Note: This won't work perfectly for existing images, but it's better than nothing
      // In a real app, you might want to handle this differently
    }
  }

  const cancelEdit = (): void => {
    resetForm()
  }

  const startLinking = async (meal: Meal): Promise<void> => {
    setLinkingMeal(meal)
    await fetchLinkedEntries(meal.id)
    setShowLinkingModal(true)
  }

  const closeLinkingModal = (): void => {
    setShowLinkingModal(false)
    setLinkingMeal(null)
  }

  const handleLinkEntry = async (entryId: string): Promise<void> => {
    if (linkingMeal == null) return
    try {
      await linkEntry(linkingMeal.id, entryId)
    } catch (error) {
      setLinkingError(error instanceof Error ? error.message : 'Failed to link entry')
    }
  }

  const handleUnlinkEntry = async (entryId: string): Promise<void> => {
    if (linkingMeal == null) return
    try {
      await unlinkEntry(linkingMeal.id, entryId)
    } catch (error) {
      setLinkingError(error instanceof Error ? error.message : 'Failed to unlink entry')
    }
  }

  return (
    <div className="min-h-screen bg-amber-50">
      <nav className="bg-white shadow-sm border-b border-amber-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Logo />
            <h1 className="text-2xl font-bold text-amber-800">Meals</h1>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {error !== '' && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {success !== '' && (
          <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
            {success}
          </div>
        )}

        <MealForm
          formData={formData}
          setFormData={setFormData}
          onSubmit={handleSubmit}
          loading={loading}
          editingMeal={editingMeal}
          onCancel={cancelEdit}
          selectedImage={selectedImage}
          imagePreview={imagePreview}
          onImageChange={handleImageChange}
          onRemoveImage={removeImage}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {meals.map((meal) => (
            <MealCard
              key={meal.id}
              meal={meal}
              onEdit={startEdit}
              onDelete={deleteMeal}
              onLink={startLinking}
            />
          ))}
        </div>

        {meals.length === 0 && !loading && (
          <div className="text-center py-8 text-gray-500">
            No meals added yet. Add your first meal above!
          </div>
        )}

        <EntryLinkingModal
          show={showLinkingModal}
          onClose={closeLinkingModal}
          meal={linkingMeal}
          availableEntries={availableEntries}
          linkedEntries={linkedEntries}
          onLinkEntry={handleLinkEntry}
          onUnlinkEntry={handleUnlinkEntry}
        />
      </main>
    </div>
  )
}

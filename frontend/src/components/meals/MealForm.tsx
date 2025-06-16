import React, { type ReactElement } from 'react'
import type { MealFormData, Meal } from '../types'

interface MealFormProps {
  formData: MealFormData
  setFormData: React.Dispatch<React.SetStateAction<MealFormData>>
  onSubmit: (e: React.FormEvent) => void
  loading: boolean
  editingMeal: Meal | null
  onCancel?: () => void
  selectedImage: File | null
  imagePreview: string | null
  onImageChange: (e: React.ChangeEvent<HTMLInputElement>) => void
  onRemoveImage: () => void
}

export function MealForm ({
  formData,
  setFormData,
  onSubmit,
  loading,
  editingMeal,
  onCancel,
  selectedImage,
  imagePreview,
  onImageChange,
  onRemoveImage
}: MealFormProps): ReactElement {
  return (
    <div className="bg-white p-6 rounded-lg shadow-md mb-8">
      <h2 className="text-xl font-semibold mb-4 text-gray-800">
        {editingMeal != null ? 'Edit Meal' : 'Add New Meal'}
      </h2>

      <form onSubmit={onSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Meal Name *
            </label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Category
            </label>
            <select
              value={formData.category}
              onChange={(e) => setFormData({ ...formData, category: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            >
              <option value="">Select category</option>
              <option value="breakfast">Breakfast</option>
              <option value="lunch">Lunch</option>
              <option value="dinner">Dinner</option>
              <option value="snack">Snack</option>
              <option value="drink">Drink</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Cuisine
            </label>
            <input
              type="text"
              value={formData.cuisine}
              onChange={(e) => setFormData({ ...formData, cuisine: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
              placeholder="e.g., Italian, Mexican, Asian"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Spicy Level
            </label>
            <select
              value={formData.spicyLevel}
              onChange={(e) => setFormData({ ...formData, spicyLevel: parseInt(e.target.value) })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            >
              <option value={1}>1 - Mild</option>
              <option value={2}>2 - Medium</option>
              <option value={3}>3 - Hot</option>
              <option value={4}>4 - Very Hot</option>
              <option value={5}>5 - Extremely Hot</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Description
          </label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            rows={3}
            placeholder="Describe the meal..."
          />
        </div>

        <div className="grid grid-cols-3 gap-4">
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={formData.fiberRich}
              onChange={(e) => setFormData({ ...formData, fiberRich: e.target.checked })}
              className="mr-2"
            />
            <span className="text-sm text-gray-700">High Fiber</span>
          </label>

          <label className="flex items-center">
            <input
              type="checkbox"
              checked={formData.dairy}
              onChange={(e) => setFormData({ ...formData, dairy: e.target.checked })}
              className="mr-2"
            />
            <span className="text-sm text-gray-700">Contains Dairy</span>
          </label>

          <label className="flex items-center">
            <input
              type="checkbox"
              checked={formData.gluten}
              onChange={(e) => setFormData({ ...formData, gluten: e.target.checked })}
              className="mr-2"
            />
            <span className="text-sm text-gray-700">Contains Gluten</span>
          </label>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Photo
          </label>
          <input
            type="file"
            accept="image/*"
            onChange={onImageChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
          />
          {imagePreview != null && (
            <div className="mt-2 relative inline-block">
              <img
                src={imagePreview}
                alt="Preview"
                className="h-32 w-32 object-cover rounded-md"
              />
              <button
                type="button"
                onClick={onRemoveImage}
                className="absolute -top-2 -right-2 bg-red-500 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs hover:bg-red-600"
              >
                ×
              </button>
            </div>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Notes
          </label>
          <textarea
            value={formData.notes}
            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amber-500"
            rows={2}
            placeholder="Additional notes..."
          />
        </div>

        <div className="flex gap-2">
          <button
            type="submit"
            disabled={loading}
            className="bg-amber-600 text-white px-4 py-2 rounded-md hover:bg-amber-700 disabled:opacity-50"
          >
            {loading ? 'Saving...' : editingMeal != null ? 'Update Meal' : 'Add Meal'}
          </button>
          
          {editingMeal != null && onCancel != null && (
            <button
              type="button"
              onClick={onCancel}
              className="bg-gray-500 text-white px-4 py-2 rounded-md hover:bg-gray-600"
            >
              Cancel
            </button>
          )}
        </div>
      </form>
    </div>
  )
}

import React, { type ReactElement } from 'react'
import { formatDate } from '../../utils/date'
import type { Meal } from '../../types'

interface MealCardProps {
  meal: Meal
  onEdit: (meal: Meal) => void
  onDelete: (mealId: string) => void
  onLink: (meal: Meal) => void
}

export function MealCard ({ meal, onEdit, onDelete, onLink }: MealCardProps): ReactElement {
  return (
    <div className="bg-white p-4 rounded-lg shadow-md">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-semibold text-gray-800">{meal.name}</h3>
        <div className="flex gap-1">
          <button
            onClick={() => onEdit(meal)}
            className="text-blue-600 hover:text-blue-800 text-sm"
          >
            Edit
          </button>
          <button
            onClick={() => onLink(meal)}
            className="text-green-600 hover:text-green-800 text-sm ml-2"
          >
            Link Entries
          </button>
          <button
            onClick={() => onDelete(meal.id)}
            className="text-red-600 hover:text-red-800 text-sm ml-2"
          >
            Delete
          </button>
        </div>
      </div>

      {meal.photoUrl != null && (
        <img
          src={meal.photoUrl}
          alt={meal.name}
          className="w-full h-32 object-cover rounded-md mb-2"
        />
      )}

      <div className="text-sm text-gray-600 space-y-1">
        {meal.category != null && (
          <div>
            <span className="font-medium">Category:</span> {meal.category}
          </div>
        )}
        
        {meal.cuisine != null && (
          <div>
            <span className="font-medium">Cuisine:</span> {meal.cuisine}
          </div>
        )}

        {meal.spicyLevel != null && (
          <div>
            <span className="font-medium">Spicy Level:</span> {meal.spicyLevel}/5
          </div>
        )}

        {meal.description != null && (
          <div>
            <span className="font-medium">Description:</span> {meal.description}
          </div>
        )}

        <div className="flex gap-2 flex-wrap">
          {meal.fiberRich && (
            <span className="bg-green-100 text-green-800 px-2 py-1 rounded-full text-xs">
              High Fiber
            </span>
          )}
          {meal.dairy && (
            <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded-full text-xs">
              Dairy
            </span>
          )}
          {meal.gluten && (
            <span className="bg-orange-100 text-orange-800 px-2 py-1 rounded-full text-xs">
              Gluten
            </span>
          )}
        </div>

        {meal.notes != null && (
          <div>
            <span className="font-medium">Notes:</span> {meal.notes}
          </div>
        )}

        {(meal.linkedEntries?.length ?? 0) > 0 && (
          <div>
            <span className="font-medium">Linked Entries:</span> {meal.linkedEntries?.length}
          </div>
        )}

        <div className="text-xs text-gray-500 mt-2">
          Added: {formatDate(meal.createdAt)}
        </div>
      </div>
    </div>
  )
}

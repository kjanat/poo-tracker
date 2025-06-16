import React, { type ReactElement } from 'react'
import { BristolAnalyzer } from '../../bristol/BristolAnalyzer'

interface BristolSelectorProps {
  selectedType: number
  onTypeSelect: (type: number) => void
  disabled?: boolean
}

export function BristolSelector ({ selectedType, onTypeSelect, disabled = false }: BristolSelectorProps): ReactElement {
  const bristolTypes = BristolAnalyzer.getAllTypes()

  return (
    <div className="space-y-4">
      <h3 className="text-lg font-semibold text-gray-800">Bristol Stool Chart</h3>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {bristolTypes.map((typeInfo) => (
          <button
            key={typeInfo.type}
            type="button"
            onClick={() => onTypeSelect(typeInfo.type)}
            disabled={disabled}
            className={`p-4 rounded-lg border-2 text-left transition-all ${
              selectedType === typeInfo.type
                ? 'border-amber-500 bg-amber-50'
                : 'border-gray-200 hover:border-gray-300'
            } ${disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
          >
            <div className="flex items-center gap-3">
              <div 
                className="w-4 h-4 rounded-full"
                style={{ backgroundColor: typeInfo.color }}
              />
              <div>
                <div className="font-medium">
                  Type {typeInfo.type}: {typeInfo.description}
                </div>
                <div className="text-sm text-gray-600">
                  {typeInfo.severity}
                </div>
              </div>
            </div>
          </button>
        ))}
      </div>
    </div>
  )
}

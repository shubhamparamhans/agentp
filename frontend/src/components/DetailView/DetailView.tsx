import React, { useContext } from 'react'
import { ObjectRenderer } from '../ObjectRenderer/ObjectRenderer'
import { AppContext } from '../../state/AppContext'

interface DetailViewProps {
  modelName: string
  selectedRow: Record<string, string | number> | null
  isOpen: boolean
  onClose: () => void
}

const DetailView: React.FC<DetailViewProps> = ({ modelName, selectedRow, isOpen, onClose }) => {
  const { state } = useContext(AppContext)

  if (!selectedRow) return null

  const fields = Object.entries(selectedRow)

  return (
    <>
      {/* Overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-30 transition-opacity"
          onClick={onClose}
        />
      )}

      {/* Slide-in Panel */}
      <div
        className={`fixed top-0 right-0 h-full w-96 bg-gray-800 shadow-2xl border-l border-cyan-600 transform transition-transform duration-300 ease-out z-40 ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        } overflow-y-auto`}
      >
        {/* Header */}
        <div className="sticky top-0 bg-gray-800 border-b border-gray-700 px-6 py-4 flex items-center justify-between">
          <h2 className="text-2xl font-bold text-white capitalize flex items-center gap-2">
            <span className="text-cyan-400">📋</span>
            {modelName} Details
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white font-bold text-2xl transition"
          >
            ✕
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-4">
          {fields.map(([key, value]) => (
            <div
              key={key}
              className="bg-gray-700 rounded-lg p-4 border border-gray-600 hover:border-cyan-600 transition"
            >
              <label className="block text-sm font-bold text-cyan-400 mb-2 uppercase tracking-wide">
                {key.replace(/_/g, ' ')}
              </label>
              <div className="text-white text-base break-words">
                <ObjectRenderer 
                  value={value} 
                  isMongoDb={state.databaseType === 'mongo'}
                />
              </div>
            </div>
          ))}
        </div>

        {/* Footer */}
        <div className="sticky bottom-0 bg-gray-800 border-t border-gray-700 px-6 py-4">
          <button
            onClick={onClose}
            className="w-full px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-all font-bold"
          >
            Close
          </button>
        </div>
      </div>
    </>
  )
}

export default DetailView

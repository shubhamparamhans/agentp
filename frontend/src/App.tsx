import { useState, useEffect } from 'react'
import { ListView } from './components/ListView/ListView'
import { FilterBuilder } from './components/FilterBuilder/FilterBuilder'
import { GroupView } from './components/GroupView/GroupView'
import DetailView from './components/DetailView/DetailView'
import { AppProvider } from './state/AppContext'
import { fetchModels, type Model } from './api/client'

interface Filter {
  id: string
  field: string
  operator: string
  value: string
}

function AppContent() {
  const [selectedModel, setSelectedModel] = useState<string | null>(null)
  const [models, setModels] = useState<Model[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [filters, setFilters] = useState<Filter[]>([])
  const [groupByField, setGroupByField] = useState<string | null>(null)
  const [showGroupView, setShowGroupView] = useState(false)
  const [showFilterModal, setShowFilterModal] = useState(false)
  const [showGroupModal, setShowGroupModal] = useState(false)
  const [selectedRow, setSelectedRow] = useState<Record<string, any> | null>(null)
  const [showDetailView, setShowDetailView] = useState(false)

  // Fetch models on component mount
  useEffect(() => {
    const loadModels = async () => {
      try {
        setLoading(true)
        setError(null)
        const fetchedModels = await fetchModels()
        setModels(fetchedModels)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load models')
        console.error('Error loading models:', err)
      } finally {
        setLoading(false)
      }
    }

    loadModels()
  }, [])

  const handleSelectModel = (modelName: string) => {
    setSelectedModel(modelName)
    setFilters([])
    setGroupByField(null)
    setShowGroupView(false)
    setShowFilterModal(false)
    setShowGroupModal(false)
    setSelectedRow(null)
    setShowDetailView(false)
  }

  const handleAddFilter = (filter: Omit<Filter, 'id'>) => {
    setFilters([...filters, { ...filter, id: Date.now().toString() }])
  }

  const handleRemoveFilter = (filterId: string) => {
    setFilters(filters.filter((f) => f.id !== filterId))
  }

  const handleRowClick = (row: Record<string, any>) => {
    setSelectedRow(row)
    setShowDetailView(true)
  }

  const currentModel = models.find((m) => m.name === selectedModel)

  // Extract field names from model
  const currentModelFields = currentModel?.fields?.map((f) => f.name) || []

  return (
    <div className="min-h-screen bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-gray-800 border-b border-gray-700 shadow-lg">
        <div className="max-w-full mx-auto px-6 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white">Universal Data Viewer</h1>
            <p className="text-gray-400 mt-1">Explore and analyze your data with ease</p>
            {error && <p className="text-red-400 text-sm mt-1">⚠️ {error}</p>}
          </div>
          {selectedModel && (
            <div className="flex gap-3">
              <button
                onClick={() => setShowFilterModal(!showFilterModal)}
                className={`px-4 py-2 rounded-lg font-semibold flex items-center gap-2 transition-all ${
                  showFilterModal || filters.length > 0
                    ? 'bg-cyan-600 text-white shadow-lg'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                }`}
              >
                🔍 Filters {filters.length > 0 && `(${filters.length})`}
              </button>
              <button
                onClick={() => setShowGroupModal(!showGroupModal)}
                className={`px-4 py-2 rounded-lg font-semibold flex items-center gap-2 transition-all ${
                  showGroupModal || groupByField
                    ? 'bg-purple-600 text-white shadow-lg'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                }`}
              >
                ⊕ Group By {groupByField && `(${groupByField})`}
              </button>
            </div>
          )}
        </div>
      </header>

      {/* Loading state */}
      {loading && (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center">
            <div className="inline-block animate-spin text-4xl mb-4">⚙️</div>
            <p className="text-xl text-gray-300">Loading models...</p>
          </div>
        </div>
      )}

      {/* Error state */}
      {error && !loading && (
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center bg-red-900 bg-opacity-30 border border-red-700 rounded-lg p-8 max-w-md">
            <div className="text-5xl mb-4">❌</div>
            <p className="text-xl text-red-300 font-semibold">{error}</p>
            <p className="text-gray-400 mt-2">Make sure the backend server is running on http://localhost:8080</p>
          </div>
        </div>
      )}

      {/* Main Content */}
      {!loading && !error && (
        <main className="flex-1 overflow-hidden flex">
          <div className="w-full h-full grid grid-cols-5 gap-0">
            {/* Left Sidebar - Models */}
            <aside className="col-span-1 bg-gray-800 border-r border-gray-700 overflow-y-auto">
              <div className="p-6">
                <h2 className="text-lg font-bold text-white mb-4 flex items-center gap-2">
                  <span className="text-cyan-400">📊</span>
                  Models ({models.length})
                </h2>
                <div className="space-y-2">
                  {models.map((model) => (
                    <button
                      key={model.name}
                      onClick={() => handleSelectModel(model.name)}
                      className={`w-full px-4 py-3 text-left rounded-lg transition-all font-medium ${
                        selectedModel === model.name
                          ? 'bg-cyan-600 text-white shadow-lg'
                          : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                      }`}
                    >
                      {model.name.charAt(0).toUpperCase() + model.name.slice(1)}
                    </button>
                  ))}
                </div>
              </div>
            </aside>

            {/* Right Content Area */}
            <div className="col-span-4 bg-gray-900 overflow-hidden flex flex-col">
              {selectedModel ? (
                <>
                  {/* Top Section - Model Info */}
                  <div className="border-b border-gray-700 bg-gray-800 p-6">
                    <h2 className="text-3xl font-bold text-white capitalize">{selectedModel}</h2>
                    <p className="text-gray-400 mt-2 flex gap-4">
                      <span>
                        <span className="font-semibold text-white">Table:</span> {currentModel?.table}
                      </span>
                      <span>
                        <span className="font-semibold text-white">Primary Key:</span> {currentModel?.primary_key}
                      </span>
                      {filters.length > 0 && (
                        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-cyan-900 text-cyan-300">
                          {filters.length} filter(s)
                        </span>
                      )}
                      {groupByField && (
                        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-purple-900 text-purple-300">
                          Grouped by {groupByField}
                        </span>
                      )}
                    </p>
                  </div>

                  {/* Data Display */}
                  <div className="flex-1 overflow-y-auto p-6">
                    {showGroupView && groupByField ? (
                      <div className="bg-gray-800 rounded-lg shadow-lg overflow-hidden border border-gray-700">
                        <GroupView
                          modelName={selectedModel}
                          groupByField={groupByField}
                          filters={filters}
                        />
                      </div>
                    ) : (
                      <div className="bg-gray-800 rounded-lg shadow-lg overflow-hidden border border-gray-700">
                        <ListView
                          modelName={selectedModel}
                          filters={filters}
                          modelFields={currentModelFields}
                          onRowClick={handleRowClick}
                        />
                      </div>
                    )}
                  </div>
                </>
              ) : (
                <div className="flex-1 flex items-center justify-center">
                  <div className="text-center">
                    <div className="text-6xl mb-4">📊</div>
                    <p className="text-2xl text-gray-400 font-semibold">Select a model to get started</p>
                    <p className="text-gray-500 mt-2">Choose from the available models on the left</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </main>
      )}

      {/* Filter Modal Overlay */}
      {showFilterModal && selectedModel && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl shadow-2xl border border-cyan-600 max-w-md w-full">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between bg-gray-800">
              <h3 className="text-xl font-bold text-white flex items-center gap-2">
                <span className="text-cyan-400">🔍</span>
                Filters
              </h3>
              <button
                onClick={() => setShowFilterModal(false)}
                className="text-gray-400 hover:text-white font-bold text-xl transition"
              >
                ✕
              </button>
            </div>
            <div className="p-6 space-y-4 max-h-96 overflow-y-auto">
              <FilterBuilder
                fields={currentModelFields}
                onAddFilter={handleAddFilter}
              />

              {filters.length > 0 && (
                <>
                  <div className="border-t border-gray-700 pt-4">
                    <h4 className="text-sm font-bold text-gray-300 mb-3 uppercase tracking-wide">Applied Filters</h4>
                    <div className="space-y-2">
                      {filters.map((filter) => (
                        <div
                          key={filter.id}
                          className="p-3 bg-gray-700 border border-gray-600 rounded-lg flex justify-between items-start hover:border-cyan-500 transition"
                        >
                          <div className="text-sm flex-1">
                            <p className="font-medium text-white">
                              {filter.field} <span className="text-cyan-400">{filter.operator}</span>
                            </p>
                            <p className="text-gray-400 text-xs mt-1">{filter.value}</p>
                          </div>
                          <button
                            onClick={() => handleRemoveFilter(filter.id)}
                            className="ml-2 text-red-400 hover:text-red-300 font-bold hover:bg-red-900 hover:bg-opacity-30 px-2 py-1 rounded transition"
                          >
                            ✕
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Group By Modal Overlay */}
      {showGroupModal && selectedModel && (
        <div className="fixed inset-0 bg-black bg-opacity-70 flex items-center justify-center z-50 p-4">
          <div className="bg-gray-800 rounded-xl shadow-2xl border border-purple-600 max-w-md w-full">
            <div className="border-b border-gray-700 px-6 py-4 flex items-center justify-between bg-gray-800">
              <h3 className="text-xl font-bold text-white flex items-center gap-2">
                <span className="text-purple-400">⊕</span>
                Group By
              </h3>
              <button
                onClick={() => setShowGroupModal(false)}
                className="text-gray-400 hover:text-white font-bold text-xl transition"
              >
                ✕
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-bold text-gray-300 mb-2 uppercase tracking-wide">Select Field</label>
                <select
                  value={groupByField || ''}
                  onChange={(e) => {
                    setGroupByField(e.target.value || null)
                    if (e.target.value) {
                      setShowGroupView(true)
                    }
                  }}
                  className="w-full px-4 py-2 bg-gray-700 border border-gray-600 text-white rounded-lg focus:border-purple-500 focus:ring-2 focus:ring-purple-500 focus:ring-opacity-20 font-medium"
                >
                  <option value="">No grouping</option>
                  {currentModelFields.map((field) => (
                    <option key={field} value={field}>
                      {field.replace('_', ' ')}
                    </option>
                  ))}
                </select>
              </div>

              {groupByField && (
                <button
                  onClick={() => setShowGroupView(!showGroupView)}
                  className="w-full px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all font-bold shadow-lg"
                >
                  {showGroupView ? '📊 Show Table View' : '📈 Show Group View'}
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Detail View Slide Panel */}
      <DetailView
        modelName={selectedModel || ''}
        selectedRow={selectedRow}
        isOpen={showDetailView}
        onClose={() => setShowDetailView(false)}
      />
    </div>
  )
}

function App() {
  return (
    <AppProvider>
      <AppContent />
    </AppProvider>
  )
}

export default App

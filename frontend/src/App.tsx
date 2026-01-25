import { useContext, useState } from 'react'
import { ModelExplorer } from './components/ModelExplorer/ModelExplorer'
import { ListView } from './components/ListView/ListView'
import { FilterBuilder } from './components/FilterBuilder/FilterBuilder'
import { GroupView } from './components/GroupView/GroupView'
import { AppContext, AppProvider } from './state/AppContext'

interface Filter {
  id: string
  field: string
  operator: string
  value: string
}

function AppContent() {
  const { state, setState } = useContext(AppContext)
  const [selectedModel, setSelectedModel] = useState<string | null>(null)
  const [filters, setFilters] = useState<Filter[]>([])
  const [groupByField, setGroupByField] = useState<string | null>(null)
  const [showGroupView, setShowGroupView] = useState(false)

  const models = [
    {
      name: 'users',
      table: 'users',
      primaryKey: 'id',
      fields: ['id', 'name', 'email', 'created_at'],
    },
    {
      name: 'orders',
      table: 'orders',
      primaryKey: 'id',
      fields: ['id', 'user_id', 'total', 'created_at'],
    },
    {
      name: 'products',
      table: 'products',
      primaryKey: 'id',
      fields: ['id', 'name', 'price', 'stock', 'created_at'],
    },
  ]

  const handleSelectModel = (modelName: string) => {
    setSelectedModel(modelName)
    setFilters([])
    setGroupByField(null)
    setShowGroupView(false)
  }

  const handleAddFilter = (filter: Omit<Filter, 'id'>) => {
    setFilters([...filters, { ...filter, id: Date.now().toString() }])
  }

  const handleRemoveFilter = (filterId: string) => {
    setFilters(filters.filter((f) => f.id !== filterId))
  }

  const currentModel = models.find((m) => m.name === selectedModel)

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50 flex flex-col">
      {/* Header */}
      <header className="bg-gradient-to-r from-blue-600 to-indigo-600 shadow-lg">
        <div className="max-w-full mx-auto px-6 py-5">
          <h1 className="text-4xl font-bold text-white">Universal Data Viewer</h1>
          <p className="text-blue-100 mt-1">Explore and analyze your data with ease</p>
        </div>
      </header>

      {/* Main Content - Two Column Layout */}
      <main className="flex-1 overflow-hidden">
        <div className="grid grid-cols-5 h-full gap-0">
          {/* Left Sidebar - Models */}
          <aside className="col-span-1 bg-white border-r-2 border-blue-100 overflow-y-auto shadow-sm">
            <div className="p-6">
              <h2 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                <span className="w-1 h-6 bg-gradient-to-b from-blue-500 to-indigo-600 rounded"></span>
                Models
              </h2>
              <div className="space-y-2">
                {models.map((model) => (
                  <button
                    key={model.name}
                    onClick={() => handleSelectModel(model.name)}
                    className={`w-full px-4 py-3 text-left rounded-lg transition-all font-medium ${
                      selectedModel === model.name
                        ? 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white shadow-md'
                        : 'bg-gray-50 text-gray-700 hover:bg-blue-50'
                    }`}
                  >
                    {model.name.charAt(0).toUpperCase() + model.name.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          </aside>

          {selectedModel && (
            /* Right Content Area - Filters, Group By, and Data */
            <div className="col-span-4 bg-white overflow-hidden flex flex-col">
              {/* Top Section - Filters & Group By */}
              <div className="border-b-2 border-blue-100 bg-gradient-to-r from-blue-50 to-indigo-50 p-6 overflow-y-auto max-h-64">
                <div className="grid grid-cols-2 gap-8">
                  {/* Filters Column */}
                  <div>
                    <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <span className="w-1 h-6 bg-blue-500 rounded"></span>
                      Filters
                    </h3>
                    <FilterBuilder
                      fields={currentModel?.fields || []}
                      onAddFilter={handleAddFilter}
                    />
                    {filters.length > 0 && (
                      <div className="mt-4 space-y-2 max-h-40 overflow-y-auto">
                        {filters.map((filter) => (
                          <div
                            key={filter.id}
                            className="p-3 bg-gradient-to-r from-blue-100 to-indigo-100 border border-blue-300 rounded-lg flex justify-between items-start hover:shadow-md transition-shadow"
                          >
                            <div className="text-sm flex-1">
                              <p className="font-medium text-gray-900">
                                {filter.field} <span className="text-blue-600">{filter.operator}</span>
                              </p>
                              <p className="text-gray-700">{filter.value}</p>
                            </div>
                            <button
                              onClick={() => handleRemoveFilter(filter.id)}
                              className="ml-2 text-red-500 hover:text-red-700 font-bold hover:bg-red-50 px-2 py-1 rounded transition-colors"
                            >
                              ✕
                            </button>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Group By Column */}
                  <div>
                    <h3 className="text-lg font-bold text-gray-900 mb-4 flex items-center gap-2">
                      <span className="w-1 h-6 bg-indigo-500 rounded"></span>
                      Group By
                    </h3>
                    <select
                      value={groupByField || ''}
                      onChange={(e) => {
                        setGroupByField(e.target.value || null)
                        if (e.target.value) {
                          setShowGroupView(true)
                        }
                      }}
                      className="w-full px-4 py-2 border-2 border-blue-200 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white text-gray-900 font-medium"
                    >
                      <option value="">No grouping</option>
                      {currentModel?.fields.map((field) => (
                        <option key={field} value={field}>
                          {field.replace('_', ' ')}
                        </option>
                      ))}
                    </select>

                    {groupByField && (
                      <button
                        onClick={() => setShowGroupView(!showGroupView)}
                        className="w-full mt-3 px-4 py-2 bg-gradient-to-r from-green-500 to-emerald-600 text-white rounded-lg hover:from-green-600 hover:to-emerald-700 transition-all font-medium shadow-md hover:shadow-lg"
                      >
                        {showGroupView ? '📊 Show Table View' : '📈 Show Group View'}
                      </button>
                    )}
                  </div>
                </div>
              </div>

              {/* Bottom Section - Data Display */}
              <div className="flex-1 overflow-y-auto p-6">
                {/* Selected Model Heading */}
                <div className="mb-6">
                  <h2 className="text-3xl font-bold text-gray-900 capitalize">{selectedModel}</h2>
                  <p className="text-gray-600 mt-2 flex gap-4">
                    <span>
                      <span className="font-semibold text-gray-900">Table:</span> {currentModel?.table}
                    </span>
                    {filters.length > 0 && (
                      <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                        {filters.length} filter(s)
                      </span>
                    )}
                    {groupByField && (
                      <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                        Grouped by {groupByField}
                      </span>
                    )}
                  </p>
                </div>

                {/* Data Display */}
                {showGroupView && groupByField ? (
                  <div className="bg-white rounded-lg shadow-lg overflow-hidden">
                    <GroupView
                      modelName={selectedModel}
                      groupByField={groupByField}
                      filters={filters}
                    />
                  </div>
                ) : (
                  <div className="bg-white rounded-lg shadow-lg overflow-hidden">
                    <ListView
                      modelName={selectedModel}
                      filters={filters}
                    />
                  </div>
                )}
              </div>
            </div>
          )}

          {!selectedModel && (
            <div className="col-span-4 flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-50">
              <div className="text-center">
                <div className="text-6xl mb-4">📊</div>
                <p className="text-2xl text-gray-600 font-semibold">Select a model to get started</p>
                <p className="text-gray-500 mt-2">Choose from Users, Orders, or Products</p>
              </div>
            </div>
          )}
        </div>
      </main>
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

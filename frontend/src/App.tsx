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
    <div className="min-h-screen bg-gray-100 flex flex-col">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-full mx-auto px-6 py-4">
          <h1 className="text-4xl font-bold text-gray-900">Universal Data Viewer</h1>
        </div>
      </header>

      {/* Main Content - Two Column Layout */}
      <main className="flex-1 overflow-hidden">
        <div className="grid grid-cols-5 h-full gap-0">
          {/* Left Sidebar - Models */}
          <aside className="col-span-1 bg-white border-r border-gray-200 overflow-y-auto">
            <div className="p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Models</h2>
              <div className="space-y-2">
                {models.map((model) => (
                  <button
                    key={model.name}
                    onClick={() => handleSelectModel(model.name)}
                    className={`w-full px-4 py-3 text-left rounded-lg transition-colors font-medium ${
                      selectedModel === model.name
                        ? 'bg-blue-500 text-white'
                        : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    {model.name.charAt(0).toUpperCase() + model.name.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          </aside>

          {selectedModel && (
            <>
              {/* Middle - Filters & Group By */}
              <aside className="col-span-1 bg-white border-r border-gray-200 overflow-y-auto">
                <div className="p-6 space-y-6">
                  {/* Filters Section */}
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Filters</h3>
                    <FilterBuilder
                      fields={currentModel?.fields || []}
                      onAddFilter={handleAddFilter}
                    />
                    {filters.length > 0 && (
                      <div className="mt-4 space-y-2">
                        {filters.map((filter) => (
                          <div
                            key={filter.id}
                            className="p-3 bg-blue-50 border border-blue-200 rounded-lg flex justify-between items-start"
                          >
                            <div className="text-sm flex-1">
                              <p className="font-medium text-gray-900">
                                {filter.field} {filter.operator}
                              </p>
                              <p className="text-gray-600">{filter.value}</p>
                            </div>
                            <button
                              onClick={() => handleRemoveFilter(filter.id)}
                              className="ml-2 text-red-500 hover:text-red-700 font-medium"
                            >
                              ✕
                            </button>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Group By Section */}
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-4">Group By</h3>
                    <select
                      value={groupByField || ''}
                      onChange={(e) => {
                        setGroupByField(e.target.value || null)
                        if (e.target.value) {
                          setShowGroupView(true)
                        }
                      }}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500"
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
                        className="w-full mt-3 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-medium"
                      >
                        {showGroupView ? 'Show Table View' : 'Show Group View'}
                      </button>
                    )}
                  </div>
                </div>
              </aside>

              {/* Right Content Area - Data Display */}
              <div className="col-span-3 bg-gray-50 overflow-y-auto flex flex-col">
                <div className="flex-1 overflow-y-auto p-8">
                  {/* Selected Model Heading */}
                  <div className="mb-6">
                    <h2 className="text-3xl font-bold text-gray-900 capitalize">{selectedModel}</h2>
                    <p className="text-gray-600 mt-1">
                      {currentModel?.table}
                      {filters.length > 0 && ` • ${filters.length} filter(s) applied`}
                      {groupByField && ` • Grouped by ${groupByField}`}
                    </p>
                  </div>

                  {/* Data Display */}
                  {showGroupView && groupByField ? (
                    <div className="bg-white rounded-lg shadow">
                      <GroupView
                        modelName={selectedModel}
                        groupByField={groupByField}
                        filters={filters}
                      />
                    </div>
                  ) : (
                    <div className="bg-white rounded-lg shadow">
                      <ListView
                        modelName={selectedModel}
                        filters={filters}
                      />
                    </div>
                  )}
                </div>
              </div>
            </>
          )}

          {!selectedModel && (
            <div className="col-span-4 flex items-center justify-center">
              <div className="text-center">
                <p className="text-xl text-gray-500">Select a model from the left to view data</p>
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

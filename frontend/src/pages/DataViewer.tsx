// DataViewer Page Component (extracted from App.tsx)
import { useState, useEffect } from 'react'
import { ListView } from '../components/ListView/ListView'
import { FilterBuilder } from '../components/FilterBuilder/FilterBuilder'
import { GroupView } from '../components/GroupView/GroupView'
import DetailView from '../components/DetailView/DetailView'
import { SearchBar } from '../components/SearchBar/SearchBar'
import { CreateForm } from '../components/CreateForm/CreateForm'
import { EditForm } from '../components/EditForm/EditForm'
import { DeleteConfirm } from '../components/DeleteConfirm/DeleteConfirm'
import { fetchModels, getSearchableFields, type Model } from '../api/client'
import { useDebounce } from '../hooks/useDebounce'
import { useAuth } from '../contexts/AuthContext'

interface Filter {
  id: string
  field: string
  operator: string
  value: string
}

export function DataViewer() {
  const { logout, user } = useAuth()
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
  const [searchInput, setSearchInput] = useState<string>('')
  const [searchMode, setSearchMode] = useState<'global' | 'column'>('global')
  const [selectedSearchColumn, setSelectedSearchColumn] = useState<string>('')
  const [columnSearchInput, setColumnSearchInput] = useState<string>('')
  const debouncedSearchQuery = useDebounce(searchInput, 500)
  const debouncedColumnSearchQuery = useDebounce(columnSearchInput, 500)
  
  // CRUD state
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [showEditForm, setShowEditForm] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)
  const [recordToEdit, setRecordToEdit] = useState<Record<string, any> | null>(null)
  const [recordToDelete, setRecordToDelete] = useState<Record<string, any> | null>(null)
  const [refreshTrigger, setRefreshTrigger] = useState(0)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

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
    setSearchInput('') // Reset search when model changes
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

  const handleClearSearch = () => {
    setSearchInput('')
  }

  // CRUD handlers
  const handleCreate = () => {
    setShowCreateForm(true)
  }

  const handleCreateSuccess = () => {
    setShowCreateForm(false)
    setSuccessMessage('Record created successfully!')
    setRefreshTrigger((prev) => prev + 1)
    setTimeout(() => setSuccessMessage(null), 3000)
  }

  const handleEdit = (row: Record<string, any>) => {
    setRecordToEdit(row)
    setShowEditForm(true)
  }

  const handleEditSuccess = () => {
    setShowEditForm(false)
    setRecordToEdit(null)
    setSuccessMessage('Record updated successfully!')
    setRefreshTrigger((prev) => prev + 1)
    setTimeout(() => setSuccessMessage(null), 3000)
  }

  const handleDelete = (row: Record<string, any>) => {
    setRecordToDelete(row)
    setShowDeleteConfirm(true)
  }

  const handleDeleteSuccess = () => {
    setShowDeleteConfirm(false)
    setRecordToDelete(null)
    setSuccessMessage('Record deleted successfully!')
    setRefreshTrigger((prev) => prev + 1)
    setTimeout(() => setSuccessMessage(null), 3000)
  }

  const currentModel = models.find((m) => m.name === selectedModel)

  // Extract field names from model
  const currentModelFields = currentModel?.fields?.map((f) => f.name) || []

  // Get searchable fields (string/text types only)
  const searchableFields = currentModel ? getSearchableFields(currentModel) : []

  return (
    <div className="min-h-screen bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-gray-800 border-b border-gray-700 shadow-lg">
        <div className="max-w-full mx-auto px-6 py-4">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold text-white">Agent P</h1>
              <p className="text-gray-400 mt-1">Universal Data Viewer</p>
              {error && <p className="text-red-400 text-sm mt-1">⚠️ {error}</p>}
            </div>
            <div className="flex items-center gap-4">
              {user && (
                <div className="flex items-center gap-2 text-gray-300">
                  <span className="text-sm">👤 {user}</span>
                </div>
              )}
              {selectedModel && (
                <div className="flex gap-3">
                  <button
                    onClick={handleCreate}
                    className="px-4 py-2 rounded-lg font-semibold flex items-center gap-2 transition-all bg-green-600 text-white hover:bg-green-700 shadow-lg"
                    title="Create new record"
                  >
                    ➕ Create
                  </button>
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
              <button
                onClick={logout}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-lg font-semibold transition-all"
              >
                Logout
              </button>
            </div>
          </div>
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
                    {/* Search Bar Above Table */}
                    {selectedModel && (
                      <div className="mb-4 bg-gray-800 rounded-lg p-4 border border-gray-700">
                        <div className="flex items-center gap-4">
                          <div className="flex items-center gap-2">
                            <label className="text-sm text-gray-300 font-medium">Search Mode:</label>
                            <select
                              value={searchMode}
                              onChange={(e) => {
                                setSearchMode(e.target.value as 'global' | 'column')
                                if (e.target.value === 'global') {
                                  setColumnSearchInput('')
                                  setSelectedSearchColumn('')
                                } else {
                                  setSearchInput('')
                                }
                              }}
                              className="px-3 py-1 bg-gray-700 border border-gray-600 text-white rounded-lg text-sm 
                                focus:border-cyan-500 focus:ring-2 focus:ring-cyan-500 focus:ring-opacity-20"
                            >
                              <option value="global">Global (All Fields)</option>
                              <option value="column">Column Specific</option>
                            </select>
                          </div>

                          {searchMode === 'global' ? (
                            <div className="flex-1 max-w-md">
                              <SearchBar
                                value={searchInput}
                                onChange={setSearchInput}
                                onClear={handleClearSearch}
                                placeholder={`Search across all fields in ${selectedModel}...`}
                                disabled={!selectedModel}
                                loading={false}
                              />
                            </div>
                          ) : (
                            <div className="flex-1 flex items-center gap-2">
                              <select
                                value={selectedSearchColumn}
                                onChange={(e) => {
                                  setSelectedSearchColumn(e.target.value)
                                  setColumnSearchInput('')
                                }}
                                className="px-3 py-2 bg-gray-700 border border-gray-600 text-white rounded-lg text-sm 
                                  focus:border-cyan-500 focus:ring-2 focus:ring-cyan-500 focus:ring-opacity-20"
                              >
                                <option value="">Select column...</option>
                                {currentModelFields.map((field) => (
                                  <option key={field} value={field}>
                                    {field.replace('_', ' ')}
                                  </option>
                                ))}
                              </select>
                              {selectedSearchColumn && (
                                <div className="flex-1 max-w-md">
                                  <SearchBar
                                    value={columnSearchInput}
                                    onChange={setColumnSearchInput}
                                    onClear={() => {
                                      setColumnSearchInput('')
                                      setSelectedSearchColumn('')
                                    }}
                                    placeholder={`Search in ${selectedSearchColumn.replace('_', ' ')}...`}
                                    disabled={!selectedSearchColumn}
                                    loading={false}
                                  />
                                </div>
                              )}
                            </div>
                          )}
                        </div>
                      </div>
                    )}

                    {showGroupView && groupByField ? (
                      <div className="bg-gray-800 rounded-lg shadow-lg overflow-hidden border border-gray-700">
                        <GroupView
                          modelName={selectedModel}
                          groupByField={groupByField}
                          filters={filters}
                          searchQuery={searchMode === 'global' ? debouncedSearchQuery : ''}
                          columnSearchQuery={searchMode === 'column' && selectedSearchColumn ? debouncedColumnSearchQuery : ''}
                          columnSearchField={searchMode === 'column' ? selectedSearchColumn : ''}
                        />
                      </div>
                    ) : (
                      <div className="bg-gray-800 rounded-lg shadow-lg overflow-hidden border border-gray-700">
                        <ListView
                          modelName={selectedModel}
                          filters={filters}
                          modelFields={currentModelFields}
                          onRowClick={handleRowClick}
                          searchQuery={searchMode === 'global' ? debouncedSearchQuery : ''}
                          searchFields={searchableFields}
                          columnSearchQuery={searchMode === 'column' && selectedSearchColumn ? debouncedColumnSearchQuery : ''}
                          columnSearchField={searchMode === 'column' ? selectedSearchColumn : ''}
                          onEdit={handleEdit}
                          onDelete={handleDelete}
                          key={refreshTrigger}
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
              <FilterBuilder fields={currentModelFields} onAddFilter={handleAddFilter} />

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

      {/* Success Message Notification */}
      {successMessage && (
        <div className="fixed top-4 right-4 bg-green-600 text-white px-6 py-4 rounded-lg shadow-xl z-50 flex items-center gap-3 animate-slide-in">
          <span className="text-2xl">✅</span>
          <span className="font-semibold">{successMessage}</span>
        </div>
      )}

      {/* Create Form Modal */}
      {showCreateForm && selectedModel && (
        <CreateForm
          modelName={selectedModel}
          onSuccess={handleCreateSuccess}
          onCancel={() => setShowCreateForm(false)}
        />
      )}

      {/* Edit Form Modal */}
      {showEditForm && selectedModel && recordToEdit && currentModel && (
        <EditForm
          modelName={selectedModel}
          recordId={recordToEdit[currentModel.primary_key] || recordToEdit.id}
          initialData={recordToEdit}
          onSuccess={handleEditSuccess}
          onCancel={() => {
            setShowEditForm(false)
            setRecordToEdit(null)
          }}
        />
      )}

      {/* Delete Confirmation Modal */}
      {showDeleteConfirm && selectedModel && recordToDelete && currentModel && (
        <DeleteConfirm
          modelName={selectedModel}
          recordId={recordToDelete[currentModel.primary_key] || recordToDelete.id}
          recordData={recordToDelete}
          onSuccess={handleDeleteSuccess}
          onCancel={() => {
            setShowDeleteConfirm(false)
            setRecordToDelete(null)
          }}
        />
      )}
    </div>
  )
}


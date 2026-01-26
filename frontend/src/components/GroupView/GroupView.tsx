// GroupView Component
import { useState, useEffect } from 'react'
import type React from 'react'
import { executeQuery, buildDSLQuery, type QueryResponse } from '../../api/client'

interface Filter {
  id: string
  field: string
  operator: string
  value: string
}

interface GroupViewProps {
  modelName?: string
  groupByField?: string
  filters?: Filter[]
}

const mockData: Record<string, any[]> = {
  users: [
    { id: 1, name: 'John Doe', email: 'john@example.com', created_at: '2024-01-15' },
    { id: 2, name: 'Jane Smith', email: 'jane@example.com', created_at: '2024-01-16' },
    { id: 3, name: 'Bob Johnson', email: 'bob@example.com', created_at: '2024-01-17' },
    { id: 4, name: 'Alice Williams', email: 'alice@example.com', created_at: '2024-01-18' },
    { id: 5, name: 'Charlie Brown', email: 'charlie@example.com', created_at: '2024-01-19' },
  ],
  orders: [
    { id: 101, user_id: 1, total: '$250.00', created_at: '2024-01-20' },
    { id: 102, user_id: 2, total: '$150.00', created_at: '2024-01-20' },
    { id: 103, user_id: 1, total: '$500.00', created_at: '2024-01-21' },
    { id: 104, user_id: 3, total: '$300.00', created_at: '2024-01-21' },
    { id: 105, user_id: 4, total: '$450.00', created_at: '2024-01-22' },
  ],
  products: [
    { id: 1, name: 'Laptop', price: '$999.99', stock: 15, created_at: '2023-06-01' },
    { id: 2, name: 'Mouse', price: '$29.99', stock: 100, created_at: '2023-06-02' },
    { id: 3, name: 'Keyboard', price: '$79.99', stock: 50, created_at: '2023-06-03' },
    { id: 4, name: 'Monitor', price: '$299.99', stock: 25, created_at: '2023-06-04' },
    { id: 5, name: 'USB Cable', price: '$9.99', stock: 200, created_at: '2023-06-05' },
  ],
}

function applyFilter(row: any, filter: Filter): boolean {
  const fieldValue = String(row[filter.field]).toLowerCase()
  const filterValue = filter.value.toLowerCase()

  switch (filter.operator) {
    case 'equals':
      return fieldValue === filterValue
    case 'contains':
      return fieldValue.includes(filterValue)
    case 'startswith':
      return fieldValue.startsWith(filterValue)
    case 'endswith':
      return fieldValue.endsWith(filterValue)
    case 'gt':
      return parseFloat(fieldValue) > parseFloat(filterValue)
    case 'lt':
      return parseFloat(fieldValue) < parseFloat(filterValue)
    case 'gte':
      return parseFloat(fieldValue) >= parseFloat(filterValue)
    case 'lte':
      return parseFloat(fieldValue) <= parseFloat(filterValue)
    default:
      return true
  }
}

// Helper function to detect and render different data types
function renderCellValue(value: any): React.ReactNode {
  if (typeof value === 'string' && value.includes('%')) {
    const percentage = parseFloat(value)
    return (
      <div className="flex items-center gap-2">
        <div className="w-20 h-2 bg-gray-700 rounded-full overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-cyan-500 to-purple-600 rounded-full"
            style={{ width: `${percentage}%` }}
          />
        </div>
        <span className="text-sm text-gray-300">{value}</span>
      </div>
    )
  }

  if (typeof value === 'string' && (value === 'Sent' || value === 'Draft' || value === 'Pending')) {
    const badgeColors: Record<string, string> = {
      Sent: 'bg-green-900 text-green-300',
      Draft: 'bg-blue-900 text-blue-300',
      Pending: 'bg-yellow-900 text-yellow-300',
    }
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-semibold ${badgeColors[value] || 'bg-gray-700 text-gray-300'}`}>
        {value}
      </span>
    )
  }

  return String(value)
}

function convertOperator(operator: string): string {
  const mapping: Record<string, string> = {
    equals: '=',
    contains: 'like',
    startswith: 'starts_with',
    endswith: 'ends_with',
    gt: '>',
    lt: '<',
    gte: '>=',
    lte: '<=',
  }
  return mapping[operator] || operator
}

export function GroupView({ modelName = 'users', groupByField = '', filters = [] }: GroupViewProps) {
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set())
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Fetch grouped data when model, filters, or groupByField change
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)

        // Build DSL query with GROUP BY
        const dslFilters = filters.map((f) => ({
          field: f.field,
          op: convertOperator(f.operator),
          value: isNaN(Number(f.value)) ? f.value : Number(f.value),
        }))

        const query = buildDSLQuery(
          modelName,
          undefined,
          dslFilters.length > 0 ? dslFilters : undefined,
          groupByField,
          100,
          0
        )

        // Execute query via backend
        const response: QueryResponse = await executeQuery(query)

        if (response.error) {
          setError(response.error)
          // Fallback to mock data
          setData(mockData[modelName] || [])
        } else {
          // Use real data if available, otherwise fallback to mock data
          if (response.data && response.data.length > 0) {
            setData(response.data)
            console.log('Grouped data from backend:', response.data)
          } else {
            setData(mockData[modelName] || [])
            console.log('Using mock data (no backend results)')
          }
          console.log('Generated GROUP BY SQL:', response.sql)
          console.log('Parameters:', response.params)
        }
      } catch (err) {
        console.error('Error fetching grouped data:', err)
        setError(err instanceof Error ? err.message : 'Failed to load grouped data')
        // Fallback to mock data
        setData(mockData[modelName] || [])
      } finally {
        setLoading(false)
      }
    }

    if (groupByField) {
      fetchData()
    }
  }, [modelName, groupByField, filters])

  let displayData = data

  // Apply filters if not using backend (fallback)
  if (displayData.length === 0 && filters.length > 0) {
    displayData = mockData[modelName]?.filter((row) =>
      filters.every((filter) => applyFilter(row, filter))
    ) || []
  }

  // Group data
  const grouped: Record<string, any[]> = {}
  displayData.forEach((row) => {
    const key = String(row[groupByField] || 'Other')
    if (!grouped[key]) {
      grouped[key] = []
    }
    grouped[key].push(row)
  })

  const groups = Object.entries(grouped).sort(([keyA], [keyB]) => keyA.localeCompare(keyB))

  const toggleGroup = (key: string) => {
    const newExpanded = new Set(expandedGroups)
    if (newExpanded.has(key)) {
      newExpanded.delete(key)
    } else {
      newExpanded.add(key)
    }
    setExpandedGroups(newExpanded)
  }

  return (
    <div className="w-full">
      {loading && (
        <div className="text-center py-8 text-gray-400">
          <span className="inline-block animate-spin mr-2">⚙️</span>
          Loading grouped data...
        </div>
      )}

      {error && (
        <div className="p-4 bg-red-900 bg-opacity-20 border border-red-700 text-red-300 rounded">
          ⚠️ {error}
        </div>
      )}

      {!loading && !error && (
        <>
          {groups.map(([groupKey, groupData]) => (
            <div key={groupKey} className="mb-1">
              {/* Group Header */}
              <button
                onClick={() => toggleGroup(groupKey)}
                className="w-full bg-gray-800 hover:bg-gray-750 px-6 py-4 flex items-center gap-3 border-b border-gray-700 transition cursor-pointer group"
              >
                <span className={`text-cyan-400 text-lg transition-transform ${expandedGroups.has(groupKey) ? 'rotate-90' : ''}`}>
                  ▶
                </span>
                <h3 className="text-base font-bold text-white capitalize flex-1 text-left">
                  {groupByField}: <span className="text-cyan-300">{groupKey}</span>
                </h3>
                <span className="text-xs font-semibold text-gray-400 bg-gray-700 px-2 py-1 rounded-full group-hover:text-cyan-300">
                  {groupData.length}
                </span>
              </button>

              {/* Group Content */}
              {expandedGroups.has(groupKey) && (
                <div className="bg-gray-900 border-b border-gray-700">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="bg-gray-800 border-b border-gray-700">
                        {Object.keys(groupData[0] || {}).map((field) => (
                          <th
                            key={field}
                            className="px-4 py-3 text-left text-xs font-bold text-cyan-400 uppercase tracking-wider"
                          >
                            {field.replace('_', ' ')}
                          </th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {groupData.map((row, idx) => (
                        <tr
                          key={idx}
                          className={`border-b border-gray-700 hover:bg-gray-800 transition ${
                            idx % 2 === 0 ? 'bg-gray-900' : 'bg-gray-850'
                          }`}
                        >
                          {Object.keys(row).map((field) => (
                            <td key={`${idx}-${field}`} className="px-4 py-3 text-gray-200">
                              {renderCellValue(row[field])}
                            </td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          ))}

          {groups.length === 0 && (
            <div className="text-center py-16 bg-gray-800 rounded-lg border border-gray-700">
              <p className="text-xl text-gray-400 font-semibold">
                {filters.length > 0 ? '🔍 No data matches the applied filters' : '📭 No data available'}
              </p>
            </div>
          )}
        </>
      )}
    </div>
  )
}

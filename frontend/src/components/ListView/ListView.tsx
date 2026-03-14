// ListView Component
import React, { useState, useEffect, useContext } from 'react'
import {
  executeQuery,
  buildDSLQuery,
  buildSearchQuery,
  type QueryResponse,
  type Sort,
} from '../../api/client'
import { Pagination } from '../Pagination/Pagination'
import { ObjectRenderer } from '../ObjectRenderer/ObjectRenderer'
import { AppContext } from '../../state/AppContext'

interface Filter {
  id: string
  field: string
  operator: string
  value: string
}

interface ListViewProps {
  modelName?: string
  filters?: Filter[]
  modelFields?: string[]
  onRowClick?: (row: Record<string, any>) => void
  searchQuery?: string
  searchFields?: string[]
  columnSearchQuery?: string
  columnSearchField?: string
  onEdit?: (row: Record<string, any>) => void
  onDelete?: (row: Record<string, any>) => void
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

export function ListView({
  modelName = 'users',
  filters = [],
  modelFields = [],
  onRowClick,
  searchQuery = '',
  searchFields = [],
  columnSearchQuery = '',
  columnSearchField = '',
  onEdit,
  onDelete,
}: ListViewProps) {
  const { state } = useContext(AppContext)
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [pageSize, setPageSize] = useState(25)
  const [totalCount, setTotalCount] = useState(0)
  const [sort, setSort] = useState<Sort | null>(null)

  // Reset to page 1 when model, filters, sort, or search changes
  useEffect(() => {
    setCurrentPage(1)
  }, [modelName, filters, sort, searchQuery])

  // Fetch data when model, filters, pagination, sort, or search changes
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)

        // Build DSL query from filters
        const dslFilters = filters.map((f) => ({
          field: f.field,
          op: convertOperator(f.operator),
          value: isNaN(Number(f.value)) ? f.value : Number(f.value),
        }))

        // Build search filter - either global or column-specific
        let searchFilter = null
        if (columnSearchQuery && columnSearchField) {
          // Column-specific search
          searchFilter = {
            field: columnSearchField,
            op: 'contains',
            value: columnSearchQuery.trim(),
          }
        } else if (searchQuery && searchFields.length > 0) {
          // Global search across multiple fields
          searchFilter = buildSearchQuery(searchQuery, searchFields, 'contains')
        }

        // Calculate offset from current page
        const offset = (currentPage - 1) * pageSize

        // Build sort array if sort is set
        const sortArray = sort ? [sort] : undefined

        const query = buildDSLQuery(
          modelName,
          modelFields.length > 0 ? modelFields : undefined,
          dslFilters.length > 0 ? dslFilters : undefined,
          undefined,
          pageSize,
          offset,
          sortArray,
          searchFilter
        )

        // Execute query via backend
        const response: QueryResponse = await executeQuery(query)

        if (response.error) {
          setError(response.error)
          // Fallback to mock data on error
          const mockDataArray = mockData[modelName] || []
          setData(mockDataArray)
          setTotalCount(mockDataArray.length)
        } else {
          // Use backend data (filtered on server)
          if (response.data && response.data.length > 0) {
            setData(response.data)
            // Use meta.total if available, otherwise estimate from data length
            if (response.meta?.total !== undefined) {
              setTotalCount(response.meta.total)
            } else if (response.total !== undefined) {
              setTotalCount(response.total)
            } else {
              // Estimate: if we got a full page, there might be more
              setTotalCount(response.data.length === pageSize ? response.data.length + 1 : response.data.length)
            }
            console.log('Data from backend:', response.data)
          } else {
            // No results from backend query (filters applied server-side)
            setData([])
            setTotalCount(0)
            console.log('No results from backend query')
          }
          console.log('Generated SQL:', response.sql)
          console.log('Parameters:', response.params)
        }
      } catch (err) {
        console.error('Error fetching data:', err)
        setError(err instanceof Error ? err.message : 'Failed to load data')
        // Fallback to mock data
        const mockDataArray = mockData[modelName] || []
        setData(mockDataArray)
        setTotalCount(mockDataArray.length)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [modelName, filters, modelFields, currentPage, pageSize, sort, searchQuery, searchFields, columnSearchQuery, columnSearchField])

  const handleSort = (field: string) => {
    if (sort && sort.field === field) {
      // Toggle direction if same field
      setSort(sort.direction === 'asc' ? { field, direction: 'desc' } : null)
    } else {
      // New field, default to ascending
      setSort({ field, direction: 'asc' })
    }
  }

  const handlePageChange = (page: number) => {
    setCurrentPage(page)
    // Scroll to top of table
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const handlePageSizeChange = (size: number) => {
    setPageSize(size)
    setCurrentPage(1) // Reset to first page when page size changes
  }

  const columns = data.length > 0 ? Object.keys(data[0]) : modelFields.length > 0 ? modelFields : []

  const getSortIndicator = (column: string) => {
    if (sort && sort.field === column) {
      return sort.direction === 'asc' ? ' ↑' : ' ↓'
    }
    return ''
  }

  // Highlight search matches in text
  const highlightMatch = (text: string, searchTerm: string): React.ReactNode => {
    if (!searchTerm || !text) return String(text)

    const escapeRegex = (str: string) => str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
    const regex = new RegExp(`(${escapeRegex(searchTerm)})`, 'gi')
    const parts = String(text).split(regex)

    return parts.map((part, i) =>
      regex.test(part) ? (
        <mark key={i} className="bg-yellow-600 text-yellow-100 px-1 rounded">
          {part}
        </mark>
      ) : (
        <span key={i}>{part}</span>
      )
    )
  }

  return (
    <div className="flex flex-col">
      <div className="overflow-x-auto">
        {loading && (
          <div className="text-center py-8 text-gray-400">
            <span className="inline-block animate-spin mr-2">⚙️</span>
            Loading data...
          </div>
        )}

        {error && (
          <div className="p-4 bg-red-900 bg-opacity-20 border border-red-700 text-red-300 rounded">
            ⚠️ {error}
          </div>
        )}

        {!loading && !error && (
          <table className="w-full divide-y divide-gray-700">
            <thead className="bg-gray-800 border-b-2 border-cyan-600">
              <tr>
                {columns.map((column) => (
                  <th
                    key={column}
                    onClick={() => handleSort(column)}
                    className="px-6 py-4 text-left text-sm font-bold text-cyan-400 capitalize cursor-pointer hover:bg-gray-700 transition-colors select-none"
                    title="Click to sort"
                  >
                    <span className="flex items-center gap-2">
                      {column.replace('_', ' ')}
                      <span className="text-xs">
                        {getSortIndicator(column) || (
                          <span className="text-gray-600 opacity-50">↕</span>
                        )}
                      </span>
                    </span>
                  </th>
                ))}
                {(onEdit || onDelete) && (
                  <th className="px-6 py-4 text-left text-sm font-bold text-cyan-400">
                    Actions
                  </th>
                )}
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {data.map((row, idx) => (
                <tr
                  key={idx}
                  className={`transition-colors ${
                    idx % 2 === 0 ? 'bg-gray-800' : 'bg-gray-750'
                  } hover:bg-gray-700`}
                >
                  {columns.map((column) => (
                    <td
                      key={`${idx}-${column}`}
                      onClick={() => onRowClick?.(row)}
                      className="px-6 py-4 text-sm text-gray-200 cursor-pointer"
                    >
                      {(searchQuery && searchFields.includes(column)) || 
                       (columnSearchQuery && columnSearchField === column)
                        ? highlightMatch(
                            String(row[column] ?? ''),
                            searchQuery || columnSearchQuery
                          )
                        : <ObjectRenderer 
                            value={row[column] ?? ''} 
                            isMongoDb={state.databaseType === 'mongo'}
                          />}
                    </td>
                  ))}
                  {(onEdit || onDelete) && (
                    <td
                      className="px-6 py-4 text-sm"
                      onClick={(e) => e.stopPropagation()}
                    >
                      <div className="flex items-center gap-2">
                        {onEdit && (
                          <button
                            onClick={() => onEdit(row)}
                            className="px-3 py-1 bg-cyan-600 text-white rounded hover:bg-cyan-700 transition-colors text-xs font-semibold"
                            title="Edit record"
                          >
                            ✏️ Edit
                          </button>
                        )}
                        {onDelete && (
                          <button
                            onClick={() => onDelete(row)}
                            className="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700 transition-colors text-xs font-semibold"
                            title="Delete record"
                          >
                            🗑️ Delete
                          </button>
                        )}
                      </div>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        )}

        {!loading && !error && data.length === 0 && (
          <div className="text-center py-12 bg-gray-800">
            <p className="text-gray-400 text-lg">
              {filters.length > 0 ? '🔍 No data matches the applied filters' : '📭 No data available'}
            </p>
          </div>
        )}
      </div>

      {/* Pagination */}
      {!loading && !error && totalCount > 0 && (
        <Pagination
          currentPage={currentPage}
          pageSize={pageSize}
          totalCount={totalCount}
          onPageChange={handlePageChange}
          onPageSizeChange={handlePageSizeChange}
        />
      )}
    </div>
  )
}

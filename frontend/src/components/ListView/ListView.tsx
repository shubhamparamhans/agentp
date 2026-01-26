// ListView Component
import { useState, useEffect } from 'react'
import { executeQuery, buildDSLQuery, type QueryResponse } from '../../api/client'

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

export function ListView({ modelName = 'users', filters = [], modelFields = [], onRowClick }: ListViewProps) {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Fetch data when model or filters change
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

        const query = buildDSLQuery(
          modelName,
          modelFields.length > 0 ? modelFields : undefined,
          dslFilters.length > 0 ? dslFilters : undefined,
          undefined,
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
            console.log('Data from backend:', response.data)
          } else {
            setData(mockData[modelName] || [])
            console.log('Using mock data (no backend results)')
          }
          console.log('Generated SQL:', response.sql)
          console.log('Parameters:', response.params)
        }
      } catch (err) {
        console.error('Error fetching data:', err)
        setError(err instanceof Error ? err.message : 'Failed to load data')
        // Fallback to mock data
        setData(mockData[modelName] || [])
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [modelName, filters, modelFields])

  const columns = data.length > 0 ? Object.keys(data[0]) : modelFields.length > 0 ? modelFields : []

  return (
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
                  className="px-6 py-4 text-left text-sm font-bold text-cyan-400 capitalize"
                >
                  {column.replace('_', ' ')}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-700">
            {data.map((row, idx) => (
              <tr
                key={idx}
                onClick={() => onRowClick?.(row)}
                className={`transition-colors cursor-pointer ${
                  idx % 2 === 0 ? 'bg-gray-800' : 'bg-gray-750'
                } hover:bg-gray-700 hover:border-l-4 hover:border-cyan-600`}
              >
                {columns.map((column) => (
                  <td
                    key={`${idx}-${column}`}
                    className="px-6 py-4 text-sm text-gray-200"
                  >
                    {String(row[column])}
                  </td>
                ))}
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
  )
}

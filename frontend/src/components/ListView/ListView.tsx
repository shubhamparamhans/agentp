// ListView Component
interface Filter {
  id: string
  field: string
  operator: string
  value: string
}

interface ListViewProps {
  modelName?: string
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

export function ListView({ modelName = 'users', filters = [] }: ListViewProps) {
  let data = mockData[modelName] || []

  // Apply filters
  if (filters.length > 0) {
    data = data.filter((row) => filters.every((filter) => applyFilter(row, filter)))
  }

  const columns = data.length > 0 ? Object.keys(data[0]) : []

  return (
    <div className="overflow-x-auto">
      <table className="w-full divide-y divide-blue-100">
        <thead className="bg-gradient-to-r from-blue-50 to-indigo-50 border-b-2 border-blue-200">
          <tr>
            {columns.map((column) => (
              <th
                key={column}
                className="px-6 py-4 text-left text-sm font-bold text-gray-900 capitalize"
              >
                {column.replace('_', ' ')}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-blue-100">
          {data.map((row, idx) => (
            <tr
              key={idx}
              className={`transition-colors ${
                idx % 2 === 0 ? 'bg-white' : 'bg-blue-50'
              } hover:bg-blue-100`}
            >
              {columns.map((column) => (
                <td
                  key={`${idx}-${column}`}
                  className="px-6 py-4 text-sm text-gray-700"
                >
                  {String(row[column])}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>

      {data.length === 0 && (
        <div className="text-center py-12 bg-blue-50">
          <p className="text-gray-600 text-lg">
            {filters.length > 0 ? '🔍 No data matches the applied filters' : '📭 No data available'}
          </p>
        </div>
      )}
    </div>
  )
}

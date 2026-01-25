// GroupView Component
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

export function GroupView({ modelName = 'users', groupByField = '', filters = [] }: GroupViewProps) {
  let data = mockData[modelName] || []

  // Apply filters
  if (filters.length > 0) {
    data = data.filter((row) => filters.every((filter) => applyFilter(row, filter)))
  }

  // Group data
  const grouped: Record<string, any[]> = {}
  data.forEach((row) => {
    const key = String(row[groupByField] || 'Other')
    if (!grouped[key]) {
      grouped[key] = []
    }
    grouped[key].push(row)
  })

  const groups = Object.entries(grouped).sort(([keyA], [keyB]) => keyA.localeCompare(keyB))

  return (
    <div className="space-y-6 p-6">
      {groups.map(([groupKey, groupData]) => (
        <div key={groupKey} className="bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-6 rounded-xl border-2 border-blue-200 shadow-md hover:shadow-lg transition-shadow">
          <div className="flex items-center justify-between mb-5">
            <h3 className="text-lg font-bold text-gray-900 capitalize">
              {groupByField}: <span className="bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">{groupKey}</span>
            </h3>
            <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-bold bg-gradient-to-r from-blue-100 to-indigo-100 text-blue-800 border border-blue-300">
              {groupData.length} item{groupData.length !== 1 ? 's' : ''}
            </span>
          </div>

          {/* Summary Stats */}
          <div className="grid grid-cols-2 gap-4 mb-5">
            <div className="bg-gradient-to-br from-white to-blue-50 p-4 rounded-lg border border-blue-200">
              <p className="text-sm text-gray-600 font-semibold">Total Records</p>
              <p className="text-3xl font-bold text-blue-600 mt-1">{groupData.length}</p>
            </div>
            <div className="bg-gradient-to-br from-white to-indigo-50 p-4 rounded-lg border border-indigo-200">
              <p className="text-sm text-gray-600 font-semibold">Fields</p>
              <p className="text-3xl font-bold text-indigo-600 mt-1">{Object.keys(groupData[0] || {}).length}</p>
            </div>
          </div>

          {/* Group Data Table */}
          <div className="bg-white rounded-lg overflow-hidden border border-blue-100">
            <table className="w-full divide-y divide-blue-100 text-sm">
              <thead className="bg-gradient-to-r from-blue-100 to-indigo-100 border-b-2 border-blue-300">
                <tr>
                  {Object.keys(groupData[0] || {}).map((field) => (
                    <th
                      key={field}
                      className="px-4 py-3 text-left text-xs font-bold text-gray-800 uppercase tracking-wider capitalize"
                    >
                      {field.replace('_', ' ')}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-blue-100">
                {groupData.map((row, idx) => (
                  <tr
                    key={idx}
                    className={`transition-colors ${
                      idx % 2 === 0 ? 'bg-white' : 'bg-blue-50'
                    } hover:bg-blue-100`}
                  >
                    {Object.keys(row).map((field) => (
                      <td key={`${idx}-${field}`} className="px-4 py-3 text-gray-700 font-medium">
                        {String(row[field])}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      ))}

      {groups.length === 0 && (
        <div className="text-center py-16 bg-gradient-to-br from-blue-50 to-indigo-50 rounded-lg border-2 border-blue-200">
          <p className="text-xl text-gray-600 font-semibold">
            {filters.length > 0 ? '🔍 No data matches the applied filters' : '📭 No data available'}
          </p>
        </div>
      )}
    </div>
  )
}

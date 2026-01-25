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
    <div className="space-y-6">
      {groups.map(([groupKey, groupData]) => (
        <div key={groupKey} className="bg-gradient-to-r from-blue-50 to-indigo-50 p-6 rounded-lg border border-blue-200">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900 capitalize">
              {groupByField}: <span className="text-blue-600">{groupKey}</span>
            </h3>
            <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
              {groupData.length} item{groupData.length !== 1 ? 's' : ''}
            </span>
          </div>

          {/* Summary Stats */}
          <div className="grid grid-cols-2 gap-4 mb-4">
            <div className="bg-white p-3 rounded-lg">
              <p className="text-sm text-gray-600">Total Records</p>
              <p className="text-2xl font-bold text-gray-900">{groupData.length}</p>
            </div>
            <div className="bg-white p-3 rounded-lg">
              <p className="text-sm text-gray-600">Details</p>
              <p className="text-lg font-semibold text-gray-700">
                {Object.keys(groupData[0] || {}).length} fields
              </p>
            </div>
          </div>

          {/* Group Data Table */}
          <div className="bg-white rounded-lg overflow-hidden">
            <table className="w-full divide-y divide-gray-200 text-sm">
              <thead className="bg-gray-50">
                <tr>
                  {Object.keys(groupData[0] || {}).map((field) => (
                    <th
                      key={field}
                      className="px-4 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider capitalize"
                    >
                      {field.replace('_', ' ')}
                    </th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                {groupData.map((row, idx) => (
                  <tr key={idx} className="hover:bg-gray-50">
                    {Object.keys(row).map((field) => (
                      <td key={`${idx}-${field}`} className="px-4 py-3 text-gray-700">
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
        <div className="text-center py-12">
          <p className="text-gray-500">
            {filters.length > 0 ? 'No data matches the applied filters' : 'No data available'}
          </p>
        </div>
      )}
    </div>
  )
}

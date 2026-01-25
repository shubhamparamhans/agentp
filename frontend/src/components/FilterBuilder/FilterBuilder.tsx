import { useState } from 'react'

interface FilterBuilderProps {
  fields: string[]
  onAddFilter: (filter: { field: string; operator: string; value: string }) => void
}

export function FilterBuilder({ fields, onAddFilter }: FilterBuilderProps) {
  const [field, setField] = useState<string>(fields[0] || '')
  const [operator, setOperator] = useState<string>('equals')
  const [value, setValue] = useState<string>('')

  const operators = [
    { value: 'equals', label: 'Equals' },
    { value: 'contains', label: 'Contains' },
    { value: 'startswith', label: 'Starts With' },
    { value: 'endswith', label: 'Ends With' },
    { value: 'gt', label: 'Greater Than' },
    { value: 'lt', label: 'Less Than' },
    { value: 'gte', label: 'Greater or Equal' },
    { value: 'lte', label: 'Less or Equal' },
  ]

  const handleAddFilter = () => {
    if (field && operator && value) {
      onAddFilter({ field, operator, value })
      setValue('')
      setField(fields[0] || '')
      setOperator('equals')
    }
  }

  return (
    <div className="space-y-3">
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Field</label>
        <select
          value={field}
          onChange={(e) => setField(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 text-sm"
        >
          {fields.map((f) => (
            <option key={f} value={f}>
              {f.replace('_', ' ')}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Operator</label>
        <select
          value={operator}
          onChange={(e) => setOperator(e.target.value)}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 text-sm"
        >
          {operators.map((op) => (
            <option key={op.value} value={op.value}>
              {op.label}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Value</label>
        <input
          type="text"
          placeholder="Enter filter value"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyPress={(e) => {
            if (e.key === 'Enter') {
              handleAddFilter()
            }
          }}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-blue-500 focus:border-blue-500 text-sm"
        />
      </div>

      <button
        onClick={handleAddFilter}
        disabled={!field || !operator || !value}
        className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed text-sm"
      >
        Add Filter
      </button>
    </div>
  )
}

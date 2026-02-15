import { useState } from 'react'

interface ObjectRendererProps {
  value: any
  isMongoDb?: boolean
  level?: number
}

export function ObjectRenderer({ value, isMongoDb = false, level = 0 }: ObjectRendererProps) {
  const [expandedKeys, setExpandedKeys] = useState<Set<string>>(new Set())

  // Check if value is an object (but not null, date, or array of primitives)
  const isObject = (val: any): boolean => {
    if (val === null || val === undefined) return false
    if (val instanceof Date) return false
    if (typeof val !== 'object') return false
    if (Array.isArray(val)) {
      // Only treat array as object if it contains objects
      return val.some((item) => typeof item === 'object' && item !== null)
    }
    return true
  }

  const toggleExpand = (key: string) => {
    const newSet = new Set(expandedKeys)
    if (newSet.has(key)) {
      newSet.delete(key)
    } else {
      newSet.add(key)
    }
    setExpandedKeys(newSet)
  }

  // For non-MongoDB databases, don't render objects as expandable
  if (!isMongoDb && isObject(value)) {
    return <span className="text-gray-500 italic">[nested data]</span>
  }

  // For primitive values
  if (typeof value !== 'object' || value === null) {
    if (value instanceof Date) {
      return <span className="text-blue-300">{value.toISOString()}</span>
    }
    if (typeof value === 'boolean') {
      return (
        <span className={value ? 'text-green-400' : 'text-red-400'}>
          {value ? 'true' : 'false'}
        </span>
      )
    }
    if (typeof value === 'number') {
      return <span className="text-yellow-300">{value.toLocaleString()}</span>
    }
    return <span className="text-gray-300">{String(value)}</span>
  }

  // For arrays
  if (Array.isArray(value)) {
    if (value.length === 0) {
      return <span className="text-gray-500">[]</span>
    }

    const arrayKey = `array-${level}-${Math.random()}`
    const isExpanded = expandedKeys.has(arrayKey)

    return (
      <div className="inline-block">
        <button
          onClick={() => toggleExpand(arrayKey)}
          className="inline-flex items-center gap-1 cursor-pointer text-cyan-300 hover:text-cyan-100 transition-colors"
        >
          <span className={`text-sm transition-transform ${isExpanded ? 'rotate-90' : ''}`}>
            ▶
          </span>
          <span className="text-gray-400">Array({value.length})</span>
        </button>

        {isExpanded && (
          <div className="ml-4 mt-2 border-l border-gray-600 pl-3 space-y-2">
            {value.map((item, idx) => (
              <div key={idx} className="text-gray-300">
                <span className="text-gray-500">[{idx}]:</span>{' '}
                <ObjectRenderer value={item} isMongoDb={isMongoDb} level={level + 1} />
              </div>
            ))}
          </div>
        )}
      </div>
    )
  }

  // For objects
  const keys = Object.keys(value)
  if (keys.length === 0) {
    return <span className="text-gray-500">{'{}'}</span>
  }

  const objectKey = `obj-${level}-${Math.random()}`
  const isExpanded = expandedKeys.has(objectKey)

  return (
    <div className="inline-block">
      <button
        onClick={() => toggleExpand(objectKey)}
        className="inline-flex items-center gap-1 cursor-pointer text-cyan-300 hover:text-cyan-100 transition-colors"
      >
        <span className={`text-sm transition-transform ${isExpanded ? 'rotate-90' : ''}`}>
          ▶
        </span>
        <span className="text-gray-400">Object({keys.length})</span>
      </button>

      {isExpanded && (
        <div className="ml-4 mt-2 border-l border-gray-600 pl-3 space-y-2">
          {keys.map((key) => (
            <div key={key} className="text-gray-300 break-words">
              <span className="text-purple-300">{key}:</span>{' '}
              <ObjectRenderer value={value[key]} isMongoDb={isMongoDb} level={level + 1} />
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

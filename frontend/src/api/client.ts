// API client
const API_BASE = (window as any).REACT_APP_API_URL || 'http://localhost:8080'

export interface Model {
  name: string
  table: string
  primary_key: string
  fields: Array<{
    name: string
    type: string
  }>
}

export interface QueryResponse {
  sql: string
  params: any[]
  data?: any[]
  error?: string
}

export async function fetchModels(): Promise<Model[]> {
  const response = await fetch(`${API_BASE}/models`)
  if (!response.ok) {
    throw new Error(`Failed to fetch models: ${response.statusText}`)
  }
  return response.json()
}

export async function executeQuery(query: unknown): Promise<QueryResponse> {
  const response = await fetch(`${API_BASE}/query`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(query),
  })
  if (!response.ok) {
    throw new Error(`Failed to execute query: ${response.statusText}`)
  }
  return response.json()
}

export interface Filter {
  field: string
  op: string
  value: any
}

export function buildDSLQuery(
  modelName: string,
  fields?: string[],
  filters?: Filter[],
  groupByField?: string,
  limit: number = 100,
  offset: number = 0
): any {
  const query: any = {
    model: modelName,
    pagination: { limit, offset },
  }

  if (fields && fields.length > 0) {
    query.fields = fields
  }

  if (filters && filters.length > 0) {
    if (filters.length === 1) {
      query.filters = {
        field: filters[0].field,
        op: filters[0].op,
        value: filters[0].value,
      }
    } else {
      query.filters = {
        and: filters.map((f) => ({
          field: f.field,
          op: f.op,
          value: f.value,
        })),
      }
    }
  }

  if (groupByField) {
    query.group_by = [groupByField]
    query.aggregates = [
      { fn: 'count', field: '', alias: 'count' },
      { fn: 'count', field: 'id', alias: 'total_rows' },
    ]
  }

  return query
}

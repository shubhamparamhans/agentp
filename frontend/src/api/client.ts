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
  affected_rows?: number // For DELETE operations
  error?: string
  total?: number
  meta?: {
    total: number
    limit: number
    offset: number
  }
}

export interface InfoResponse {
  database_type?: 'mongo' | 'postgres'
  status?: string
}

export async function fetchDatabaseInfo(): Promise<InfoResponse> {
  try {
    const response = await fetch(`${API_BASE}/info`)
    if (!response.ok) {
      return { database_type: 'postgres' }
    }
    return response.json()
  } catch {
    return { database_type: 'postgres' }
  }
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

export interface Sort {
  field: string
  direction: 'asc' | 'desc'
}

/**
 * Builds a search filter (OR filter across multiple fields)
 */
export function buildSearchQuery(
  searchTerm: string,
  searchFields: string[],
  operator: string = 'contains'
): any | null {
  if (!searchTerm || searchTerm.trim() === '' || searchFields.length === 0) {
    return null
  }

  // Build OR filter across all search fields
  const orFilters = searchFields.map((field) => ({
    field,
    op: operator,
    value: searchTerm.trim(),
  }))

  return {
    or: orFilters,
  }
}

/**
 * Gets searchable fields from model (string/text types only)
 */
export function getSearchableFields(model: Model): string[] {
  return model.fields
    .filter((field) => {
      const type = field.type.toLowerCase()
      return (
        type === 'string' ||
        type === 'text' ||
        type === 'varchar' ||
        type === 'char' ||
        type.includes('text')
      )
    })
    .map((field) => field.name)
}

export function buildDSLQuery(
  modelName: string,
  fields?: string[],
  filters?: Filter[],
  groupByField?: string,
  limit: number = 100,
  offset: number = 0,
  sort?: Sort[],
  searchFilter?: any // OR filter from search
): any {
  const query: any = {
    model: modelName,
    pagination: { limit, offset },
  }

  if (fields && fields.length > 0) {
    query.fields = fields
  }

  // Combine search filter with existing filters
  const allFilters: any[] = []

  // Add search filter if exists (can be OR filter for global or single field filter for column)
  if (searchFilter) {
    allFilters.push(searchFilter)
  }

  // Add existing filters
  if (filters && filters.length > 0) {
    if (filters.length === 1) {
      allFilters.push({
        field: filters[0].field,
        op: filters[0].op,
        value: filters[0].value,
      })
    } else {
      allFilters.push({
        and: filters.map((f) => ({
          field: f.field,
          op: f.op,
          value: f.value,
        })),
      })
    }
  }

  // Set filters in query
  if (allFilters.length > 0) {
    if (allFilters.length === 1) {
      query.filters = allFilters[0]
    } else {
      query.filters = {
        and: allFilters,
      }
    }
  }

  if (sort && sort.length > 0) {
    query.sort = sort.map((s) => ({
      field: s.field,
      direction: s.direction,
    }))
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

/**
 * Creates a new record
 * @param modelName - The model name
 * @param data - The data to insert
 * @returns The created record
 */
export async function createRecord(
  modelName: string,
  data: Record<string, any>
): Promise<QueryResponse> {
  const query = {
    operation: 'create',
    model: modelName,
    data: data,
  }
  return executeQuery(query)
}

/**
 * Updates an existing record
 * @param modelName - The model name
 * @param id - The record ID
 * @param data - The data to update (partial)
 * @returns The updated record
 */
export async function updateRecord(
  modelName: string,
  id: string | number,
  data: Record<string, any>
): Promise<QueryResponse> {
  const query = {
    operation: 'update',
    model: modelName,
    id: id,
    data: data,
  }
  return executeQuery(query)
}

/**
 * Deletes a record
 * @param modelName - The model name
 * @param id - The record ID
 * @returns The number of affected rows
 */
export async function deleteRecord(
  modelName: string,
  id: string | number
): Promise<QueryResponse> {
  const query = {
    operation: 'delete',
    model: modelName,
    id: id,
  }
  return executeQuery(query)
}

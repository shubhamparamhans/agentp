// Application state types
export interface Model {
  name: string
  table: string
  primaryKey: string
  fields: Field[]
}

export interface Field {
  name: string
  type: string
  nullable: boolean
}

export interface AppState {
  models: Model[]
  selectedModel: Model | null
  query: QueryState
  databaseType?: 'mongo' | 'postgres'
}

export interface QueryState {
  filters: unknown[]
  groupBy: string[]
  sort: unknown[]
}

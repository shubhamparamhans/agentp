import { createContext, useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import type { AppState } from '../types'
import { fetchDatabaseInfo } from '../api/client'

const initialState: AppState = {
  models: [],
  selectedModel: null,
  query: {
    filters: [],
    groupBy: [],
    sort: [],
  },
  databaseType: 'postgres',
}

export const AppContext = createContext<{
  state: AppState
  setState: (state: AppState) => void
}>({ state: initialState, setState: () => {} })

export function AppProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AppState>(initialState)

  // Fetch database type on mount
  useEffect(() => {
    const initializeDatabase = async () => {
      try {
        const info = await fetchDatabaseInfo()
        setState((prev) => ({
          ...prev,
          databaseType: info.database_type || 'postgres',
        }))
      } catch {
        // Keep default if fetch fails
        setState((prev) => ({
          ...prev,
          databaseType: 'postgres',
        }))
      }
    }

    initializeDatabase()
  }, [])

  return (
    <AppContext.Provider value={{ state, setState }}>
      {children}
    </AppContext.Provider>
  )
}

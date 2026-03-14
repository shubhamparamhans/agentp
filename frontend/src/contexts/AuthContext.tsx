// Authentication Context
import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

interface AuthContextType {
  isAuthenticated: boolean
  user: string | null
  login: (username: string, password: string) => Promise<boolean>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    // Check localStorage on mount
    return localStorage.getItem('isAuthenticated') === 'true'
  })
  const [user, setUser] = useState<string | null>(() => {
    return localStorage.getItem('user') || null
  })

  const login = async (username: string, password: string): Promise<boolean> => {
    // TODO: Replace with actual API call
    // For now, accept any credentials (demo mode)
    if (username && password) {
      setIsAuthenticated(true)
      setUser(username)
      localStorage.setItem('isAuthenticated', 'true')
      localStorage.setItem('user', username)
      return true
    }
    return false
  }

  const logout = () => {
    setIsAuthenticated(false)
    setUser(null)
    localStorage.removeItem('isAuthenticated')
    localStorage.removeItem('user')
  }

  // Optional: Check authentication status on mount
  useEffect(() => {
    const stored = localStorage.getItem('isAuthenticated')
    if (stored === 'true') {
      setIsAuthenticated(true)
      setUser(localStorage.getItem('user'))
    }
  }, [])

  return (
    <AuthContext.Provider value={{ isAuthenticated, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}


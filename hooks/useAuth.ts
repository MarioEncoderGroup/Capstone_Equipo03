// MisViáticos - useAuth Hook

'use client'

import { useRouter } from 'next/navigation'
import { useCallback, useEffect, useState } from 'react'
import { AuthService } from '@/services/authService'
import { TenantService } from '@/services/tenantService'
import type { LoginRequest, RegisterRequest } from '@/types/api'
import { getErrorMessage } from '@/lib/api/errorHandler'

interface UseAuthReturn {
  isAuthenticated: boolean
  isLoading: boolean
  error: string | null
  login: (data: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<string> // Retorna email_token
  logout: () => void
  clearError: () => void
}

export function useAuth(): UseAuthReturn {
  const router = useRouter()
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = useCallback(() => {
    const authenticated = AuthService.isAuthenticated()
    setIsAuthenticated(authenticated)
    setIsLoading(false)
  }, [])

  const login = useCallback(
    async (data: LoginRequest) => {
      setIsLoading(true)
      setError(null)

      try {
        await AuthService.login(data)

        // Verificar estado del tenant
        const tenantStatus = await TenantService.getStatus()

        if (!tenantStatus.has_tenants) {
          // No tiene tenant - redirigir a creación
          router.push('/tenant/create')
        } else {
          // Tiene uno o más tenants - redirigir a selección
          // El usuario debe seleccionar manualmente en la página
          router.push('/tenant/select')
        }

        setIsAuthenticated(true)
      } catch (err) {
        setError(getErrorMessage(err))
        throw err
      } finally {
        setIsLoading(false)
      }
    },
    [router]
  )

  const register = useCallback(async (data: RegisterRequest) => {
    setIsLoading(true)
    setError(null)

    try {
      const result = await AuthService.register(data)
      return result.email_token
    } catch (err) {
      setError(getErrorMessage(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const logout = useCallback(() => {
    AuthService.logout()
    setIsAuthenticated(false)
    router.push('/auth/login')
  }, [router])

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  return {
    isAuthenticated,
    isLoading,
    error,
    login,
    register,
    logout,
    clearError,
  }
}

// MisViáticos - useTenant Hook

'use client'

import { useCallback, useState } from 'react'
import { useRouter } from 'next/navigation'
import { TenantService } from '@/services/tenantService'
import type { CreateTenantRequest, TenantData } from '@/types/api'
import { getErrorMessage } from '@/lib/api/errorHandler'

interface UseTenantReturn {
  isLoading: boolean
  error: string | null
  createTenant: (data: CreateTenantRequest) => Promise<void>
  selectTenant: (tenantId: string) => Promise<void>
  getTenantStatus: () => Promise<{
    has_tenants: boolean
    tenants: TenantData[]
    tenant_count: number
  }>
  clearError: () => void
}

export function useTenant(): UseTenantReturn {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const createTenant = useCallback(
    async (data: CreateTenantRequest) => {
      setIsLoading(true)
      setError(null)

      try {
        const tenant = await TenantService.create(data)

        // Automáticamente seleccionar el tenant recién creado
        await TenantService.select(tenant.id)

        // Esperar un momento para que las cookies se actualicen
        await new Promise(resolve => setTimeout(resolve, 100))

        // Refrescar el router para que el middleware lea las cookies actualizadas
        router.refresh()

        // Redirigir al dashboard
        router.push('/dashboard')
      } catch (err) {
        setError(getErrorMessage(err))
        throw err
      } finally {
        setIsLoading(false)
      }
    },
    [router]
  )

  const selectTenant = useCallback(
    async (tenantId: string) => {
      setIsLoading(true)
      setError(null)

      try {
        await TenantService.select(tenantId)

        // Esperar un momento para que las cookies se actualicen
        await new Promise(resolve => setTimeout(resolve, 100))

        // Refrescar el router para que el middleware lea las cookies actualizadas
        router.refresh()

        // Redirigir al dashboard
        router.push('/dashboard')
      } catch (err) {
        setError(getErrorMessage(err))
        throw err
      } finally {
        setIsLoading(false)
      }
    },
    [router]
  )

  const getTenantStatus = useCallback(async () => {
    setIsLoading(true)
    setError(null)

    try {
      const status = await TenantService.getStatus()
      return {
        has_tenants: status.has_tenants,
        tenants: status.tenants,
        tenant_count: status.tenant_count,
      }
    } catch (err) {
      setError(getErrorMessage(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  return {
    isLoading,
    error,
    createTenant,
    selectTenant,
    getTenantStatus,
    clearError,
  }
}

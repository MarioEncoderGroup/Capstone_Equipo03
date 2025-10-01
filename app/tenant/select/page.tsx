'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { BuildingOfficeIcon, CheckIcon } from '@heroicons/react/24/outline'
import { useTenant } from '@/hooks/useTenant'
import { Button } from '@/components/ui/Button'
import { Alert } from '@/components/ui/Alert'
import type { TenantData } from '@/types/api'

export default function SelectTenantPage() {
  const router = useRouter()
  const { selectTenant, getTenantStatus, isLoading, error, clearError } =
    useTenant()

  const [tenants, setTenants] = useState<TenantData[]>([])
  const [selectedTenantId, setSelectedTenantId] = useState<string | null>(null)
  const [loadingTenants, setLoadingTenants] = useState(true)

  useEffect(() => {
    loadTenants()
  }, [])

  const loadTenants = async () => {
    try {
      const status = await getTenantStatus()

      if (!status.has_tenants) {
        // No tiene tenants, redirigir a creación
        router.push('/tenant/create')
        return
      }

      // Mostrar TODOS los tenants visualmente, sin auto-seleccionar
      // El usuario debe hacer clic manualmente en "Seleccionar Empresa"
      setTenants(status.tenants)
    } catch (err) {
      console.error('Error cargando tenants:', err)
    } finally {
      setLoadingTenants(false)
    }
  }

  const handleSelectTenant = async (tenantId: string) => {
    setSelectedTenantId(tenantId)

    try {
      await selectTenant(tenantId)
      // useTenant ya redirige al dashboard después del éxito
    } catch (err) {
      // El error ya está manejado por useTenant
      console.error('Error seleccionando tenant:', err)
      setSelectedTenantId(null)
    }
  }

  if (loadingTenants) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mb-4" />
          <p className="text-gray-600">Cargando empresas...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Selecciona tu Empresa
          </h1>
          <p className="text-gray-600">
            Tienes acceso a múltiples empresas. Selecciona una para continuar.
          </p>
        </div>

        {error && (
          <div className="mb-6">
            <Alert variant="error" onClose={clearError}>
              {error}
            </Alert>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {tenants.map((tenant) => (
            <div
              key={tenant.id}
              className={`
                relative bg-white p-6 rounded-lg shadow-md border-2 transition-all duration-200
                ${
                  selectedTenantId === tenant.id
                    ? 'border-purple-600 ring-2 ring-purple-200'
                    : 'border-gray-200'
                }
              `}
            >
              <div className="flex items-start space-x-4 mb-4">
                <div className="flex-shrink-0 w-12 h-12 rounded-lg bg-purple-100 flex items-center justify-center">
                  <BuildingOfficeIcon className="w-6 h-6 text-purple-600" />
                </div>

                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-gray-900 mb-1">
                    {tenant.business_name}
                  </h3>
                  <p className="text-sm text-gray-600 mb-2">
                    RUT: {tenant.rut}
                  </p>
                  <p className="text-sm text-gray-600 mb-2">
                    Email: {tenant.email}
                  </p>
                  <p className="text-sm text-gray-600 mb-2">
                    Teléfono: {tenant.phone}
                  </p>
                  <span
                    className={`
                    inline-block px-2 py-1 rounded-full text-xs
                    ${
                      tenant.status === 'active'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-gray-100 text-gray-800'
                    }
                  `}
                  >
                    {tenant.status === 'active' ? 'Activo' : 'Inactivo'}
                  </span>
                </div>
              </div>

              <Button
                variant="primary"
                fullWidth
                isLoading={selectedTenantId === tenant.id && isLoading}
                disabled={isLoading && selectedTenantId !== tenant.id}
                onClick={() => handleSelectTenant(tenant.id)}
              >
                {selectedTenantId === tenant.id && isLoading
                  ? 'Seleccionando...'
                  : 'Seleccionar Empresa'}
              </Button>

              {selectedTenantId === tenant.id && isLoading && (
                <div className="absolute inset-0 bg-white/50 rounded-lg flex items-center justify-center">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600" />
                </div>
              )}
            </div>
          ))}
        </div>

        <div className="mt-8 text-center space-y-4">
          <Link href="/tenant/create">
            <Button variant="secondary">Crear Nueva Empresa</Button>
          </Link>

          <div>
            <Link
              href="/auth/login"
              className="text-sm text-purple-600 hover:text-purple-700"
            >
              ← Volver al login
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

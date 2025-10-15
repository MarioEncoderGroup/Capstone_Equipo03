'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
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
    <div className="min-h-screen relative overflow-hidden bg-gradient-to-br from-purple-50 via-white to-indigo-50">
      {/* Background decorations */}
      <div className="absolute inset-0 overflow-hidden">
        {/* Large gradient orbs */}
        <div className="absolute -top-24 -right-24 w-96 h-96 rounded-full bg-gradient-to-br from-purple-400/20 to-indigo-400/20 blur-3xl animate-pulse" />
        <div className="absolute top-1/4 -left-32 w-80 h-80 rounded-full bg-gradient-to-br from-indigo-400/15 to-purple-400/15 blur-3xl" />
        <div className="absolute bottom-1/4 right-1/4 w-64 h-64 rounded-full bg-gradient-to-br from-purple-300/10 to-pink-300/10 blur-2xl animate-pulse" />
        
      </div>

      <div className="relative z-10 pb-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            {/* Logo corporativo */}
            <div className="flex justify-center mb-4 -mt-4">
              <Image
                src="/icon-mv/Assets MV_Logo2.svg"
                alt="MisViáticos"
                width={240}
                height={240}
                className="w-60 h-60"
              />
            </div>
            
            <h1 className="text-4xl md:text-5xl font-black text-purple-600 mb-4 tracking-tight">
              Selecciona tu Empresa
            </h1>
            <p className="text-lg text-gray-600 font-medium max-w-2xl mx-auto leading-relaxed">
              Tienes acceso a múltiples empresas. Selecciona una para continuar tu gestión de viáticos.
            </p>
            
            {/* Decorative line */}
            <div className="flex justify-center mt-6">
              <div className="w-32 h-1 bg-purple-500 rounded-full" />
            </div>
          </div>

        {error && (
          <div className="mb-6">
            <Alert variant="error" onClose={clearError}>
              {error}
            </Alert>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {tenants.map((tenant) => (
            <div
              key={tenant.id}
              className={`
                relative bg-white rounded-2xl shadow-lg border transition-all duration-300 hover:shadow-xl transform hover:-translate-y-1
                ${
                  selectedTenantId === tenant.id
                    ? 'border-purple-500 ring-4 ring-purple-200/50 shadow-purple-200/50'
                    : 'border-gray-200 hover:border-purple-300'
                }
              `}
            >
              {/* Header con color sólido */}
              <div className="relative h-20 bg-purple-600 rounded-t-2xl overflow-hidden">
                <div className="absolute inset-0 bg-black/10" />
                <div className="absolute top-4 left-6 flex items-center space-x-3">
                  {/* Logo MisViáticos */}
                  <div className="w-12 h-12 bg-white rounded-xl flex items-center justify-center shadow-lg">
                    <Image
                      src="/icon-mv/Assets MV_Isotipo2.svg"
                      alt="MisViáticos"
                      width={32}
                      height={32}
                      className="w-8 h-8"
                    />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg leading-tight">
                      {tenant.business_name}
                    </h3>
                    <p className="text-purple-100 text-sm font-medium">
                      RUT: {tenant.rut}
                    </p>
                  </div>
                </div>
              </div>

              {/* Contenido */}
              <div className="p-6">
                <div className="space-y-3 mb-6">
                  <div className="flex items-center space-x-3">
                    <div className="w-2 h-2 bg-purple-500 rounded-full" />
                    <span className="text-gray-600 text-sm">
                      <span className="font-medium text-gray-800">Email:</span> {tenant.email}
                    </span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <div className="w-2 h-2 bg-purple-500 rounded-full" />
                    <span className="text-gray-600 text-sm">
                      <span className="font-medium text-gray-800">Teléfono:</span> {tenant.phone}
                    </span>
                  </div>
                </div>

                <Button
                  variant="primary"
                  fullWidth
                  isLoading={selectedTenantId === tenant.id && isLoading}
                  disabled={isLoading && selectedTenantId !== tenant.id}
                  onClick={() => handleSelectTenant(tenant.id)}
                  className="bg-purple-600 hover:bg-purple-700 shadow-lg hover:shadow-xl transition-all duration-200 py-3 rounded-xl font-semibold"
                >
                  {selectedTenantId === tenant.id && isLoading
                    ? 'Seleccionando...'
                    : 'Seleccionar Empresa'}
                </Button>
              </div>

              {selectedTenantId === tenant.id && isLoading && (
                <div className="absolute inset-0 bg-white/80 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                  <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-4 border-purple-200 border-t-purple-600 mx-auto mb-3" />
                    <p className="text-purple-600 font-medium">Conectando...</p>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

          <div className="mt-12 text-center space-y-6">
            <Link href="/tenant/create">
              <Button 
                variant="secondary"
                className="bg-white/80 backdrop-blur-sm border-2 border-purple-200 hover:border-purple-300 text-purple-700 font-semibold px-8 py-3 rounded-xl shadow-lg hover:shadow-xl transition-all duration-200"
              >
                + Crear Nueva Empresa
              </Button>
            </Link>

            <div>
              <Link
                href="/auth/login"
                className="inline-flex items-center text-purple-600 hover:text-purple-700 font-medium text-sm transition-colors duration-200"
              >
                ← Volver al login
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

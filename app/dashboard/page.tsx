'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { TokenManager } from '@/lib/auth/tokenManager'
import { BuildingOfficeIcon, UserIcon } from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'

interface User {
  id: string
  email: string
  username: string
  full_name: string
  phone?: string
  is_active: boolean
  last_login?: string
}

export default function DashboardPage() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [tenantId, setTenantId] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    checkAuthAndTenant()
  }, [])

  const checkAuthAndTenant = async () => {
    try {
      // 1. Verificar autenticación
      if (!TokenManager.isAuthenticated()) {
        router.push('/auth/login')
        return
      }

      // 2. Verificar que el token tenga tenant_id
      const currentTenantId = TokenManager.getTenantId()
      if (!currentTenantId) {
        // No tiene tenant_id en el token - redirigir a selección
        router.push('/tenant/select')
        return
      }

      // 3. Obtener datos del usuario
      const userData = TokenManager.getUserData()
      setUser(userData as User)
      setTenantId(currentTenantId)
    } catch (error) {
      console.error('Error verificando autenticación:', error)
      router.push('/auth/login')
    } finally {
      setIsLoading(false)
    }
  }

  const handleLogout = () => {
    TokenManager.clearTokens()
    router.push('/auth/login')
  }

  const handleChangeTenant = () => {
    // Redirigir a selección de tenant
    router.push('/tenant/select')
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mb-4" />
          <p className="text-gray-600">Cargando dashboard...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div className="flex items-center">
              <img
                src="/icon-mv/Assets MV_Logo2.svg"
                alt="MisViáticos"
                className="h-12 w-auto"
              />
            </div>

            <div className="flex items-center space-x-4">
              <span className="text-gray-700">
                Hola, {user?.full_name || user?.username}
              </span>

              <Button variant="secondary" onClick={handleChangeTenant}>
                <BuildingOfficeIcon className="w-5 h-5 mr-2 inline" />
                Cambiar Empresa
              </Button>

              <Button variant="danger" onClick={handleLogout}>
                Cerrar Sesión
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* User Info Card */}
        <div className="bg-white shadow rounded-lg p-6 mb-6">
          <div className="flex items-center space-x-4">
            <div className="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center">
              <UserIcon className="w-8 h-8 text-purple-600" />
            </div>
            <div>
              <h2 className="text-xl font-semibold text-gray-900">
                {user?.full_name || user?.username || 'Usuario'}
              </h2>
              <p className="text-gray-600">{user?.email}</p>
            </div>
          </div>
        </div>

        {/* Tenant Info */}
        <div className="bg-white shadow rounded-lg p-6 mb-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Información de Empresa
          </h3>
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Tenant ID:</span>
              <span className="font-mono text-sm bg-gray-100 px-3 py-1 rounded">
                {tenantId}
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600">Estado:</span>
              <span className="px-3 py-1 bg-green-100 text-green-800 rounded-full text-sm font-medium">
                Activo ✓
              </span>
            </div>
          </div>
        </div>

        {/* Welcome Message */}
        <div className="bg-gradient-to-r from-purple-600 to-violet-600 rounded-lg p-8 text-white">
          <h3 className="text-2xl font-bold mb-2">
            ¡Bienvenido a MisViáticos! 🎉
          </h3>
          <p className="text-purple-100 mb-4">
            Tu cuenta está completamente configurada y lista para usar. Ahora
            puedes gestionar los viáticos de tu empresa de forma simple y
            eficiente.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
            <div className="bg-white/10 rounded-lg p-4">
              <h4 className="font-semibold mb-1">Paso 1</h4>
              <p className="text-sm text-purple-100">
                Configura las categorías de gastos
              </p>
            </div>
            <div className="bg-white/10 rounded-lg p-4">
              <h4 className="font-semibold mb-1">Paso 2</h4>
              <p className="text-sm text-purple-100">Invita a tu equipo</p>
            </div>
            <div className="bg-white/10 rounded-lg p-4">
              <h4 className="font-semibold mb-1">Paso 3</h4>
              <p className="text-sm text-purple-100">
                Comienza a registrar gastos
              </p>
            </div>
          </div>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6">
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow">
            <h3 className="text-lg font-medium text-gray-900 mb-2">Gastos</h3>
            <p className="text-gray-600">Gestiona tus gastos de viáticos</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow">
            <h3 className="text-lg font-medium text-gray-900 mb-2">Reportes</h3>
            <p className="text-gray-600">Genera reportes detallados</p>
          </div>
          <div className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow">
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              Configuración
            </h3>
            <p className="text-gray-600">Ajusta tu perfil y preferencias</p>
          </div>
        </div>

        {/* Debug Info (solo en desarrollo) */}
        {process.env.NODE_ENV === 'development' && (
          <div className="mt-6 bg-gray-100 rounded-lg p-4">
            <h3 className="font-semibold text-gray-900 mb-2">
              🔧 Debug Info (solo desarrollo)
            </h3>
            <pre className="text-xs text-gray-600 overflow-auto">
              {JSON.stringify(
                {
                  user,
                  tenantId,
                  tokenHasTenant: TokenManager.hasTenantId(),
                  isAuthenticated: TokenManager.isAuthenticated(),
                },
                null,
                2
              )}
            </pre>
          </div>
        )}
      </main>
    </div>
  )
}

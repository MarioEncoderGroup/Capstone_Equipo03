'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { TokenManager } from '@/lib/auth/tokenManager'
import { TenantService } from '@/services/tenantService'
import { BuildingOfficeIcon, UserIcon } from '@heroicons/react/24/outline'
import { Button } from '@/components/ui/Button'
import type { TenantData } from '@/types/api'

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
  const [tenantData, setTenantData] = useState<TenantData | null>(null)
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

      // 4. Obtener datos del tenant actual
      try {
        const tenantStatus = await TenantService.getStatus()
        const currentTenant = tenantStatus.tenants.find(t => t.id === currentTenantId)
        if (currentTenant) {
          setTenantData(currentTenant)
        }
      } catch (tenantError) {
        console.error('Error obteniendo datos del tenant:', tenantError)
        // No es crítico, continuar sin datos del tenant
      }
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

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Welcome Message */}
        <div className="bg-purple-600 rounded-lg p-8 text-white mb-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold mb-2">
                ¡Bienvenido, {user?.full_name || user?.username}! 👋
              </h1>
              <p className="text-purple-100 text-lg">
                Estás conectado como <span className="font-medium">{user?.email}</span>
              </p>
            </div>
            <div className="w-20 h-20 bg-white/20 rounded-full flex items-center justify-center">
              <UserIcon className="w-10 h-10 text-white" />
            </div>
          </div>
        </div>

        {/* Company Information Card */}
        <div className="bg-white shadow-lg rounded-xl p-8 mb-8 border border-gray-100">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 bg-white rounded-lg flex items-center justify-center p-2 border border-purple-100">
                <img
                  src="/icon-mv/Assets MV_Isotipo2.svg"
                  alt="MisViáticos Isotipo"
                  className="w-full h-full object-contain"
                />
              </div>
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Información de Empresa</h2>
                <p className="text-gray-600">Detalles de tu organización actual</p>
              </div>
            </div>
            <Button 
              variant="secondary" 
              onClick={handleChangeTenant}
              className="flex items-center space-x-2 px-6 py-3"
            >
              <BuildingOfficeIcon className="w-5 h-5" />
              <span>Cambiar Empresa</span>
            </Button>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="space-y-2">
                <span className="text-gray-600 text-sm font-medium">Nombre de la Empresa</span>
                <p className="text-gray-900 font-semibold text-lg">
                  {tenantData?.business_name || 'Cargando...'}
                </p>
              </div>
            </div>
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="space-y-2">
                <span className="text-gray-600 text-sm font-medium">RUT</span>
                <p className="text-gray-900 font-semibold text-lg">
                  {tenantData?.rut || 'Cargando...'}
                </p>
              </div>
            </div>
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="space-y-2">
                <span className="text-gray-600 text-sm font-medium">Dirección</span>
                <p className="text-gray-900 font-medium">
                  {tenantData?.address || 'Cargando...'}
                </p>
              </div>
            </div>
            <div className="bg-gray-50 rounded-lg p-4">
              <div className="space-y-2">
                <span className="text-gray-600 text-sm font-medium">Email</span>
                <p className="text-gray-900 font-medium">
                  {tenantData?.email || 'Cargando...'}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Expense Process Guide */}
        <div className="bg-white shadow-lg rounded-xl p-8 mb-8 border border-gray-100">
          <div className="flex items-center mb-6">
            <div className="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center p-2">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div className="ml-3">
              <h2 className="text-2xl font-bold text-gray-900">Cómo Rendir un Gasto</h2>
              <p className="text-gray-600">Sigue estos pasos para rendir tus gastos de manera eficiente</p>
            </div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Paso 1 */}
            <div className="relative bg-purple-50 rounded-xl p-6 border border-purple-100">
              <div className="flex items-center justify-between mb-4">
                <div className="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center">
                  <span className="text-white font-bold text-sm">1</span>
                </div>
                <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
                  <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </div>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Captura el Comprobante</h3>
              <p className="text-gray-600 text-sm">Toma una foto clara del boleta, factura o recibo inmediatamente después de realizar la compra.</p>
            </div>

            {/* Paso 2 */}
            <div className="relative bg-purple-50 rounded-xl p-6 border border-purple-100">
              <div className="flex items-center justify-between mb-4">
                <div className="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center">
                  <span className="text-white font-bold text-sm">2</span>
                </div>
                <div className="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
                  <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </div>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Completa los Datos</h3>
              <p className="text-gray-600 text-sm">Ingresa la información del gasto: monto, fecha, categoría y descripción detallada del motivo.</p>
            </div>

            {/* Paso 3 */}
            <div className="relative bg-violet-50 rounded-xl p-6 border border-violet-100">
              <div className="flex items-center justify-between mb-4">
                <div className="w-8 h-8 bg-violet-600 rounded-full flex items-center justify-center">
                  <span className="text-white font-bold text-sm">3</span>
                </div>
                <div className="w-10 h-10 bg-violet-100 rounded-lg flex items-center justify-center">
                  <svg className="w-6 h-6 text-violet-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
                  </svg>
                </div>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Envía y Rastrea</h3>
              <p className="text-gray-600 text-sm">Envía tu solicitud de reembolso y monitorea su estado hasta la aprobación y pago.</p>
            </div>
          </div>

          {/* Tips adicionales */}
          <div className="mt-6 bg-yellow-50 rounded-lg p-4 border border-yellow-100">
            <div className="flex items-start">
              <div className="w-6 h-6 bg-yellow-100 rounded-full flex items-center justify-center mt-0.5">
                <svg className="w-4 h-4 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div className="ml-3">
                <h4 className="text-sm font-medium text-yellow-800">💡 Consejo Importante</h4>
                <p className="text-sm text-yellow-700 mt-1">
                  Asegúrate de que todos los comprobantes estén legibles y completos. Los gastos sin documentación adecuada pueden ser rechazados.
                </p>
              </div>
            </div>
          </div>
        </div>

      </main>
    </div>
  )
}

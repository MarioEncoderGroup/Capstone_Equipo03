'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'

interface User {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  full_name: string;
  phone: string;
  email_verified: boolean;
  is_active: boolean;
  last_login: string;
}

export default function DashboardPage() {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('auth_token')
    const userData = localStorage.getItem('user_data')

    if (!token) {
      router.push('/auth/login')
      return
    }

    if (userData) {
      setUser(JSON.parse(userData))
    }
    
    setLoading(false)
  }, [router])

  const handleLogout = () => {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_data')
    router.push('/')
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Cargando...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <img 
                src="/icon-mv/Assets MV_Elemento1.svg" 
                alt="MisViáticos" 
                className="h-8 w-auto"
              />
              <span className="ml-2 text-xl font-bold text-gray-900">MisViáticos</span>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">Hola, {user?.full_name || user?.username}</span>
              <button
                onClick={handleLogout}
                className="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
              >
                Cerrar Sesión
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="border-4 border-dashed border-gray-200 rounded-lg p-8">
            <div className="text-center">
              <h1 className="text-3xl font-bold text-gray-900 mb-4">
                ¡Bienvenido a MisViáticos!
              </h1>
              <p className="text-lg text-gray-600 mb-8">
                Panel de gestión de viáticos empresariales
              </p>

              {/* User Info Card */}
              <div className="bg-white rounded-lg shadow p-6 max-w-md mx-auto mb-8">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Información de Usuario</h2>
                <div className="space-y-3 text-left">
                  <div>
                    <span className="font-medium text-gray-700">Nombre:</span>
                    <span className="ml-2 text-gray-600">{user?.full_name}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Email:</span>
                    <span className="ml-2 text-gray-600">{user?.email}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Usuario:</span>
                    <span className="ml-2 text-gray-600">{user?.username}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Teléfono:</span>
                    <span className="ml-2 text-gray-600">{user?.phone || 'No especificado'}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-gray-700">Email verificado:</span>
                    <span className={`ml-2 px-2 py-1 rounded-full text-xs font-medium ${
                      user?.email_verified 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {user?.email_verified ? 'Sí' : 'No'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-gray-700">Estado:</span>
                    <span className={`ml-2 px-2 py-1 rounded-full text-xs font-medium ${
                      user?.is_active 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {user?.is_active ? 'Activo' : 'Inactivo'}
                    </span>
                  </div>
                </div>
              </div>

              {/* Features Coming Soon */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-white rounded-lg shadow p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-2">Gastos</h3>
                  <p className="text-gray-600">Gestiona tus gastos de viáticos</p>
                </div>
                <div className="bg-white rounded-lg shadow p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-2">Reportes</h3>
                  <p className="text-gray-600">Genera reportes detallados</p>
                </div>
                <div className="bg-white rounded-lg shadow p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-2">Configuración</h3>
                  <p className="text-gray-600">Ajusta tu perfil y preferencias</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

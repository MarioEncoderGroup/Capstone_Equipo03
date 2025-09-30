'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import RegisterHeader from './components/RegisterHeader'
import RegisterForm from './components/RegisterForm'
import { RegisterService } from './utils/api'
import type { RegisterFormData } from './types'

export default function RegisterPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  const handleRegister = async (data: RegisterFormData) => {
    setIsLoading(true)
    setError(null)
    
    try {
      const result = await RegisterService.register(data)
      
      if (result.success && result.token) {
        // Save token and redirect
        localStorage.setItem('auth_token', result.token)
        if (result.user) {
          localStorage.setItem('user_data', JSON.stringify(result.user))
        }
        router.push('/dashboard') // or wherever new users should go
      } else {
        setError(result.error || 'Error al crear la cuenta')
      }
    } catch (err) {
      setError('Ha ocurrido un error inesperado')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <RegisterHeader />
        
        <div className="bg-white py-8 px-6 shadow rounded-lg">
          {error && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}
          
          <RegisterForm onSubmit={handleRegister} isLoading={isLoading} />
        </div>

        <div className="text-center space-y-2">
          <p className="text-sm text-gray-600">
            ¿Ya tienes una cuenta?{' '}
            <Link href="/auth/login" className="font-medium text-purple-600 hover:text-purple-500">
              Inicia sesión aquí
            </Link>
          </p>
          <Link href="/" className="block text-purple-600 hover:text-purple-500 text-sm font-medium">
            ← Volver al inicio
          </Link>
        </div>
      </div>
    </div>
  )
}

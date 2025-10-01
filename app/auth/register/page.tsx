'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import RegisterHeader from './components/RegisterHeader'
import RegisterForm from './components/RegisterForm'
import { useAuth } from '@/hooks/useAuth'
import { Alert } from '@/components/ui/Alert'
import type { RegisterFormData } from './types'

export default function RegisterPage() {
  const router = useRouter()
  const { register, isLoading, error, clearError } = useAuth()
  const [emailToken, setEmailToken] = useState<string | null>(null)

  const handleRegister = async (data: RegisterFormData) => {
    try {
      // El formulario ya envía full_name directamente
      const token = await register(data)

      // Guardar el email_token y redirigir a verificación
      setEmailToken(token)
      router.push(`/auth/verify-email?token=${token}`)
    } catch (err) {
      // El error ya está manejado por useAuth
      console.error('Error en registro:', err)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <RegisterHeader />
        
        <div className="bg-white py-8 px-6 shadow rounded-lg">
          {error && (
            <div className="mb-4">
              <Alert variant="error" onClose={clearError}>
                {error}
              </Alert>
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

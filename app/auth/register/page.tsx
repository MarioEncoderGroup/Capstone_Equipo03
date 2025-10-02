'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import RegisterHeader from './components/RegisterHeader'
import RegisterForm from './components/RegisterForm'
import RegisterBenefitsSection from './components/RegisterBenefitsSection'
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
    <div className="min-h-screen bg-white flex">
      {/* Left side - Register Form */}
      <div className="flex-1 lg:flex-none lg:w-1/2 flex items-center justify-center px-4 sm:px-6 lg:px-20 xl:px-24">
        <div className="w-full max-w-sm lg:w-96">
          <div className="mb-8">
            <Link href="/" className="inline-block">
              <div className="flex items-center">
                <img 
                  src="/icon-mv/Assets MV_Elemento1.svg" 
                  alt="MisViáticos" 
                  className="h-10 w-auto"
                />
                <span className="ml-2 text-xl font-bold text-gray-900">MisViáticos</span>
              </div>
            </Link>
          </div>

          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">
              Inicia sesión o regístrate para gestionar
            </h2>
            <p className="text-gray-600 text-sm mb-8">
              ¿Tienes cuenta en MisViáticos? Usa el mismo correo y contraseña para que te podamos reconocer o{' '}
              <Link href="/auth/login" className="text-purple-600 hover:text-purple-500 font-medium">
                inicia sesión
              </Link>
              .
            </p>
          </div>

          {error && (
            <div className="mb-4">
              <Alert variant="error" onClose={clearError}>
                {error}
              </Alert>
            </div>
          )}

          <RegisterForm onSubmit={handleRegister} isLoading={isLoading} />

          <div className="mt-8 text-center space-y-2">
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

      {/* Right side - Benefits */}
      <RegisterBenefitsSection />
    </div>
  )
}

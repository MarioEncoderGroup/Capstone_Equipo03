'use client'

import { Suspense, useEffect } from 'react'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import LoginHeader from './components/LoginHeader'
import LoginForm from './components/LoginForm'
import BenefitsSection from './components/BenefitsSection'
import { useAuth } from '@/hooks/useAuth'
import { Alert } from '@/components/ui/Alert'
import type { LoginFormData } from './types'

function LoginContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { login, isLoading, error, clearError } = useAuth()

  const verified = searchParams.get('verified')

  useEffect(() => {
    // Limpiar error cuando se monta el componente
    clearError()
  }, [clearError])

  const handleLogin = async (data: LoginFormData) => {
    try {
      // useAuth ya maneja todo el flujo: login → tenant status → redirect
      await login(data)
    } catch (err) {
      // El error ya está manejado por useAuth
      console.error('Error en login:', err)
    }
  }

  return (
    <div className="min-h-screen bg-white flex">
      {/* Left side - Login Form */}
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
              <Link href="/auth/register" className="text-purple-600 hover:text-purple-500 font-medium">
                inicia sesión
              </Link>
              .
            </p>
          </div>

          {verified && (
            <div className="mb-4">
              <Alert variant="success">
                ¡Email verificado exitosamente! Ahora puedes iniciar sesión.
              </Alert>
            </div>
          )}

          {error && (
            <div className="mb-4">
              <Alert variant="error" onClose={clearError}>
                {error}
              </Alert>
            </div>
          )}

          <LoginForm onSubmit={handleLogin} isLoading={isLoading} />
          
          <div className="mt-6 flex items-center justify-between">
            <div className="flex items-center">
              <input
                id="remember-me"
                name="remember-me"
                type="checkbox"
                className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
              />
              <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-900">
                Recordarme
              </label>
            </div>

            <div className="text-sm">
              <Link href="/auth/reset-password" className="text-purple-600 hover:text-purple-500 text-sm">
                ¿Olvidaste tu contraseña?
              </Link>
            </div>
          </div>
          

          <div className="mt-8 text-center space-y-2">
            <p className="text-sm text-gray-600">
              ¿No tienes una cuenta?{' '}
              <Link href="/auth/register" className="font-medium text-purple-600 hover:text-purple-500">
                Regístrate aquí
              </Link>
            </p>
            <Link href="/" className="block text-purple-600 hover:text-purple-500 text-sm font-medium">
              ← Volver al inicio
            </Link>
          </div>
        </div>
      </div>

      {/* Right side - Benefits */}
      <BenefitsSection />
    </div>
  )
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
          <div className="max-w-md w-full bg-white py-8 px-6 shadow rounded-lg">
            <div className="text-center">
              <div className="inline-block animate-spin rounded-full h-16 w-16 border-b-2 border-purple-600 mb-4" />
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Cargando...
              </h2>
            </div>
          </div>
        </div>
      }
    >
      <LoginContent />
    </Suspense>
  )
}

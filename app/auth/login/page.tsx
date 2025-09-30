'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import LoginHeader from './components/LoginHeader'
import LoginForm from './components/LoginForm'
import SocialLogin from './components/SocialLogin'
import { AuthService } from './utils/api'
import type { LoginFormData } from './types'

export default function LoginPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()

  const handleLogin = async (data: LoginFormData) => {
    setIsLoading(true)
    setError(null)
    
    try {
      const result = await AuthService.login(data)
      
      if (result.success && result.token) {
        // Check if user data exists
        if (!result.user) {
          setError('Error en los datos del usuario')
          return
        }
        
        // Check if user is active
        if (!result.user.is_active) {
          setError('Cuenta desactivada. Contacta al administrador')
          return
        }
        
        // Check if email is verified
        if (!result.user.email_verified) {
          setError('Tu email no ha sido verificado. Revisa tu bandeja de entrada y haz clic en el enlace de verificación.')
          return
        }
        
        // Save token and redirect
        localStorage.setItem('auth_token', result.token)
        if (result.refreshToken) {
          localStorage.setItem('refresh_token', result.refreshToken)
        }
        if (result.user) {
          localStorage.setItem('user_data', JSON.stringify(result.user))
        }
        router.push('/dashboard')
      } else {
        setError(result.error || 'Error al iniciar sesión')
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
        <LoginHeader />
        
        <div className="bg-white py-8 px-6 shadow rounded-lg">
          {error && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">{error}</p>
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
          
          <div className="mt-6">
            <SocialLogin />
          </div>
        </div>

        <div className="text-center space-y-2">
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
  )
}

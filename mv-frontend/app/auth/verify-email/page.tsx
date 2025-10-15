'use client'

import { Suspense, useEffect, useState, useRef, useCallback } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline'
import { AuthService } from '@/services/authService'
import { Button } from '@/components/ui/Button'
import { Alert } from '@/components/ui/Alert'

type VerificationStatus = 'loading' | 'success' | 'error' | 'no-token'

function VerifyEmailContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = searchParams.get('token')

  const [status, setStatus] = useState<VerificationStatus>('loading')
  const [message, setMessage] = useState('')
  const [countdown, setCountdown] = useState(3)

  // useRef para prevenir llamadas duplicadas (React StrictMode ejecuta efectos 2 veces)
  const hasVerified = useRef(false)

  /**
   * Función de verificación memoizada con useCallback
   * Previene recreación en cada render
   */
  const verifyEmail = useCallback(async () => {
    if (!token) return

    // Prevenir ejecuciones duplicadas
    if (hasVerified.current) {
      console.log('⚠️ Verificación ya ejecutada, saltando...')
      return
    }

    hasVerified.current = true

    try {
      console.log('🔄 Iniciando verificación de email...')
      await AuthService.verifyEmail(token)

      setStatus('success')
      setMessage('¡Email verificado exitosamente! Tu cuenta ha sido activada.')
      console.log('✅ Verificación exitosa')
    } catch (error) {
      setStatus('error')
      setMessage(
        error instanceof Error
          ? error.message
          : 'Error al verificar el email. El token puede haber expirado.'
      )
      console.error('❌ Error en verificación:', error)
    }
  }, [token])

  /**
   * useEffect para ejecutar verificación solo una vez
   * Dependencias correctas: token y verifyEmail
   */
  useEffect(() => {
    if (!token) {
      setStatus('no-token')
      setMessage('No se encontró el token de verificación')
      return
    }

    // Solo ejecutar si no se ha verificado aún
    if (!hasVerified.current) {
      verifyEmail()
    }
  }, [token, verifyEmail])

  /**
   * useEffect para countdown y redirección
   */
  useEffect(() => {
    if (status === 'success' && countdown > 0) {
      const timer = setTimeout(() => {
        setCountdown(countdown - 1)
      }, 1000)

      return () => clearTimeout(timer)
    }

    if (status === 'success' && countdown === 0) {
      router.push('/auth/login?verified=true')
    }
  }, [status, countdown, router])

  if (status === 'loading') {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
        <div className="max-w-md w-full bg-white py-8 px-6 shadow rounded-lg">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-16 w-16 border-b-2 border-purple-600 mb-4" />
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              Verificando email...
            </h2>
            <p className="text-gray-600">Por favor espera un momento</p>
          </div>
        </div>
      </div>
    )
  }

  if (status === 'no-token') {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
        <div className="max-w-md w-full space-y-6">
          <div className="bg-white py-8 px-6 shadow rounded-lg">
            <div className="text-center">
              <XCircleIcon className="h-16 w-16 text-red-600 mx-auto mb-4" />
              <h2 className="text-2xl font-bold text-red-600 mb-4">
                Token no encontrado
              </h2>
              <p className="text-gray-600 mb-6">{message}</p>

              <Link href="/auth/register">
                <Button variant="primary" fullWidth>
                  Volver a registro
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (status === 'success') {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
        <div className="max-w-md w-full space-y-6">
          <div className="bg-white py-8 px-6 shadow rounded-lg">
            <div className="text-center">
              <CheckCircleIcon className="h-16 w-16 text-green-600 mx-auto mb-4" />
              <h2 className="text-2xl font-bold text-green-600 mb-4">
                ¡Verificación exitosa!
              </h2>

              <Alert variant="success">
                {message}
              </Alert>

              <p className="text-sm text-gray-500 mt-4">
                Redirigiendo al login en {countdown} segundo{countdown !== 1 ? 's' : ''}
                ...
              </p>

              <div className="mt-6">
                <Link href="/auth/login?verified=true">
                  <Button variant="primary" fullWidth>
                    Ir al login ahora
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // Error state
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
      <div className="max-w-md w-full space-y-6">
        <div className="bg-white py-8 px-6 shadow rounded-lg">
          <div className="text-center">
            <XCircleIcon className="h-16 w-16 text-red-600 mx-auto mb-4" />
            <h2 className="text-2xl font-bold text-red-600 mb-4">
              Error de verificación
            </h2>

            <div className="mb-6">
              <Alert variant="error">{message}</Alert>
            </div>

            <div className="space-y-3">
              <Button
                variant="primary"
                fullWidth
                onClick={() => router.push('/auth/login')}
              >
                Ir al login
              </Button>

              <Button
                variant="secondary"
                fullWidth
                onClick={() => window.location.reload()}
              >
                Intentar de nuevo
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function VerifyEmailPage() {
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
      <VerifyEmailContent />
    </Suspense>
  )
}

'use client'

import { Suspense, useState, useEffect } from 'react'
import Link from 'next/link'
import { useSearchParams, useRouter } from 'next/navigation'
import ResetPasswordHeader from './components/ResetPasswordHeader'
import ResetPasswordForm from './components/ResetPasswordForm'
import SuccessMessage from './components/SuccessMessage'
import { ResetPasswordService } from './utils/api'
import type { ResetPasswordFormData } from './types'

function ResetPasswordContent() {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const [mode, setMode] = useState<'request' | 'reset'>('request')
  const [email, setEmail] = useState('')

  const searchParams = useSearchParams()
  const router = useRouter()

  // Check if there's a token in the URL
  useEffect(() => {
    const token = searchParams.get('token')
    if (token) {
      setMode('reset')
      // Pre-fill token if available
    }
  }, [searchParams])

  const handleSubmit = async (data: ResetPasswordFormData) => {
    setIsLoading(true)
    setError(null)
    
    try {
      if (mode === 'request') {
        const result = await ResetPasswordService.requestPasswordReset(data.email)
        
        if (result.success) {
          setEmail(data.email)
          setSuccess(true)
        } else {
          setError(result.error || 'Error al enviar instrucciones')
        }
      } else {
        const result = await ResetPasswordService.resetPassword(data.token, data.new_password)
        
        if (result.success) {
          setSuccess(true)
        } else {
          setError(result.error || 'Error al actualizar contraseña')
        }
      }
    } catch (err) {
      setError('Ha ocurrido un error inesperado')
    } finally {
      setIsLoading(false)
    }
  }

  if (success) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <SuccessMessage 
            mode={mode === 'request' ? 'email-sent' : 'password-updated'} 
            email={email} 
          />
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <ResetPasswordHeader mode={mode} />
        
        <div className="bg-white py-8 px-6 shadow rounded-lg">
          {error && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}
          
          <ResetPasswordForm onSubmit={handleSubmit} isLoading={isLoading} mode={mode} />
        </div>

        <div className="text-center space-y-2">
          <p className="text-sm text-gray-600">
            ¿Recordaste tu contraseña?{' '}
            <Link href="/auth/login" className="font-medium text-blue-600 hover:text-blue-500">
              Inicia sesión aquí
            </Link>
          </p>
          <Link href="/" className="block text-blue-600 hover:text-blue-500 text-sm font-medium">
            ← Volver al inicio
          </Link>
        </div>
      </div>
    </div>
  )
}

export default function ResetPasswordPage() {
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
      <ResetPasswordContent />
    </Suspense>
  )
}

import React from 'react'
import Link from 'next/link'
import { CheckCircleIcon } from '@heroicons/react/24/outline'

interface SuccessMessageProps {
  mode: 'email-sent' | 'password-updated';
  email?: string;
}

export default function SuccessMessage({ mode, email }: SuccessMessageProps) {
  return (
    <div className="text-center">
      <div className="mb-6">
        <CheckCircleIcon className="mx-auto h-16 w-16 text-green-600" />
      </div>
      
      {mode === 'email-sent' ? (
        <>
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            Revisa tu email
          </h2>
          <p className="text-gray-600 mb-6">
            Hemos enviado las instrucciones de recuperación a{' '}
            <span className="font-semibold text-gray-900">{email}</span>
          </p>
          <div className="space-y-4">
            <p className="text-sm text-gray-500">
              No recibiste el email? Revisa tu carpeta de spam o solicita un nuevo enlace en unos minutos.
            </p>
            <Link
              href="/auth/login"
              className="inline-block text-blue-600 hover:text-blue-500 font-medium"
            >
              Volver al login
            </Link>
          </div>
        </>
      ) : (
        <>
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            ¡Contraseña actualizada!
          </h2>
          <p className="text-gray-600 mb-6">
            Tu contraseña ha sido actualizada exitosamente. Ya puedes iniciar sesión con tu nueva contraseña.
          </p>
          <Link
            href="/auth/login"
            className="inline-block bg-blue-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-blue-700 transition-colors"
          >
            Iniciar sesión
          </Link>
        </>
      )}
    </div>
  )
}

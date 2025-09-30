'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { CheckCircleIcon, XCircleIcon, ClockIcon } from '@heroicons/react/24/outline';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

interface VerificationResponse {
  success: boolean;
  message: string;
  error?: string;
}

export default function VerifyEmailPage() {
  const params = useParams();
  const router = useRouter();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');

  const token = params.token as string;

  useEffect(() => {
    if (!token) {
      setStatus('error');
      setMessage('Token de verificación no válido');
      return;
    }

    verifyEmail();
  }, [token]);

  const verifyEmail = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/verify-email`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token }),
      });

      const data: VerificationResponse = await response.json();

      if (data.success) {
        setStatus('success');
        setMessage('¡Email verificado exitosamente! Tu cuenta ha sido activada.');
        
        // Redirect to login after 3 seconds
        setTimeout(() => {
          router.push('/auth/login?verified=true');
        }, 3000);
      } else {
        setStatus('error');
        setMessage(data.error || data.message || 'Error al verificar el email');
      }
    } catch (error) {
      setStatus('error');
      setMessage('Error de conexión. Inténtalo de nuevo más tarde.');
      console.error('Error verifying email:', error);
    }
  };

  const getIcon = () => {
    switch (status) {
      case 'loading':
        return <ClockIcon className="h-16 w-16 text-purple-600 animate-spin" />;
      case 'success':
        return <CheckCircleIcon className="h-16 w-16 text-green-600" />;
      case 'error':
        return <XCircleIcon className="h-16 w-16 text-red-600" />;
    }
  };

  const getTitle = () => {
    switch (status) {
      case 'loading':
        return 'Verificando email...';
      case 'success':
        return '¡Verificación exitosa!';
      case 'error':
        return 'Error de verificación';
    }
  };

  const getTextColor = () => {
    switch (status) {
      case 'loading':
        return 'text-purple-600';
      case 'success':
        return 'text-green-600';
      case 'error':
        return 'text-red-600';
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        <div className="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
          <div className="text-center">
            <div className="flex justify-center mb-6">
              {getIcon()}
            </div>
            
            <h2 className={`text-2xl font-bold mb-4 ${getTextColor()}`}>
              {getTitle()}
            </h2>
            
            <p className="text-gray-600 mb-6">
              {message}
            </p>

            {status === 'success' && (
              <p className="text-sm text-gray-500">
                Serás redirigido al login en unos segundos...
              </p>
            )}

            {status === 'error' && (
              <div className="space-y-4">
                <button
                  onClick={() => router.push('/auth/login')}
                  className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
                >
                  Ir al login
                </button>
                
                <button
                  onClick={() => window.location.reload()}
                  className="w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
                >
                  Intentar de nuevo
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

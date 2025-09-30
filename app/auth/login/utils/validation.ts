// MisViáticos Login - Validation Utilities

import type { LoginFormData, AuthError } from '../types'

export const validateEmail = (email: string): string | null => {
  if (!email) {
    return 'El email es requerido'
  }
  
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email)) {
    return 'Por favor ingresa un email válido'
  }
  
  return null
}

export const validatePassword = (password: string): string | null => {
  if (!password) {
    return 'La contraseña es requerida'
  }
  
  if (password.length < 6) {
    return 'La contraseña debe tener al menos 6 caracteres'
  }
  
  return null
}

export const validateLoginForm = (data: LoginFormData): Record<keyof LoginFormData, string | null> => {
  return {
    email: validateEmail(data.email),
    password: validatePassword(data.password)
  }
}

export const hasValidationErrors = (errors: Record<string, string | null>): boolean => {
  return Object.values(errors).some(error => error !== null)
}

export const formatAuthError = (error: AuthError): string => {
  const errorMessages: Record<string, string> = {
    'invalid-credentials': 'Email o contraseña incorrectos',
    'user-not-found': 'No existe una cuenta con este email',
    'account-disabled': 'Tu cuenta ha sido deshabilitada',
    'too-many-requests': 'Demasiados intentos fallidos. Intenta más tarde',
    'network-error': 'Error de conexión. Verifica tu internet',
    'server-error': 'Error del servidor. Intenta más tarde'
  }
  
  return errorMessages[error.code] || error.message || 'Ha ocurrido un error inesperado'
}

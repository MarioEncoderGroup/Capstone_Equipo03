// MisViáticos Register - Validation Utilities

import { validatePhone } from '@/lib/utils/validators'
import type { RegisterFormData, RegisterError, PasswordStrength } from '../types'

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
  
  if (password.length < 8) {
    return 'La contraseña debe tener al menos 8 caracteres'
  }
  
  const hasUppercase = /[A-Z]/.test(password)
  const hasLowercase = /[a-z]/.test(password)
  const hasNumbers = /\d/.test(password)
  const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password)
  
  if (!hasUppercase || !hasLowercase || !hasNumbers) {
    return 'La contraseña debe contener mayúsculas, minúsculas y números'
  }
  
  return null
}

export const validatePasswordConfirmation = (password: string, confirmPassword: string): string | null => {
  if (!confirmPassword) {
    return 'Confirma tu contraseña'
  }
  
  if (password !== confirmPassword) {
    return 'Las contraseñas no coinciden'
  }
  
  return null
}

export const validateName = (name: string, fieldName: string): string | null => {
  if (!name) {
    return `${fieldName} es requerido`
  }
  
  if (name.length < 2) {
    return `${fieldName} debe tener al menos 2 caracteres`
  }
  
  if (name.length > 50) {
    return `${fieldName} no puede exceder 50 caracteres`
  }
  
  return null
}

export const validateTermsAcceptance = (accepted: boolean): string | null => {
  if (!accepted) {
    return 'Debes aceptar los términos y condiciones'
  }
  
  return null
}

export const validateRegisterForm = (data: RegisterFormData): Record<keyof RegisterFormData, string | null> => {
  return {
    full_name: validateName(data.full_name, 'El nombre'),
    email: validateEmail(data.email),
    phone: validatePhone(data.phone),
    password: validatePassword(data.password),
    password_confirm: validatePasswordConfirmation(data.password, data.password_confirm)
  }
}

export const calculatePasswordStrength = (password: string): PasswordStrength => {
  let score = 0
  const feedback: string[] = []
  
  if (password.length >= 8) score += 1
  else feedback.push('Mínimo 8 caracteres')
  
  if (/[a-z]/.test(password)) score += 1
  else feedback.push('Incluye minúsculas')
  
  if (/[A-Z]/.test(password)) score += 1
  else feedback.push('Incluye mayúsculas')
  
  if (/\d/.test(password)) score += 1
  else feedback.push('Incluye números')
  
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) score += 1
  else feedback.push('Incluye caracteres especiales')
  
  return {
    score,
    feedback,
    isValid: score >= 4
  }
}

export const hasValidationErrors = (errors: Record<string, string | null>): boolean => {
  return Object.values(errors).some(error => error !== null)
}

export const formatRegisterError = (error: RegisterError): string => {
  const errorMessages: Record<string, string> = {
    'email-already-exists': 'Ya existe una cuenta con este email',
    'invalid-email': 'El formato del email no es válido',
    'weak-password': 'La contraseña es muy débil',
    'server-error': 'Error del servidor. Intenta más tarde',
    'network-error': 'Error de conexión. Verifica tu internet'
  }
  
  return errorMessages[error.code] || error.message || 'Ha ocurrido un error inesperado'
}

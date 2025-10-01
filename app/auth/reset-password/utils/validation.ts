// MisViáticos Reset Password - Validation Utilities

import type { ResetPasswordFormData, ResetPasswordError } from '../types'

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

export const validateNewPassword = (password: string): string | null => {
  if (!password) {
    return 'La nueva contraseña es requerida'
  }
  
  if (password.length < 8) {
    return 'La contraseña debe tener al menos 8 caracteres'
  }
  
  const hasUppercase = /[A-Z]/.test(password)
  const hasLowercase = /[a-z]/.test(password)
  const hasNumbers = /\d/.test(password)
  
  if (!hasUppercase || !hasLowercase || !hasNumbers) {
    return 'La contraseña debe contener mayúsculas, minúsculas y números'
  }
  
  return null
}

export const validatePasswordConfirmation = (password: string, confirmPassword: string): string | null => {
  if (!confirmPassword) {
    return 'Confirma tu nueva contraseña'
  }
  
  if (password !== confirmPassword) {
    return 'Las contraseñas no coinciden'
  }
  
  return null
}

export const validateResetToken = (token: string): string | null => {
  if (!token) {
    return 'Token de recuperación requerido'
  }
  
  // Basic token format validation (should be a UUID or similar)
  if (token.length < 20) {
    return 'Token de recuperación inválido'
  }
  
  return null
}

export const validateResetPasswordForm = (
  data: ResetPasswordFormData, 
  mode: 'request' | 'reset'
): Record<string, string | null> => {
  if (mode === 'request') {
    return {
      email: validateEmail(data.email),
      newPassword: null,
      confirmPassword: null,
      token: null
    }
  }
  
  return {
    email: null,
    new_password: validateNewPassword(data.new_password),
    confirm_password: validatePasswordConfirmation(data.new_password, data.confirm_password),
    token: validateResetToken(data.token)
  }
}

export const hasValidationErrors = (errors: Record<string, string | null>): boolean => {
  return Object.values(errors).some(error => error !== null)
}

export const formatResetPasswordError = (error: ResetPasswordError): string => {
  const errorMessages: Record<string, string> = {
    'invalid-email': 'El email ingresado no es válido',
    'email-not-found': 'No existe una cuenta con este email',
    'invalid-token': 'El enlace de recuperación es inválido o ha expirado',
    'token-expired': 'El enlace de recuperación ha expirado',
    'password-too-weak': 'La contraseña no cumple con los requisitos de seguridad',
    'rate-limit-exceeded': 'Demasiados intentos. Intenta más tarde',
    'server-error': 'Error del servidor. Intenta más tarde',
    'network-error': 'Error de conexión. Verifica tu internet'
  }
  
  return errorMessages[error.code] || error.message || 'Ha ocurrido un error inesperado'
}

export const getPasswordStrengthIndicator = (password: string) => {
  if (!password) return { strength: 0, message: '', color: 'gray' }
  
  let score = 0
  const checks = [
    { test: password.length >= 8, message: 'Al menos 8 caracteres' },
    { test: /[a-z]/.test(password), message: 'Minúsculas' },
    { test: /[A-Z]/.test(password), message: 'Mayúsculas' },
    { test: /\d/.test(password), message: 'Números' },
    { test: /[!@#$%^&*(),.?":{}|<>]/.test(password), message: 'Símbolos' }
  ]
  
  score = checks.filter(check => check.test).length
  
  const strengthLevels = [
    { min: 0, max: 1, message: 'Muy débil', color: 'red' },
    { min: 2, max: 2, message: 'Débil', color: 'orange' },
    { min: 3, max: 3, message: 'Regular', color: 'yellow' },
    { min: 4, max: 4, message: 'Buena', color: 'blue' },
    { min: 5, max: 5, message: 'Excelente', color: 'green' }
  ]
  
  const level = strengthLevels.find(l => score >= l.min && score <= l.max) || strengthLevels[0]
  
  return {
    strength: score,
    message: level.message,
    color: level.color,
    requirements: checks.map(check => ({
      ...check,
      passed: check.test
    }))
  }
}

// MisViáticos - Users Feature Validation

import type { UserFormData } from '../types'

/**
 * Validation errors for user form
 */
export interface UserValidationErrors {
  full_name?: string
  email?: string
  phone?: string
  password?: string
  password_confirm?: string
}

/**
 * Validate user form data
 */
export function validateUserForm(
  formData: UserFormData,
  isEdit: boolean
): UserValidationErrors {
  const errors: UserValidationErrors = {}

  // Full name validation
  if (!formData.full_name.trim()) {
    errors.full_name = 'El nombre es requerido'
  } else if (formData.full_name.trim().length < 3) {
    errors.full_name = 'El nombre debe tener al menos 3 caracteres'
  } else if (formData.full_name.trim().length > 150) {
    errors.full_name = 'El nombre no puede exceder 150 caracteres'
  }

  // Email validation
  if (!formData.email.trim()) {
    errors.email = 'El email es requerido'
  } else if (!isValidEmail(formData.email)) {
    errors.email = 'Email inválido'
  }

  // Phone validation
  if (!formData.phone.trim()) {
    errors.phone = 'El teléfono es requerido'
  } else if (!isValidChileanPhone(formData.phone)) {
    errors.phone = 'Teléfono inválido (formato: +56 9 XXXX XXXX)'
  }

  // Password validation (only for new users or when changing password)
  if (!isEdit) {
    if (!formData.password) {
      errors.password = 'La contraseña es requerida'
    } else if (formData.password.length < 8) {
      errors.password = 'La contraseña debe tener al menos 8 caracteres'
    } else if (!isStrongPassword(formData.password)) {
      errors.password = 'La contraseña debe contener mayúsculas, minúsculas y números'
    }

    if (formData.password !== formData.password_confirm) {
      errors.password_confirm = 'Las contraseñas no coinciden'
    }
  } else if (formData.password) {
    // Validar solo si se está cambiando la contraseña
    if (formData.password.length < 8) {
      errors.password = 'La contraseña debe tener al menos 8 caracteres'
    } else if (!isStrongPassword(formData.password)) {
      errors.password = 'La contraseña debe contener mayúsculas, minúsculas y números'
    }

    if (formData.password !== formData.password_confirm) {
      errors.password_confirm = 'Las contraseñas no coinciden'
    }
  }

  return errors
}

/**
 * Check if email is valid
 */
function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Check if Chilean phone number is valid
 */
function isValidChileanPhone(phone: string): boolean {
  // Formato: +56 9 XXXX XXXX or +56 9 XXXXXXXX
  const phoneRegex = /^\+56\s?9\s?\d{4}\s?\d{4}$/
  return phoneRegex.test(phone)
}

/**
 * Check if password is strong enough
 */
function isStrongPassword(password: string): boolean {
  const hasUpperCase = /[A-Z]/.test(password)
  const hasLowerCase = /[a-z]/.test(password)
  const hasNumbers = /\d/.test(password)
  return hasUpperCase && hasLowerCase && hasNumbers
}

/**
 * Format phone number to Chilean format
 */
export function formatChileanPhone(phone: string): string {
  // Remove all non-digit characters except +
  const cleaned = phone.replace(/[^\d+]/g, '')

  // If doesn't start with +56, add it
  if (!cleaned.startsWith('+56')) {
    if (cleaned.startsWith('56')) {
      return `+${cleaned}`
    }
    return `+56${cleaned}`
  }

  return cleaned
}

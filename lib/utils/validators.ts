// MisViáticos - Validation Utilities

/**
 * Valida formato de email
 */
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

/**
 * Valida contraseña (mínimo 8 caracteres)
 */
export const validatePassword = (password: string): string | null => {
  if (!password) {
    return 'La contraseña es requerida'
  }

  if (password.length < 8) {
    return 'La contraseña debe tener al menos 8 caracteres'
  }

  return null
}

/**
 * Valida que las contraseñas coincidan
 */
export const validatePasswordConfirm = (
  password: string,
  passwordConfirm: string
): string | null => {
  if (!passwordConfirm) {
    return 'Confirma tu contraseña'
  }

  if (password !== passwordConfirm) {
    return 'Las contraseñas no coinciden'
  }

  return null
}

/**
 * Valida nombre completo (mínimo 2 caracteres)
 */
export const validateFullName = (name: string): string | null => {
  if (!name) {
    return 'El nombre es requerido'
  }

  if (name.trim().length < 2) {
    return 'El nombre debe tener al menos 2 caracteres'
  }

  return null
}

/**
 * Valida teléfono chileno (+56912345678 o 912345678)
 */
export const validatePhone = (phone: string): string | null => {
  if (!phone) {
    return 'El teléfono es requerido'
  }

  // Permitir formato +56912345678 o 912345678
  const phoneRegex = /^(\+56)?[2-9]\d{8}$/
  const cleanPhone = phone.replace(/\s|-/g, '')

  if (!phoneRegex.test(cleanPhone)) {
    return 'Ingresa un teléfono válido (ej: +56912345678)'
  }

  return null
}

/**
 * Valida RUT chileno
 */
export const validateRUT = (rut: string): string | null => {
  if (!rut) {
    return 'El RUT es requerido'
  }

  // Remover puntos y guión
  const cleanRUT = rut.replace(/\./g, '').replace(/-/g, '')

  // Formato: 12345678K o 123456789
  const rutRegex = /^(\d{7,8})([0-9Kk])$/
  if (!rutRegex.test(cleanRUT)) {
    return 'Ingresa un RUT válido (ej: 12.345.678-9)'
  }

  // Validar dígito verificador
  const match = cleanRUT.match(rutRegex)
  if (!match) return 'RUT inválido'

  const num = match[1]
  const dv = match[2].toUpperCase()

  let sum = 0
  let multiplier = 2

  for (let i = num.length - 1; i >= 0; i--) {
    sum += Number.parseInt(num[i]) * multiplier
    multiplier = multiplier === 7 ? 2 : multiplier + 1
  }

  const expectedDV = 11 - (sum % 11)
  const calculatedDV =
    expectedDV === 11 ? '0' : expectedDV === 10 ? 'K' : expectedDV.toString()

  if (dv !== calculatedDV) {
    return 'Dígito verificador del RUT inválido'
  }

  return null
}

/**
 * Valida que un campo no esté vacío
 */
export const validateRequired = (
  value: string,
  fieldName: string
): string | null => {
  if (!value || value.trim().length === 0) {
    return `${fieldName} es requerido`
  }
  return null
}

/**
 * Valida formato de URL
 */
export const validateWebsite = (url: string): string | null => {
  if (!url) {
    return 'El sitio web es requerido'
  }

  try {
    const urlObj = new URL(url)
    if (urlObj.protocol !== 'http:' && urlObj.protocol !== 'https:') {
      return 'La URL debe comenzar con http:// o https://'
    }
    return null
  } catch {
    return 'Ingresa una URL válida (ej: https://ejemplo.cl)'
  }
}

/**
 * Valida longitud máxima de un campo
 */
export const validateMaxLength = (
  value: string,
  maxLength: number,
  fieldName: string
): string | null => {
  if (value.length > maxLength) {
    return `${fieldName} no puede tener más de ${maxLength} caracteres`
  }
  return null
}

// MisViáticos - Input Sanitization Utilities

/**
 * Sanitiza un input de texto general
 * Previene XSS básico y limita longitud
 */
export const sanitizeInput = (input: string, maxLength = 255): string => {
  return input
    .trim()
    .replace(/[<>]/g, '') // Prevenir XSS básico
    .substring(0, maxLength)
}

/**
 * Sanitiza y normaliza un email
 */
export const sanitizeEmail = (email: string): string => {
  return email.toLowerCase().trim()
}

/**
 * Sanitiza un número de teléfono
 * Remueve espacios y caracteres no numéricos excepto + y -
 */
export const sanitizePhone = (phone: string): string => {
  return phone.trim().replace(/[^\d+\s-]/g, '')
}

/**
 * Sanitiza un RUT chileno
 * Formato esperado: 12.345.678-9
 */
export const sanitizeRUT = (rut: string): string => {
  return rut
    .trim()
    .replace(/[^\dKk.-]/g, '') // Solo dígitos, K, puntos y guiones
    .toUpperCase()
}

/**
 * Sanitiza un nombre completo
 */
export const sanitizeName = (name: string): string => {
  return name
    .trim()
    .replace(/[<>]/g, '')
    .replace(/\s+/g, ' ') // Reemplazar múltiples espacios por uno
}

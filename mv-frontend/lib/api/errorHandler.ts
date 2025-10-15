// MisViáticos - API Error Handler

import type { ApiErrorResponse, ValidationError } from '@/types/api'

export class ApiError extends Error {
  constructor(
    public statusCode: number,
    public code: string,
    message: string,
    public validationErrors?: ValidationError[]
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

/**
 * Maneja la respuesta de la API y lanza errores si es necesario
 */
export async function handleApiResponse<T>(response: Response): Promise<T> {
  const contentType = response.headers.get('content-type')
  const isJson = contentType?.includes('application/json')

  if (!response.ok) {
    let errorData: ApiErrorResponse | undefined

    if (isJson) {
      errorData = await response.json().catch(() => undefined)
    }

    throw new ApiError(
      response.status,
      errorData?.error || 'UNKNOWN_ERROR',
      errorData?.message || `Error ${response.status}: ${response.statusText}`,
      errorData?.data
    )
  }

  if (isJson) {
    return response.json()
  }

  // Si no es JSON, retornar como texto
  return (await response.text()) as T
}

/**
 * Formatea errores de validación para mostrar al usuario
 */
export function formatValidationErrors(errors: ValidationError[]): string {
  if (errors.length === 0) return 'Error de validación'

  return errors.map((err) => `${err.Field}: ${err.Message}`).join(', ')
}

/**
 * Obtiene mensaje de error user-friendly
 */
export function getErrorMessage(error: unknown): string {
  console.log('getErrorMessage called with:', error)
  
  if (error instanceof ApiError) {
    console.log('ApiError detected:', {
      statusCode: error.statusCode,
      code: error.code,
      message: error.message,
      validationErrors: error.validationErrors
    })
    
    if (error.validationErrors && error.validationErrors.length > 0) {
      return formatValidationErrors(error.validationErrors)
    }

    // Mensajes específicos por código de error
    const errorMessages: Record<string, string> = {
      EMAIL_NOT_VERIFIED:
        'Por favor verifica tu email antes de iniciar sesión',
      ACCOUNT_DEACTIVATED: 'Tu cuenta ha sido desactivada. Contacta al soporte',
      INVALID_CREDENTIALS: 'Email o contraseña incorrectos',
      TOKEN_EXPIRED: 'Tu sesión ha expirado. Por favor inicia sesión nuevamente',
      TENANT_NOT_FOUND: 'Empresa no encontrada',
      TENANT_ALREADY_EXISTS: 'Ya existe una empresa registrada con este RUT',
      RUT_ALREADY_EXISTS: 'Ya existe una empresa registrada con este RUT',
      DUPLICATE_RUT: 'Ya existe una empresa registrada con este RUT',
      CONFLICT: 'Ya existe una empresa registrada con este RUT',
      'ya existe un tenant con el RUT proporcionado': 'Ya existe una empresa registrada con este RUT',
      UNAUTHORIZED: 'No tienes autorización para realizar esta acción',
    }

    // Verificar también el mensaje de error directamente para RUT duplicado
    if (error.message && error.message.toLowerCase().includes('rut')) {
      return 'Ya existe una empresa registrada con este RUT'
    }

    // Verificar status code 409 (Conflict) típico para duplicados
    if (error.statusCode === 409) {
      return 'Ya existe una empresa registrada con este RUT'
    }

    const mappedMessage = errorMessages[error.code]
    console.log(`Mapped message for code "${error.code}":`, mappedMessage)
    
    return mappedMessage || error.message
  }

  if (error instanceof Error) {
    console.log('Regular Error:', error.message)
    return error.message
  }

  console.log('Unknown error type:', error)
  return 'Ha ocurrido un error inesperado'
}

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
  if (error instanceof ApiError) {
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
      UNAUTHORIZED: 'No tienes autorización para realizar esta acción',
    }

    return errorMessages[error.code] || error.message
  }

  if (error instanceof Error) {
    return error.message
  }

  return 'Ha ocurrido un error inesperado'
}

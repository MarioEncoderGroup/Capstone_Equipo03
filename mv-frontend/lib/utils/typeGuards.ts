// MisViáticos - TypeScript Type Guards

import type {
  ApiResponse,
  UserData,
  TenantData,
  JWTPayload,
  ApiErrorResponse,
} from '@/types/api'

/**
 * Type guard para verificar si es una respuesta exitosa de la API
 */
export function isApiSuccess<T>(
  response: unknown
): response is ApiResponse<T> & { success: true; data: T } {
  return (
    typeof response === 'object' &&
    response !== null &&
    'success' in response &&
    response.success === true &&
    'data' in response
  )
}

/**
 * Type guard para verificar si es un error de la API
 */
export function isApiError(error: unknown): error is ApiErrorResponse {
  return (
    typeof error === 'object' &&
    error !== null &&
    'success' in error &&
    error.success === false &&
    'message' in error
  )
}

/**
 * Type guard para verificar si es un UserData válido
 */
export function isUser(value: unknown): value is UserData {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'email' in value &&
    'full_name' in value &&
    'is_active' in value
  )
}

/**
 * Type guard para verificar si es un TenantData válido
 */
export function isTenant(value: unknown): value is TenantData {
  return (
    typeof value === 'object' &&
    value !== null &&
    'id' in value &&
    'rut' in value &&
    'business_name' in value &&
    'status' in value
  )
}

/**
 * Type guard para verificar si es un JWT válido
 */
export function isJWTPayload(value: unknown): value is JWTPayload {
  return (
    typeof value === 'object' &&
    value !== null &&
    'user_id' in value &&
    'type' in value &&
    'iat' in value &&
    'exp' in value
  )
}

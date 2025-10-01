// MisViáticos - Auth Service

import { API_BASE_URL, API_ENDPOINTS } from '@/constants/api'
import { handleApiResponse } from '@/lib/api/errorHandler'
import { TokenManager } from '@/lib/auth/tokenManager'
import {
  sanitizeEmail,
  sanitizeName,
  sanitizePhone,
} from '@/lib/utils/sanitize'
import type {
  ApiResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  VerifyEmailRequest,
} from '@/types/api'

export class AuthService {
  /**
   * Registra un nuevo usuario
   */
  static async register(data: RegisterRequest): Promise<RegisterResponse> {
    // Sanitizar inputs
    const sanitizedData: RegisterRequest = {
      full_name: sanitizeName(data.full_name),
      email: sanitizeEmail(data.email),
      phone: sanitizePhone(data.phone),
      password: data.password,
      password_confirm: data.password_confirm,
    }

    const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.AUTH_REGISTER}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(sanitizedData),
    })

    const result = await handleApiResponse<ApiResponse<RegisterResponse>>(
      response
    )

    if (!result.data) {
      throw new Error('No se recibieron datos en la respuesta')
    }

    return result.data
  }

  /**
   * Verifica el email con el token
   */
  static async verifyEmail(token: string): Promise<void> {
    const data: VerifyEmailRequest = { token }

    const response = await fetch(
      `${API_BASE_URL}${API_ENDPOINTS.AUTH_VERIFY_EMAIL}`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      }
    )

    await handleApiResponse<ApiResponse<void>>(response)
  }

  /**
   * Inicia sesión con email y contraseña
   */
  static async login(data: LoginRequest): Promise<LoginResponse> {
    // Sanitizar email
    const sanitizedData: LoginRequest = {
      email: sanitizeEmail(data.email),
      password: data.password,
    }

    const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.AUTH_LOGIN}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(sanitizedData),
    })

    const result = await handleApiResponse<ApiResponse<LoginResponse>>(response)

    if (!result.data) {
      throw new Error('No se recibieron datos en la respuesta')
    }

    // Guardar tokens y datos de usuario
    TokenManager.setTokens(result.data.access_token, result.data.refresh_token)
    TokenManager.setUserData(result.data.user)

    return result.data
  }

  /**
   * Cierra sesión del usuario
   */
  static logout(): void {
    TokenManager.clearTokens()
  }

  /**
   * Verifica si el usuario está autenticado
   */
  static isAuthenticated(): boolean {
    return TokenManager.isAuthenticated()
  }

  /**
   * Obtiene los datos del usuario actual
   */
  static getCurrentUser(): unknown {
    return TokenManager.getUserData()
  }
}

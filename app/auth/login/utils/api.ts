// MisViáticos Login - API Utilities

import type { LoginFormData, LoginResponse, AuthError } from '../types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export class AuthService {
  static async login(data: LoginFormData): Promise<LoginResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      })

      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.message || 'Error en el login')
      }

      return {
        success: true,
        token: result.data.access_token,
        user: result.data.user,
        refreshToken: result.data.refresh_token,
        expiresIn: result.data.expires_in
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error desconocido'
      }
    }
  }

  static async loginWithGoogle(): Promise<LoginResponse> {
    // TODO: Implementar autenticación con Google
    return {
      success: false,
      error: 'Login con Google no implementado aún'
    }
  }

  static async loginWithMicrosoft(): Promise<LoginResponse> {
    // TODO: Implementar autenticación con Microsoft
    return {
      success: false,
      error: 'Login con Microsoft no implementado aún'
    }
  }

  static async logout(): Promise<void> {
    // TODO: Implementar logout
    localStorage.removeItem('auth_token')
    localStorage.removeItem('user_data')
  }
}

export const saveAuthToken = (token: string): void => {
  localStorage.setItem('auth_token', token)
}

export const getAuthToken = (): string | null => {
  return localStorage.getItem('auth_token')
}

export const isAuthenticated = (): boolean => {
  return !!getAuthToken()
}

// MisViáticos Reset Password - API Utilities

import type { ResetPasswordFormData, ResetPasswordResponse, TokenValidationResponse } from '../types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export class ResetPasswordService {
  static async requestPasswordReset(email: string): Promise<ResetPasswordResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/forgot-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.message || 'Error al solicitar recuperación')
      }

      return {
        success: true,
        message: result.message || 'Instrucciones enviadas por email'
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error desconocido'
      }
    }
  }

  static async validateResetToken(token: string): Promise<TokenValidationResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/reset-password/validate`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token }),
      })

      const result = await response.json()

      if (!response.ok) {
        return {
          isValid: false,
          error: result.message || 'Token inválido'
        }
      }

      return {
        isValid: true,
        email: result.email,
        expiresAt: result.expiresAt
      }
    } catch (error) {
      return {
        isValid: false,
        error: error instanceof Error ? error.message : 'Error al validar token'
      }
    }
  }

  static async resetPassword(token: string, newPassword: string): Promise<ResetPasswordResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/reset-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          token,
          newPassword 
        }),
      })

      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.message || 'Error al actualizar contraseña')
      }

      return {
        success: true,
        message: result.message || 'Contraseña actualizada correctamente'
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error desconocido'
      }
    }
  }

  static async resendResetEmail(email: string): Promise<ResetPasswordResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/reset-password/resend`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      const result = await response.json()

      return {
        success: response.ok,
        message: result.message,
        error: response.ok ? undefined : result.message
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error al reenviar email'
      }
    }
  }
}

export const extractTokenFromUrl = (): string | null => {
  if (typeof window === 'undefined') return null
  
  const urlParams = new URLSearchParams(window.location.search)
  return urlParams.get('token')
}

export const buildResetPasswordUrl = (token: string): string => {
  const baseUrl = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'
  return `${baseUrl}/auth/reset-password?token=${token}`
}

export const isValidTokenFormat = (token: string): boolean => {
  // Basic validation for token format (UUID-like pattern)
  const tokenRegex = /^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$/i
  return tokenRegex.test(token) || token.length >= 32
}

export const getTokenExpirationTime = (expiresAt: string): number => {
  return Math.floor((new Date(expiresAt).getTime() - Date.now()) / 1000 / 60) // minutes remaining
}

// MisViáticos Register - API Utilities

import type { RegisterFormData, RegisterResponse, RegisterError } from '../types'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export class RegisterService {
  static async register(data: RegisterFormData): Promise<RegisterResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          firstname: data.firstname,
          lastname: data.lastname,
          email: data.email,
          phone: data.phone,
          password: data.password,
          password_confirm: data.password_confirm
        }),
      })

      const result = await response.json()

      if (!response.ok) {
        throw new Error(result.message || 'Error en el registro')
      }

      return {
        success: true,
        token: result.token,
        user: result.user
      }
    } catch (error) {
      return {
        success: false,
        error: error instanceof Error ? error.message : 'Error desconocido'
      }
    }
  }

  // TODO: Implementar cuando el backend tenga el endpoint /auth/check-email
  static async checkEmailAvailability(email: string): Promise<boolean> {
    // Temporalmente deshabilitado - endpoint no existe en backend
    console.warn('checkEmailAvailability: Endpoint no implementado en backend')
    return true // Assume available if check fails
    
    /* 
    try {
      const response = await fetch(`${API_BASE_URL}/auth/check-email`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      const result = await response.json()
      return result.available
    } catch (error) {
      console.error('Error checking email availability:', error)
      return true // Assume available if check fails
    }
    */
  }

  static async sendVerificationEmail(email: string): Promise<boolean> {
    try {
      const response = await fetch(`${API_BASE_URL}/auth/resend-verification`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      })

      return response.ok
    } catch (error) {
      console.error('Error sending verification email:', error)
      return false
    }
  }
}

export const saveUserData = (user: any, token: string): void => {
  localStorage.setItem('auth_token', token)
  localStorage.setItem('user_data', JSON.stringify(user))
}

export const getUserData = (): any | null => {
  const userData = localStorage.getItem('user_data')
  return userData ? JSON.parse(userData) : null
}

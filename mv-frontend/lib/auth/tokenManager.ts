// MisViáticos - Token Manager (Secure Token Handling)

import type { JWTPayload } from '@/types/api'
import { isJWTPayload } from '@/lib/utils/typeGuards'

export class TokenManager {
  private static readonly TOKEN_KEY = 'auth_token'
  private static readonly REFRESH_KEY = 'refresh_token'
  private static readonly USER_KEY = 'user_data'

  /**
   * Guarda los tokens de autenticación
   * Almacena en localStorage Y en cookies para compatibilidad con middleware
   */
  static setTokens(accessToken: string, refreshToken: string): void {
    if (typeof window === 'undefined') return

    // Guardar en localStorage (acceso desde cliente)
    localStorage.setItem(this.TOKEN_KEY, accessToken)
    localStorage.setItem(this.REFRESH_KEY, refreshToken)

    // Guardar en cookies (acceso desde middleware)
    // Decodificar token para obtener la expiración
    const payload = this.decodeJWT(accessToken)
    const maxAge = payload?.exp ? payload.exp - Math.floor(Date.now() / 1000) : 60 * 60 * 24 * 7 // 7 días por defecto

    // Set cookie con flags de seguridad
    document.cookie = `${this.TOKEN_KEY}=${accessToken}; path=/; max-age=${maxAge}; SameSite=Lax`
    document.cookie = `${this.REFRESH_KEY}=${refreshToken}; path=/; max-age=${maxAge}; SameSite=Lax`
  }

  /**
   * Obtiene un valor de cookie por nombre
   */
  private static getCookie(name: string): string | null {
    if (typeof window === 'undefined') return null

    const value = `; ${document.cookie}`
    const parts = value.split(`; ${name}=`)
    if (parts.length === 2) {
      return parts.pop()?.split(';').shift() || null
    }
    return null
  }

  /**
   * Obtiene el access token
   * Prioriza localStorage, fallback a cookies
   */
  static getAccessToken(): string | null {
    if (typeof window === 'undefined') return null

    // Intentar desde localStorage primero
    const localToken = localStorage.getItem(this.TOKEN_KEY)
    if (localToken) return localToken

    // Fallback a cookies
    return this.getCookie(this.TOKEN_KEY)
  }

  /**
   * Obtiene el refresh token
   * Prioriza localStorage, fallback a cookies
   */
  static getRefreshToken(): string | null {
    if (typeof window === 'undefined') return null

    // Intentar desde localStorage primero
    const localToken = localStorage.getItem(this.REFRESH_KEY)
    if (localToken) return localToken

    // Fallback a cookies
    return this.getCookie(this.REFRESH_KEY)
  }

  /**
   * Guarda los datos del usuario
   */
  static setUserData(user: unknown): void {
    if (typeof window === 'undefined') return
    localStorage.setItem(this.USER_KEY, JSON.stringify(user))
  }

  /**
   * Obtiene los datos del usuario
   */
  static getUserData(): unknown {
    if (typeof window === 'undefined') return null

    try {
      const data = localStorage.getItem(this.USER_KEY)
      return data ? JSON.parse(data) : null
    } catch {
      return null
    }
  }

  /**
   * Limpia todos los tokens y datos
   * Elimina de localStorage Y de cookies
   */
  static clearTokens(): void {
    if (typeof window === 'undefined') return

    // Limpiar localStorage
    localStorage.removeItem(this.TOKEN_KEY)
    localStorage.removeItem(this.REFRESH_KEY)
    localStorage.removeItem(this.USER_KEY)

    // Limpiar cookies (expirar inmediatamente)
    document.cookie = `${this.TOKEN_KEY}=; path=/; max-age=0`
    document.cookie = `${this.REFRESH_KEY}=; path=/; max-age=0`
  }

  /**
   * Decodifica un JWT y retorna el payload
   */
  static decodeJWT(token: string): JWTPayload | null {
    try {
      const base64Url = token.split('.')[1]
      if (!base64Url) return null

      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      const jsonPayload = decodeURIComponent(
        atob(base64)
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join('')
      )

      const payload = JSON.parse(jsonPayload)

      return isJWTPayload(payload) ? payload : null
    } catch {
      return null
    }
  }

  /**
   * Verifica si un token está expirado
   */
  static isTokenExpired(token: string): boolean {
    const payload = this.decodeJWT(token)
    if (!payload) return true

    // Verificar si el token expira en los próximos 60 segundos
    return payload.exp * 1000 < Date.now() + 60000
  }

  /**
   * Verifica si el token actual tiene tenant_id
   */
  static hasTenantId(): boolean {
    const token = this.getAccessToken()
    if (!token) return false

    const payload = this.decodeJWT(token)
    return payload?.tenant_id !== undefined
  }

  /**
   * Obtiene el tenant_id del token actual
   */
  static getTenantId(): string | null {
    const token = this.getAccessToken()
    if (!token) return null

    const payload = this.decodeJWT(token)
    return payload?.tenant_id || null
  }

  /**
   * Verifica si el usuario está autenticado
   */
  static isAuthenticated(): boolean {
    try {
      console.log('🔍 TokenManager: Verificando autenticación...')
      const token = this.getAccessToken()
      console.log('📝 TokenManager: Token obtenido:', token ? 'Sí' : 'No')
      
      if (!token) {
        console.log('❌ TokenManager: No hay token')
        return false
      }

      const isExpired = this.isTokenExpired(token)
      console.log('⏰ TokenManager: Token expirado:', isExpired)
      
      const isAuth = !isExpired
      console.log('✅ TokenManager: Resultado autenticación:', isAuth)
      return isAuth
    } catch (error) {
      console.error('💥 TokenManager: Error en isAuthenticated:', error)
      return false
    }
  }
}

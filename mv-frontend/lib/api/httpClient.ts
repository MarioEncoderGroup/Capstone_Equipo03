// MisViáticos - HTTP Client with Token Management

import { API_BASE_URL } from '@/constants/api'
import { TokenManager } from '@/lib/auth/tokenManager'
import type { ApiResponse } from '@/types'

/**
 * HTTP Client configuration
 */
interface RequestConfig extends RequestInit {
  skipAuth?: boolean
}

/**
 * HTTP Client class with automatic token management
 */
export class HttpClient {
  private static isRefreshing = false
  private static refreshPromise: Promise<string> | null = null

  /**
   * Refresh access token using refresh token
   */
  private static async refreshAccessToken(): Promise<string> {
    const refreshToken = TokenManager.getRefreshToken()
    if (!refreshToken) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      })

      if (!response.ok) {
        throw new Error('Failed to refresh token')
      }

      const result: ApiResponse<{ access_token: string; refresh_token: string }> =
        await response.json()

      if (!result.data?.access_token) {
        throw new Error('Invalid refresh response')
      }

      // Update tokens
      TokenManager.setTokens(result.data.access_token, result.data.refresh_token)

      return result.data.access_token
    } catch (error) {
      // Clear tokens and redirect to login
      TokenManager.clearTokens()
      if (typeof window !== 'undefined') {
        window.location.href = '/auth/login?reason=session_expired'
      }
      throw error
    }
  }

  /**
   * Get valid access token (refreshes if needed)
   */
  private static async getValidToken(): Promise<string | null> {
    console.log('🔍 HttpClient: Getting valid token...')
    const token = TokenManager.getAccessToken()

    if (!token) {
      console.log('❌ HttpClient: No token found')
      return null
    }

    console.log('✅ HttpClient: Token found:', token.substring(0, 20) + '...')

    // Decode and log token payload for debugging
    const payload = TokenManager.decodeJWT(token)
    console.log('📋 HttpClient: Token payload:', {
      tenant_id: payload?.tenant_id,
      roles: payload?.roles,
      permissions: payload?.permissions,
      exp: payload?.exp ? new Date(payload.exp * 1000).toISOString() : 'N/A',
    })

    // Check if token is expired or about to expire
    const isExpired = TokenManager.isTokenExpired(token)
    console.log('⏰ HttpClient: Token expired?', isExpired)

    if (isExpired) {
      console.log('🔄 HttpClient: Token expired, attempting refresh...')
      // If already refreshing, wait for the promise
      if (this.isRefreshing && this.refreshPromise) {
        console.log('⏳ HttpClient: Already refreshing, waiting...')
        return await this.refreshPromise
      }

      // Start refresh process
      this.isRefreshing = true
      this.refreshPromise = this.refreshAccessToken()

      try {
        const newToken = await this.refreshPromise
        console.log('✅ HttpClient: Token refreshed successfully')
        return newToken
      } catch (error) {
        console.error('❌ HttpClient: Token refresh failed:', error)
        throw error
      } finally {
        this.isRefreshing = false
        this.refreshPromise = null
      }
    }

    console.log('✅ HttpClient: Token is valid, using it')
    return token
  }

  /**
   * Make HTTP request with automatic token management
   */
  static async request<T = any>(
    url: string,
    config: RequestConfig = {}
  ): Promise<ApiResponse<T>> {
    const { skipAuth = false, headers = {}, ...rest } = config

    console.log('🌐 HttpClient: Making request to:', url)
    console.log('🔐 HttpClient: Skip auth?', skipAuth)

    // Prepare headers
    const requestHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(headers as Record<string, string>),
    }

    // Add authorization header if not skipping auth
    if (!skipAuth) {
      const token = await this.getValidToken()
      if (token) {
        requestHeaders['Authorization'] = `Bearer ${token}`
        console.log('✅ HttpClient: Authorization header added')
      } else {
        console.error('❌ HttpClient: No token available, redirecting to login')
        // No token available - redirect to login
        if (typeof window !== 'undefined') {
          window.location.href = '/auth/login?reason=no_token'
        }
        throw new Error('No authentication token available')
      }
    }

    // Make request
    const fullUrl = url.startsWith('http') ? url : `${API_BASE_URL}${url}`
    console.log('📡 HttpClient: Full URL:', fullUrl)

    try {
      const response = await fetch(fullUrl, {
        ...rest,
        headers: requestHeaders,
      })

      console.log('📥 HttpClient: Response status:', response.status)

      // Handle 401 Unauthorized
      if (response.status === 401) {
        console.warn('⚠️ HttpClient: Got 401, attempting token refresh...')
        // Try to refresh token once
        if (!this.isRefreshing) {
          const newToken = await this.getValidToken()
          if (newToken) {
            console.log('🔄 HttpClient: Retrying with refreshed token...')
            // Retry request with new token
            requestHeaders['Authorization'] = `Bearer ${newToken}`
            const retryResponse = await fetch(fullUrl, {
              ...rest,
              headers: requestHeaders,
            })

            console.log('📥 HttpClient: Retry response status:', retryResponse.status)

            if (!retryResponse.ok) {
              throw new Error(`HTTP ${retryResponse.status}: ${retryResponse.statusText}`)
            }

            return await retryResponse.json()
          }
        }

        // If refresh failed or already refreshing, redirect to login
        console.error('❌ HttpClient: Token refresh failed, clearing tokens and redirecting')
        TokenManager.clearTokens()
        if (typeof window !== 'undefined') {
          window.location.href = '/auth/login?reason=unauthorized'
        }
        throw new Error('Unauthorized')
      }

      // Handle other errors
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        console.error('❌ HttpClient: Request failed:', response.status, errorData)
        throw new Error(errorData.message || `HTTP ${response.status}`)
      }

      console.log('✅ HttpClient: Request successful')
      return await response.json()
    } catch (error) {
      console.error('💥 HttpClient Error:', error)
      throw error
    }
  }

  /**
   * Convenience methods
   */
  static async get<T = any>(url: string, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'GET' })
  }

  static async post<T = any>(
    url: string,
    data?: any,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    return this.request<T>(url, {
      ...config,
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  static async put<T = any>(
    url: string,
    data?: any,
    config?: RequestConfig
  ): Promise<ApiResponse<T>> {
    return this.request<T>(url, {
      ...config,
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  static async delete<T = any>(url: string, config?: RequestConfig): Promise<ApiResponse<T>> {
    return this.request<T>(url, { ...config, method: 'DELETE' })
  }
}

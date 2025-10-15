// MisViáticos - Tenant Service

import { API_BASE_URL, API_ENDPOINTS } from '@/constants/api'
import { handleApiResponse } from '@/lib/api/errorHandler'
import { TokenManager } from '@/lib/auth/tokenManager'
import {
  sanitizeEmail,
  sanitizeInput,
  sanitizePhone,
  sanitizeRUT,
} from '@/lib/utils/sanitize'
import type {
  ApiResponse,
  CreateTenantRequest,
  SelectTenantResponse,
  TenantData,
  TenantStatusResponse,
} from '@/types/api'

export class TenantService {
  /**
   * Obtiene el estado de los tenants del usuario
   */
  static async getStatus(): Promise<TenantStatusResponse> {
    const token = TokenManager.getAccessToken()
    if (!token) {
      throw new Error('No hay token de autenticación')
    }

    const response = await fetch(
      `${API_BASE_URL}${API_ENDPOINTS.TENANT_STATUS}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    )

    const result = await handleApiResponse<ApiResponse<TenantStatusResponse>>(
      response
    )

    if (!result.data) {
      throw new Error('No se recibieron datos en la respuesta')
    }

    return result.data
  }

  /**
   * Crea un nuevo tenant
   */
  static async create(data: CreateTenantRequest): Promise<TenantData> {
    const token = TokenManager.getAccessToken()
    if (!token) {
      throw new Error('No hay token de autenticación')
    }

    // Sanitizar inputs
    const sanitizedData: CreateTenantRequest = {
      rut: sanitizeRUT(data.rut),
      business_name: sanitizeInput(data.business_name, 150),
      email: sanitizeEmail(data.email),
      phone: sanitizePhone(data.phone),
      address: sanitizeInput(data.address, 200),
      website: sanitizeInput(data.website, 150),
      region_id: data.region_id,
      commune_id: data.commune_id,
      country_id: data.country_id,
    }

    // Agregar logo si existe
    if (data.logo) {
      sanitizedData.logo = sanitizeInput(data.logo)
    }

    const response = await fetch(
      `${API_BASE_URL}${API_ENDPOINTS.TENANT_CREATE}`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(sanitizedData),
      }
    )

    const result = await handleApiResponse<ApiResponse<TenantData>>(response)

    if (!result.data) {
      throw new Error('No se recibieron datos en la respuesta')
    }

    return result.data
  }

  /**
   * Selecciona un tenant y obtiene nuevo token con tenant_id
   */
  static async select(tenantId: string): Promise<SelectTenantResponse> {
    const token = TokenManager.getAccessToken()
    if (!token) {
      throw new Error('No hay token de autenticación')
    }

    const response = await fetch(
      `${API_BASE_URL}${API_ENDPOINTS.TENANT_SELECT(tenantId)}`,
      {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    )

    const result = await handleApiResponse<ApiResponse<SelectTenantResponse>>(
      response
    )

    if (!result.data) {
      throw new Error('No se recibieron datos en la respuesta')
    }

    // IMPORTANTE: Reemplazar tokens con los nuevos que incluyen tenant_id
    TokenManager.setTokens(
      result.data.access_token,
      result.data.refresh_token
    )
    TokenManager.setUserData(result.data.user)

    return result.data
  }

  /**
   * Verifica si el usuario tiene tenant_id en su token
   */
  static hasTenant(): boolean {
    return TokenManager.hasTenantId()
  }

  /**
   * Obtiene el tenant_id del token actual
   */
  static getCurrentTenantId(): string | null {
    return TokenManager.getTenantId()
  }
}

// MisViáticos - Region Service
import { API_BASE_URL, API_ENDPOINTS } from '@/constants/api'
import { handleApiResponse, getErrorMessage } from '@/lib/api/errorHandler'
import type {
  ApiResponse,
  RegionsResponse,
  CommunesResponse,
  Region,
  Commune
} from '@/types/api'

/**
 * Service for handling region and commune operations
 */
export class RegionService {
  /**
   * Get all regions of Chile
   * @returns Array of regions
   */
  static async getRegions(): Promise<Region[]> {
    try {
      const response = await fetch(`${API_BASE_URL}${API_ENDPOINTS.REGIONS}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      const data = await handleApiResponse<ApiResponse<RegionsResponse>>(response)

      if (!data.success || !data.data) {
        throw new Error(data.message || 'Error al obtener regiones')
      }

      return data.data.regions
    } catch (error) {
      throw new Error(getErrorMessage(error))
    }
  }

  /**
   * Get communes by region
   * @param regionId - Region ID (e.g., "RM", "VA")
   * @returns Array of communes
   */
  static async getCommunes(regionId: string): Promise<Commune[]> {
    try {
      const response = await fetch(
        `${API_BASE_URL}${API_ENDPOINTS.COMMUNES(regionId)}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )

      const data = await handleApiResponse<ApiResponse<CommunesResponse>>(response)

      if (!data.success || !data.data) {
        throw new Error(data.message || 'Error al obtener comunas')
      }

      return data.data.communes
    } catch (error) {
      throw new Error(getErrorMessage(error))
    }
  }

  /**
   * Get all communes (without region filter)
   * @returns Array of all communes
   */
  static async getAllCommunes(): Promise<Commune[]> {
    try {
      const response = await fetch(
        `${API_BASE_URL}${API_ENDPOINTS.COMMUNES()}`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )

      const data = await handleApiResponse<ApiResponse<CommunesResponse>>(response)

      if (!data.success || !data.data) {
        throw new Error(data.message || 'Error al obtener comunas')
      }

      return data.data.communes
    } catch (error) {
      throw new Error(getErrorMessage(error))
    }
  }
}

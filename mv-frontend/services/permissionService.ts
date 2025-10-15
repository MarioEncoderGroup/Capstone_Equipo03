// MisViáticos - Permission Service

import { HttpClient } from '@/lib/api/httpClient'
import type {
  Permission,
  PermissionsResponse,
  GroupedPermission,
} from '@/types'

export class PermissionService {
  private static readonly ENDPOINT = '/permissions'

  /**
   * Get all permissions with pagination
   */
  static async getAll(limit: number = 100, page: number = 1): Promise<Permission[]> {
    const result = await HttpClient.get<PermissionsResponse>(
      `${this.ENDPOINT}?limit=${limit}&page=${page}`
    )
    return result.data?.permissions || []
  }

  /**
   * Get permissions grouped by section
   * Backend returns: { success, message, data: GroupedPermission[] }
   */
  static async getGrouped(): Promise<GroupedPermission[]> {
    const result = await HttpClient.get<GroupedPermission[]>(
      `${this.ENDPOINT}/grouped`
    )
    console.log('📊 PermissionService.getGrouped() response:', result)
    // El backend devuelve el array directamente en data
    return result.data || []
  }

  /**
   * Get permission by ID
   */
  static async getById(id: string): Promise<Permission> {
    const result = await HttpClient.get<Permission>(`${this.ENDPOINT}/${id}`)
    if (!result.data) throw new Error('Permiso no encontrado')
    return result.data
  }
}

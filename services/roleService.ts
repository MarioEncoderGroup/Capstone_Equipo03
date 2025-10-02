// MisViáticos - Role Service

import { HttpClient } from '@/lib/api/httpClient'
import { sanitizeInput } from '@/lib/utils/sanitize'
import type {
  Role,
  CreateRoleRequest,
  UpdateRoleRequest,
  RolesResponse,
  SyncPermissionsRequest,
} from '@/types'

export class RoleService {
  private static readonly ENDPOINT = '/roles'

  /**
   * Get all roles with pagination
   * Backend returns: { success, message, data: Role[], pagination }
   */
  static async getAll(limit: number = 100, page: number = 1): Promise<Role[]> {
    const result = await HttpClient.get<Role[]>(
      `${this.ENDPOINT}?limit=${limit}&page=${page}`
    )
    console.log('🎭 RoleService.getAll() response:', result)
    // El backend devuelve el array directamente en data
    return result.data || []
  }

  /**
   * Get role by ID
   */
  static async getById(id: string): Promise<Role> {
    const result = await HttpClient.get<Role>(`${this.ENDPOINT}/${id}`)
    if (!result.data) throw new Error('Rol no encontrado')
    return result.data
  }

  /**
   * Create role
   */
  static async create(data: CreateRoleRequest): Promise<Role> {
    const sanitizedData: CreateRoleRequest = {
      name: sanitizeInput(data.name, 100),
      description: sanitizeInput(data.description, 255),
      permission_ids: data.permission_ids,
    }

    const result = await HttpClient.post<Role>(this.ENDPOINT, sanitizedData)
    if (!result.data) throw new Error('Error al crear rol')
    return result.data
  }

  /**
   * Update role
   */
  static async update(id: string, data: UpdateRoleRequest): Promise<Role> {
    const sanitizedData: UpdateRoleRequest = {
      ...(data.name && { name: sanitizeInput(data.name, 100) }),
      ...(data.description && { description: sanitizeInput(data.description, 255) }),
    }

    const result = await HttpClient.put<Role>(`${this.ENDPOINT}/${id}`, sanitizedData)
    if (!result.data) throw new Error('Error al actualizar rol')
    return result.data
  }

  /**
   * Sync role permissions (replace all permissions) - SINGLE CALL
   */
  static async syncPermissions(roleId: string, permissionIds: string[]): Promise<void> {
    const payload: SyncPermissionsRequest = {
      permission_ids: permissionIds,
    }
    await HttpClient.put<void>(`/role-permissions/roles/${roleId}/sync`, payload)
  }

  /**
   * Sync role users (replace all users)
   */
  static async syncUsers(roleId: string, userIds: string[]): Promise<void> {
    await HttpClient.put<void>(`/user-roles/roles/${roleId}/sync`, { user_ids: userIds })
  }

  /**
   * Delete role (logical)
   */
  static async delete(id: string): Promise<void> {
    await HttpClient.delete<void>(`${this.ENDPOINT}/${id}`)
  }

  /**
   * Restore role (undo logical delete)
   */
  static async restore(id: string): Promise<Role> {
    const result = await HttpClient.post<Role>(`${this.ENDPOINT}/${id}/restore`)
    if (!result.data) throw new Error('Error al restaurar rol')
    return result.data
  }
}

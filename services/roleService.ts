// MisViáticos - Role Service

import { HttpClient } from '@/lib/api/httpClient'
import { sanitizeInput } from '@/lib/utils/sanitize'
import type {
  Role,
  CreateRoleRequest,
  UpdateRoleRequest,
  RolesResponse,
  Permission,
  PermissionsResponse,
} from '@/types'

export class RoleService {
  private static readonly ENDPOINT = '/roles'
  private static readonly PERMISSIONS_ENDPOINT = '/permissions'

  /**
   * Get all roles
   */
  static async getAll(): Promise<Role[]> {
    const result = await HttpClient.get<RolesResponse>(this.ENDPOINT)
    return result.data?.roles || []
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
   * Sync role permissions (replace all permissions)
   */
  static async syncPermissions(roleId: string, permissionIds: string[]): Promise<void> {
    await HttpClient.put<void>(`/role-permissions/roles/${roleId}/sync`, {
      permission_ids: permissionIds,
    })
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

  /**
   * Get all permissions
   */
  static async getAllPermissions(): Promise<Permission[]> {
    const result = await HttpClient.get<PermissionsResponse>(this.PERMISSIONS_ENDPOINT)
    return result.data?.permissions || []
  }
}

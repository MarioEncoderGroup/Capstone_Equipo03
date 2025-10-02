// MisViáticos - User Service

import { HttpClient } from '@/lib/api/httpClient'
import { sanitizeInput, sanitizeEmail, sanitizePhone } from '@/lib/utils/sanitize'
import type {
  User,
  CreateUserRequest,
  UpdateUserRequest,
  UsersResponse,
} from '@/app/dashboard/users/types'

export class UserService {
  private static readonly ENDPOINT = '/users'

  /**
   * Get all users
   * Backend returns: { success, message, data: User[], pagination }
   */
  static async getAll(): Promise<User[]> {
    const result = await HttpClient.get<User[]>(`/admin${this.ENDPOINT}`)
    console.log('👥 UserService.getAll() response:', result)
    // El backend devuelve el array directamente en data
    return result.data || []
  }

  /**
   * Get user by ID
   */
  static async getById(id: string): Promise<User> {
    const result = await HttpClient.get<User>(`/admin${this.ENDPOINT}/${id}`)
    if (!result.data) throw new Error('Usuario no encontrado')
    return result.data
  }

  /**
   * Create user
   */
  static async create(data: CreateUserRequest): Promise<User> {
    const sanitizedData: CreateUserRequest = {
      full_name: sanitizeInput(data.full_name, 150),
      email: sanitizeEmail(data.email),
      phone: sanitizePhone(data.phone),
      password: data.password, // No sanitizar password
      role_ids: data.role_ids,
      permission_ids: data.permission_ids,
    }

    const result = await HttpClient.post<User>(`${this.ENDPOINT}`, sanitizedData)
    if (!result.data) throw new Error('Error al crear usuario')
    return result.data
  }

  /**
   * Update user
   */
  static async update(id: string, data: UpdateUserRequest): Promise<User> {
    const sanitizedData: UpdateUserRequest = {
      ...(data.full_name && { full_name: sanitizeInput(data.full_name, 150) }),
      ...(data.email && { email: sanitizeEmail(data.email) }),
      ...(data.phone && { phone: sanitizePhone(data.phone) }),
      ...(data.password && { password: data.password }),
      ...(data.is_active !== undefined && { is_active: data.is_active }),
    }

    const result = await HttpClient.put<User>(`/admin${this.ENDPOINT}/${id}`, sanitizedData)
    if (!result.data) throw new Error('Error al actualizar usuario')
    return result.data
  }

  /**
   * Assign a single role to a user
   */
  static async assignRole(userId: string, roleId: string): Promise<void> {
    await HttpClient.post<void>(`/user-roles/assign`, {
      user_id: userId,
      role_id: roleId,
    })
  }

  /**
   * Assign multiple roles to a user (using Promise.all for parallel execution)
   */
  static async assignRoles(userId: string, roleIds: string[]): Promise<void> {
    if (roleIds.length === 0) return

    await Promise.all(roleIds.map((roleId) => this.assignRole(userId, roleId)))
  }

  /**
   * Sync user roles (replace all roles)
   */
  static async syncRoles(userId: string, roleIds: string[]): Promise<void> {
    await HttpClient.put<void>(`/user-roles/users/${userId}/sync`, { role_ids: roleIds })
  }

  /**
   * Delete user (logical)
   */
  static async delete(id: string): Promise<void> {
    await HttpClient.delete<void>(`/admin${this.ENDPOINT}/${id}`)
  }

  /**
   * Restore user (undo logical delete)
   */
  static async restore(id: string): Promise<User> {
    const result = await HttpClient.post<User>(`${this.ENDPOINT}/${id}/restore`)
    if (!result.data) throw new Error('Error al restaurar usuario')
    return result.data
  }
}

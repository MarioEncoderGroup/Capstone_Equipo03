// MisViáticos - Users Feature Types

import type { Role, Permission } from '@/types'

/**
 * User entity
 */
export interface User {
  id: string
  full_name: string
  email: string
  phone?: string
  is_active: boolean
  tenant_id?: string
  roles?: Role[]
  permissions?: Permission[]
  last_login?: string
  created_at?: string
  updated_at?: string
  deleted_at?: string | null
}

/**
 * Request to create a new user
 */
export interface CreateUserRequest {
  full_name: string
  email: string
  phone: string
  password: string
  role_ids: string[]
  permission_ids: string[]
}

/**
 * Request to update an existing user
 */
export interface UpdateUserRequest {
  full_name?: string
  email?: string
  phone?: string
  password?: string
  is_active?: boolean
  role_ids?: string[]
  permission_ids?: string[]
}

/**
 * Response with users list
 */
export interface UsersResponse {
  users: User[]
  total: number
  page: number
  limit: number
}

/**
 * Form data for user modal
 */
export interface UserFormData {
  full_name: string
  email: string
  phone: string
  password: string
  password_confirm: string
  is_active: boolean
  role_ids: string[]
  permission_ids: string[]
}

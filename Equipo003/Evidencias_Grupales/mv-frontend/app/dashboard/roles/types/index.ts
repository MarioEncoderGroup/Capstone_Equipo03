// MisViáticos - Roles Feature Types

import type { Permission } from '@/types'

/**
 * Role with users assigned
 */
export interface RoleWithUsers {
  id: string
  name: string
  description: string
  tenant_id: string
  permissions: Permission[]
  users: RoleUser[]
  created_at?: string
  updated_at?: string
  deleted_at?: string | null
}

/**
 * User assigned to a role (simplified)
 */
export interface RoleUser {
  id: string
  full_name: string
  email: string
}

/**
 * Form data for role modal
 */
export interface RoleFormData {
  name: string
  description: string
  permission_ids: string[]
  user_ids: string[]
}

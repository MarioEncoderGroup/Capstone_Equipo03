// MisViáticos - Role Types (Shared)

import type { Permission } from './permission'

export interface Role {
  id: string
  name: string
  description: string
  tenant_id: string
  permissions: Permission[]
  created_at?: string
  updated_at?: string
  deleted_at?: string | null
}

export interface CreateRoleRequest {
  name: string
  description: string
  permission_ids: string[]
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
  permission_ids?: string[]
  user_ids?: string[]
}

export interface RolesResponse {
  roles: Role[]
  total: number
}

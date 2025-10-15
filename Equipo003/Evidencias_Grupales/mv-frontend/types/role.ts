// MisViáticos - Role Types (Shared)

import type { Permission } from './permission'

export interface Role {
  id: string
  name: string
  description: string
  tenant_id: string | null
  permissions?: Permission[]
  created?: string
  updated?: string
  deleted_at?: string | null
}

export interface CreateRoleRequest {
  name: string
  description: string
  permission_ids?: string[]
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
}

export interface RolesResponse {
  roles: Role[]
  pagination?: {
    page: number
    page_size: number
    total_pages: number
    total_records: number
    has_next: boolean
    has_previous: boolean
  }
}

export interface SyncPermissionsRequest {
  permission_ids: string[]
}

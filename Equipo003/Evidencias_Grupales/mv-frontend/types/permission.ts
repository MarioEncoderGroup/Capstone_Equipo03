// MisViáticos - Permission Types

export interface Permission {
  id: string
  name: string
  description: string
  section: string
  created?: string
  updated?: string
}

export interface PermissionsResponse {
  permissions: Permission[]
  pagination?: {
    page: number
    page_size: number
    total_pages: number
    total_records: number
    has_next: boolean
    has_previous: boolean
  }
}

export interface GroupedPermission {
  section: string
  permissions: Permission[]
}

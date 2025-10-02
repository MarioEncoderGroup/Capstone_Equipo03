// MisViáticos - Permission Types

export interface Permission {
  id: string
  name: string
  description: string
  resource: string
  action: string
  created_at?: string
}

export interface PermissionsResponse {
  permissions: Permission[]
  total: number
}

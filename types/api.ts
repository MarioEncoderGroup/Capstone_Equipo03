// MisViáticos - API Types

// ========== COMMON ==========
export interface ApiResponse<T> {
  success: boolean
  message: string
  data?: T
  error?: string
}

export interface ValidationError {
  Field: string
  Message: string
}

export interface ApiErrorResponse {
  success: false
  message: string
  data?: ValidationError[]
  error?: string
}

// ========== AUTH ==========
export interface RegisterRequest {
  full_name: string
  email: string
  phone: string
  password: string
  password_confirm: string
}

export interface RegisterResponse {
  id: string
  full_name: string
  email: string
  phone: string
  email_token: string
  requires_email_verification: boolean
}

export interface VerifyEmailRequest {
  token: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: 'Bearer'
  user: UserData
}

export interface UserData {
  id: string
  username: string
  full_name: string
  email: string
  is_active: boolean
  last_login?: string
}

// ========== TENANT ==========
export interface TenantData {
  id: string
  rut: string
  business_name: string
  email: string
  phone: string
  address: string
  region_id: string
  commune_id: string
  country_id: string
  status: 'active' | 'inactive'
  node_number: number
  tenant_name: string
  created_by?: string
  created?: string
}

export interface TenantStatusResponse {
  has_tenants: boolean
  tenants: TenantData[]
  requires_tenant_creation: boolean
  tenant_count: number
}

export interface CreateTenantRequest {
  rut: string
  business_name: string
  email: string
  phone: string
  address: string
  website: string
  region_id: string
  commune_id: string
  country_id: string
  logo?: string
}

export interface SelectTenantResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserData
  tenant: {
    id: string
    rut: string
    business_name: string
    status: string
  }
}

// ========== REGIONS & COMMUNES ==========
export interface Region {
  id: string
  number: number
  roman_number: string
  name: string
}

export interface RegionsResponse {
  regions: Region[]
  total: number
}

export interface Commune {
  id: string
  region_id: string
  name: string
}

export interface CommunesResponse {
  communes: Commune[]
}

// ========== JWT ==========
export interface JWTPayload {
  user_id: string
  tenant_id?: string
  type: 'access' | 'refresh'
  iat: number
  exp: number
  iss: string
}

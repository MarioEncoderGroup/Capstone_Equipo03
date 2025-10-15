// MisViáticos - API Constants
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

export const API_ENDPOINTS = {
  // Auth
  AUTH_REGISTER: '/auth/register',
  AUTH_LOGIN: '/auth/login',
  AUTH_VERIFY_EMAIL: '/auth/verify-email',
  AUTH_VERIFY_EMAIL_GET: (token: string) => `/auth/verify-email/${token}`,

  // Tenant
  TENANT_STATUS: '/tenant/status',
  TENANT_CREATE: '/tenant/create',
  TENANT_SELECT: (tenantId: string) => `/tenant/select/${tenantId}`,

  // Regions & Communes
  REGIONS: '/regions',
  COMMUNES: (regionId?: string) => regionId ? `/communes?region_id=${regionId}` : '/communes',
} as const

export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  INTERNAL_SERVER_ERROR: 500,
} as const

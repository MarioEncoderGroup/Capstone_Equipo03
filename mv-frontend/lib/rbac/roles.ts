// MisViáticos - Role-Based Access Control (RBAC) Constants

/**
 * System roles
 * Define los roles disponibles en el sistema
 */
export const ROLES = {
  ADMINISTRATOR: 'administrator',
  APPROVER: 'approver',
  EMPLOYEE: 'employee',
  ACCOUNTANT: 'accountant',
} as const

/**
 * Role type
 */
export type Role = (typeof ROLES)[keyof typeof ROLES]

/**
 * Role display names
 * Mapeo de roles a nombres legibles en español
 */
export const ROLE_NAMES: Record<Role, string> = {
  [ROLES.ADMINISTRATOR]: 'Administrador',
  [ROLES.APPROVER]: 'Aprobador',
  [ROLES.EMPLOYEE]: 'Empleado',
  [ROLES.ACCOUNTANT]: 'Contador',
}

/**
 * Verifica si un usuario tiene un rol específico
 */
export function hasRole(userRoles: string[] | undefined, role: Role): boolean {
  if (!userRoles || userRoles.length === 0) return false
  return userRoles.includes(role)
}

/**
 * Verifica si un usuario tiene al menos uno de los roles especificados
 */
export function hasAnyRole(userRoles: string[] | undefined, roles: Role[]): boolean {
  if (!userRoles || userRoles.length === 0) return false
  return roles.some((role) => userRoles.includes(role))
}

/**
 * Verifica si un usuario tiene todos los roles especificados
 */
export function hasAllRoles(userRoles: string[] | undefined, roles: Role[]): boolean {
  if (!userRoles || userRoles.length === 0) return false
  return roles.every((role) => userRoles.includes(role))
}

/**
 * Verifica si un usuario es administrador
 */
export function isAdmin(userRoles: string[] | undefined): boolean {
  return hasRole(userRoles, ROLES.ADMINISTRATOR)
}

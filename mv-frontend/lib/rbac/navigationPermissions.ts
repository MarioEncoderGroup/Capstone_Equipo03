// MisViáticos - Navigation Permissions
// Define qué rutas/páginas son accesibles según el rol del usuario

import { ROLES, type Role } from './roles'

/**
 * Configuración de permisos para cada ruta del sidebar
 * Si allowedRoles está vacío, la ruta es pública para todos los usuarios autenticados
 */
export const NAVIGATION_PERMISSIONS = {
  DASHBOARD: {
    path: '/dashboard',
    allowedRoles: [], // Todos los usuarios autenticados
  },
  USERS: {
    path: '/dashboard/users',
    allowedRoles: [ROLES.ADMINISTRATOR], // Solo administradores
  },
  ROLES: {
    path: '/dashboard/roles',
    allowedRoles: [ROLES.ADMINISTRATOR], // Solo administradores
  },
  SETTINGS: {
    path: '/dashboard/settings',
    allowedRoles: [], // Todos los usuarios autenticados
  },
} as const

/**
 * Verifica si un usuario puede acceder a una ruta
 *
 * @param path - Ruta a verificar
 * @param userRoles - Roles del usuario
 * @returns true si el usuario puede acceder, false si no
 */
export function canAccessRoute(path: string, userRoles: string[] | undefined): boolean {
  // Si no hay roles definidos, no puede acceder
  if (!userRoles || userRoles.length === 0) {
    return false
  }

  // Buscar la configuración de permisos para esta ruta
  const permission = Object.values(NAVIGATION_PERMISSIONS).find((p) => p.path === path)

  // Si no hay configuración, denegar por defecto
  if (!permission) {
    return false
  }

  // Si no hay roles requeridos, todos los autenticados pueden acceder
  if (permission.allowedRoles.length === 0) {
    return true
  }

  // Verificar si el usuario tiene alguno de los roles permitidos
  return permission.allowedRoles.some((role) => userRoles.includes(role))
}

/**
 * Filtra los items de navegación según los roles del usuario
 *
 * @param items - Items de navegación con sus roles permitidos
 * @param userRoles - Roles del usuario
 * @returns Items filtrados que el usuario puede ver
 */
export function filterNavigationItems<T extends { href: string; allowedRoles?: Role[] }>(
  items: T[],
  userRoles: string[] | undefined
): T[] {
  if (!userRoles || userRoles.length === 0) {
    return []
  }

  return items.filter((item) => {
    // Si no tiene roles definidos, es accesible para todos
    if (!item.allowedRoles || item.allowedRoles.length === 0) {
      return true
    }

    // Verificar si el usuario tiene alguno de los roles permitidos
    return item.allowedRoles.some((role) => userRoles.includes(role))
  })
}

// MisViáticos - User Roles Hook
'use client'

import { useState, useEffect } from 'react'
import { TokenManager } from '@/lib/auth/tokenManager'
import type { Role } from '@/lib/rbac/roles'

/**
 * Hook para obtener los roles del usuario desde el JWT
 *
 * @returns {Object} roles - Array de roles del usuario
 * @returns {boolean} isLoading - Estado de carga
 * @returns {boolean} isAdmin - Si el usuario es administrador
 */
export function useUserRoles() {
  const [roles, setRoles] = useState<string[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const loadRoles = () => {
      try {
        const token = TokenManager.getAccessToken()

        if (!token) {
          setRoles([])
          setIsLoading(false)
          return
        }

        const payload = TokenManager.decodeJWT(token)

        if (payload && payload.roles) {
          setRoles(payload.roles)
        } else {
          setRoles([])
        }
      } catch (error) {
        console.error('Error loading user roles:', error)
        setRoles([])
      } finally {
        setIsLoading(false)
      }
    }

    loadRoles()

    // Escuchar cambios en localStorage (login/logout/tenant selection)
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'auth_token') {
        loadRoles()
      }
    }

    window.addEventListener('storage', handleStorageChange)
    return () => window.removeEventListener('storage', handleStorageChange)
  }, [])

  const isAdmin = roles.includes('administrator')

  return {
    roles,
    isLoading,
    isAdmin,
  }
}

/**
 * Hook para verificar si el usuario tiene un rol específico
 */
export function useHasRole(role: Role): boolean {
  const { roles } = useUserRoles()
  return roles.includes(role)
}

/**
 * Hook para verificar si el usuario tiene al menos uno de los roles
 */
export function useHasAnyRole(requiredRoles: Role[]): boolean {
  const { roles } = useUserRoles()
  return requiredRoles.some((role) => roles.includes(role))
}

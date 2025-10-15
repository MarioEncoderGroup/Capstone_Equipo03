// MisViáticos - Role Guard Component
// Protege rutas según el rol del usuario
'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useUserRoles } from '@/hooks/useUserRoles'
import { hasAnyRole, type Role } from '@/lib/rbac/roles'

interface RoleGuardProps {
  children: React.ReactNode
  allowedRoles: Role[]
  fallbackUrl?: string
  showAccessDenied?: boolean
}

/**
 * Component wrapper para proteger rutas según rol del usuario
 *
 * @param children - Contenido a proteger
 * @param allowedRoles - Roles permitidos para acceder
 * @param fallbackUrl - URL a redirigir si no tiene permisos (default: /dashboard)
 * @param showAccessDenied - Si mostrar mensaje de acceso denegado antes de redirigir
 */
export function RoleGuard({
  children,
  allowedRoles,
  fallbackUrl = '/dashboard',
  showAccessDenied = false,
}: RoleGuardProps) {
  const router = useRouter()
  const { roles, isLoading } = useUserRoles()

  const hasAccess = hasAnyRole(roles, allowedRoles)

  useEffect(() => {
    if (!isLoading && !hasAccess) {
      console.warn(
        `🚫 Acceso denegado. Roles requeridos: ${allowedRoles.join(', ')}. Roles del usuario: ${roles.join(', ') || 'ninguno'}`
      )

      if (showAccessDenied) {
        alert('No tienes permisos para acceder a esta página')
      }

      router.push(fallbackUrl)
    }
  }, [isLoading, hasAccess, allowedRoles, roles, router, fallbackUrl, showAccessDenied])

  // Mostrar loading mientras verifica roles
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600" />
      </div>
    )
  }

  // Si no tiene acceso, no renderizar nada (ya se está redirigiendo)
  if (!hasAccess) {
    return null
  }

  // Si tiene acceso, renderizar children
  return <>{children}</>
}

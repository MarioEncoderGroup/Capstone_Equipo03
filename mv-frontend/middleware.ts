// MisViáticos - Next.js Middleware (Auth Protection)

import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Rutas que requieren autenticación
const protectedRoutes = ['/dashboard', '/expenses', '/reports', '/settings']

// Rutas que requieren tenant_id en el token
const tenantRequiredRoutes = ['/dashboard', '/expenses', '/reports']

// Rutas de auth que redirigen si ya está autenticado
const authRoutes = ['/auth/login', '/auth/register']

/**
 * Decodifica un JWT básico (solo para leer, NO para validar)
 * La validación real se hace en el backend
 */
function decodeJWT(token: string): any {
  try {
    const base64Url = token.split('.')[1]
    if (!base64Url) return null

    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )

    return JSON.parse(jsonPayload)
  } catch {
    return null
  }
}

/**
 * Verifica si un token ha expirado
 */
function isTokenExpired(payload: any): boolean {
  if (!payload || !payload.exp) return true
  return payload.exp * 1000 < Date.now()
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl

  // Obtener token de las cookies o del header (Next.js usa cookies por defecto)
  const token = request.cookies.get('auth_token')?.value

  // Verificar rutas protegidas
  const isProtectedRoute = protectedRoutes.some((route) =>
    pathname.startsWith(route)
  )
  const isAuthRoute = authRoutes.some((route) => pathname.startsWith(route))
  const isTenantRequiredRoute = tenantRequiredRoutes.some((route) =>
    pathname.startsWith(route)
  )

  // Si es ruta protegida y no hay token, redirigir a login
  if (isProtectedRoute && !token) {
    const loginUrl = new URL('/auth/login', request.url)
    loginUrl.searchParams.set('from', pathname)
    return NextResponse.redirect(loginUrl)
  }

  // Si hay token, verificar su validez
  if (token) {
    const payload = decodeJWT(token)

    // Token inválido o expirado
    if (!payload || isTokenExpired(payload)) {
      // Limpiar cookie solo si NO está en una ruta de auth
      // (para evitar loops infinitos)
      if (!isAuthRoute) {
        const response = NextResponse.redirect(new URL('/auth/login', request.url))
        response.cookies.delete('auth_token')
        response.cookies.delete('refresh_token')
        return response
      }
    }

    // Verificar tenant_id en rutas que lo requieren
    if (isTenantRequiredRoute && !payload.tenant_id) {
      // Tiene token pero sin tenant_id - redirigir a selección de tenant
      // El componente TenantSelect verificará el status y redirigirá apropiadamente
      return NextResponse.redirect(new URL('/tenant/select', request.url))
    }

    // NO redirigir automáticamente desde rutas de auth
    // Dejar que los componentes manejen la lógica de redirección
    // Esto permite que useAuth verifique el tenant status correctamente
  }

  // Continuar con la request normalmente
  return NextResponse.next()
}

// Configurar qué rutas procesará el middleware
export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    '/((?!api|_next/static|_next/image|favicon.ico|icon-mv|.*\\.svg|.*\\.png|.*\\.jpg).*)',
  ],
}

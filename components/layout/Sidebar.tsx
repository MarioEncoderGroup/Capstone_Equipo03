// MisViáticos - Sidebar Component
'use client'

import Link from 'next/link'
import Image from 'next/image'
import { usePathname } from 'next/navigation'
import {
  HomeIcon,
  UserGroupIcon,
  ShieldCheckIcon,
  Cog6ToothIcon,
  ArrowRightOnRectangleIcon,
} from '@heroicons/react/24/outline'
import { useAuth } from '@/hooks/useAuth'
import { useUserRoles } from '@/hooks/useUserRoles'
import { filterNavigationItems } from '@/lib/rbac/navigationPermissions'
import { ROLES, type Role } from '@/lib/rbac/roles'

interface SidebarItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
  allowedRoles?: Role[]
}

const navigationItems: SidebarItem[] = [
  { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
  {
    name: 'Usuarios',
    href: '/dashboard/users',
    icon: UserGroupIcon,
    allowedRoles: [ROLES.ADMINISTRATOR], // Solo administradores
  },
  {
    name: 'Roles',
    href: '/dashboard/roles',
    icon: ShieldCheckIcon,
    allowedRoles: [ROLES.ADMINISTRATOR], // Solo administradores
  },
  { name: 'Configuración', href: '/dashboard/settings', icon: Cog6ToothIcon },
]

export function Sidebar() {
  const pathname = usePathname()
  const { logout } = useAuth()
  const { roles } = useUserRoles()

  // Filtrar items de navegación según roles del usuario
  const visibleItems = filterNavigationItems(navigationItems, roles)

  return (
    <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
      <div className="flex flex-col flex-grow bg-white border-r border-gray-200 overflow-y-auto">
        {/* Logo */}
        <div className="py-4">
          <Image
            src="/icon-mv/Assets-MV_Logo2.png"
            alt="MisViáticos"
            width={14400}
            height={4800}
            className="w-auto h-[30px] object-contain block mx-auto leading-none"
          />
        </div>

        {/* Navigation */}
        <nav className="mt-8 flex-1 px-2 space-y-1">
          {visibleItems.map((item) => {
            const isActive = pathname === item.href
            const Icon = item.icon
            const isConfigItem = item.name === 'Configuración'

            return (
              <div key={item.name}>
                {isConfigItem && (
                  <div className="border-t border-gray-300 my-4 mx-2"></div>
                )}
                <Link
                  href={item.href}
                  className={`
                    group flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors
                    ${
                      isActive
                        ? 'bg-purple-600 text-white'
                        : 'text-gray-700 hover:bg-purple-50 hover:text-purple-600'
                    }
                  `}
                >
                  <Icon
                    className={`
                      mr-3 flex-shrink-0 h-6 w-6
                      ${isActive ? 'text-white' : 'text-gray-500 group-hover:text-purple-600'}
                    `}
                  />
                  {item.name}
                </Link>
              </div>
            )
          })}
        </nav>

        {/* Logout */}
        <div className="flex-shrink-0 px-2 pb-4">
          <button
            onClick={logout}
            className="group flex items-center w-full px-3 py-2 text-sm font-medium text-gray-700 rounded-md hover:bg-purple-50 hover:text-purple-600 transition-colors"
          >
            <ArrowRightOnRectangleIcon className="mr-3 flex-shrink-0 h-6 w-6 text-gray-500 group-hover:text-purple-600" />
            Cerrar Sesión
          </button>
        </div>
      </div>
    </div>
  )
}

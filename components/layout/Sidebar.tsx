// MisViáticos - Sidebar Component
'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  HomeIcon,
  UserGroupIcon,
  ShieldCheckIcon,
  Cog6ToothIcon,
  ArrowRightOnRectangleIcon,
  BuildingOfficeIcon,
} from '@heroicons/react/24/outline'
import { useAuth } from '@/hooks/useAuth'

interface SidebarItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}

const navigationItems: SidebarItem[] = [
  { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
  { name: 'Usuarios', href: '/dashboard/users', icon: UserGroupIcon },
  { name: 'Roles', href: '/dashboard/roles', icon: ShieldCheckIcon },
  { name: 'Configuración', href: '/dashboard/settings', icon: Cog6ToothIcon },
]

export function Sidebar() {
  const pathname = usePathname()
  const { logout } = useAuth()

  return (
    <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
      <div className="flex flex-col flex-grow bg-purple-700 pt-5 pb-4 overflow-y-auto">
        {/* Logo */}
        <div className="flex items-center flex-shrink-0 px-4">
          <BuildingOfficeIcon className="h-8 w-8 text-white" />
          <span className="ml-2 text-white text-xl font-bold">MisViáticos</span>
        </div>

        {/* Navigation */}
        <nav className="mt-8 flex-1 px-2 space-y-1">
          {navigationItems.map((item) => {
            const isActive = pathname === item.href
            const Icon = item.icon

            return (
              <Link
                key={item.name}
                href={item.href}
                className={`
                  group flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors
                  ${
                    isActive
                      ? 'bg-purple-800 text-white'
                      : 'text-purple-100 hover:bg-purple-600 hover:text-white'
                  }
                `}
              >
                <Icon
                  className={`
                    mr-3 flex-shrink-0 h-6 w-6
                    ${isActive ? 'text-white' : 'text-purple-300 group-hover:text-white'}
                  `}
                />
                {item.name}
              </Link>
            )
          })}
        </nav>

        {/* Logout */}
        <div className="flex-shrink-0 px-2">
          <button
            onClick={logout}
            className="group flex items-center w-full px-3 py-2 text-sm font-medium text-purple-100 rounded-md hover:bg-purple-600 hover:text-white transition-colors"
          >
            <ArrowRightOnRectangleIcon className="mr-3 flex-shrink-0 h-6 w-6 text-purple-300 group-hover:text-white" />
            Cerrar Sesión
          </button>
        </div>
      </div>
    </div>
  )
}

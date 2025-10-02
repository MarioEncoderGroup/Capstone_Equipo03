// MisViáticos - Users Management Page
'use client'

import { useState, useEffect, useRef } from 'react'
import { PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'
import { UserService } from '@/services/userService'
import { RoleService } from '@/services/roleService'
import { Button } from '@/components/ui/Button'
import { Alert } from '@/components/ui/Alert'
import { DeleteConfirmModal } from '@/components/shared/DeleteConfirmModal'
import { UserModal } from './components/UserModal'
import type { Role, Permission } from '@/types'
import type { User } from './types'

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [roles, setRoles] = useState<Role[]>([])
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Modal states
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const [userToDelete, setUserToDelete] = useState<User | null>(null)

  // Ref para evitar múltiples llamadas simultáneas
  const isLoadingRef = useRef(false)
  const abortControllerRef = useRef<AbortController | null>(null)

  useEffect(() => {
    // Cleanup function para cancelar requests en progreso
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, [])

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    // Prevenir múltiples llamadas simultáneas
    if (isLoadingRef.current) {
      console.log('⏳ Ya hay una carga en progreso, ignorando...')
      return
    }

    // Cancelar request anterior si existe
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    isLoadingRef.current = true
    setIsLoading(true)
    setError(null)

    try {
      console.log('🔄 Cargando datos de usuarios, roles y permisos...')
      
      const [usersData, rolesData, permissionsData] = await Promise.all([
        UserService.getAll(),
        RoleService.getAll(),
        RoleService.getAllPermissions(),
      ])

      setUsers(usersData)
      setRoles(rolesData)
      setPermissions(permissionsData)
      console.log('✅ Datos cargados exitosamente')
    } catch (err) {
      // Ignorar errores de abort
      if (err instanceof Error && err.name === 'AbortError') {
        console.log('🚫 Request cancelado')
        return
      }
      
      console.error('❌ Error cargando datos:', err)
      setError(err instanceof Error ? err.message : 'Error al cargar datos')
    } finally {
      setIsLoading(false)
      isLoadingRef.current = false
    }
  }

  const handleCreate = () => {
    setSelectedUser(null)
    setIsModalOpen(true)
  }

  const handleEdit = (user: User) => {
    setSelectedUser(user)
    setIsModalOpen(true)
  }

  const handleDelete = (user: User) => {
    setUserToDelete(user)
    setIsDeleteModalOpen(true)
  }

  const confirmDelete = async () => {
    if (!userToDelete) return

    try {
      await UserService.delete(userToDelete.id)
      await loadData()
      setIsDeleteModalOpen(false)
      setUserToDelete(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al eliminar usuario')
    }
  }

  const handleModalClose = async (success: boolean) => {
    setIsModalOpen(false)
    setSelectedUser(null)

    if (success) {
      await loadData()
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Usuarios</h1>
          <p className="mt-1 text-sm text-gray-600">
            Gestiona los usuarios del sistema
          </p>
        </div>
        <Button variant="primary" onClick={handleCreate}>
          <PlusIcon className="h-5 w-5 mr-2" />
          Nuevo Usuario
        </Button>
      </div>

      {/* Error Alert */}
      {error && (
        <Alert variant="error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Users Table */}
      <div className="bg-white shadow-md rounded-lg overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Usuario
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Email
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Teléfono
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Roles
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Estado
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {users.map((user) => (
              <tr key={user.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">
                    {user.full_name}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm text-gray-500">{user.email}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm text-gray-500">{user.phone}</div>
                </td>
                <td className="px-6 py-4">
                  <div className="flex flex-wrap gap-1">
                    {user.roles.map((role) => (
                      <span
                        key={role.id}
                        className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800"
                      >
                        {role.name}
                      </span>
                    ))}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      user.is_active
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {user.is_active ? 'Activo' : 'Inactivo'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    onClick={() => handleEdit(user)}
                    className="text-purple-600 hover:text-purple-900 mr-4"
                  >
                    <PencilIcon className="h-5 w-5" />
                  </button>
                  <button
                    onClick={() => handleDelete(user)}
                    className="text-red-600 hover:text-red-900"
                  >
                    <TrashIcon className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {users.length === 0 && (
          <div className="text-center py-12">
            <p className="text-gray-500">No hay usuarios registrados</p>
          </div>
        )}
      </div>

      {/* Modals */}
      {isModalOpen && (
        <UserModal
          user={selectedUser}
          roles={roles}
          permissions={permissions}
          onClose={handleModalClose}
        />
      )}

      {isDeleteModalOpen && userToDelete && (
        <DeleteConfirmModal
          title="Eliminar Usuario"
          message={`¿Estás seguro de que deseas eliminar a ${userToDelete.full_name}? Esta acción se puede revertir.`}
          onConfirm={confirmDelete}
          onCancel={() => {
            setIsDeleteModalOpen(false)
            setUserToDelete(null)
          }}
        />
      )}
    </div>
  )
}

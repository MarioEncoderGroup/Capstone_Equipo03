// MisViáticos - Roles Management Page
'use client'

import { useState, useEffect } from 'react'
import { PlusIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline'
import { RoleService } from '@/services/roleService'
import { UserService } from '@/services/userService'
import { Button } from '@/components/ui/Button'
import { Alert } from '@/components/ui/Alert'
import { DeleteConfirmModal } from '@/components/shared/DeleteConfirmModal'
import { RoleModal } from './components/RoleModal'
import type { Role, Permission } from '@/types'
import type { User } from '@/app/dashboard/users/types'

export default function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([])
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Modal states
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [selectedRole, setSelectedRole] = useState<Role | null>(null)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const [roleToDelete, setRoleToDelete] = useState<Role | null>(null)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const [rolesData, permissionsData, usersData] = await Promise.all([
        RoleService.getAll(),
        RoleService.getAllPermissions(),
        UserService.getAll(),
      ])

      setRoles(rolesData)
      setPermissions(permissionsData)
      setUsers(usersData)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al cargar datos')
    } finally {
      setIsLoading(false)
    }
  }

  const handleCreate = () => {
    setSelectedRole(null)
    setIsModalOpen(true)
  }

  const handleEdit = (role: Role) => {
    setSelectedRole(role)
    setIsModalOpen(true)
  }

  const handleDelete = (role: Role) => {
    setRoleToDelete(role)
    setIsDeleteModalOpen(true)
  }

  const confirmDelete = async () => {
    if (!roleToDelete) return

    try {
      await RoleService.delete(roleToDelete.id)
      await loadData()
      setIsDeleteModalOpen(false)
      setRoleToDelete(null)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al eliminar rol')
    }
  }

  const handleModalClose = async (success: boolean) => {
    setIsModalOpen(false)
    setSelectedRole(null)

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
          <h1 className="text-3xl font-bold text-gray-900">Roles</h1>
          <p className="mt-1 text-sm text-gray-600">
            Gestiona los roles y permisos del sistema
          </p>
        </div>
        <Button variant="primary" onClick={handleCreate}>
          <PlusIcon className="h-5 w-5 mr-2" />
          Nuevo Rol
        </Button>
      </div>

      {/* Error Alert */}
      {error && (
        <Alert variant="error" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Roles Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {roles.map((role) => (
          <div
            key={role.id}
            className="bg-white shadow-md rounded-lg p-6 hover:shadow-lg transition-shadow"
          >
            {/* Role Header */}
            <div className="flex items-start justify-between mb-4">
              <div>
                <h3 className="text-lg font-bold text-gray-900">{role.name}</h3>
                <p className="text-sm text-gray-500 mt-1">{role.description}</p>
              </div>
              <div className="flex space-x-2">
                <button
                  onClick={() => handleEdit(role)}
                  className="text-purple-600 hover:text-purple-900"
                >
                  <PencilIcon className="h-5 w-5" />
                </button>
                <button
                  onClick={() => handleDelete(role)}
                  className="text-red-600 hover:text-red-900"
                >
                  <TrashIcon className="h-5 w-5" />
                </button>
              </div>
            </div>

            {/* Permissions */}
            <div className="border-t pt-4">
              <h4 className="text-sm font-medium text-gray-700 mb-2">
                Permisos ({role.permissions.length})
              </h4>
              <div className="flex flex-wrap gap-2">
                {role.permissions.slice(0, 3).map((permission) => (
                  <span
                    key={permission.id}
                    className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800"
                  >
                    {permission.name}
                  </span>
                ))}
                {role.permissions.length > 3 && (
                  <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600">
                    +{role.permissions.length - 3} más
                  </span>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      {roles.length === 0 && (
        <div className="text-center py-12 bg-white rounded-lg">
          <p className="text-gray-500">No hay roles registrados</p>
        </div>
      )}

      {/* Modals */}
      {isModalOpen && (
        <RoleModal
          role={selectedRole}
          permissions={permissions}
          users={users}
          onClose={handleModalClose}
        />
      )}

      {isDeleteModalOpen && roleToDelete && (
        <DeleteConfirmModal
          title="Eliminar Rol"
          message={`¿Estás seguro de que deseas eliminar el rol "${roleToDelete.name}"? Esta acción se puede revertir.`}
          onConfirm={confirmDelete}
          onCancel={() => {
            setIsDeleteModalOpen(false)
            setRoleToDelete(null)
          }}
        />
      )}
    </div>
  )
}

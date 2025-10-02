// MisViáticos - Role Modal Component
'use client'

import { useState } from 'react'
import { XMarkIcon } from '@heroicons/react/24/outline'
import { RoleService } from '@/services/roleService'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Alert } from '@/components/ui/Alert'
import type { Role, GroupedPermission } from '@/types'
import type { User } from '@/app/dashboard/users/types'
import type { RoleFormData } from '../types'
import { validateRoleForm } from '../utils/validation'

interface RoleModalProps {
  role: Role | null
  groupedPermissions: GroupedPermission[]
  users: User[]
  onClose: (success: boolean) => void
}

export function RoleModal({ role, groupedPermissions, users, onClose }: RoleModalProps) {
  const isEdit = !!role

  const [formData, setFormData] = useState<RoleFormData>({
    name: role?.name || '',
    description: role?.description || '',
    permission_ids: role?.permissions?.map(p => p.id) || [],
    user_ids: [],
  })

  const [errors, setErrors] = useState<{ name?: string; description?: string }>({})
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const validate = (): boolean => {
    const validationErrors = validateRoleForm(formData)
    setErrors(validationErrors)
    return Object.keys(validationErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) return

    setIsLoading(true)
    setError(null)

    try {
      if (isEdit && role) {
        // Update role basic info
        await RoleService.update(role.id, {
          name: formData.name,
          description: formData.description,
        })

        // Sync permissions
        await RoleService.syncPermissions(role.id, formData.permission_ids)

        // Sync users if any selected
        if (formData.user_ids.length > 0) {
          await RoleService.syncUsers(role.id, formData.user_ids)
        }
      } else {
        // Create role
        const newRole = await RoleService.create({
          name: formData.name,
          description: formData.description,
          permission_ids: formData.permission_ids,
        })

        // Sync permissions after creation
        if (formData.permission_ids.length > 0) {
          await RoleService.syncPermissions(newRole.id, formData.permission_ids)
        }
      }

      onClose(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar rol')
    } finally {
      setIsLoading(false)
    }
  }

  const togglePermission = (permissionId: string) => {
    setFormData(prev => ({
      ...prev,
      permission_ids: prev.permission_ids.includes(permissionId)
        ? prev.permission_ids.filter(id => id !== permissionId)
        : [...prev.permission_ids, permissionId]
    }))
  }

  const toggleUser = (userId: string) => {
    setFormData(prev => ({
      ...prev,
      user_ids: prev.user_ids.includes(userId)
        ? prev.user_ids.filter(id => id !== userId)
        : [...prev.user_ids, userId]
    }))
  }

  return (
    <div className="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-2xl font-bold text-gray-900">
            {isEdit ? 'Editar Rol' : 'Nuevo Rol'}
          </h2>
          <button
            onClick={() => onClose(false)}
            className="text-gray-400 hover:text-gray-500"
          >
            <XMarkIcon className="h-6 w-6" />
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit} className="p-6 space-y-6">
          {error && (
            <Alert variant="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          {/* Basic Info */}
          <div className="grid grid-cols-1 gap-4">
            <Input
              label="Nombre del Rol"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              error={errors.name}
              placeholder="Ej: Administrador, Usuario, etc."
            />
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Descripción
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={3}
                className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-purple-500 ${
                  errors.description ? 'border-red-500' : 'border-gray-300'
                }`}
                placeholder="Describe las responsabilidades de este rol"
              />
              {errors.description && (
                <p className="mt-1 text-sm text-red-600">{errors.description}</p>
              )}
            </div>
          </div>

          {/* Permissions - Grouped by Section */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              Permisos ({formData.permission_ids.length} seleccionados)
            </label>
            <div className="max-h-96 overflow-y-auto p-4 border rounded-lg space-y-6">
              {groupedPermissions.map((group) => (
                <div key={group.section} className="space-y-3">
                  {/* Section Header */}
                  <h4 className="text-sm font-semibold text-gray-900 capitalize border-b pb-2">
                    {group.section}
                  </h4>
                  
                  {/* Permissions in this section */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    {group.permissions.map((permission) => (
                      <label
                        key={permission.id}
                        className="flex items-start space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-purple-50 transition-colors"
                      >
                        <input
                          type="checkbox"
                          checked={formData.permission_ids.includes(permission.id)}
                          onChange={() => togglePermission(permission.id)}
                          className="mt-1 h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                        />
                        <div className="flex-1 min-w-0">
                          <span className="text-sm font-medium text-gray-900 block">
                            {permission.name}
                          </span>
                          <p className="text-xs text-gray-500 mt-0.5">{permission.description}</p>
                        </div>
                      </label>
                    ))}
                  </div>
                </div>
              ))}
            </div>
            {formData.permission_ids.length === 0 && (
              <p className="mt-2 text-sm text-yellow-600">⚠️ Se recomienda seleccionar al menos un permiso</p>
            )}
          </div>

          {/* Users (only in edit mode) */}
          {isEdit && users.length > 0 && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-3">
                Asignar Usuarios (Opcional)
              </label>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-3 max-h-60 overflow-y-auto p-4 border rounded-lg">
                {users.map((user) => (
                  <label
                    key={user.id}
                    className="flex items-center space-x-2 p-3 border rounded-lg cursor-pointer hover:bg-gray-50"
                  >
                    <input
                      type="checkbox"
                      checked={formData.user_ids.includes(user.id)}
                      onChange={() => toggleUser(user.id)}
                      className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                    />
                    <div>
                      <span className="text-sm font-medium text-gray-900">{user.full_name}</span>
                      <p className="text-xs text-gray-500">{user.email}</p>
                    </div>
                  </label>
                ))}
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex justify-end space-x-3 pt-4 border-t">
            <Button
              type="button"
              variant="secondary"
              onClick={() => onClose(false)}
            >
              Cancelar
            </Button>
            <Button type="submit" variant="primary" isLoading={isLoading}>
              {isEdit ? 'Actualizar' : 'Crear'} Rol
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}

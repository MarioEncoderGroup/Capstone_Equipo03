// MisViáticos - User Modal Component
'use client'

import { useState } from 'react'
import { XMarkIcon } from '@heroicons/react/24/outline'
import { UserService } from '@/services/userService'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Alert } from '@/components/ui/Alert'
import type { Role, Permission } from '@/types'
import type { User, CreateUserRequest, UpdateUserRequest, UserFormData } from '../types'
import { validateUserForm } from '../utils/validation'

interface UserModalProps {
  user: User | null
  roles: Role[]
  permissions: Permission[]
  onClose: (success: boolean) => void
}

export function UserModal({ user, roles, permissions, onClose }: UserModalProps) {
  const isEdit = !!user

  const [formData, setFormData] = useState<UserFormData>({
    full_name: user?.full_name || '',
    email: user?.email || '',
    phone: user?.phone || '+56 9 ',
    password: '',
    password_confirm: '',
    is_active: user?.is_active ?? true,
    role_ids: user?.roles.map(r => r.id) || [],
    permission_ids: user?.permissions.map(p => p.id) || [],
  })

  const [errors, setErrors] = useState<{
    full_name?: string
    email?: string
    phone?: string
    password?: string
    password_confirm?: string
  }>({})
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const validate = (): boolean => {
    const validationErrors = validateUserForm(formData, isEdit)
    setErrors(validationErrors)
    return Object.keys(validationErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validate()) return

    setIsLoading(true)
    setError(null)

    try {
      if (isEdit && user) {
        // Update user basic info
        const updateData: UpdateUserRequest = {
          full_name: formData.full_name,
          email: formData.email,
          phone: formData.phone,
          is_active: formData.is_active,
        }

        if (formData.password) {
          updateData.password = formData.password
        }

        await UserService.update(user.id, updateData)

        // Sync roles
        if (formData.role_ids.length > 0) {
          await UserService.syncRoles(user.id, formData.role_ids)
        }
      } else {
        // Create user
        const createData: CreateUserRequest = {
          full_name: formData.full_name,
          email: formData.email,
          phone: formData.phone,
          password: formData.password,
          role_ids: formData.role_ids,
          permission_ids: formData.permission_ids,
        }

        const newUser = await UserService.create(createData)

        // Sync roles after creation
        if (formData.role_ids.length > 0) {
          await UserService.syncRoles(newUser.id, formData.role_ids)
        }
      }

      onClose(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error al guardar usuario')
    } finally {
      setIsLoading(false)
    }
  }

  const toggleRole = (roleId: string) => {
    setFormData(prev => ({
      ...prev,
      role_ids: prev.role_ids.includes(roleId)
        ? prev.role_ids.filter(id => id !== roleId)
        : [...prev.role_ids, roleId]
    }))
  }

  const togglePermission = (permissionId: string) => {
    setFormData(prev => ({
      ...prev,
      permission_ids: prev.permission_ids.includes(permissionId)
        ? prev.permission_ids.filter(id => id !== permissionId)
        : [...prev.permission_ids, permissionId]
    }))
  }

  return (
    <div className="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-2xl font-bold text-gray-900">
            {isEdit ? 'Editar Usuario' : 'Nuevo Usuario'}
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
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label="Nombre Completo"
              value={formData.full_name}
              onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
              error={errors.full_name}
            />
            <Input
              label="Email"
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              error={errors.email}
            />
            <Input
              label="Teléfono"
              type="tel"
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              error={errors.phone}
            />
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Estado
              </label>
              <select
                value={formData.is_active ? 'true' : 'false'}
                onChange={(e) => setFormData({ ...formData, is_active: e.target.value === 'true' })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
              >
                <option value="true">Activo</option>
                <option value="false">Inactivo</option>
              </select>
            </div>
          </div>

          {/* Password */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label={isEdit ? 'Nueva Contraseña (opcional)' : 'Contraseña'}
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              error={errors.password}
              helperText={isEdit ? 'Dejar en blanco para mantener la actual' : undefined}
            />
            <Input
              label="Confirmar Contraseña"
              type="password"
              value={formData.password_confirm}
              onChange={(e) => setFormData({ ...formData, password_confirm: e.target.value })}
              error={errors.password_confirm}
            />
          </div>

          {/* Roles */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              Roles
            </label>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
              {roles.map((role) => (
                <label
                  key={role.id}
                  className="flex items-center space-x-2 p-3 border rounded-lg cursor-pointer hover:bg-gray-50"
                >
                  <input
                    type="checkbox"
                    checked={formData.role_ids.includes(role.id)}
                    onChange={() => toggleRole(role.id)}
                    className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                  />
                  <span className="text-sm text-gray-900">{role.name}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Permissions */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              Permisos Adicionales
            </label>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-3 max-h-60 overflow-y-auto">
              {permissions.map((permission) => (
                <label
                  key={permission.id}
                  className="flex items-center space-x-2 p-3 border rounded-lg cursor-pointer hover:bg-gray-50"
                >
                  <input
                    type="checkbox"
                    checked={formData.permission_ids.includes(permission.id)}
                    onChange={() => togglePermission(permission.id)}
                    className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                  />
                  <span className="text-sm text-gray-900">{permission.name}</span>
                </label>
              ))}
            </div>
          </div>

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
              {isEdit ? 'Actualizar' : 'Crear'} Usuario
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}

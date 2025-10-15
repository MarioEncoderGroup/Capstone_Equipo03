// MisViáticos - Roles Feature Validation

import type { RoleFormData } from '../types'

/**
 * Validation errors for role form
 */
export interface RoleValidationErrors {
  name?: string
  description?: string
}

/**
 * Validate role form data
 */
export function validateRoleForm(formData: RoleFormData): RoleValidationErrors {
  const errors: RoleValidationErrors = {}

  // Name validation
  if (!formData.name.trim()) {
    errors.name = 'El nombre es requerido'
  } else if (formData.name.trim().length < 3) {
    errors.name = 'El nombre debe tener al menos 3 caracteres'
  } else if (formData.name.trim().length > 100) {
    errors.name = 'El nombre no puede exceder 100 caracteres'
  }

  // Description validation
  if (!formData.description.trim()) {
    errors.description = 'La descripción es requerida'
  } else if (formData.description.trim().length < 10) {
    errors.description = 'La descripción debe tener al menos 10 caracteres'
  } else if (formData.description.trim().length > 255) {
    errors.description = 'La descripción no puede exceder 255 caracteres'
  }

  return errors
}

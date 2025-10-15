'use client'

import React, { useState } from 'react'
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'
import type { ResetPasswordFormData, ResetPasswordFormProps } from '../types'

export default function ResetPasswordForm({ onSubmit, isLoading, mode = 'request' }: ResetPasswordFormProps) {
  const [formData, setFormData] = useState<ResetPasswordFormData>({
    email: '',
    token: '',
    new_password: '',
    confirm_password: ''
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [errors, setErrors] = useState<Partial<ResetPasswordFormData>>({})

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    const newErrors: Partial<ResetPasswordFormData> = {}
    
    if (mode === 'request') {
      if (!formData.email) newErrors.email = 'El email es requerido'
    } else {
      if (!formData.token) newErrors.token = 'Token de recuperación requerido'
      if (!formData.new_password) newErrors.new_password = 'La nueva contraseña es requerida'
      if (!formData.confirm_password) newErrors.confirm_password = 'Confirma tu nueva contraseña'
      if (formData.new_password !== formData.confirm_password) {
        newErrors.confirm_password = 'Las contraseñas no coinciden'
      }
    }
    
    setErrors(newErrors)
    
    if (Object.keys(newErrors).length === 0) {
      await onSubmit(formData)
    }
  }

  const handleChange = (field: keyof ResetPasswordFormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: undefined }))
    }
  }

  if (mode === 'request') {
    return (
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
            Email
          </label>
          <input
            id="email"
            type="email"
            value={formData.email}
            onChange={(e) => handleChange('email', e.target.value)}
            className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.email ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="tu@email.com"
          />
          {errors.email && (
            <p className="mt-1 text-sm text-red-600">{errors.email}</p>
          )}
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full bg-purple-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-purple-700 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {isLoading ? 'Enviando...' : 'Enviar instrucciones'}
        </button>
      </form>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div>
        <label htmlFor="token" className="block text-sm font-medium text-gray-700 mb-2">
          Token de Recuperación
        </label>
        <input
          id="token"
          type="text"
          value={formData.token}
          onChange={(e) => handleChange('token', e.target.value)}
          className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
            errors.token ? 'border-red-500' : 'border-gray-300'
          }`}
          placeholder="Token del email de recuperación"
        />
        {errors.token && (
          <p className="mt-1 text-sm text-red-600">{errors.token}</p>
        )}
      </div>

      <div>
        <label htmlFor="new_password" className="block text-sm font-medium text-gray-700 mb-2">
          Nueva Contraseña
        </label>
        <div className="relative">
          <input
            id="new_password"
            type={showPassword ? 'text' : 'password'}
            value={formData.new_password}
            onChange={(e) => handleChange('new_password', e.target.value)}
            className={`w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.new_password ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="Mínimo 8 caracteres"
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
          >
            {showPassword ? (
              <EyeSlashIcon className="w-5 h-5" />
            ) : (
              <EyeIcon className="w-5 h-5" />
            )}
          </button>
        </div>
        {errors.new_password && (
          <p className="mt-1 text-sm text-red-600">{errors.new_password}</p>
        )}
      </div>

      <div>
        <label htmlFor="confirm_password" className="block text-sm font-medium text-gray-700 mb-2">
          Confirmar Nueva Contraseña
        </label>
        <div className="relative">
          <input
            id="confirm_password"
            type={showConfirmPassword ? 'text' : 'password'}
            value={formData.confirm_password}
            onChange={(e) => handleChange('confirm_password', e.target.value)}
            className={`w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.confirm_password ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="Repite tu nueva contraseña"
          />
          <button
            type="button"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
          >
            {showConfirmPassword ? (
              <EyeSlashIcon className="w-5 h-5" />
            ) : (
              <EyeIcon className="w-5 h-5" />
            )}
          </button>
        </div>
        {errors.confirm_password && (
          <p className="mt-1 text-sm text-red-600">{errors.confirm_password}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={isLoading}
        className="w-full bg-purple-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-purple-700 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {isLoading ? 'Actualizando...' : 'Actualizar Contraseña'}
      </button>
    </form>
  )
}

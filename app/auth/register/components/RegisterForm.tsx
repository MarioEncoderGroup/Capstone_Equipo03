'use client'

import React, { useState } from 'react'
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'
import type { RegisterFormData, RegisterFormProps } from '../types'

export default function RegisterForm({ onSubmit, isLoading }: RegisterFormProps) {
  const [formData, setFormData] = useState<RegisterFormData>({
    firstname: '',
    lastname: '',
    email: '',
    phone: '',
    password: '',
    password_confirm: ''
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [errors, setErrors] = useState<Partial<RegisterFormData>>({})

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // Basic validation
    const newErrors: Partial<RegisterFormData> = {}
    if (!formData.firstname) newErrors.firstname = 'El nombre es requerido'
    if (!formData.lastname) newErrors.lastname = 'El apellido es requerido'
    if (!formData.email) newErrors.email = 'El email es requerido'
    if (!formData.phone) newErrors.phone = 'El teléfono es requerido'
    if (!formData.password) newErrors.password = 'La contraseña es requerida'
    if (!formData.password_confirm) newErrors.password_confirm = 'Confirma tu contraseña'
    if (formData.password !== formData.password_confirm) {
      newErrors.password_confirm = 'Las contraseñas no coinciden'
    }
    
    setErrors(newErrors)
    
    if (Object.keys(newErrors).length === 0) {
      await onSubmit(formData)
    }
  }

  const handleChange = (field: keyof RegisterFormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: undefined }))
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label htmlFor="firstname" className="block text-sm font-medium text-gray-700 mb-2">
            Nombre
          </label>
          <input
            id="firstname"
            type="text"
            value={formData.firstname}
            onChange={(e) => handleChange('firstname', e.target.value)}
            className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.firstname ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="Tu nombre"
          />
          {errors.firstname && (
            <p className="mt-1 text-sm text-red-600">{errors.firstname}</p>
          )}
        </div>

        <div>
          <label htmlFor="lastname" className="block text-sm font-medium text-gray-700 mb-2">
            Apellido
          </label>
          <input
            id="lastname"
            type="text"
            value={formData.lastname}
            onChange={(e) => handleChange('lastname', e.target.value)}
            className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.lastname ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="Tu apellido"
          />
          {errors.lastname && (
            <p className="mt-1 text-sm text-red-600">{errors.lastname}</p>
          )}
        </div>
      </div>

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

      <div>
        <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-2">
          Teléfono
        </label>
        <input
          id="phone"
          type="tel"
          value={formData.phone}
          onChange={(e) => handleChange('phone', e.target.value)}
          className={`w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
            errors.phone ? 'border-red-500' : 'border-gray-300'
          }`}
          placeholder="+56 9 1234 5678"
        />
        {errors.phone && (
          <p className="mt-1 text-sm text-red-600">{errors.phone}</p>
        )}
      </div>

      <div>
        <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
          Contraseña
        </label>
        <div className="relative">
          <input
            id="password"
            type={showPassword ? 'text' : 'password'}
            value={formData.password}
            onChange={(e) => handleChange('password', e.target.value)}
            className={`w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.password ? 'border-red-500' : 'border-gray-300'
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
        {errors.password && (
          <p className="mt-1 text-sm text-red-600">{errors.password}</p>
        )}
      </div>

      <div>
        <label htmlFor="password_confirm" className="block text-sm font-medium text-gray-700 mb-2">
          Confirmar Contraseña
        </label>
        <div className="relative">
          <input
            id="password_confirm"
            type={showConfirmPassword ? 'text' : 'password'}
            value={formData.password_confirm}
            onChange={(e) => handleChange('password_confirm', e.target.value)}
            className={`w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
              errors.password_confirm ? 'border-red-500' : 'border-gray-300'
            }`}
            placeholder="Repite tu contraseña"
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
        {errors.password_confirm && (
          <p className="mt-1 text-sm text-red-600">{errors.password_confirm}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={isLoading}
        className="w-full bg-purple-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-purple-700 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {isLoading ? 'Creando cuenta...' : 'Crear Cuenta'}
      </button>
    </form>
  )
}

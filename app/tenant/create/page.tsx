'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useTenant } from '@/hooks/useTenant'
import { useRegions } from '@/hooks/useRegions'
import { Button } from '@/components/ui/Button'
import { Input } from '@/components/ui/Input'
import { Alert } from '@/components/ui/Alert'
import { useFormField } from '@/hooks/useFormField'
import {
  validateEmail,
  validatePhone,
  validateRequired,
  validateRUT,
  validateWebsite,
} from '@/lib/utils/validators'
import type { CreateTenantRequest } from '@/types/api'

export default function CreateTenantPage() {
  const router = useRouter()
  const { createTenant, isLoading, error, clearError } = useTenant()
  const {
    regions,
    communes,
    isLoadingRegions,
    isLoadingCommunes,
    error: regionError,
    selectedRegionId,
    selectedCommuneId,
    setSelectedRegionId,
    setSelectedCommuneId,
    clearError: clearRegionError,
  } = useRegions()

  // Form fields con validación
  const rut = useFormField<string>('', validateRUT)
  const businessName = useFormField<string>('', (value) =>
    validateRequired(value, 'Razón social')
  )
  const email = useFormField<string>('', validateEmail)
  const phone = useFormField<string>('', validatePhone)
  const address = useFormField<string>('', (value) =>
    validateRequired(value, 'Dirección')
  )
  const website = useFormField<string>('', validateWebsite)

  // Country ID (Chile - hardcoded)
  const [countryId] = useState('550e8400-e29b-41d4-a716-446655440000')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    // Validar todos los campos
    const isValid =
      rut.validate() &&
      businessName.validate() &&
      email.validate() &&
      phone.validate() &&
      address.validate() &&
      website.validate()

    if (!isValid) return

    // Validar región y comuna
    if (!selectedRegionId) {
      clearError()
      alert('Por favor selecciona una región')
      return
    }

    if (!selectedCommuneId) {
      clearError()
      alert('Por favor selecciona una comuna')
      return
    }

    try {
      const tenantData: CreateTenantRequest = {
        rut: rut.value,
        business_name: businessName.value,
        email: email.value,
        phone: phone.value,
        address: address.value,
        website: website.value,
        region_id: selectedRegionId,
        commune_id: selectedCommuneId,
        country_id: countryId,
      }

      await createTenant(tenantData)
      // useTenant ya redirige al dashboard después del éxito
    } catch (err) {
      // El error ya está manejado por useTenant
      console.error('Error creando tenant:', err)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-2xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Crea tu Empresa
          </h1>
          <p className="text-gray-600">
            Ingresa los datos de tu empresa para comenzar a gestionar tus
            viáticos
          </p>
        </div>

        <div className="bg-white py-8 px-6 shadow rounded-lg">
          {error && (
            <div className="mb-6">
              <Alert variant="error" onClose={clearError}>
                {error}
              </Alert>
            </div>
          )}

          {regionError && (
            <div className="mb-6">
              <Alert variant="warning" onClose={clearRegionError}>
                {regionError}
              </Alert>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <Input
              label="RUT de la Empresa"
              type="text"
              placeholder="12.345.678-9"
              value={rut.value}
              onChange={(e) => rut.setValue(e.target.value)}
              error={rut.error}
              helperText="Ingresa el RUT sin puntos, con guión"
            />

            <Input
              label="Razón Social"
              type="text"
              placeholder="Empresa Demo SpA"
              value={businessName.value}
              onChange={(e) => businessName.setValue(e.target.value)}
              error={businessName.error}
            />

            <Input
              label="Email de la Empresa"
              type="email"
              placeholder="contacto@empresa.cl"
              value={email.value}
              onChange={(e) => email.setValue(e.target.value)}
              error={email.error}
            />

            <Input
              label="Teléfono"
              type="tel"
              placeholder="+56 9 1234 5678"
              value={phone.value}
              onChange={(e) => phone.setValue(e.target.value)}
              error={phone.error}
            />

            <Input
              label="Dirección"
              type="text"
              placeholder="Av. Libertador 1234, Oficina 567"
              value={address.value}
              onChange={(e) => address.setValue(e.target.value)}
              error={address.error}
            />

            <Input
              label="Sitio Web"
              type="url"
              placeholder="https://ejemplo.cl"
              value={website.value}
              onChange={(e) => website.setValue(e.target.value)}
              error={website.error}
              helperText="Debe comenzar con http:// o https://"
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Región *
                </label>
                {isLoadingRegions ? (
                  <div className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-500">
                    Cargando regiones...
                  </div>
                ) : (
                  <select
                    value={selectedRegionId}
                    onChange={(e) => setSelectedRegionId(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors"
                    required
                  >
                    <option value="">Seleccione una región</option>
                    {regions.map((region) => (
                      <option key={region.id} value={region.id}>
                        {region.roman_number} - {region.name}
                      </option>
                    ))}
                  </select>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Comuna *
                </label>
                {isLoadingCommunes ? (
                  <div className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-500">
                    Cargando comunas...
                  </div>
                ) : (
                  <select
                    value={selectedCommuneId}
                    onChange={(e) => setSelectedCommuneId(e.target.value)}
                    disabled={!selectedRegionId || communes.length === 0}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors disabled:bg-gray-100 disabled:cursor-not-allowed"
                    required
                  >
                    <option value="">
                      {!selectedRegionId
                        ? 'Seleccione primero una región'
                        : 'Seleccione una comuna'}
                    </option>
                    {communes.map((commune) => (
                      <option key={commune.id} value={commune.id}>
                        {commune.name}
                      </option>
                    ))}
                  </select>
                )}
                {selectedRegionId && communes.length === 0 && !isLoadingCommunes && (
                  <p className="text-xs text-gray-500 mt-1">
                    No hay comunas disponibles para esta región
                  </p>
                )}
              </div>
            </div>

            <div className="pt-4">
              <Button
                type="submit"
                variant="primary"
                fullWidth
                isLoading={isLoading}
              >
                Crear Empresa
              </Button>
            </div>
          </form>

          <div className="mt-6 text-center">
            <Link
              href="/auth/login"
              className="text-sm text-purple-600 hover:text-purple-700"
            >
              ← Volver al login
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

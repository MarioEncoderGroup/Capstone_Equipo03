'use client'

import React from 'react'
import {
  ChartBarIcon,
  DocumentTextIcon,
  ShieldCheckIcon,
  CogIcon,
  MapPinIcon,
  DevicePhoneMobileIcon,
} from '@heroicons/react/24/outline'

interface Benefit {
  icon: React.ReactNode
  title: string
  description: string
}

const benefits: Benefit[] = [
  {
    icon: <ChartBarIcon className="w-6 h-6 text-purple-600" />,
    title: 'Gestión centralizada',
    description: 'Controla todos los viáticos de tu empresa desde un solo lugar',
  },
  {
    icon: <DocumentTextIcon className="w-6 h-6 text-purple-600" />,
    title: 'Reportes automáticos',
    description: 'Genera reportes detallados y exporta información financiera',
  },
  {
    icon: <ShieldCheckIcon className="w-6 h-6 text-purple-600" />,
    title: 'Seguridad empresarial',
    description: 'Datos protegidos con los más altos estándares de seguridad',
  },
  {
    icon: <CogIcon className="w-6 h-6 text-purple-600" />,
    title: 'Configuración flexible',
    description: 'Adapta las políticas de viáticos según tu empresa',
  },
]

const features = [
  {
    icon: <MapPinIcon className="w-5 h-5 text-white" />,
    text: 'Geolocalización automática de gastos',
  },
  {
    icon: <DevicePhoneMobileIcon className="w-5 h-5 text-white" />,
    text: 'App móvil para registro en tiempo real',
  },
]

export default function TenantBenefitsSection() {
  return (
    <div className="hidden lg:block lg:flex-1 relative min-h-screen lg:min-h-0 overflow-hidden">
      <div className="absolute inset-0 bg-gradient-to-br from-purple-600 via-purple-700 to-indigo-800"></div>
      
      {/* Background decoration - positioned to not interfere with text */}
      <div className="absolute inset-0 bg-purple-900 bg-opacity-20"></div>
      <div className="absolute top-0 right-0 w-1/4 h-1/4 bg-white bg-opacity-6 transform rotate-12 translate-x-8 -translate-y-8"></div>
      <div className="absolute bottom-0 left-0 w-1/5 h-1/5 bg-white bg-opacity-6 transform -rotate-12 -translate-x-4 translate-y-4"></div>
      
      <div className="relative h-full flex flex-col justify-center p-12">
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-white mb-4 drop-shadow-lg">
            Configura tu Empresa
          </h2>
          <div className="bg-purple-800 bg-opacity-30 backdrop-blur-sm rounded-lg p-4 inline-block">
            <p className="text-white text-lg font-medium">
              Comienza a gestionar los viáticos de tu empresa de forma profesional
            </p>
          </div>
        </div>

        <div className="space-y-6 mb-8">
          {benefits.map((benefit, index) => (
            <div key={index} className="flex items-start space-x-4">
              <div className="flex-shrink-0 w-12 h-12 bg-white bg-opacity-90 rounded-lg flex items-center justify-center shadow-lg">
                {benefit.icon}
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white mb-1 drop-shadow">
                  {benefit.title}
                </h3>
                <p className="text-purple-100 text-sm leading-relaxed drop-shadow-sm">
                  {benefit.description}
                </p>
              </div>
            </div>
          ))}
        </div>

        <div className="bg-purple-800 bg-opacity-40 backdrop-blur-sm rounded-lg p-6">
          <h3 className="text-xl font-bold text-white mb-4">
            Características destacadas
          </h3>
          <div className="space-y-3">
            {features.map((feature, index) => (
              <div key={index} className="flex items-center space-x-3">
                <div className="flex-shrink-0 w-8 h-8 bg-purple-600 bg-opacity-80 rounded-full flex items-center justify-center">
                  {feature.icon}
                </div>
                <span className="text-white text-sm font-medium">
                  {feature.text}
                </span>
              </div>
            ))}
          </div>
          
          <div className="mt-6 pt-4 border-t border-purple-300 border-opacity-30">
            <div className="bg-purple-700 bg-opacity-50 backdrop-blur-sm rounded-lg p-4">
              <p className="text-white text-sm font-semibold text-center">
                ✓ Configuración completada en menos de 5 minutos
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

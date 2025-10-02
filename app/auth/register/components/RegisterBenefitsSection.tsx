import React from 'react'

const benefits = [
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
      </svg>
    ),
    title: 'Control total de gastos',
    description: 'Monitorea cada peso gastado en viáticos y reembolsos'
  },
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
    ),
    title: 'Documentación digital',
    description: 'Digitaliza y organiza todos tus comprobantes automáticamente'
  },
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
      </svg>
    ),
    title: 'Ahorro comprobado',
    description: 'Reduce hasta un 40% los costos administrativos de viáticos'
  },
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4" />
      </svg>
    ),
    title: 'Configuración flexible',
    description: 'Adapta las políticas de viáticos según tu empresa'
  }
]

const features = [
  {
    icon: (
      <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
    ),
    title: 'Geolocalización',
    description: 'automática de gastos'
  },
  {
    icon: (
      <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
      </svg>
    ),
    title: 'App móvil',
    description: 'para registro en tiempo real'
  }
]

export default function RegisterBenefitsSection() {
  return (
    <div className="hidden lg:block lg:w-1/2 bg-gradient-to-br from-purple-500 via-purple-400 to-indigo-500 relative overflow-hidden">
      {/* Background decoration - positioned to not interfere with text */}
      <div className="absolute inset-0 bg-purple-900 bg-opacity-20"></div>
      <div className="absolute top-0 right-0 w-1/3 h-1/3 bg-white bg-opacity-8 transform rotate-12 translate-x-12 -translate-y-12"></div>
      <div className="absolute bottom-0 left-0 w-1/4 h-1/4 bg-white bg-opacity-8 transform -rotate-12 -translate-x-6 translate-y-6"></div>
      
      <div className="relative h-full flex flex-col justify-center p-12">
        <div className="mb-8">
          <h2 className="text-3xl font-bold text-white mb-4 drop-shadow-lg">
            Beneficios MisViáticos
          </h2>
          <div className="bg-purple-800 bg-opacity-30 backdrop-blur-sm rounded-lg p-4 inline-block">
            <p className="text-white text-lg font-medium">
              Únete a más de 1000+ empresas que ya optimizan sus viáticos
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
                <p className="text-white text-sm leading-relaxed font-medium drop-shadow-sm">
                  {benefit.description}
                </p>
              </div>
            </div>
          ))}
        </div>

        <div className="bg-purple-800 bg-opacity-30 backdrop-blur-sm rounded-lg p-6 mt-8">
          <h3 className="text-xl font-semibold text-white mb-4">
            Características destacadas
          </h3>
          <div className="grid grid-cols-1 gap-4 mb-6">
            {features.map((feature, index) => (
              <div key={index} className="flex items-center space-x-3">
                <div className="flex-shrink-0 w-8 h-8 bg-white bg-opacity-90 rounded-full flex items-center justify-center shadow-lg">
                  {feature.icon}
                </div>
                <div>
                  <span className="text-white font-semibold text-sm">{feature.title}</span>
                  <span className="text-white text-sm ml-1 font-medium">{feature.description}</span>
                </div>
              </div>
            ))}
          </div>
          
          <div className="border-t border-white border-opacity-30 pt-4">
            <div className="flex items-center space-x-2 text-white">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <span className="text-sm font-semibold">Implementación en menos de 24 horas</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

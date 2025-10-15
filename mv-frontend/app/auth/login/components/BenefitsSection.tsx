import React from 'react'

const benefits = [
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    ),
    title: 'Gestión automatizada',
    description: 'de viáticos y gastos empresariales en tiempo real'
  },
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      </svg>
    ),
    title: 'Reportes detallados',
    description: 'para mejor control financiero y toma de decisiones'
  },
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
      </svg>
    ),
    title: 'Seguridad garantizada',
    description: 'con encriptación de datos y respaldos automáticos'
  },
  {
    icon: (
      <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
      </svg>
    ),
    title: 'Aprobaciones rápidas',
    description: 'con flujo de trabajo optimizado y notificaciones'
  }
]

export default function BenefitsSection() {
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
          <p className="text-white text-lg font-medium drop-shadow">
            La plataforma integral para la gestión de viáticos empresariales
          </p>
        </div>

        <div className="space-y-6">
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

        <div className="mt-12">
          <div className="bg-purple-800 bg-opacity-30 backdrop-blur-sm rounded-lg p-4">
            <div className="flex items-center space-x-2 text-white">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <span className="text-sm font-semibold">Más de 1000+ empresas confían en nosotros</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function Features() {
  return (
    <section className="py-20 lg:py-32 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center space-y-4 mb-16">
          <h2 className="text-3xl lg:text-5xl font-bold text-gray-900">
            Simplifica la gestión de viáticos empresariales
          </h2>
          <p className="text-lg lg:text-xl text-gray-600 max-w-3xl mx-auto">
            Automatiza procesos, reduce errores y ahorra tiempo con nuestra plataforma integral
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-12">
          {/* Feature 1 - App Móvil */}
          <div className="group text-center space-y-6 p-8 rounded-2xl hover:bg-gray-50 transition-all duration-300">
            <div className="mx-auto w-16 h-16 bg-purple-50 rounded-2xl flex items-center justify-center transition-colors duration-300">
              <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
            </div>
            
            <div className="space-y-4">
              <h3 className="text-xl lg:text-2xl font-bold text-gray-900">App Móvil</h3>
              <p className="text-gray-600 leading-relaxed">
                Captura gastos al instante con tu smartphone. Fotografía recibos y carga información automáticamente.
              </p>
            </div>

            {/* Feature Benefits */}
            <div className="space-y-2 pt-4 border-t border-gray-100">
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Escaneo automático de recibos</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Modo offline disponible</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Geolocalización automática</span>
              </div>
            </div>
          </div>

          {/* Feature 2 - Aprobación Rápida */}
          <div className="group text-center space-y-6 p-8 rounded-2xl hover:bg-gray-50 transition-all duration-300">
            <div className="mx-auto w-16 h-16 bg-purple-50 rounded-2xl flex items-center justify-center transition-colors duration-300">
              <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            
            <div className="space-y-4">
              <h3 className="text-xl lg:text-2xl font-bold text-gray-900">Aprobación Rápida</h3>
              <p className="text-gray-600 leading-relaxed">
                Flujos de aprobación automatizados que aceleran el proceso y mantienen el control financiero.
              </p>
            </div>

            {/* Feature Benefits */}
            <div className="space-y-2 pt-4 border-t border-gray-100">
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Aprobaciones por niveles</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Notificaciones en tiempo real</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Políticas personalizables</span>
              </div>
            </div>
          </div>

          {/* Feature 3 - Reportes Inteligentes */}
          <div className="group text-center space-y-6 p-8 rounded-2xl hover:bg-gray-50 transition-all duration-300">
            <div className="mx-auto w-16 h-16 bg-purple-50 rounded-2xl flex items-center justify-center transition-colors duration-300">
              <svg className="w-8 h-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
            </div>
            
            <div className="space-y-4">
              <h3 className="text-xl lg:text-2xl font-bold text-gray-900">Reportes Inteligentes</h3>
              <p className="text-gray-600 leading-relaxed">
                Genera reportes detallados y análisis de gastos para optimizar el presupuesto de viajes.
              </p>
            </div>

            {/* Feature Benefits */}
            <div className="space-y-2 pt-4 border-t border-gray-100">
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Dashboard en tiempo real</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Exportación a Excel/PDF</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-500">
                <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
                <span>Analytics predictivo</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

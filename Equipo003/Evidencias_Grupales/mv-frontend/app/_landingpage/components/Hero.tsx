export default function Hero() {
  return (
    <section className="relative overflow-hidden bg-gray-50">
      {/* Background Banner */}
      <div className="bg-gradient-to-r from-purple-600 to-violet-600 text-white py-3">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-center space-x-4">
            <div className="flex items-center space-x-2">
              <img 
                src="/icon-mv/Assets MV_Elemento6.svg" 
                alt="MisViáticos" 
                className="w-5 h-5"
              />
              <span className="text-sm font-medium">Revoluciona la gestión de viáticos en tu empresa - Prueba GRATIS por 14 días sin compromiso</span>
            </div>
            <button className="bg-white/20 hover:bg-white/30 text-white px-4 py-1 rounded-full text-sm font-medium transition-colors duration-200">
              Prueba gratis
            </button>
          </div>
        </div>
      </div>

      {/* Main Hero Content */}
      <div 
        className="relative overflow-hidden"
        style={{
          background: `
            radial-gradient(circle at 20% 30%, rgba(147, 51, 234, 0.08) 0%, transparent 50%),
            radial-gradient(circle at 80% 70%, rgba(124, 58, 237, 0.06) 0%, transparent 50%),
            linear-gradient(135deg, rgba(147, 51, 234, 0.02) 0%, transparent 30%, rgba(168, 85, 247, 0.02) 70%, transparent 100%),
            #fefefe
          `
        }}
      >
        {/* Geometric accent elements */}
        <div className="absolute top-10 left-10 w-20 h-20 bg-gradient-to-br from-purple-100 to-violet-50 rounded-full blur-xl opacity-40"></div>
        <div className="absolute bottom-20 right-20 w-32 h-32 bg-gradient-to-tl from-violet-100 to-purple-50 rounded-full blur-2xl opacity-30"></div>
        <div className="absolute top-1/2 left-1/4 w-2 h-2 bg-purple-200 rounded-full opacity-60"></div>
        <div className="absolute top-1/3 right-1/3 w-1 h-1 bg-violet-300 rounded-full opacity-50"></div>
        
        <div className="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 lg:py-24">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-start lg:items-center isolation-auto">
          {/* Left Column - Content */}
          <div className="space-y-8">
            <div className="space-y-6">
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-black leading-tight tracking-tight">
                <span className="text-slate-800 drop-shadow-sm">Simplifica los </span>
                <span className="bg-gradient-to-r from-purple-600 via-violet-600 to-purple-700 bg-clip-text text-transparent drop-shadow-lg">
                  viáticos empresariales
                </span>
              </h1>
              
              <p className="text-lg lg:text-xl text-slate-600 leading-relaxed max-w-xl font-medium">
                No hay razón para continuar llenando manualmente planillas o reportes.{' '}
                <span className="text-slate-700 font-semibold">MisViáticos digitaliza el 100%</span> de los gastos de viaje en tu empresa, 
                haciéndolos simples y rápidos tanto para quienes los rinden como para quienes los revisan.
              </p>
            </div>

            {/* CTA Buttons */}
            <div className="isolate relative z-50">
              <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
                <button className="w-full sm:w-auto bg-gradient-to-r from-purple-600 to-violet-600 text-white px-8 py-4 rounded-lg font-semibold text-lg transition-colors duration-200 will-change-auto transform-gpu">
                  Solicita tu demo
                </button>
                
                <button className="w-full sm:w-auto group flex items-center justify-center space-x-2 text-purple-600 hover:text-purple-700 px-8 py-4 rounded-lg font-semibold text-lg transition-colors duration-200 border border-purple-200 hover:border-purple-300 hover:bg-purple-50 will-change-auto transform-gpu">
                  <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h8m2-8V6a2 2 0 00-2-2H6a2 2 0 00-2 2v4a2 2 0 002 2h16l-2-2z" />
                  </svg>
                  <span>Cómo funciona</span>
                </button>
              </div>
            </div>
          </div>

          {/* Right Column - Phone Mockup */}
          <div className="relative flex justify-center isolate z-10">
            <div className="mockup-container w-72 h-96 bg-gray-800 rounded-3xl p-2 shadow-2xl transform rotate-12 hover:rotate-6 transition-transform duration-300 will-change-transform">
              <div className="w-full h-full bg-white rounded-2xl overflow-hidden border border-gray-200">
                {/* Header */}
                <div className="bg-gray-100 px-3 py-1 border-b border-gray-200 h-12 overflow-hidden">
                  <div className="flex items-center justify-start h-full">
                    <img 
                      src="/icon-mv/Assets MV_Logo2.svg" 
                      alt="MisViáticos" 
                      className="h-16 w-auto"
                    />
                  </div>
                </div>
                
                {/* Content */}
                <div className="p-4 space-y-4">
                  <div className="flex items-center space-x-2">
                    <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                      <svg className="w-5 h-5 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    </div>
                    <div>
                      <div className="text-sm font-medium text-gray-900">Tu gasto de viaje</div>
                      <div className="text-xs text-green-600">fue aprobado</div>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="bg-gray-50 rounded-lg p-3">
                      <div className="text-sm font-medium text-gray-900">Almuerzo - Cliente</div>
                      <div className="text-xs text-gray-500">Restaurante Plaza</div>
                      <div className="text-lg font-bold text-purple-600">$25.000</div>
                    </div>
                    
                    <div className="bg-gray-50 rounded-lg p-3">
                      <div className="text-sm font-medium text-gray-900">Taxi al aeropuerto</div>
                      <div className="text-xs text-gray-500">Transporte</div>
                      <div className="text-lg font-bold text-purple-600">$15.500</div>
                    </div>
                  </div>
                  
                  <button className="w-full bg-gradient-to-r from-purple-600 to-violet-600 text-white py-3 rounded-lg text-sm font-medium">
                    Enviar reporte
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        </div>
      </div>
    </section>
  )
}

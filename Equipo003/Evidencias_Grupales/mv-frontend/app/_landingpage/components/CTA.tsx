export default function CTA() {
  return (
    <section className="py-20 lg:py-32 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center space-y-4 mb-16">
          <h2 className="text-3xl lg:text-4xl font-black text-slate-800">
            Beneficios medibles y cuantificables
          </h2>
          <p className="text-lg text-slate-600 max-w-2xl mx-auto">
            Resultados tangibles que transformarán la gestión de gastos de tu empresa
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-16">
          {/* Automatización */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 text-center space-y-4">
            <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-lg">Automatización del</h3>
              <span className="text-2xl font-black text-purple-600">95%</span>
            </div>
            <p className="text-sm text-slate-600 leading-relaxed">
              Reduce la intervención manual drásticamente mediante nuestros algoritmos avanzados
            </p>
          </div>

          {/* Precisión */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 text-center space-y-4">
            <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-lg">Precisión</h3>
              <span className="text-2xl font-black text-purple-600">garantizada</span>
            </div>
            <p className="text-sm text-slate-600 leading-relaxed">
              La IA elimina los errores humanos comunes en la gestión de gastos
            </p>
          </div>

          {/* Ahorro de tiempo */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 text-center space-y-4">
            <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-lg">Ahorro de tiempo</h3>
              <span className="text-2xl font-black text-purple-600">real</span>
            </div>
            <p className="text-sm text-slate-600 leading-relaxed">
              Reduce el procesamiento de horas a minutos con la automatización inteligente
            </p>
          </div>

          {/* Compliance */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 text-center space-y-4">
            <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-lg">Compliance sin</h3>
              <span className="text-2xl font-black text-purple-600">esfuerzo</span>
            </div>
            <p className="text-sm text-slate-600 leading-relaxed">
              Cumplimiento automático de regulaciones fiscales locales e internacionales
            </p>
          </div>

          {/* Detección de fraude */}
          <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 text-center space-y-4">
            <div className="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center mx-auto">
              <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <div>
              <h3 className="font-bold text-slate-800 text-lg">Detección de fraude</h3>
            </div>
            <p className="text-sm text-slate-600 leading-relaxed">
              Algoritmos avanzados que protegen tu empresa contra gastos fraudulentos
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}

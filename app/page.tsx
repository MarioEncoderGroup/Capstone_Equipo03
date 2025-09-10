import Link from 'next/link'
import { CONTACT_INFO, SOCIAL_LINKS } from './_landingpage/utils/constants'

export default function LandingPage() {
  const footerSections = [
    {
      title: 'Producto',
      links: [
        { label: 'Características', href: '/producto/caracteristicas' },
        { label: 'Integraciones', href: '/producto/integraciones' },
        { label: 'Seguridad', href: '/producto/seguridad' },
        { label: 'API', href: '/producto/api' }
      ]
    },
    {
      title: 'Empresa',
      links: [
        { label: 'Sobre nosotros', href: '/nosotros' },
        { label: 'Blog', href: '/blog' },
        { label: 'Carreras', href: '/carreras' },
        { label: 'Prensa', href: '/prensa' }
      ]
    },
    {
      title: 'Soporte',
      links: [
        { label: 'Centro de ayuda', href: '/ayuda' },
        { label: 'Contacto', href: '/contacto' },
        { label: 'Estado del servicio', href: '/estado' },
        { label: 'Documentación', href: '/docs' }
      ]
    }
  ]
  return (
    <div className="min-h-screen bg-white">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-100 h-16 overflow-hidden">
        <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-full">
          <div className="flex justify-between items-center h-full">
            {/* Logo */}
            <div className="flex items-center">
              <img 
                src="/icon_mv/Assets MV_Logo2.svg" 
                alt="MisViáticos" 
                className="h-24 w-auto sm:h-28 md:h-32 lg:h-40"
              />
            </div>

            {/* Navigation Menu */}
            <div className="hidden md:flex items-center space-x-8">
              <div className="relative group">
                <button className="text-gray-700 hover:text-purple-600 transition-colors duration-200 flex items-center space-x-1">
                  <span>Producto</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
              
              <div className="relative group">
                <button className="text-gray-700 hover:text-purple-600 transition-colors duration-200 flex items-center space-x-1">
                  <span>Soluciones</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
              
              <Link href="/pricing" className="text-gray-700 hover:text-purple-600 transition-colors duration-200">
                Precios
              </Link>
              
              <div className="relative group">
                <button className="text-gray-700 hover:text-purple-600 transition-colors duration-200 flex items-center space-x-1">
                  <span>Clientes</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
              
              <div className="relative group">
                <button className="text-gray-700 hover:text-purple-600 transition-colors duration-200 flex items-center space-x-1">
                  <span>Nosotros</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex items-center space-x-4">
              <Link href="/login" className="text-gray-700 hover:text-purple-600 transition-colors duration-200">
                Iniciar sesión
              </Link>
              
              <button className="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg font-medium transition-colors duration-200 flex items-center space-x-2">
                <span>Contactar ventas</span>
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
              </button>

              {/* Language Selector */}
              <div className="flex items-center space-x-1 text-gray-500">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m-9 9a9 9 0 019-9" />
                </svg>
                <span className="text-sm">ES</span>
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>

            {/* Mobile menu button */}
            <div className="md:hidden">
              <button className="text-gray-700 hover:text-purple-600 transition-colors duration-200">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </button>
            </div>
          </div>
        </nav>
      </header>

      {/* Hero Section */}
      <section className="relative overflow-hidden bg-gray-50">
        {/* Background Banner */}
        <div className="bg-gradient-to-r from-purple-600 to-violet-600 text-white py-3">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex items-center justify-center space-x-4">
              <div className="flex items-center space-x-2">
                <img 
                  src="/icon_mv/Assets MV_Elemento6.svg" 
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
                        src="/icon_mv/Assets MV_Logo2.svg" 
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

      {/* Features Section */}
      <section className="py-20 lg:py-32 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center space-y-4 mb-16">
            <h2 className="text-3xl lg:text-5xl font-bold text-gray-900">
              Simplifica la gestión de viáticos empresariales
            </h2>
            <p className="text-lg lg:text-xl text-gray-600 max-w-3xl mx-auto">
              Automatiza procesos, reduce errores y ahorra tiempo con nuestra plataforma integral
            </p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-12">
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

      {/* CTA Section */}
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

      {/* Footer */}
      <footer className="bg-gray-50 border-t border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="py-12 sm:py-16">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-8 lg:gap-12">
              <div className="sm:col-span-2 lg:col-span-2 space-y-6">
                <div className="flex items-center justify-center sm:justify-start">
                  <img 
                    src="/icon_mv/Assets MV_Logo2.svg" 
                    alt="MisViáticos" 
                    className="h-36 w-auto sm:h-48 md:h-60 lg:h-72"
                  />
                </div>
                <p className="text-gray-600 max-w-md leading-relaxed text-center sm:text-left mx-auto sm:mx-0">
                  La plataforma líder en gestión de viáticos empresariales en Chile. 
                  Digitaliza el 100% de los gastos de viaje de tu empresa.
                </p>
                <div className="flex space-x-4 justify-center sm:justify-start">
                  <Link href={SOCIAL_LINKS.linkedin} className="w-10 h-10 bg-gray-200 hover:bg-purple-600 rounded-lg flex items-center justify-center transition-colors duration-200 group">
                    <svg className="w-5 h-5 text-gray-600 group-hover:text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.338 16.338H13.67V12.16c0-.995-.017-2.277-1.387-2.277-1.39 0-1.601 1.086-1.601 2.207v4.248H8.014v-8.59h2.559v1.174h.037c.356-.675 1.227-1.387 2.526-1.387 2.703 0 3.203 1.778 3.203 4.092v4.711zM5.005 6.575a1.548 1.548 0 11-.003-3.096 1.548 1.548 0 01.003 3.096zm-1.337 9.763H6.34v-8.59H3.667v8.59zM17.668 1H2.328C1.595 1 1 1.581 1 2.298v15.403C1 18.418 1.595 19 2.328 19h15.34c.734 0 1.332-.582 1.332-1.299V2.298C19 1.581 18.402 1 17.668 1z" clipRule="evenodd" />
                    </svg>
                  </Link>
                  
                  <Link href={SOCIAL_LINKS.twitter} className="w-10 h-10 bg-gray-200 hover:bg-purple-600 rounded-lg flex items-center justify-center transition-colors duration-200 group">
                    <svg className="w-5 h-5 text-gray-600 group-hover:text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M6.29 18.251c7.547 0 11.675-6.253 11.675-11.675 0-.178 0-.355-.012-.53A8.348 8.348 0 0020 3.92a8.19 8.19 0 01-2.357.646 4.118 4.118 0 001.804-2.27 8.224 8.224 0 01-2.605.996 4.107 4.107 0 00-6.993 3.743 11.65 11.65 0 01-8.457-4.287 4.106 4.106 0 001.27 5.477A4.073 4.073 0 01.8 7.713v.052a4.105 4.105 0 003.292 4.022 4.095 4.095 0 01-1.853.07 4.108 4.108 0 003.834 2.85A8.233 8.233 0 010 16.407a11.616 11.616 0 006.29 1.84" />
                    </svg>
                  </Link>
                  
                  <Link href={SOCIAL_LINKS.facebook} className="w-10 h-10 bg-gray-200 hover:bg-purple-600 rounded-lg flex items-center justify-center transition-colors duration-200 group">
                    <svg className="w-5 h-5 text-gray-600 group-hover:text-white" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M20 10C20 4.477 15.523 0 10 0S0 4.477 0 10c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V10h2.54V7.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V10h2.773l-.443 2.89h-2.33v6.988C16.343 19.128 20 14.991 20 10z" clipRule="evenodd" />
                    </svg>
                  </Link>
                </div>
              </div>
              {footerSections.map((section) => (
                <div key={section.title} className="space-y-4 text-center sm:text-left">
                  <h4 className="font-semibold text-gray-900 text-lg">{section.title}</h4>
                  <ul className="space-y-3">
                    {section.links.map((link) => (
                      <li key={link.label}>
                        <Link href={link.href} className="text-gray-600 hover:text-purple-600 transition-colors duration-200 text-sm">
                          {link.label}
                        </Link>
                      </li>
                    ))}
                  </ul>
                </div>
              ))}
            </div>
          </div>
          <div className="py-6 sm:py-8 border-t border-gray-200">
            <div className="flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0 gap-4">
              {/* Copyright */}
              <div className="text-gray-500 text-sm text-center sm:text-left order-2 sm:order-1">
                © 2024 MisViáticos. Todos los derechos reservados.
              </div>

              {/* Legal Links */}
              <div className="flex flex-wrap items-center justify-center gap-4 sm:gap-6 text-sm order-1 sm:order-2">
                <Link href="/privacidad" className="text-gray-500 hover:text-purple-600 transition-colors duration-200 whitespace-nowrap">
                  Política de Privacidad
                </Link>
                <Link href="/terminos" className="text-gray-500 hover:text-purple-600 transition-colors duration-200 whitespace-nowrap">
                  Términos de Servicio
                </Link>
                <Link href="/cookies" className="text-gray-500 hover:text-purple-600 transition-colors duration-200 whitespace-nowrap">
                  Cookies
                </Link>
              </div>

              {/* Contact Info */}
              <div className="flex items-center justify-center sm:justify-end text-sm text-gray-500 order-3">
                <div className="flex items-center space-x-1">
                  <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                  <span className="break-all sm:break-normal">{CONTACT_INFO.email}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}

import Link from 'next/link'
import { CONTACT_INFO, SOCIAL_LINKS } from '../utils/constants'

export default function Footer() {
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
    <footer className="bg-gray-50 border-t border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="py-12 sm:py-16">
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-8 lg:gap-12">
            <div className="sm:col-span-2 lg:col-span-2 space-y-6">
              <div className="flex items-center justify-center sm:justify-start">
                <img 
                  src="/icon-mv/Assets MV_Logo2.svg" 
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
  )
}

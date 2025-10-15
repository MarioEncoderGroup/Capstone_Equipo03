// MisViáticos Landing Page - Constants and Configuration

export const NAVIGATION_ITEMS = [
  {
    label: 'Producto',
    href: '/producto',
    hasDropdown: true,
    dropdownItems: [
      { label: 'Características', href: '/producto/caracteristicas' },
      { label: 'Integrations', href: '/producto/integraciones' },
      { label: 'Seguridad', href: '/producto/seguridad' }
    ]
  },
  {
    label: 'Soluciones',
    href: '/soluciones',
    hasDropdown: true,
    dropdownItems: [
      { label: 'Empresas', href: '/soluciones/empresas' },
      { label: 'Startups', href: '/soluciones/startups' },
      { label: 'Consultoras', href: '/soluciones/consultoras' }
    ]
  },
  {
    label: 'Precios',
    href: '/precios',
    hasDropdown: false
  },
  {
    label: 'Clientes',
    href: '/clientes',
    hasDropdown: true,
    dropdownItems: [
      { label: 'Casos de Éxito', href: '/clientes/casos-exito' },
      { label: 'Testimonios', href: '/clientes/testimonios' }
    ]
  },
  {
    label: 'Nosotros',
    href: '/nosotros',
    hasDropdown: true,
    dropdownItems: [
      { label: 'Sobre MisViáticos', href: '/nosotros/about' },
      { label: 'Blog', href: '/blog' },
      { label: 'Contacto', href: '/contacto' }
    ]
  }
];

export const BRAND_COLORS = {
  primary: 'var(--misviaticos-primary)',
  secondary: 'var(--misviaticos-secondary)',
  accent: 'var(--misviaticos-accent)',
  gradientStart: 'var(--misviaticos-gradient-start)',
  gradientEnd: 'var(--misviaticos-gradient-end)'
};

export const HERO_CONTENT = {
  title: {
    line1: 'Gestiona tus viáticos',
    line2: 'en segundos con',
    line3: 'MisViáticos'
  },
  subtitle: 'No hay razón para continuar llenando manualmente planillas o reportes. MisViáticos digitaliza el 100% de los gastos de viaje en tu empresa, haciéndolos simples y rápidos tanto para quienes los rinden como para quienes los revisan.',
  ctaButtons: {
    primary: 'Solicita tu demo',
    secondary: 'Cómo funciona'
  }
};

export const FEATURES_DATA = [
  {
    id: 'app-movil',
    title: 'App Móvil',
    description: 'Captura gastos al instante con tu smartphone. Fotografía recibos y carga información automáticamente.',
    benefits: [
      'Escaneo automático de recibos',
      'Modo offline disponible',
      'Geolocalización automática'
    ]
  },
  {
    id: 'aprobacion-rapida',
    title: 'Aprobación Rápida',
    description: 'Flujos de aprobación automatizados que aceleran el proceso y mantienen el control financiero.',
    benefits: [
      'Aprobaciones por niveles',
      'Notificaciones en tiempo real',
      'Políticas personalizables'
    ]
  },
  {
    id: 'reportes-inteligentes',
    title: 'Reportes Inteligentes',
    description: 'Genera reportes detallados y análisis de gastos para optimizar el presupuesto de viajes.',
    benefits: [
      'Dashboard en tiempo real',
      'Exportación a Excel/PDF',
      'Analytics predictivo'
    ]
  }
];

export const ADDITIONAL_FEATURES = [
  { label: 'Seguridad bancaria', icon: 'shield' },
  { label: 'Múltiples monedas', icon: 'currency' },
  { label: 'Integración API', icon: 'api' },
  { label: 'Soporte 24/7', icon: 'support' }
];

export const CONTACT_INFO = {
  email: 'contacto@misviaticos.cl',
  phone: '+56 2 2XXX XXXX',
  address: 'Santiago, Chile'
};

export const SOCIAL_LINKS = {
  linkedin: 'https://linkedin.com/company/misviaticos',
  twitter: 'https://twitter.com/misviaticos',
  facebook: 'https://facebook.com/misviaticos'
};

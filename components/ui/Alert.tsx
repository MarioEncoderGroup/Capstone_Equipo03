// MisViáticos - Alert Component (Design System)

import type { ReactNode } from 'react'

interface AlertProps {
  children: ReactNode
  variant?: 'success' | 'error' | 'warning' | 'info'
  onClose?: () => void
}

export function Alert({
  children,
  variant = 'info',
  onClose,
}: AlertProps) {
  const variants = {
    success: 'bg-green-50 border-green-200 text-green-800',
    error: 'bg-red-50 border-red-200 text-red-600',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-800',
    info: 'bg-blue-50 border-blue-200 text-blue-800',
  }

  const icons = {
    success: '✓',
    error: '✕',
    warning: '⚠',
    info: 'ℹ',
  }

  return (
    <div
      className={`p-4 border rounded-lg flex items-start justify-between ${variants[variant]}`}
    >
      <div className="flex items-start space-x-3">
        <span className="text-xl font-bold">{icons[variant]}</span>
        <div className="text-sm flex-1">{children}</div>
      </div>
      {onClose && (
        <button
          onClick={onClose}
          className="text-current opacity-60 hover:opacity-100 transition-opacity"
          aria-label="Cerrar"
        >
          ✕
        </button>
      )}
    </div>
  )
}

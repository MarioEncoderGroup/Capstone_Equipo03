import React from 'react'
import Link from 'next/link'

export default function LoginHeader() {
  return (
    <div className="text-center">
      <div className="mb-8">
        <Link href="/" className="inline-block">
          <div className="flex items-center justify-center">
            <img 
              src="/icon-mv/Assets MV_Elemento1.svg" 
              alt="MisViáticos" 
              className="h-12 w-auto"
            />
            <span className="ml-2 text-2xl font-bold text-gray-900">MisViáticos</span>
          </div>
        </Link>
      </div>
      
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Bienvenido de vuelta
        </h1>
        <p className="text-gray-600">
          Inicia sesión en tu cuenta de MisViáticos
        </p>
      </div>
    </div>
  )
}

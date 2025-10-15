// MisViáticos - useFormField Hook

'use client'

import { useState, useCallback } from 'react'

interface UseFormFieldReturn<T> {
  value: T
  error: string | null
  setValue: (value: T) => void
  setError: (error: string | null) => void
  validate: () => boolean
  reset: () => void
}

/**
 * Hook reutilizable para manejar campos de formulario con validación
 * 
 * @template T - Tipo del valor del campo (debe ser string, number o boolean)
 * @param initialValue - Valor inicial del campo
 * @param validator - Función opcional de validación que retorna string de error o null
 */
export function useFormField<T extends string | number | boolean = string>(
  initialValue: T,
  validator?: (value: T) => string | null
): UseFormFieldReturn<T> {
  const [value, setValue] = useState<T>(initialValue)
  const [error, setError] = useState<string | null>(null)

  const handleChange = useCallback(
    (newValue: T) => {
      setValue(newValue)
      // Limpiar error cuando el usuario empieza a escribir
      if (error && validator) {
        const validationError = validator(newValue)
        if (!validationError) {
          setError(null)
        }
      }
    },
    [error, validator]
  )

  const validate = useCallback(() => {
    if (validator) {
      const validationError = validator(value)
      setError(validationError)
      return validationError === null
    }
    return true
  }, [value, validator])

  const reset = useCallback(() => {
    setValue(initialValue)
    setError(null)
  }, [initialValue])

  return {
    value,
    error,
    setValue: handleChange,
    setError,
    validate,
    reset,
  }
}

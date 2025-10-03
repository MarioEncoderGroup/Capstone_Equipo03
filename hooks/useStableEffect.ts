// MisViáticos - Stable Effect Hook
// Prevents duplicate API calls in React StrictMode (development)

'use client'

import React, { useEffect, useRef, useState, DependencyList, EffectCallback } from 'react'

/**
 * Custom hook that prevents duplicate API calls in React StrictMode
 * Similar to useEffect but only runs once even in StrictMode
 *
 * @param effect - Effect function to run
 * @param deps - Dependency array
 */
export function useStableEffect(effect: EffectCallback, deps?: DependencyList) {
  const hasRun = useRef(false)
  const cleanup = useRef<void | (() => void)>(undefined)

  useEffect(() => {
    // In production or if already run, skip
    if (hasRun.current) {
      return cleanup.current
    }

    // Mark as run and execute effect
    hasRun.current = true
    cleanup.current = effect()

    // Return cleanup function
    return () => {
      if (cleanup.current) {
        cleanup.current()
      }
    }
  }, deps)
}

/**
 * Custom hook for stable data fetching that prevents duplicate calls
 * Automatically handles loading and error states
 *
 * @param fetchFn - Async function to fetch data
 * @param deps - Dependency array
 */
export function useStableFetch<T>(
  fetchFn: () => Promise<T>,
  deps: DependencyList = []
): {
  data: T | null
  isLoading: boolean
  error: Error | null
  refetch: () => Promise<void>
} {
  const [data, setData] = useState<T | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  const hasFetched = useRef(false)
  const abortControllerRef = useRef<AbortController | null>(null)

  const fetch = async () => {
    // Cancel previous request
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    abortControllerRef.current = new AbortController()

    setIsLoading(true)
    setError(null)

    try {
      const result = await fetchFn()
      setData(result)
    } catch (err) {
      // Ignore abort errors
      if (err instanceof Error && err.name === 'AbortError') {
        return
      }
      setError(err instanceof Error ? err : new Error('Unknown error'))
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    // Prevent duplicate fetches in StrictMode
    if (hasFetched.current) {
      return
    }

    hasFetched.current = true
    fetch()

    // Cleanup
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
    }
  }, deps)

  return {
    data,
    isLoading,
    error,
    refetch: fetch,
  }
}

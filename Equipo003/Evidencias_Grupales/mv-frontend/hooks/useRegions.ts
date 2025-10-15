// MisViáticos - useRegions Hook
'use client'

import { useState, useEffect, useCallback } from 'react'
import { RegionService } from '@/services/regionService'
import type { Region, Commune } from '@/types/api'

interface UseRegionsReturn {
  regions: Region[]
  communes: Commune[]
  isLoadingRegions: boolean
  isLoadingCommunes: boolean
  error: string | null
  selectedRegionId: string
  selectedCommuneId: string
  setSelectedRegionId: (regionId: string) => void
  setSelectedCommuneId: (communeId: string) => void
  loadCommunes: (regionId: string) => Promise<void>
  clearError: () => void
}

/**
 * Custom hook for managing regions and communes
 * Provides cascade functionality: when region changes, communes are loaded
 */
export function useRegions(): UseRegionsReturn {
  const [regions, setRegions] = useState<Region[]>([])
  const [communes, setCommunes] = useState<Commune[]>([])
  const [isLoadingRegions, setIsLoadingRegions] = useState(false)
  const [isLoadingCommunes, setIsLoadingCommunes] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedRegionId, setSelectedRegionId] = useState('')
  const [selectedCommuneId, setSelectedCommuneId] = useState('')

  /**
   * Load all regions on mount
   */
  useEffect(() => {
    const fetchRegions = async () => {
      setIsLoadingRegions(true)
      setError(null)

      try {
        const data = await RegionService.getRegions()
        setRegions(data)
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Error al cargar regiones'
        setError(errorMessage)
      } finally {
        setIsLoadingRegions(false)
      }
    }

    fetchRegions()
  }, [])

  /**
   * Load communes when region is selected
   */
  const loadCommunes = useCallback(async (regionId: string) => {
    if (!regionId) {
      setCommunes([])
      return
    }

    setIsLoadingCommunes(true)
    setError(null)

    try {
      const data = await RegionService.getCommunes(regionId)
      setCommunes(data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Error al cargar comunas'
      setError(errorMessage)
      setCommunes([])
    } finally {
      setIsLoadingCommunes(false)
    }
  }, [])

  /**
   * Handle region selection change
   * Automatically loads communes and clears selected commune
   */
  const handleSetSelectedRegionId = useCallback(
    (regionId: string) => {
      setSelectedRegionId(regionId)
      setSelectedCommuneId('') // Clear commune when region changes
      if (regionId) {
        loadCommunes(regionId)
      } else {
        setCommunes([])
      }
    },
    [loadCommunes]
  )

  /**
   * Clear error message
   */
  const clearError = useCallback(() => {
    setError(null)
  }, [])

  return {
    regions,
    communes,
    isLoadingRegions,
    isLoadingCommunes,
    error,
    selectedRegionId,
    selectedCommuneId,
    setSelectedRegionId: handleSetSelectedRegionId,
    setSelectedCommuneId,
    loadCommunes,
    clearError,
  }
}

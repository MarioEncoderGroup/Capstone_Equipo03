package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/core/region/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
)

// RegionController maneja las operaciones de regiones
type RegionController struct {
	regionService ports.RegionService
}

// NewRegionController crea una nueva instancia del controller de regiones
func NewRegionController(regionService ports.RegionService) *RegionController {
	return &RegionController{
		regionService: regionService,
	}
}

// GetAllRegions maneja GET /regions - Lista todas las regiones
func (rc *RegionController) GetAllRegions(c *fiber.Ctx) error {
	regions, err := rc.regionService.GetAllRegions(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo regiones",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Regiones obtenidas exitosamente",
		Data: fiber.Map{
			"regions": regions,
			"total":   len(regions),
		},
	})
}

// GetRegionByID maneja GET /regions/:id - Obtiene una región por ID
func (rc *RegionController) GetRegionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de región es requerido",
		})
	}

	region, err := rc.regionService.GetRegionByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo región",
			Error:   err.Error(),
		})
	}

	if region == nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Región no encontrada",
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Región obtenida exitosamente",
		Data:    region,
	})
}

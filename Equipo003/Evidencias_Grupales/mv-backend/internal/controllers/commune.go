package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/core/commune/ports"
	"github.com/JoseLuis21/mv-backend/internal/shared/types"
)

// CommuneController maneja las operaciones de comunas
type CommuneController struct {
	communeService ports.CommuneService
}

// NewCommuneController crea una nueva instancia del controller de comunas
func NewCommuneController(communeService ports.CommuneService) *CommuneController {
	return &CommuneController{
		communeService: communeService,
	}
}

// GetCommunes maneja GET /communes - Lista comunas (todas o filtradas por región)
func (cc *CommuneController) GetCommunes(c *fiber.Ctx) error {
	regionID := c.Query("region_id")

	var communes interface{}
	var err error

	if regionID != "" {
		// Obtener comunas de una región específica
		communes, err = cc.communeService.GetCommunesByRegion(c.Context(), regionID)
	} else {
		// Obtener todas las comunas
		communes, err = cc.communeService.GetAllCommunes(c.Context())
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo comunas",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Comunas obtenidas exitosamente",
		Data: fiber.Map{
			"communes": communes,
		},
	})
}

// GetCommuneByID maneja GET /communes/:id - Obtiene una comuna por ID
func (cc *CommuneController) GetCommuneByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(types.APIResponse{
			Success: false,
			Message: "ID de comuna es requerido",
		})
	}

	commune, err := cc.communeService.GetCommuneByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(types.APIResponse{
			Success: false,
			Message: "Error obteniendo comuna",
			Error:   err.Error(),
		})
	}

	if commune == nil {
		return c.Status(fiber.StatusNotFound).JSON(types.APIResponse{
			Success: false,
			Message: "Comuna no encontrada",
		})
	}

	return c.Status(fiber.StatusOK).JSON(types.APIResponse{
		Success: true,
		Message: "Comuna obtenida exitosamente",
		Data:    commune,
	})
}

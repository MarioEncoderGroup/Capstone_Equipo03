package services

import (
	"context"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/region/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/region/ports"
)

// regionService implementa RegionService
type regionService struct {
	regionRepo ports.RegionRepository
}

// NewRegionService crea una nueva instancia del servicio de regiones
func NewRegionService(regionRepo ports.RegionRepository) ports.RegionService {
	return &regionService{
		regionRepo: regionRepo,
	}
}

// GetAllRegions obtiene todas las regiones
func (s *regionService) GetAllRegions(ctx context.Context) ([]domain.Region, error) {
	regions, err := s.regionRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo regiones: %w", err)
	}
	return regions, nil
}

// GetRegionByID obtiene una región por su ID
func (s *regionService) GetRegionByID(ctx context.Context, id string) (*domain.Region, error) {
	region, err := s.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo región: %w", err)
	}
	if region == nil {
		return nil, fmt.Errorf("región no encontrada")
	}
	return region, nil
}

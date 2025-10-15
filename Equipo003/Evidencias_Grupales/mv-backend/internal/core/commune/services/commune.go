package services

import (
	"context"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/commune/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/commune/ports"
)

// communeService implementa CommuneService
type communeService struct {
	communeRepo ports.CommuneRepository
}

// NewCommuneService crea una nueva instancia del servicio de comunas
func NewCommuneService(communeRepo ports.CommuneRepository) ports.CommuneService {
	return &communeService{
		communeRepo: communeRepo,
	}
}

// GetAllCommunes obtiene todas las comunas
func (s *communeService) GetAllCommunes(ctx context.Context) ([]domain.Commune, error) {
	communes, err := s.communeRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo comunas: %w", err)
	}
	return communes, nil
}

// GetCommunesByRegion obtiene las comunas de una región específica
func (s *communeService) GetCommunesByRegion(ctx context.Context, regionID string) ([]domain.Commune, error) {
	if regionID == "" {
		return nil, fmt.Errorf("region_id es requerido")
	}
	
	communes, err := s.communeRepo.GetByRegionID(ctx, regionID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo comunas de la región: %w", err)
	}
	return communes, nil
}

// GetCommuneByID obtiene una comuna por su ID
func (s *communeService) GetCommuneByID(ctx context.Context, id string) (*domain.Commune, error) {
	commune, err := s.communeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo comuna: %w", err)
	}
	if commune == nil {
		return nil, fmt.Errorf("comuna no encontrada")
	}
	return commune, nil
}

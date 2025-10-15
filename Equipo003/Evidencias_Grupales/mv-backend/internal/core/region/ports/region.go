package ports

import (
	"context"
	"github.com/JoseLuis21/mv-backend/internal/core/region/domain"
)

// RegionRepository define las operaciones de persistencia para regiones
type RegionRepository interface {
	GetAll(ctx context.Context) ([]domain.Region, error)
	GetByID(ctx context.Context, id string) (*domain.Region, error)
}

// RegionService define las operaciones de negocio para regiones
type RegionService interface {
	GetAllRegions(ctx context.Context) ([]domain.Region, error)
	GetRegionByID(ctx context.Context, id string) (*domain.Region, error)
}

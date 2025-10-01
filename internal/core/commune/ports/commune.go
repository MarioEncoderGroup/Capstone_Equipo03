package ports

import (
	"context"
	"github.com/JoseLuis21/mv-backend/internal/core/commune/domain"
)

// CommuneRepository define las operaciones de persistencia para comunas
type CommuneRepository interface {
	GetAll(ctx context.Context) ([]domain.Commune, error)
	GetByRegionID(ctx context.Context, regionID string) ([]domain.Commune, error)
	GetByID(ctx context.Context, id string) (*domain.Commune, error)
}

// CommuneService define las operaciones de negocio para comunas
type CommuneService interface {
	GetAllCommunes(ctx context.Context) ([]domain.Commune, error)
	GetCommunesByRegion(ctx context.Context, regionID string) ([]domain.Commune, error)
	GetCommuneByID(ctx context.Context, id string) (*domain.Commune, error)
}

package adapters

import (
	"context"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/region/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/region/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// PostgreSQLRegionRepository implementa RegionRepository usando PostgreSQL
type PostgreSQLRegionRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLRegionRepository crea una nueva instancia del repositorio
func NewPostgreSQLRegionRepository(client *postgresql.PostgresqlClient) ports.RegionRepository {
	return &PostgreSQLRegionRepository{
		client: client,
	}
}

// GetAll obtiene todas las regiones ordenadas por número
func (r *PostgreSQLRegionRepository) GetAll(ctx context.Context) ([]domain.Region, error) {
	query := `
		SELECT id, number, roman_number, name
		FROM region
		ORDER BY number ASC
	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando consulta: %w", err)
	}
	defer rows.Close()

	var regions []domain.Region
	for rows.Next() {
		var region domain.Region
		if err := rows.Scan(&region.ID, &region.Number, &region.RomanNumber, &region.Name); err != nil {
			return nil, fmt.Errorf("error escaneando región: %w", err)
		}
		regions = append(regions, region)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return regions, nil
}

// GetByID obtiene una región por su ID
func (r *PostgreSQLRegionRepository) GetByID(ctx context.Context, id string) (*domain.Region, error) {
	query := `
		SELECT id, number, roman_number, name
		FROM region
		WHERE id = $1
	`

	var region domain.Region
	row := r.client.QueryRow(ctx, query, id)
	err := row.Scan(&region.ID, &region.Number, &region.RomanNumber, &region.Name)
	
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("error obteniendo región: %w", err)
	}

	return &region, nil
}

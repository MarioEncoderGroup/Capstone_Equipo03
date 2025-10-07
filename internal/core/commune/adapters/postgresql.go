package adapters

import (
	"context"
	"fmt"

	"github.com/JoseLuis21/mv-backend/internal/core/commune/domain"
	"github.com/JoseLuis21/mv-backend/internal/core/commune/ports"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
)

// PostgreSQLCommuneRepository implementa CommuneRepository usando PostgreSQL
type PostgreSQLCommuneRepository struct {
	client *postgresql.PostgresqlClient
}

// NewPostgreSQLCommuneRepository crea una nueva instancia del repositorio
func NewPostgreSQLCommuneRepository(client *postgresql.PostgresqlClient) ports.CommuneRepository {
	return &PostgreSQLCommuneRepository{
		client: client,
	}
}

// GetAll obtiene todas las comunas ordenadas por nombre
func (r *PostgreSQLCommuneRepository) GetAll(ctx context.Context) ([]domain.Commune, error) {
	query := `
		SELECT id, region_id, name
		FROM commune
		ORDER BY name ASC
	`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando consulta: %w", err)
	}
	defer rows.Close()

	var communes []domain.Commune
	for rows.Next() {
		var commune domain.Commune
		if err := rows.Scan(&commune.ID, &commune.RegionID, &commune.Name); err != nil {
			return nil, fmt.Errorf("error escaneando comuna: %w", err)
		}
		communes = append(communes, commune)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return communes, nil
}

// GetByRegionID obtiene todas las comunas de una región específica
func (r *PostgreSQLCommuneRepository) GetByRegionID(ctx context.Context, regionID string) ([]domain.Commune, error) {
	query := `
		SELECT id, region_id, name
		FROM commune
		WHERE region_id = $1
		ORDER BY name ASC
	`

	rows, err := r.client.Query(ctx, query, regionID)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando consulta: %w", err)
	}
	defer rows.Close()

	var communes []domain.Commune
	for rows.Next() {
		var commune domain.Commune
		if err := rows.Scan(&commune.ID, &commune.RegionID, &commune.Name); err != nil {
			return nil, fmt.Errorf("error escaneando comuna: %w", err)
		}
		communes = append(communes, commune)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return communes, nil
}

// GetByID obtiene una comuna por su ID
func (r *PostgreSQLCommuneRepository) GetByID(ctx context.Context, id string) (*domain.Commune, error) {
	query := `
		SELECT id, region_id, name
		FROM commune
		WHERE id = $1
	`

	var commune domain.Commune
	row := r.client.QueryRow(ctx, query, id)
	err := row.Scan(&commune.ID, &commune.RegionID, &commune.Name)
	
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("error obteniendo comuna: %w", err)
	}

	return &commune, nil
}

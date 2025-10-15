package domain

// Commune representa una comuna de Chile
type Commune struct {
	ID       string `json:"id" db:"id"`               // ID de la comuna (slug)
	RegionID string `json:"region_id" db:"region_id"` // ID de la región padre
	Name     string `json:"name" db:"name"`           // Nombre de la comuna
}

// CommuneListResponse representa la respuesta de lista de comunas
type CommuneListResponse struct {
	Communes []Commune `json:"communes"`
	Total    int       `json:"total"`
}

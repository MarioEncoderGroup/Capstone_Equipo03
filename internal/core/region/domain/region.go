package domain

// Region representa una región de Chile
type Region struct {
	ID          string `json:"id" db:"id"`                       // Código ISO (AP, TA, ...)
	Number      int    `json:"number" db:"number"`               // Número de región (1-16)
	RomanNumber string `json:"roman_number" db:"roman_number"`   // Número romano (I, II, ...)
	Name        string `json:"name" db:"name"`                   // Nombre de la región
}

// RegionListResponse representa la respuesta de lista de regiones
type RegionListResponse struct {
	Regions []Region `json:"regions"`
	Total   int      `json:"total"`
}

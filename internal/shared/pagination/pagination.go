package pagination

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// PaginationRequest representa los parámetros de paginación desde query parameters
type PaginationRequest struct {
	Page     int    `query:"page" validate:"min=1"`
	PageSize int    `query:"page_size" validate:"min=1,max=100"`
	SortBy   string `query:"sort_by"`
	SortDir  string `query:"sort_dir" validate:"oneof=asc desc ASC DESC"`
	Search   string `query:"search"`
}

// PaginationResponse representa la información de paginación en las respuestas
type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// PaginatedResponse estructura genérica para respuestas paginadas
type PaginatedResponse struct {
	Data       interface{}         `json:"data"`
	Pagination *PaginationResponse `json:"pagination"`
	Success    bool               `json:"success"`
	Message    string             `json:"message"`
}

// DefaultPageSize constante para el tamaño de página por defecto
const DefaultPageSize = 20

// MaxPageSize constante para el tamaño máximo de página
const MaxPageSize = 100

// ParsePaginationFromContext extrae y valida parámetros de paginación desde el contexto Fiber
func ParsePaginationFromContext(c *fiber.Ctx) (*PaginationRequest, error) {
	req := &PaginationRequest{
		Page:     1,
		PageSize: DefaultPageSize,
		SortBy:   "created",
		SortDir:  "DESC",
		Search:   "",
	}

	// Parse page
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return nil, fmt.Errorf("página inválida: debe ser un número mayor a 0")
		}
		req.Page = page
	}

	// Parse page_size
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 || pageSize > MaxPageSize {
			return nil, fmt.Errorf("tamaño de página inválido: debe ser un número entre 1 y %d", MaxPageSize)
		}
		req.PageSize = pageSize
	}

	// Parse sort_by (validar que sea un campo válido)
	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = sanitizeSortField(sortBy)
	}

	// Parse sort_dir
	if sortDir := c.Query("sort_dir"); sortDir != "" {
		sortDir = strings.ToUpper(sortDir)
		if sortDir != "ASC" && sortDir != "DESC" {
			return nil, fmt.Errorf("dirección de ordenamiento inválida: debe ser ASC o DESC")
		}
		req.SortDir = sortDir
	}

	// Parse search
	if search := c.Query("search"); search != "" {
		req.Search = strings.TrimSpace(search)
	}

	return req, nil
}

// GetOffset calcula el offset para la consulta SQL basado en la página y tamaño
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit devuelve el límite para la consulta SQL
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

// BuildOrderBy construye la cláusula ORDER BY para SQL
func (p *PaginationRequest) BuildOrderBy() string {
	return fmt.Sprintf("%s %s", p.SortBy, p.SortDir)
}

// BuildSearchCondition construye condiciones de búsqueda para campos específicos
func (p *PaginationRequest) BuildSearchCondition(searchFields []string) (string, []interface{}) {
	if p.Search == "" || len(searchFields) == 0 {
		return "", nil
	}

	var conditions []string
	var args []interface{}
	argIndex := 1

	for _, field := range searchFields {
		conditions = append(conditions, fmt.Sprintf("%s ILIKE $%d", field, argIndex))
		args = append(args, "%"+p.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " OR ")
	return "(" + whereClause + ")", args
}

// CalculatePagination calcula la información de paginación basada en el total de registros
func (p *PaginationRequest) CalculatePagination(totalRecords int64) *PaginationResponse {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(p.PageSize)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &PaginationResponse{
		CurrentPage:  p.Page,
		PageSize:     p.PageSize,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		HasNext:      p.Page < totalPages,
		HasPrevious:  p.Page > 1,
	}
}

// CreatePaginatedResponse crea una respuesta paginada estándar
func CreatePaginatedResponse(data interface{}, pagination *PaginationResponse, message string) *PaginatedResponse {
	return &PaginatedResponse{
		Data:       data,
		Pagination: pagination,
		Success:    true,
		Message:    message,
	}
}

// sanitizeSortField limpia y valida el campo de ordenamiento para prevenir SQL injection
func sanitizeSortField(field string) string {
	// Remover caracteres peligrosos y espacios
	field = strings.TrimSpace(field)
	field = strings.ReplaceAll(field, ";", "")
	field = strings.ReplaceAll(field, "--", "")
	field = strings.ReplaceAll(field, "/*", "")
	field = strings.ReplaceAll(field, "*/", "")

	// Lista de campos permitidos por defecto
	allowedFields := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created":    true,
		"updated":    true,
		"full_name":  true,
		"phone":      true,
		"is_active":  true,
		"deleted_at": true,
	}

	// Validar que el campo esté en la lista permitida
	if !allowedFields[field] {
		return "created" // Campo por defecto seguro
	}

	return field
}

// ValidateSortField valida que un campo específico esté permitido para ordenamiento
func ValidateSortField(field string, allowedFields []string) bool {
	for _, allowedField := range allowedFields {
		if field == allowedField {
			return true
		}
	}
	return false
}

// SetCustomSortField establece un campo de ordenamiento personalizado con validación
func (p *PaginationRequest) SetCustomSortField(field string, allowedFields []string) error {
	if !ValidateSortField(field, allowedFields) {
		return fmt.Errorf("campo de ordenamiento no permitido: %s", field)
	}
	p.SortBy = field
	return nil
}
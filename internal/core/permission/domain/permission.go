package domain

import (
	"time"

	"github.com/google/uuid"
)

// Permission representa la entidad principal de permiso en el dominio
// Mapea directamente con la tabla 'permissions' en la base de datos control
type Permission struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Section     string     `json:"section"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// PermissionSection representa las secciones de permisos del sistema
type PermissionSection string

const (
	SectionRole                  PermissionSection = "role"
	SectionPermission           PermissionSection = "permission"
	SectionUser                 PermissionSection = "user"
	SectionTenant               PermissionSection = "tenant"
	SectionRegion               PermissionSection = "region"
	SectionCommune              PermissionSection = "commune"
	SectionCountry              PermissionSection = "country"
	SectionCurrency             PermissionSection = "currency"
	SectionBank                 PermissionSection = "bank"
	SectionPolicy               PermissionSection = "policy"
	SectionCategory             PermissionSection = "category"
	SectionPolicyApprover       PermissionSection = "policy-approver"
	SectionPolicySubmitter      PermissionSection = "policy-submitter"
	SectionPolicyReportStatus   PermissionSection = "policy-report-status"
	SectionPolicyField          PermissionSection = "policy-field"
	SectionReportExpense        PermissionSection = "report-expense"
	SectionReportExpenseDetail  PermissionSection = "report-expense-detail"
	SectionReportExpenseRecord  PermissionSection = "report-expense-record"
	SectionReportExpenseComment PermissionSection = "report-expense-comment"
	SectionReportExpenseApprover PermissionSection = "report-expense-approver"
	SectionExpenseGallery       PermissionSection = "expense-gallery"
)

// Permisos predefinidos del sistema - Roles
const (
	PermissionListRole   = "list-role"
	PermissionCreateRole = "create-role"
	PermissionUpdateRole = "update-role"
	PermissionDeleteRole = "delete-role"
)

// Permisos predefinidos del sistema - Permisos
const (
	PermissionListPermission   = "list-permission"
	PermissionCreatePermission = "create-permission"
	PermissionUpdatePermission = "update-permission"
	PermissionDeletePermission = "delete-permission"
)

// Permisos predefinidos del sistema - Usuarios
const (
	PermissionListUser   = "list-user"
	PermissionCreateUser = "create-user"
	PermissionUpdateUser = "update-user"
	PermissionDeleteUser = "delete-user"
)

// Permisos predefinidos del sistema - Tenants
const (
	PermissionListTenant   = "list-tenant"
	PermissionCreateTenant = "create-tenant"
	PermissionUpdateTenant = "update-tenant"
	PermissionDeleteTenant = "delete-tenant"
)

// NewPermission crea una nueva instancia de permiso con valores por defecto
func NewPermission(name, description, section string) *Permission {
	now := time.Now()

	// Validar descripción como puntero
	var descPtr *string
	if description != "" {
		descPtr = &description
	}

	return &Permission{
		ID:          uuid.New(),
		Name:        name,
		Description: descPtr,
		Section:     section,
		Created:     now,
		Updated:     now,
	}
}

// IsSystemPermission verifica si es uno de los permisos predefinidos del sistema
func (p *Permission) IsSystemPermission() bool {
	systemPermissions := []string{
		PermissionListRole, PermissionCreateRole, PermissionUpdateRole, PermissionDeleteRole,
		PermissionListPermission, PermissionCreatePermission, PermissionUpdatePermission, PermissionDeletePermission,
		PermissionListUser, PermissionCreateUser, PermissionUpdateUser, PermissionDeleteUser,
		PermissionListTenant, PermissionCreateTenant, PermissionUpdateTenant, PermissionDeleteTenant,
	}

	for _, sysPermission := range systemPermissions {
		if p.Name == sysPermission {
			return true
		}
	}
	return false
}

// IsValidSection verifica si la sección es válida
func (p *Permission) IsValidSection() bool {
	validSections := []PermissionSection{
		SectionRole, SectionPermission, SectionUser, SectionTenant,
		SectionRegion, SectionCommune, SectionCountry, SectionCurrency,
		SectionBank, SectionPolicy, SectionCategory, SectionPolicyApprover,
		SectionPolicySubmitter, SectionPolicyReportStatus, SectionPolicyField,
		SectionReportExpense, SectionReportExpenseDetail, SectionReportExpenseRecord,
		SectionReportExpenseComment, SectionReportExpenseApprover, SectionExpenseGallery,
	}

	for _, validSection := range validSections {
		if string(validSection) == p.Section {
			return true
		}
	}
	return false
}

// GetPermissionsBySection retorna permisos agrupados por sección
func GetPermissionsBySection(section PermissionSection) []string {
	permissionMap := map[PermissionSection][]string{
		SectionRole: {
			PermissionListRole, PermissionCreateRole, PermissionUpdateRole, PermissionDeleteRole,
		},
		SectionPermission: {
			PermissionListPermission, PermissionCreatePermission, PermissionUpdatePermission, PermissionDeletePermission,
		},
		SectionUser: {
			PermissionListUser, PermissionCreateUser, PermissionUpdateUser, PermissionDeleteUser,
		},
		SectionTenant: {
			PermissionListTenant, PermissionCreateTenant, PermissionUpdateTenant, PermissionDeleteTenant,
		},
	}

	if permissions, exists := permissionMap[section]; exists {
		return permissions
	}
	return []string{}
}

// Update actualiza los campos modificables del permiso
func (p *Permission) Update(name, description, section string) {
	p.Name = name
	if description != "" {
		p.Description = &description
	} else {
		p.Description = nil
	}
	p.Section = section
	p.Updated = time.Now()
}

// SoftDelete marca el permiso como eliminado
func (p *Permission) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.Updated = now
}

// IsActive verifica si el permiso está activo (no eliminado)
func (p *Permission) IsActive() bool {
	return p.DeletedAt == nil
}
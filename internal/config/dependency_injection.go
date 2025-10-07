package config

import (
	"context"
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/middlewares"
	"github.com/JoseLuis21/mv-backend/internal/shared/email"
	"github.com/JoseLuis21/mv-backend/internal/shared/hasher"
	"github.com/JoseLuis21/mv-backend/internal/shared/tokens"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/services"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/adapters"
	tenantServices "github.com/JoseLuis21/mv-backend/internal/core/tenant/services"
	userAdapters "github.com/JoseLuis21/mv-backend/internal/core/user/adapters"
	userServices "github.com/JoseLuis21/mv-backend/internal/core/user/services"
	roleAdapters "github.com/JoseLuis21/mv-backend/internal/core/role/adapters"
	roleServices "github.com/JoseLuis21/mv-backend/internal/core/role/services"
	permissionAdapters "github.com/JoseLuis21/mv-backend/internal/core/permission/adapters"
	permissionServices "github.com/JoseLuis21/mv-backend/internal/core/permission/services"
	userRoleAdapters "github.com/JoseLuis21/mv-backend/internal/core/user_role/adapters"
	userRoleServices "github.com/JoseLuis21/mv-backend/internal/core/user_role/services"
	rolePermissionAdapters "github.com/JoseLuis21/mv-backend/internal/core/role_permission/adapters"
	rolePermissionServices "github.com/JoseLuis21/mv-backend/internal/core/role_permission/services"
	regionAdapters "github.com/JoseLuis21/mv-backend/internal/core/region/adapters"
	regionServices "github.com/JoseLuis21/mv-backend/internal/core/region/services"
	communeAdapters "github.com/JoseLuis21/mv-backend/internal/core/commune/adapters"
	communeServices "github.com/JoseLuis21/mv-backend/internal/core/commune/services"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
)

// Dependencies contiene todas las dependencias inyectadas de la aplicación
type Dependencies struct {
	// Controllers
	AuthController       *controllers.AuthController
	TenantController     *controllers.TenantController
	UserController       *controllers.UserController
	RoleController           *controllers.RoleController
	PermissionController     *controllers.PermissionController
	UserRoleController       *controllers.UserRoleController
	RolePermissionController *controllers.RolePermissionController
	RegionController         *controllers.RegionController
	CommuneController    *controllers.CommuneController

	// Middlewares
	RBACMiddleware *middlewares.RBACMiddleware

	// Infrastructure
	DBControl *postgresql.PostgresqlClient
	Validator *validatorapi.Validator
}

// NewDependencies crea e inyecta todas las dependencias del sistema
// Sigue el patrón de Dependency Injection para la arquitectura hexagonal
func NewDependencies(dbControl *postgresql.PostgresqlClient) (*Dependencies, error) {
	// 1. Crear infraestructura
	validator := validatorapi.NewValidator()

	// 2. Crear adaptadores (outer layer)
	userRepo := userAdapters.NewPostgreSQLUserRepository(dbControl)
	tenantRepo := adapters.NewPostgreSQLTenantRepository(dbControl)
	roleRepo := roleAdapters.NewPostgreSQLRoleRepository(dbControl)
	permissionRepo := permissionAdapters.NewPgPermissionRepository(dbControl)
	userRoleRepo := userRoleAdapters.NewPgUserRoleRepository(dbControl)
	rolePermissionRepo := rolePermissionAdapters.NewPgRolePermissionRepository(dbControl)
	regionRepo := regionAdapters.NewPostgreSQLRegionRepository(dbControl)
	communeRepo := communeAdapters.NewPostgreSQLCommuneRepository(dbControl)

	// 3. Crear servicios auxiliares usando módulos genéricos
	passwordHasher := hasher.NewService()
	tokenGenerator := tokens.NewService()
	emailService := email.NewService()

	// 4. Crear servicios de dominio (core layer) siguiendo patrón de referencia
	userService := userServices.NewUserService(userRepo)
	permissionService := permissionServices.NewPermissionService(permissionRepo)
	roleService := roleServices.NewRoleService(roleRepo, permissionService, rolePermissionRepo)
	userRoleService := userRoleServices.NewUserRoleService(userRoleRepo, userService, roleService)
	rolePermissionService := rolePermissionServices.NewRolePermissionService(rolePermissionRepo, roleService, permissionService)

	authService := services.NewAuthService(
		userService,
		passwordHasher,
		tokenGenerator,
		emailService,
		roleService,
		rolePermissionService,
	)

	tenantService := tenantServices.NewTenantService(
		tenantRepo,
		userService,
		roleService,
		userRoleService,
	)

	regionService := regionServices.NewRegionService(regionRepo)
	communeService := communeServices.NewCommuneService(communeRepo)

	// 5. Inicializar roles y permisos del sistema
	// Esto debe ejecutarse al inicio de la aplicación para garantizar que existan los roles necesarios
	ctx := context.Background()
	if err := permissionService.InitializeSystemPermissions(ctx); err != nil {
		return nil, err
	}
	if err := roleService.InitializeSystemRoles(ctx); err != nil {
		return nil, err
	}

	// 6. Crear middlewares
	rbacMiddleware := middlewares.NewRBACMiddleware(roleService, permissionService)

	// 7. Crear controllers (adapter layer - entrada HTTP)
	authController := controllers.NewAuthController(authService, validator)
	tenantController := controllers.NewTenantController(tenantService, authService, validator)
	userController := controllers.NewUserController(userService, roleService, userRoleService, validator)
	roleController := controllers.NewRoleController(roleService, validator)
	permissionController := controllers.NewPermissionController(permissionService, validator)
	userRoleController := controllers.NewUserRoleController(userRoleService, validator)
	rolePermissionController := controllers.NewRolePermissionController(rolePermissionService, validator)
	regionController := controllers.NewRegionController(regionService)
	communeController := controllers.NewCommuneController(communeService)

	return &Dependencies{
		AuthController:           authController,
		TenantController:         tenantController,
		UserController:           userController,
		RoleController:           roleController,
		PermissionController:     permissionController,
		UserRoleController:       userRoleController,
		RolePermissionController: rolePermissionController,
		RegionController:         regionController,
		CommuneController:        communeController,
		RBACMiddleware:           rbacMiddleware,
		DBControl:                dbControl,
		Validator:                validator,
	}, nil
}

// ConfigureRoutes configura todas las rutas con las dependencias inyectadas
func (d *Dependencies) ConfigureRoutes(app *fiber.App) {
	// Importar routes package en runtime para evitar imports circulares
	// En una implementación real, esto se haría directamente en main.go

	// Por ahora, esta función servirá como referencia para main.go
	// El main.go debería llamar a routes.AuthRoutes(app, d.AuthController)
}

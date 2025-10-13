package config

import (
	"context"
	"os"

	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/services"
	ocrServices "github.com/JoseLuis21/mv-backend/internal/core/ocr/services"
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
	expenseAdapters "github.com/JoseLuis21/mv-backend/internal/core/expense/adapters"
	expenseServices "github.com/JoseLuis21/mv-backend/internal/core/expense/services"
	policyAdapters "github.com/JoseLuis21/mv-backend/internal/core/policy/adapters"
	policyServices "github.com/JoseLuis21/mv-backend/internal/core/policy/services"
	"github.com/JoseLuis21/mv-backend/internal/libraries/ocr"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	redis_custom "github.com/JoseLuis21/mv-backend/internal/libraries/redis"
	"github.com/JoseLuis21/mv-backend/internal/middlewares"
	"github.com/JoseLuis21/mv-backend/internal/shared/email"
	"github.com/JoseLuis21/mv-backend/internal/shared/hasher"
	"github.com/JoseLuis21/mv-backend/internal/shared/tokens"
	"github.com/JoseLuis21/mv-backend/internal/shared/utils"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
)

// Dependencies contiene todas las dependencias inyectadas de la aplicación
type Dependencies struct {
	// Controllers
	AuthController           *controllers.AuthController
	TenantController         *controllers.TenantController
	UserController           *controllers.UserController
	RoleController           *controllers.RoleController
	PermissionController     *controllers.PermissionController
	UserRoleController       *controllers.UserRoleController
	RolePermissionController *controllers.RolePermissionController
	RegionController         *controllers.RegionController
	CommuneController        *controllers.CommuneController
	ExpenseController        *controllers.ExpenseController
	OCRController            *controllers.OCRController
	PolicyController         *controllers.PolicyController

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
	expenseRepo := expenseAdapters.NewPostgreSQLExpenseRepository(dbControl)
	categoryRepo := expenseAdapters.NewPostgreSQLCategoryRepository(dbControl)
	policyRepo := policyAdapters.NewPolicyRepository(dbControl)

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
	expenseService := expenseServices.NewExpenseService(expenseRepo, categoryRepo)
	policyService := policyServices.NewPolicyService(policyRepo)
	ruleEngine := policyServices.NewRuleEngine(policyRepo)

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
	expenseController := controllers.NewExpenseController(expenseService, validator)
	policyController := controllers.NewPolicyController(policyService, ruleEngine, validator)

	// 8. Configurar OCR dependencies (opcional si credenciales no están configuradas)
	var ocrController *controllers.OCRController
	googleProjectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	googleCredsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if googleProjectID != "" && googleCredsPath != "" {
		// Crear cliente de Redis para cache
		redisStorage := redis_custom.New(redis_custom.Config{
			Host:     utils.GetEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     6379,
			Database: 0,
		})
		redisClient := redisStorage.GetClient()

		// Crear cliente Google Vision
		ctx := context.Background()
		visionClient, err := ocr.NewGoogleVisionClient(ctx, ocr.GoogleVisionConfig{
			ProjectID:           googleProjectID,
			CredentialsFilePath: googleCredsPath,
		})
		if err != nil {
			// Log warning pero continuar sin OCR
			// TODO: agregar logging estructurado
		} else {
			// Crear OCR service
			ocrService := ocrServices.NewOCRService(ocrServices.OCRServiceConfig{
				VisionClient: visionClient,
				RedisClient:  redisClient,
				CacheEnabled: true,
			})

			// Crear OCR controller
			ocrController = controllers.NewOCRController(ocrService)
		}
	}

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
		ExpenseController:        expenseController,
		OCRController:            ocrController,
		PolicyController:         policyController,
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

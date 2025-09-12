package config

import (
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/shared/email"
	"github.com/JoseLuis21/mv-backend/internal/shared/hasher"
	"github.com/JoseLuis21/mv-backend/internal/shared/tokens"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/services"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/adapters"
	tenantServices "github.com/JoseLuis21/mv-backend/internal/core/tenant/services"
	userAdapters "github.com/JoseLuis21/mv-backend/internal/core/user/adapters"
	userServices "github.com/JoseLuis21/mv-backend/internal/core/user/services"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/gofiber/fiber/v2"
)

// Dependencies contiene todas las dependencias inyectadas de la aplicación
type Dependencies struct {
	// Controllers
	AuthController   *controllers.AuthController
	TenantController *controllers.TenantController

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

	// 3. Crear servicios auxiliares usando módulos genéricos
	passwordHasher := hasher.NewService()
	tokenGenerator := tokens.NewService()
	emailService := email.NewService()

	// 4. Crear servicios de dominio (core layer) siguiendo patrón de referencia
	userService := userServices.NewUserService(userRepo)
	
	authService := services.NewAuthService(
		userService,
		passwordHasher,
		tokenGenerator,
		emailService,
	)
	
	tenantService := tenantServices.NewTenantService(
		tenantRepo,
		userService,
	)

	// 5. Crear controllers (adapter layer - entrada HTTP)
	authController := controllers.NewAuthController(authService, validator)
	tenantController := controllers.NewTenantController(tenantService, authService, validator)

	return &Dependencies{
		AuthController:   authController,
		TenantController: tenantController,
		DBControl:        dbControl,
		Validator:        validator,
	}, nil
}

// ConfigureRoutes configura todas las rutas con las dependencias inyectadas
func (d *Dependencies) ConfigureRoutes(app *fiber.App) {
	// Importar routes package en runtime para evitar imports circulares
	// En una implementación real, esto se haría directamente en main.go

	// Por ahora, esta función servirá como referencia para main.go
	// El main.go debería llamar a routes.AuthRoutes(app, d.AuthController)
}

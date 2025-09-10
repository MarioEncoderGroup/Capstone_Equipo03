package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/core/auth/services"
	"github.com/JoseLuis21/mv-backend/internal/core/tenant/adapters"
	userAdapters "github.com/JoseLuis21/mv-backend/internal/core/user/adapters"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
)

// Dependencies contiene todas las dependencias inyectadas de la aplicación
type Dependencies struct {
	// Controllers
	AuthController *controllers.AuthController

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

	// 3. Crear servicios auxiliares
	passwordHasher := services.NewPasswordHasher()
	tokenGenerator := services.NewTokenGenerator()
	emailService := services.NewEmailService()

	// 4. Crear servicios de dominio (core layer)
	registerService := services.NewRegisterService(
		userRepo,
		tenantRepo,
		passwordHasher,
		tokenGenerator,
		emailService,
	)

	// 5. Crear controllers (adapter layer - entrada HTTP)
	authController := controllers.NewAuthController(registerService, validator)

	return &Dependencies{
		AuthController: authController,
		DBControl:      dbControl,
		Validator:      validator,
	}, nil
}

// ConfigureRoutes configura todas las rutas con las dependencias inyectadas
func (d *Dependencies) ConfigureRoutes(app *fiber.App) {
	// Importar routes package en runtime para evitar imports circulares
	// En una implementación real, esto se haría directamente en main.go
	
	// Por ahora, esta función servirá como referencia para main.go
	// El main.go debería llamar a routes.AuthRoutes(app, d.AuthController)
}
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
	"github.com/JoseLuis21/mv-backend/internal/libraries/ocr"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	redis_custom "github.com/JoseLuis21/mv-backend/internal/libraries/redis"
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
	AuthController   *controllers.AuthController
	TenantController *controllers.TenantController
	OCRController    *controllers.OCRController

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

	// 6. Configurar OCR dependencies (opcional si credenciales no están configuradas)
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
		AuthController:   authController,
		TenantController: tenantController,
		OCRController:    ocrController,
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

package server

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	redis_custom "github.com/JoseLuis21/mv-backend/internal/libraries/redis"
	"github.com/JoseLuis21/mv-backend/internal/middleware"
	"github.com/JoseLuis21/mv-backend/internal/routes"
	"github.com/JoseLuis21/mv-backend/internal/shared/validatorapi"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	port      string
	host      string
	dbControl *postgresql.PostgresqlClient
	dbTenant  *postgresql.PostgresqlClient
}

func NewServer(host, port string, dbControl *postgresql.PostgresqlClient, dbTenant *postgresql.PostgresqlClient) *Server {
	return &Server{
		host:      host,
		port:      port,
		dbControl: dbControl,
		dbTenant:  dbTenant,
	}
}

func (s *Server) Start() error {
	// Configure Redis for rate limiting and caching
	store := redis_custom.New(redis_custom.Config{
		Host:      getEnvOrDefault("REDIS_HOST", "localhost"),
		Port:      atoiDefault(os.Getenv("REDIS_PORT"), 6379),
		Username:  os.Getenv("REDIS_USERNAME"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		Database:  atoiDefault(os.Getenv("REDIS_DB"), 0),
		TLSConfig: nil,
		Reset:     false,
	})

	// Configure Fiber app for MisViaticos
	app := fiber.New(fiber.Config{
		AppName:           "MisViaticos API",
		CaseSensitive:     true,
		Concurrency:       500,
		JSONDecoder:       json.Unmarshal,
		JSONEncoder:       json.Marshal,
		EnablePrintRoutes: true,
		// Enhanced limits for expense receipts
		ReadBufferSize:  64 * 1024,        // 64KB for receipt uploads
		WriteBufferSize: 64 * 1024,        // 64KB buffer
		BodyLimit:       50 * 1024 * 1024, // 50MB for receipt files
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(validatorapi.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
				Code:    "MISVIATICOS_ERROR",
			})
		},
	})

	// Recovery middleware for production stability
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Enhanced logging for expense tracking
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n",
	}))

	// CORS configuration for MisViaticos frontend
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnvOrDefault("BASE_URL_FRONTEND", "http://localhost:3000"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Tenant-ID",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	// Rate limiting: 200 requests per minute for expense operations
	app.Use(limiter.New(limiter.Config{
		Max:        200,
		Expiration: time.Minute,
		Storage:    store,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Include tenant ID in rate limit key for multi-tenant fairness
			tenantID := c.Get("X-Tenant-ID", "global")
			return c.IP() + ":" + tenantID
		},
		LimitReached: func(c *fiber.Ctx) error {
			c.Set("Retry-After", "60")
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Rate limit exceeded - MisViaticos API",
				"message": "Too many requests. Please try again in a minute.",
				"rate_limit": fiber.Map{
					"limit":               c.Get("X-RateLimit-Limit"),
					"remaining":           "0",
					"retry_after_seconds": 60,
				},
			})
		},
		// Skip rate limiting for health checks
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/" || c.Path() == "/health"
		},
	}))

	// Middleware to pass rate limit headers to context
	app.Use(func(c *fiber.Ctx) error {
		if rem := c.Get("X-RateLimit-Remaining"); rem != "" {
			c.Locals("rate_limit_remaining", rem)
		}
		if lim := c.Get("X-RateLimit-Limit"); lim != "" {
			c.Locals("rate_limit_limit", lim)
		}
		return c.Next()
	})

	// Initialize validator for expense data validation
	validatorApi := &validatorapi.XValidator{
		Validator: validator.New(),
	}

	// Core middlewares
	app.Use(middleware.ValidatorMiddleware(validatorApi))
	app.Use(middleware.DatabaseControlMiddleware(s.dbControl))
	app.Use(middleware.DatabaseTenant1Middleware(s.dbTenant))

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "MisViaticos API is running",
			"version": "1.0.0",
			"status":  "healthy",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "misviaticos-backend",
		})
	})

	// Define public routes (authentication, registration)
	routes.PublicRoutes(app)

	// Define private routes (expenses, receipts, reports)
	routes.PrivateRoutes(app, s.dbControl)

	// Start server
	return app.Listen(s.host + ":" + s.port)
}

func atoiDefault(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
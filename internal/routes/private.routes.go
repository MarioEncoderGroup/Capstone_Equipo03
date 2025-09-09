package routes

import (
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes defines all private routes for MisViaticos API
func PrivateRoutes(app *fiber.App, dbControl *postgresql.PostgresqlClient) *fiber.App {
	// Create authentication middleware
	authMiddleware := middleware.AuthMiddleware(dbControl)

	// Private API routes group with authentication
	private := app.Group("/api/v1", authMiddleware)

	// Tenant management routes
	tenant := private.Group("/tenant")
	tenant.Post("/create", controllers.CreateTenant)
	tenant.Put("/update/:tenantID", controllers.UpdateTenant)
	tenant.Post("/select", controllers.SelectTenant)
	tenant.Get("/current", controllers.GetCurrentTenant)

	// User management routes
	users := private.Group("/users")
	users.Get("/profile", controllers.GetUserProfile)
	users.Put("/profile", controllers.UpdateUserProfile)
	users.Post("/change-password", controllers.ChangePassword)

	// Expense management routes (core MisViaticos functionality)
	expenses := private.Group("/expenses")
	expenses.Get("/", controllers.GetExpenses)
	expenses.Post("/", controllers.CreateExpense)
	expenses.Get("/:expenseID", controllers.GetExpense)
	expenses.Put("/:expenseID", controllers.UpdateExpense)
	expenses.Delete("/:expenseID", controllers.DeleteExpense)
	expenses.Post("/:expenseID/submit", controllers.SubmitExpense)
	expenses.Post("/:expenseID/approve", controllers.ApproveExpense)
	expenses.Post("/:expenseID/reject", controllers.RejectExpense)

	// Receipt management routes
	receipts := private.Group("/receipts")
	receipts.Post("/upload", controllers.UploadReceipt)
	receipts.Get("/:receiptID", controllers.GetReceipt)
	receipts.Delete("/:receiptID", controllers.DeleteReceipt)

	// Category management routes
	categories := private.Group("/categories")
	categories.Get("/", controllers.GetCategories)
	categories.Post("/", controllers.CreateCategory)
	categories.Put("/:categoryID", controllers.UpdateCategory)
	categories.Delete("/:categoryID", controllers.DeleteCategory)

	// Reporting routes
	reports := private.Group("/reports")
	reports.Get("/expenses", controllers.GetExpenseReport)
	reports.Get("/summary", controllers.GetExpenseSummary)
	reports.Get("/export/excel", controllers.ExportExpensesToExcel)
	reports.Get("/export/pdf", controllers.ExportExpensesToPDF)

	// Admin routes (require admin role)
	admin := private.Group("/admin")
	// admin.Use(middleware.RequireRole("admin")) // Would be implemented later
	admin.Get("/users", controllers.GetAllUsers)
	admin.Get("/expenses", controllers.GetAllExpenses)
	admin.Get("/analytics", controllers.GetAnalytics)

	return app
}
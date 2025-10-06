package routes

import (
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/middleware"
	"github.com/JoseLuis21/mv-backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes defines all private routes for MisViaticos API
func PrivateRoutes(app *fiber.App, dbControl *postgresql.PostgresqlClient, tenantController *controllers.TenantController, userController *controllers.UserController, roleController *controllers.RoleController, permissionController *controllers.PermissionController, userRoleController *controllers.UserRoleController, rolePermissionController *controllers.RolePermissionController, rbacMiddleware *middlewares.RBACMiddleware) *fiber.App {
	// Create authentication middleware
	authMiddleware := middleware.AuthMiddleware(dbControl)
	tenantMiddleware := middleware.RequireTenantMiddleware()

	// Private API routes group with authentication
	private := app.Group("/api/v1", authMiddleware)

	// Tenant management routes (sin RequireTenant porque son para selección inicial)
	tenant := private.Group("/tenant")
	tenant.Get("/status", tenantController.GetTenantStatus)                    // No requiere tenant seleccionado
	tenant.Post("/create", tenantController.CreateTenant)                       // No requiere tenant seleccionado
	tenant.Get("/", tenantController.GetTenantsByUser)                          // No requiere tenant seleccionado
	tenant.Post("/select/:tenantId", tenantController.SelectTenant)             // No requiere tenant seleccionado
	tenant.Get("/:tenantId/profile", tenantMiddleware, tenantController.GetTenantProfile)   // Requiere tenant
	tenant.Put("/:tenantId/profile", tenantMiddleware, tenantController.UpdateTenantProfile) // Requiere tenant

	// User management routes (requieren tenant seleccionado)
	users := private.Group("/users", tenantMiddleware)
	users.Get("/profile", userController.GetProfile)
	users.Put("/profile", userController.UpdateProfile)
	users.Post("/change-password", userController.ChangePassword)

	// Admin user management routes (requieren tenant + administrator role)
	adminUsers := private.Group("/admin/users", tenantMiddleware, rbacMiddleware.RequireRole("administrator"))
	adminUsers.Get("/", userController.GetUsers)
	adminUsers.Get("/:id", userController.GetUserByID)
	adminUsers.Post("/", rbacMiddleware.RequirePermission("create-user"), userController.CreateUser)
	adminUsers.Put("/:id", rbacMiddleware.RequirePermission("update-user"), userController.UpdateUser)
	adminUsers.Delete("/:id", rbacMiddleware.RequirePermission("delete-user"), userController.DeleteUser)

	// Role management routes (requieren tenant + administrator role)
	roles := private.Group("/roles", tenantMiddleware, rbacMiddleware.RequireRole("administrator"))
	roles.Get("/", roleController.GetRoles)
	roles.Get("/:id", roleController.GetRoleByID)
	roles.Post("/", rbacMiddleware.RequirePermission("create-role"), roleController.CreateRole)
	roles.Put("/:id", rbacMiddleware.RequirePermission("update-role"), roleController.UpdateRole)
	roles.Delete("/:id", rbacMiddleware.RequirePermission("delete-role"), roleController.DeleteRole)

	// Permission management routes (requieren tenant + administrator role)
	permissions := private.Group("/permissions", tenantMiddleware, rbacMiddleware.RequireRole("administrator"))
	permissions.Get("/", permissionController.GetPermissions)
	permissions.Get("/sections", permissionController.GetAvailableSections)
	permissions.Get("/grouped", permissionController.GetPermissionsGrouped)
	permissions.Get("/:id", permissionController.GetPermissionByID)
	permissions.Post("/", rbacMiddleware.RequireSystemAdmin(), permissionController.CreatePermission)
	permissions.Put("/:id", rbacMiddleware.RequireSystemAdmin(), permissionController.UpdatePermission)
	permissions.Delete("/:id", rbacMiddleware.RequireSystemAdmin(), permissionController.DeletePermission)

	// User-Role assignment routes (requieren tenant + administrator/manager role)
	userRoles := private.Group("/user-roles", tenantMiddleware, rbacMiddleware.RequireRole("administrator", "manager"))
	userRoles.Get("/users/:userID/roles", userRoleController.GetUserRoles)
	userRoles.Get("/roles/:roleID/users", userRoleController.GetRoleUsers)
	userRoles.Post("/assign", userRoleController.CreateUserRole)
	userRoles.Delete("/unassign", userRoleController.DeleteUserRole)
	userRoles.Put("/users/:userID/sync", userRoleController.SyncUserRoles)
	userRoles.Put("/roles/:roleID/sync", userRoleController.SyncRoleUsers)

	// Role-Permission assignment routes (requieren tenant + administrator role)
	rolePermissions := private.Group("/role-permissions", tenantMiddleware, rbacMiddleware.RequireRole("administrator"))
	rolePermissions.Get("/roles/:roleID", rolePermissionController.GetRolePermissions)
	rolePermissions.Get("/permissions/:permissionID", rolePermissionController.GetPermissionRoles)
	rolePermissions.Post("/assign", rolePermissionController.CreateRolePermission)
	rolePermissions.Delete("/unassign", rolePermissionController.DeleteRolePermission)
	rolePermissions.Put("/roles/:roleID/sync", rolePermissionController.SyncRolePermissions)
	rolePermissions.Get("/check", rolePermissionController.CheckRoleHasPermission)
	rolePermissions.Delete("/roles/:roleID/all", rolePermissionController.RemoveAllPermissionsFromRole)

	// TODO: Implement these tenant controllers
	// tenant.Post("/create", controllers.CreateTenant)
	// tenant.Put("/update/:tenantID", controllers.UpdateTenant)

	// TODO: Implement all remaining controllers and uncomment these routes
	/*

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
		admin := private.Group("/admin", rbacMiddleware.RequireAdminRole())
		admin.Get("/analytics", controllers.GetAnalytics)
		// TODO: Implement expense management routes
		// admin.Get("/expenses", controllers.GetAllExpenses)
	*/

	return app
}

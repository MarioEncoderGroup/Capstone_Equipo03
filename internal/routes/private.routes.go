package routes

import (
	"github.com/JoseLuis21/mv-backend/internal/controllers"
	"github.com/JoseLuis21/mv-backend/internal/libraries/postgresql"
	"github.com/JoseLuis21/mv-backend/internal/middleware"
	"github.com/JoseLuis21/mv-backend/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes defines all private routes for MisViaticos API
func PrivateRoutes(app *fiber.App, dbControl *postgresql.PostgresqlClient, tenantController *controllers.TenantController, userController *controllers.UserController, roleController *controllers.RoleController, permissionController *controllers.PermissionController, userRoleController *controllers.UserRoleController, rolePermissionController *controllers.RolePermissionController, expenseController *controllers.ExpenseController, policyController *controllers.PolicyController, reportController *controllers.ReportController, approvalController *controllers.ApprovalController, rbacMiddleware *middlewares.RBACMiddleware) *fiber.App {
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

	// Expense management routes (core MisViaticos functionality)
	expenses := private.Group("/expenses", tenantMiddleware)
	expenses.Get("/", expenseController.GetExpenses)
	expenses.Post("/", expenseController.CreateExpense)
	expenses.Get("/:id", expenseController.GetExpenseByID)
	expenses.Put("/:id", expenseController.UpdateExpense)
	expenses.Delete("/:id", expenseController.DeleteExpense)
	expenses.Post("/:id/receipts", expenseController.UploadReceipt)
	expenses.Post("/validate", policyController.ValidateExpense)

	// Receipt management routes
	receipts := private.Group("/receipts", tenantMiddleware)
	receipts.Delete("/:id", expenseController.DeleteReceipt)

	// Expense Report management routes (core approval workflow)
	reports := private.Group("/expense-reports", tenantMiddleware)
	reports.Post("/", reportController.CreateReport)
	reports.Get("/", reportController.GetUserReports)
	reports.Get("/:id", reportController.GetReportByID)
	reports.Put("/:id", reportController.UpdateReport)
	reports.Delete("/:id", reportController.DeleteReport)
	reports.Post("/:id/submit", reportController.SubmitReport)
	reports.Post("/:id/expenses", reportController.AddExpensesToReport)
	reports.Delete("/:id/expenses/:expenseId", reportController.RemoveExpenseFromReport)
	reports.Post("/:id/comments", reportController.AddComment)
	reports.Get("/:id/comments", reportController.GetReportComments)

	// Approval management routes (approval workflow)
	approvals := private.Group("/approvals", tenantMiddleware)
	approvals.Get("/pending", approvalController.GetPendingApprovals)
	approvals.Get("/reports/:id", approvalController.GetApprovalsByReport)
	approvals.Post("/:id/approve", approvalController.ApproveReport)
	approvals.Post("/:id/reject", approvalController.RejectReport)
	approvals.Post("/:id/escalate", approvalController.EscalateApproval)
	approvals.Get("/reports/:id/history", approvalController.GetApprovalHistory)

	return app
}

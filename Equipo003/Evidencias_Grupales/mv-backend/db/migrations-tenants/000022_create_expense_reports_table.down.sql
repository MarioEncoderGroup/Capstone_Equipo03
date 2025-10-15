-- Eliminar índices
DROP INDEX IF EXISTS idx_expense_reports_user_status;
DROP INDEX IF EXISTS idx_expense_reports_deleted;
DROP INDEX IF EXISTS idx_expense_reports_submission;
DROP INDEX IF EXISTS idx_expense_reports_status;
DROP INDEX IF EXISTS idx_expense_reports_policy;
DROP INDEX IF EXISTS idx_expense_reports_user;

-- Eliminar tabla
DROP TABLE IF EXISTS expense_reports;

-- Eliminar ENUM
DROP TYPE IF EXISTS report_status;

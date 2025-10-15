-- Eliminar índices
DROP INDEX IF EXISTS idx_approvals_report_level;
DROP INDEX IF EXISTS idx_approvals_pending_user;
DROP INDEX IF EXISTS idx_approvals_decision_date;
DROP INDEX IF EXISTS idx_approvals_level;
DROP INDEX IF EXISTS idx_approvals_status;
DROP INDEX IF EXISTS idx_approvals_approver;
DROP INDEX IF EXISTS idx_approvals_report;

-- Eliminar tabla
DROP TABLE IF EXISTS approvals;

-- Eliminar ENUM
DROP TYPE IF EXISTS approval_status;

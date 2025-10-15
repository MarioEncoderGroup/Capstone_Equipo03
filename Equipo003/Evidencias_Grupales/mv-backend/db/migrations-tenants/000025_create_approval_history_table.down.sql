-- Eliminar índices
DROP INDEX IF EXISTS idx_approval_history_metadata;
DROP INDEX IF EXISTS idx_approval_history_report_timeline;
DROP INDEX IF EXISTS idx_approval_history_created;
DROP INDEX IF EXISTS idx_approval_history_action;
DROP INDEX IF EXISTS idx_approval_history_actor;
DROP INDEX IF EXISTS idx_approval_history_report;
DROP INDEX IF EXISTS idx_approval_history_approval;

-- Eliminar tabla
DROP TABLE IF EXISTS approval_history;

-- Eliminar ENUM
DROP TYPE IF EXISTS approval_action;

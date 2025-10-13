-- Eliminar índices
DROP INDEX IF EXISTS idx_expense_comments_attachments;
DROP INDEX IF EXISTS idx_expense_comments_visible;
DROP INDEX IF EXISTS idx_expense_comments_thread;
DROP INDEX IF EXISTS idx_expense_comments_deleted;
DROP INDEX IF EXISTS idx_expense_comments_created;
DROP INDEX IF EXISTS idx_expense_comments_type;
DROP INDEX IF EXISTS idx_expense_comments_parent;
DROP INDEX IF EXISTS idx_expense_comments_user;
DROP INDEX IF EXISTS idx_expense_comments_expense;
DROP INDEX IF EXISTS idx_expense_comments_report;

-- Eliminar tabla
DROP TABLE IF EXISTS expense_comments;

-- Eliminar ENUM
DROP TYPE IF EXISTS comment_type;

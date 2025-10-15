-- Eliminar índices
DROP INDEX IF EXISTS idx_expense_report_items_sequence;
DROP INDEX IF EXISTS idx_expense_report_items_expense;
DROP INDEX IF EXISTS idx_expense_report_items_report;

-- Eliminar tabla
DROP TABLE IF EXISTS expense_report_items;

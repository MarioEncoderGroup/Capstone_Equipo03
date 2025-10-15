-- Tabla de items de reporte (relación entre gastos y reportes)
CREATE TABLE expense_report_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES expense_reports(id) ON DELETE CASCADE,
    expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    sequence_number INT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Un gasto no puede estar en múltiples reportes activos simultáneamente
    CONSTRAINT unique_expense_per_report UNIQUE (expense_id)
);

-- Índices para optimizar consultas
CREATE INDEX idx_expense_report_items_report ON expense_report_items(report_id);
CREATE INDEX idx_expense_report_items_expense ON expense_report_items(expense_id);
CREATE INDEX idx_expense_report_items_sequence ON expense_report_items(report_id, sequence_number);

-- Comentarios
COMMENT ON TABLE expense_report_items IS 'Relación many-to-one entre gastos y reportes de gastos';
COMMENT ON COLUMN expense_report_items.sequence_number IS 'Orden del gasto dentro del reporte';
COMMENT ON CONSTRAINT unique_expense_per_report ON expense_report_items IS 'Un gasto solo puede pertenecer a un reporte a la vez';

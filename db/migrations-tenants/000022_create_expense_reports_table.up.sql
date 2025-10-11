-- Crear ENUM para el estado del reporte de gastos
CREATE TYPE report_status AS ENUM (
    'draft',
    'submitted',
    'under_review',
    'approved',
    'rejected',
    'paid'
);

-- Tabla de reportes de gastos (agrupación de gastos para aprobación)
CREATE TABLE expense_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    policy_id UUID REFERENCES policies(id),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    status report_status DEFAULT 'draft',
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'CLP',
    submission_date TIMESTAMP,
    approval_date TIMESTAMP,
    payment_date TIMESTAMP,
    rejection_reason TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT check_total_amount CHECK (total_amount >= 0)
);

-- Índices para optimizar consultas
CREATE INDEX idx_expense_reports_user ON expense_reports(user_id);
CREATE INDEX idx_expense_reports_policy ON expense_reports(policy_id);
CREATE INDEX idx_expense_reports_status ON expense_reports(status);
CREATE INDEX idx_expense_reports_submission ON expense_reports(submission_date);
CREATE INDEX idx_expense_reports_deleted ON expense_reports(deleted_at);

-- Índice compuesto para queries comunes
CREATE INDEX idx_expense_reports_user_status ON expense_reports(user_id, status) WHERE deleted_at IS NULL;

-- Comentarios
COMMENT ON TABLE expense_reports IS 'Reportes de gastos que agrupan múltiples gastos para aprobación';
COMMENT ON COLUMN expense_reports.total_amount IS 'Monto total calculado automáticamente de todos los gastos incluidos';
COMMENT ON COLUMN expense_reports.status IS 'Estado del flujo de aprobación del reporte';

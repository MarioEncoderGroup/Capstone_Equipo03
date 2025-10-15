-- Crear ENUM para el estado de una aprobación
CREATE TYPE approval_status AS ENUM (
    'pending',
    'approved',
    'rejected',
    'escalated'
);

-- Tabla de aprobaciones (proceso de aprobación multi-nivel)
CREATE TABLE approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES expense_reports(id) ON DELETE CASCADE,
    approver_id UUID NOT NULL,
    level INT NOT NULL CHECK (level > 0),
    status approval_status DEFAULT 'pending',
    comments TEXT,
    approved_amount DECIMAL(12,2),
    decision_date TIMESTAMP,
    escalation_date TIMESTAMP,
    escalated_to UUID,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT check_approved_amount CHECK (approved_amount IS NULL OR approved_amount >= 0)
);

-- Índices para optimizar consultas
CREATE INDEX idx_approvals_report ON approvals(report_id);
CREATE INDEX idx_approvals_approver ON approvals(approver_id);
CREATE INDEX idx_approvals_status ON approvals(status);
CREATE INDEX idx_approvals_level ON approvals(level);
CREATE INDEX idx_approvals_decision_date ON approvals(decision_date);

-- Índice compuesto para queries de aprobaciones pendientes por usuario
CREATE INDEX idx_approvals_pending_user ON approvals(approver_id, status) WHERE status = 'pending';

-- Índice para reportes con nivel de aprobación
CREATE INDEX idx_approvals_report_level ON approvals(report_id, level);

-- Comentarios
COMMENT ON TABLE approvals IS 'Proceso de aprobación multi-nivel para reportes de gastos';
COMMENT ON COLUMN approvals.level IS 'Nivel de aprobación (1 = primer nivel, 2 = segundo nivel, etc.)';
COMMENT ON COLUMN approvals.escalation_date IS 'Fecha en que la aprobación fue escalada automáticamente';
COMMENT ON COLUMN approvals.escalated_to IS 'Usuario al que se escaló la aprobación';
COMMENT ON COLUMN approvals.approved_amount IS 'Monto aprobado (puede ser menor al solicitado)';

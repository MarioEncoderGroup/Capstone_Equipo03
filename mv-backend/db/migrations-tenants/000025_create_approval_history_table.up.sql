-- Crear ENUM para el tipo de acción en el historial
CREATE TYPE approval_action AS ENUM (
    'created',
    'approved',
    'rejected',
    'escalated',
    'reassigned',
    'commented'
);

-- Tabla de historial de aprobaciones (auditoría completa)
CREATE TABLE approval_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    approval_id UUID NOT NULL REFERENCES approvals(id) ON DELETE CASCADE,
    report_id UUID NOT NULL REFERENCES expense_reports(id) ON DELETE CASCADE,
    actor_id UUID NOT NULL,
    action approval_action NOT NULL,
    previous_status approval_status,
    new_status approval_status,
    comments TEXT,
    metadata JSONB,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Al menos uno de los estados debe estar presente para cambios de estado
    CONSTRAINT check_status_change CHECK (
        action IN ('commented', 'created') OR
        (previous_status IS NOT NULL AND new_status IS NOT NULL)
    )
);

-- Índices para optimizar consultas
CREATE INDEX idx_approval_history_approval ON approval_history(approval_id);
CREATE INDEX idx_approval_history_report ON approval_history(report_id);
CREATE INDEX idx_approval_history_actor ON approval_history(actor_id);
CREATE INDEX idx_approval_history_action ON approval_history(action);
CREATE INDEX idx_approval_history_created ON approval_history(created DESC);

-- Índice compuesto para timeline de un reporte
CREATE INDEX idx_approval_history_report_timeline ON approval_history(report_id, created DESC);

-- Índice GIN para búsquedas en metadata JSONB
CREATE INDEX idx_approval_history_metadata ON approval_history USING GIN (metadata);

-- Comentarios
COMMENT ON TABLE approval_history IS 'Historial completo de todas las acciones en el proceso de aprobación';
COMMENT ON COLUMN approval_history.actor_id IS 'Usuario que realizó la acción';
COMMENT ON COLUMN approval_history.action IS 'Tipo de acción realizada';
COMMENT ON COLUMN approval_history.metadata IS 'Datos adicionales en formato JSON (información contextual)';

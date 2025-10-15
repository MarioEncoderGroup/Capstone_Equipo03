-- Crear ENUM para el tipo de comentario
CREATE TYPE comment_type AS ENUM (
    'general',
    'question',
    'clarification',
    'approval_note',
    'rejection_note',
    'system'
);

-- Tabla de comentarios para gastos y reportes
CREATE TABLE expense_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID REFERENCES expense_reports(id) ON DELETE CASCADE,
    expense_id UUID REFERENCES expenses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    comment_type comment_type DEFAULT 'general',
    content TEXT NOT NULL,
    parent_id UUID REFERENCES expense_comments(id) ON DELETE CASCADE,
    is_internal BOOLEAN DEFAULT false,
    attachments JSONB,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- El comentario debe estar asociado al menos a un reporte o gasto
    CONSTRAINT check_association CHECK (
        report_id IS NOT NULL OR expense_id IS NOT NULL
    ),

    -- Si tiene parent, debe ser del mismo reporte/gasto
    CONSTRAINT check_parent_context CHECK (
        parent_id IS NULL OR
        (report_id IS NOT NULL AND expense_id IS NULL) OR
        (expense_id IS NOT NULL AND report_id IS NULL)
    )
);

-- Índices para optimizar consultas
CREATE INDEX idx_expense_comments_report ON expense_comments(report_id) WHERE report_id IS NOT NULL;
CREATE INDEX idx_expense_comments_expense ON expense_comments(expense_id) WHERE expense_id IS NOT NULL;
CREATE INDEX idx_expense_comments_user ON expense_comments(user_id);
CREATE INDEX idx_expense_comments_parent ON expense_comments(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_expense_comments_type ON expense_comments(comment_type);
CREATE INDEX idx_expense_comments_created ON expense_comments(created DESC);
CREATE INDEX idx_expense_comments_deleted ON expense_comments(deleted_at);

-- Índice compuesto para thread de comentarios
CREATE INDEX idx_expense_comments_thread ON expense_comments(parent_id, created) WHERE deleted_at IS NULL;

-- Índice para comentarios visibles (no internos, no eliminados)
CREATE INDEX idx_expense_comments_visible ON expense_comments(report_id, expense_id, created DESC)
    WHERE is_internal = false AND deleted_at IS NULL;

-- Índice GIN para búsquedas en attachments JSONB
CREATE INDEX idx_expense_comments_attachments ON expense_comments USING GIN (attachments);

-- Comentarios
COMMENT ON TABLE expense_comments IS 'Comentarios y comunicación sobre gastos y reportes';
COMMENT ON COLUMN expense_comments.is_internal IS 'Si es true, solo visible para aprobadores y administradores';
COMMENT ON COLUMN expense_comments.attachments IS 'Array JSON de archivos adjuntos con metadata';
COMMENT ON COLUMN expense_comments.parent_id IS 'Comentario padre para respuestas en thread';
COMMENT ON COLUMN expense_comments.comment_type IS 'Tipo de comentario para categorización';

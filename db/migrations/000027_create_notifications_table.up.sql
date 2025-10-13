-- Crear tipo ENUM para tipos de notificación
CREATE TYPE notification_type AS ENUM (
    'expense_approved',
    'expense_rejected',
    'expense_submitted',
    'approval_needed',
    'approval_approved',
    'approval_rejected',
    'approval_escalated',
    'report_submitted',
    'report_approved',
    'report_rejected',
    'comment_added',
    'system_notification'
);

-- Tabla de notificaciones para usuarios
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type notification_type NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    data JSONB DEFAULT '{}'::jsonb,

    -- Información de la entidad relacionada
    related_entity_id UUID,
    related_entity_type VARCHAR(50),

    -- Estado de lectura
    read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,

    -- Auditoría
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT check_read_at CHECK (read_at IS NULL OR read = true)
);

-- Índices para optimizar búsquedas
CREATE INDEX idx_notifications_user ON notifications(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_type ON notifications(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_read ON notifications(read) WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_created ON notifications(created DESC);
CREATE INDEX idx_notifications_deleted ON notifications(deleted_at) WHERE deleted_at IS NULL;

-- Índices compuestos para queries comunes
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, created DESC)
    WHERE deleted_at IS NULL AND read = false;
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created DESC)
    WHERE deleted_at IS NULL;
CREATE INDEX idx_notifications_related_entity ON notifications(related_entity_type, related_entity_id)
    WHERE deleted_at IS NULL;

-- Índice para búsquedas por entidad relacionada y usuario
CREATE INDEX idx_notifications_user_entity ON notifications(user_id, related_entity_type, related_entity_id)
    WHERE deleted_at IS NULL;

-- Comentarios
COMMENT ON TABLE notifications IS 'Notificaciones del sistema para usuarios';
COMMENT ON COLUMN notifications.user_id IS 'Usuario destinatario de la notificación';
COMMENT ON COLUMN notifications.type IS 'Tipo de notificación (expense_approved, approval_needed, etc.)';
COMMENT ON COLUMN notifications.title IS 'Título breve de la notificación';
COMMENT ON COLUMN notifications.message IS 'Mensaje completo de la notificación';
COMMENT ON COLUMN notifications.data IS 'Datos adicionales en formato JSON';
COMMENT ON COLUMN notifications.related_entity_id IS 'ID de la entidad relacionada (expense, approval, report)';
COMMENT ON COLUMN notifications.related_entity_type IS 'Tipo de entidad relacionada (expense, approval, report)';
COMMENT ON COLUMN notifications.read IS 'Indica si la notificación ha sido leída';
COMMENT ON COLUMN notifications.read_at IS 'Fecha y hora en que se leyó la notificación';

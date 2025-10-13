-- ========== POLÍTICAS ==========
CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    policy_type VARCHAR(50) NOT NULL,  -- travel, daily, project
    is_active BOOLEAN DEFAULT true,
    config JSONB,                      -- Configuración flexible en JSON
    created_by UUID NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT check_policy_type CHECK (policy_type IN ('travel', 'daily', 'project'))
);

-- Índices
CREATE INDEX idx_policies_type ON policies(policy_type);
CREATE INDEX idx_policies_active ON policies(is_active) WHERE is_active = true;
CREATE INDEX idx_policies_created_by ON policies(created_by);
CREATE INDEX idx_policies_deleted ON policies(deleted_at) WHERE deleted_at IS NULL;

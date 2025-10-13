-- ========== REGLAS DE POLÍTICA ==========
CREATE TABLE policy_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL,
    category_id UUID NOT NULL,
    rule_type VARCHAR(50) NOT NULL,       -- limit, auto_approve, require_approval
    condition JSONB,                       -- Condiciones en JSON
    action JSONB,                          -- Acción a tomar
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_policy_rules_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE,
    CONSTRAINT check_rule_type CHECK (rule_type IN ('limit', 'auto_approve', 'require_approval', 'reject'))
);

-- Índices
CREATE INDEX idx_policy_rules_policy_id ON policy_rules(policy_id);
CREATE INDEX idx_policy_rules_category_id ON policy_rules(category_id);
CREATE INDEX idx_policy_rules_type ON policy_rules(rule_type);
CREATE INDEX idx_policy_rules_active ON policy_rules(is_active) WHERE is_active = true;
CREATE INDEX idx_policy_rules_priority ON policy_rules(priority DESC);

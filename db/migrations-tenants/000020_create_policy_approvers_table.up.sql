-- ========== APROBADORES DE POLÍTICA ==========
CREATE TABLE policy_approvers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL,
    user_id UUID NOT NULL,
    level INT NOT NULL,                    -- Nivel de aprobación (1, 2, 3...)
    amount_min DECIMAL(12,2),
    amount_max DECIMAL(12,2),
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_policy_approvers_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE,
    CONSTRAINT check_approver_level CHECK (level > 0),
    CONSTRAINT check_approver_amounts CHECK (
        (amount_min IS NULL OR amount_min >= 0) AND
        (amount_max IS NULL OR amount_max >= 0) AND
        (amount_min IS NULL OR amount_max IS NULL OR amount_max >= amount_min)
    )
);

-- Índices
CREATE INDEX idx_policy_approvers_policy_id ON policy_approvers(policy_id);
CREATE INDEX idx_policy_approvers_user_id ON policy_approvers(user_id);
CREATE INDEX idx_policy_approvers_level ON policy_approvers(level);
CREATE INDEX idx_policy_approvers_amounts ON policy_approvers(amount_min, amount_max);

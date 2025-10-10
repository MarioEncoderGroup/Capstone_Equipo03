-- ========== EMPLEADOS QUE PUEDEN USAR LA POLÍTICA ==========
CREATE TABLE policy_submitters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    policy_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role_id UUID,                          -- Por rol o usuario específico
    department VARCHAR(100),
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_policy_submitters_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE,
    CONSTRAINT check_submitter_assignment CHECK (
        (user_id IS NOT NULL OR role_id IS NOT NULL OR department IS NOT NULL)
    ),
    UNIQUE(policy_id, user_id)
);

-- Índices
CREATE INDEX idx_policy_submitters_policy_id ON policy_submitters(policy_id);
CREATE INDEX idx_policy_submitters_user_id ON policy_submitters(user_id);
CREATE INDEX idx_policy_submitters_role_id ON policy_submitters(role_id);
CREATE INDEX idx_policy_submitters_department ON policy_submitters(department);

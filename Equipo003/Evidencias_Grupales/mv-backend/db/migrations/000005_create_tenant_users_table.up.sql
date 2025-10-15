-- Crear Tabla para tenant_users
CREATE TABLE tenant_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    UNIQUE (tenant_id, user_id)
);

-- Crear un índice para mejorar la búsqueda por tenant_id
CREATE INDEX idx_tenant_users_tenant_id ON tenant_users (tenant_id);

-- Crear un índice para mejorar la búsqueda por user_id
CREATE INDEX idx_tenant_users_user_id ON tenant_users (user_id);

-- Crear un índice para mejorar la búsqueda por tenant_id y user_id
CREATE INDEX idx_tenant_users_tenant_user ON tenant_users (tenant_id, user_id);

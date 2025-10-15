-- Crear Tabla UserRoles
CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    tenant_id UUID NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    UNIQUE (user_id, role_id)
);

-- Crear un índice para mejorar la búsqueda por user_id
CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);

-- Crear un índice para mejorar la búsqueda por role_id
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

-- Crear un índice para mejorar la búsqueda por user_id y role_id
CREATE INDEX idx_user_roles_user_role ON user_roles (user_id, role_id);

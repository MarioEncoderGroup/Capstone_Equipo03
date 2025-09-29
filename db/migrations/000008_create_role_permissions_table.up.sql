-- Crear Tabla RolePermissions
CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL,
    UNIQUE (role_id, permission_id)
);

-- Crear un índice para mejorar la búsqueda por role_id
CREATE INDEX idx_role_permissions_role_id ON role_permissions (role_id);

-- Crear un índice para mejorar la búsqueda por permission_id
CREATE INDEX idx_role_permissions_permission_id ON role_permissions (permission_id);

-- Crear un índice para mejorar la búsqueda por role_id y permission_id
CREATE INDEX idx_role_permissions_role_permission ON role_permissions (role_id, permission_id);

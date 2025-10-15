-- Crear Tabla Roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    tenant_id UUID NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Crear un índice para mejorar la búsqueda por name
CREATE INDEX idx_roles_name ON roles (name);

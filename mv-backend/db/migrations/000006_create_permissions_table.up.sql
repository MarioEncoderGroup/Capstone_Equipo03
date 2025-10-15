-- Crear Tabla Permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name VARCHAR(50) NOT NULL UNIQUE,   
    description TEXT NULL,
    section VARCHAR(50) NOT NULL, -- sección a la que pertenece el permiso
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Crear un índice para mejorar la búsqueda por name
CREATE INDEX idx_permissions_name ON permissions (name);

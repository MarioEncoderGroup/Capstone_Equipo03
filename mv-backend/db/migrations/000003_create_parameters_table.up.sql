-- Crear Tabla parameters
CREATE TABLE parameters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name VARCHAR(50) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Crear un índice para mejorar la búsqueda por name
CREATE INDEX idx_parameters_name ON parameters (name);

-- Insertar parámetro inicial
INSERT INTO parameters (name, value) VALUES ('TENANT_NODE', '1');

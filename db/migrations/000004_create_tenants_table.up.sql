DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tenant_status') THEN
        CREATE TYPE tenant_status AS ENUM ('active', 'inactive', 'suspended');
    END IF;
END$$;

-- Crear Tabla para tenants
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    rut VARCHAR(20) NOT NULL,           -- RUT del negocio
    business_name VARCHAR(150) NOT NULL,  -- nombre del negocio
    email VARCHAR(150) NOT NULL,         -- correo electrónico del negocio
    phone VARCHAR(20) NOT NULL,         -- teléfono del negocio
    address VARCHAR(200) NOT NULL,       -- dirección del negocio
    website VARCHAR(150) NOT NULL,    -- sitio web del negocio
    logo TEXT NULL,         -- logo del negocio
    region_id Char(2) NOT NULL,  -- ID de la región (AP, TA, ...)
    commune_id VARCHAR(100) NOT NULL,  -- ID de la comuna (AP, TA, ...)
    country_id UUID NOT NULL,  -- país del negocio
    status      tenant_status   NOT NULL DEFAULT 'active', -- ENUM status
    node_number INT NOT NULL,  -- número de nodo
    tenant_name TEXT NOT NULL,      -- nombre de BD
    created_by UUID NOT NULL,   -- ID del usuario que creó el tenant
    updated_by UUID NOT NULL,   -- ID del usuario que actualizó el tenant
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Crear un índice para mejorar la búsqueda por business_name
CREATE INDEX idx_tenants_business_name ON tenants (business_name);

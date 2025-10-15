-- Create Table for countries
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name VARCHAR(100) NOT NULL,  -- Nombre del país
    code VARCHAR(10) NOT NULL,    -- Código del país (ISO 3166-1 alpha-2)
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
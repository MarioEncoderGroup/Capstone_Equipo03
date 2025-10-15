-- Create Table for currency
CREATE TABLE currencies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    name VARCHAR(100) NOT NULL,  -- Nombre de la moneda
    code VARCHAR(10) NOT NULL,    -- Código de la moneda (ISO 4217)
    symbol VARCHAR(10) NOT NULL,  -- Símbolo de la moneda
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
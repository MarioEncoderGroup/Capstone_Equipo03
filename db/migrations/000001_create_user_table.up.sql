-- Crear Tabla users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    username VARCHAR(50) NOT NULL UNIQUE,
    phone VARCHAR(20) NULL,
    full_name VARCHAR(200) NOT NULL,
    identification_number VARCHAR(50) NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    email_token VARCHAR(255) NULL,
    email_token_expires TIMESTAMP NULL,
    email_verified BOOLEAN DEFAULT FALSE NOT NULL,
    password VARCHAR(255) NOT NULL,
    password_reset_token VARCHAR(255) NULL,
    password_reset_expires TIMESTAMP NULL,
    last_password_change TIMESTAMP NULL,
    last_login TIMESTAMP NULL,
    bank_id UUID NULL,
    bank_account_number VARCHAR(50) NULL,
    bank_account_type VARCHAR(50) NULL, -- Cuenta corriente, Cuenta vista, Cuenta de Ahorro
    image_url TEXT NULL,
    is_active BOOLEAN DEFAULT FALSE NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

-- Crear un índice para mejorar la búsqueda por username
CREATE INDEX idx_users_username ON users (username);

-- Crear un índice para mejorar la búsqueda por email
CREATE INDEX idx_users_email ON users (email);

-- Crear un índice para mejorar la búsqueda por email_token
CREATE INDEX idx_users_email_token ON users (email_token);

-- Crear un índice para mejorar la búsqueda por password_reset_token
CREATE INDEX idx_users_password_reset_token ON users (password_reset_token);

-- Crear un índice para mejorar la búsqueda por is_active
CREATE INDEX idx_users_is_active ON users (is_active);
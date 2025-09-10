-- Crear la tabla user_refresh_tokens
CREATE TABLE user_refresh_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  token_hash TEXT NOT NULL UNIQUE,
  node_number INT NOT NULL,      -- número de nodo
  db_name TEXT NOT NULL,         -- nombre de BD
  created TIMESTAMP NOT NULL DEFAULT now(),
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP NULL
);

-- Crear un índice para mejorar la búsqueda por user_id
CREATE INDEX idx_refresh_tokens_user_id ON user_refresh_tokens (user_id);

-- Crear un índice para mejorar la búsqueda por token
CREATE INDEX idx_refresh_tokens_token ON user_refresh_tokens (token_hash);
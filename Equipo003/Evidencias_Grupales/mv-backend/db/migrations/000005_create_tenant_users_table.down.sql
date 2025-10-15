-- Eliminar índices
DROP INDEX IF EXISTS idx_tenant_users_tenant_user;
DROP INDEX IF EXISTS idx_tenant_users_user_id;
DROP INDEX IF EXISTS idx_tenant_users_tenant_id;

-- Eliminar tabla tenant_users
DROP TABLE IF EXISTS tenant_users;



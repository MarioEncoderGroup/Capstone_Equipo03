-- Eliminar índices
DROP INDEX IF EXISTS idx_user_roles_user_role;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_user_roles_user_id;

-- Eliminar tabla user_roles
DROP TABLE IF EXISTS user_roles;

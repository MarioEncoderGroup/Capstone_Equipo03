-- Eliminar índices
DROP INDEX IF EXISTS idx_role_permissions_role_permission;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;

-- Eliminar tabla role_permissions
DROP TABLE IF EXISTS role_permissions;

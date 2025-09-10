INSERT INTO role_permissions (role_id, permission_id)
SELECT '0197c20a-4578-7bd3-a497-31d5ae6e53d2', id
FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;
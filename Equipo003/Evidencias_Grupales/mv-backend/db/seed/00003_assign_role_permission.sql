-- Asignar TODOS los permisos al rol administrator
INSERT INTO role_permissions (role_id, permission_id)
SELECT '0197c20a-4578-7bd3-a497-31d5ae6e53d2', id
FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Asignar permisos específicos al rol approver (aprobador de gastos)
INSERT INTO role_permissions (role_id, permission_id)
SELECT '0197e595-e807-72f7-bd91-3b5b8c3ade98', id
FROM permissions
WHERE name IN (
    -- Gestión de usuarios (para ver quiénes son los que rinden)
    'list-user',
    
    -- Gestión de reportes de gastos (aprobar/revisar)
    'list-report-expense',
    'update-report-expense',
    
    -- Gestión de aprobadores de reportes
    'list-report-expense-approver',
    'create-report-expense-approver',
    'update-report-expense-approver',
    'delete-report-expense-approver',
    
    -- Comentarios en reportes de gastos
    'list-report-expense-comment',
    'create-report-expense-comment',
    'update-report-expense-comment',
    
    -- Detalles de reportes de gastos
    'list-report-expense-details',
    
    -- Registros de reportes de gastos
    'list-report-expense-record',
    
    -- Galerías de gastos (ver comprobantes)
    'list-expense-gallery',
    
    -- Estados de reportes de política
    'list-policy-report-status',
    'update-policy-report-status',
    
    -- Ver políticas
    'list-policy',
    
    -- Ver categorías
    'list-category'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Asignar permisos específicos al rol expense-submitter (rendidor de gastos)
INSERT INTO role_permissions (role_id, permission_id)
SELECT '0197e596-216e-74fc-8224-7692deb6295c', id
FROM permissions
WHERE name IN (
    -- Ver su propio usuario
    'list-user',
    
    -- Gestión de reportes de gastos (crear y editar sus propios reportes)
    'list-report-expense',
    'create-report-expense',
    'update-report-expense',
    'delete-report-expense',
    
    -- Detalles de reportes de gastos
    'list-report-expense-details',
    'create-report-expense-details',
    'update-report-expense-details',
    'delete-report-expense-details',
    
    -- Registros de reportes de gastos
    'list-report-expense-record',
    'create-report-expense-record',
    'update-report-expense-record',
    'delete-report-expense-record',
    
    -- Comentarios en reportes de gastos
    'list-report-expense-comment',
    'create-report-expense-comment',
    
    -- Galerías de gastos (subir comprobantes)
    'list-expense-gallery',
    'create-expense-gallery',
    'update-expense-gallery',
    'delete-expense-gallery',
    
    -- Ver políticas aplicables
    'list-policy',
    
    -- Ver postulantes de políticas
    'list-policy-submitter',
    
    -- Ver categorías
    'list-category',
    
    -- Ver estados de reportes
    'list-policy-report-status',
    
    -- Ver bancos (para datos de reembolso)
    'list-bank'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;
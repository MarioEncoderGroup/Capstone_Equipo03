INSERT INTO roles (id, name, description)
VALUES ('0197c20a-4578-7bd3-a497-31d5ae6e53d2', 'administrator', 'Rol con acceso completo a todas las funciones')
ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (id, name, description)
VALUES ('0197e595-e807-72f7-bd91-3b5b8c3ade98', 'approver', 'Rol con acceso completo a todas las funciones de aprobador')
ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (id, name, description)
VALUES ('0197e596-216e-74fc-8224-7692deb6295c', 'expense-submitter', 'Rol con acceso completo a todas las funciones de rendidor')
ON CONFLICT (name) DO NOTHING;
INSERT INTO parameters (name, value)
VALUES
  ('TENANT_NODE', '1')
ON CONFLICT (name) DO NOTHING;
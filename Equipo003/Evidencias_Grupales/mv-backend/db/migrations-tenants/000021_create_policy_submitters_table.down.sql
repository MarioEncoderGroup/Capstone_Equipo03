-- Drop indexes
DROP INDEX IF EXISTS idx_policy_submitters_department;
DROP INDEX IF EXISTS idx_policy_submitters_role_id;
DROP INDEX IF EXISTS idx_policy_submitters_user_id;
DROP INDEX IF EXISTS idx_policy_submitters_policy_id;

-- Drop table
DROP TABLE IF EXISTS policy_submitters;

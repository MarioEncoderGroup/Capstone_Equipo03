-- Drop indexes
DROP INDEX IF EXISTS idx_policy_approvers_amounts;
DROP INDEX IF EXISTS idx_policy_approvers_level;
DROP INDEX IF EXISTS idx_policy_approvers_user_id;
DROP INDEX IF EXISTS idx_policy_approvers_policy_id;

-- Drop table
DROP TABLE IF EXISTS policy_approvers;

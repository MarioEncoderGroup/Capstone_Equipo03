-- Drop indexes
DROP INDEX IF EXISTS idx_policy_rules_priority;
DROP INDEX IF EXISTS idx_policy_rules_active;
DROP INDEX IF EXISTS idx_policy_rules_type;
DROP INDEX IF EXISTS idx_policy_rules_category_id;
DROP INDEX IF EXISTS idx_policy_rules_policy_id;

-- Drop table
DROP TABLE IF EXISTS policy_rules;

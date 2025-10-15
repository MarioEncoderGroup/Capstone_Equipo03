-- Drop indexes
DROP INDEX IF EXISTS idx_policies_deleted;
DROP INDEX IF EXISTS idx_policies_created_by;
DROP INDEX IF EXISTS idx_policies_active;
DROP INDEX IF EXISTS idx_policies_type;

-- Drop table
DROP TABLE IF EXISTS policies;

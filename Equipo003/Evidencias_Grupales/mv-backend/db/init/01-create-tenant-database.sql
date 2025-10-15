-- Initialize tenant database for MisViaticos
-- This script runs automatically when PostgreSQL container starts

-- Create the first tenant database
CREATE DATABASE misviaticos_tenant_1
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

-- Grant all privileges to postgres user
GRANT ALL PRIVILEGES ON DATABASE misviaticos_tenant_1 TO postgres;

-- Log successful creation
\echo 'Successfully created misviaticos_tenant_1 database'

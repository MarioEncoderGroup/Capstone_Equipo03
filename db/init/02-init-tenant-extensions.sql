-- Initialize UUID extension and custom functions for tenant database
-- This script runs after the database creation

-- Connect to the tenant database
\c misviaticos_tenant_1

-- Create UUID extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create or replace the uuid_generate_v7 function
CREATE OR REPLACE FUNCTION uuid_generate_v7()
RETURNS uuid
AS $$
SELECT encode(
    set_bit(
        set_bit(
            overlay(
                uuid_send(gen_random_uuid())
                placing substring(
                    int8send(
                        floor(extract(epoch from clock_timestamp()) * 1000)::bigint
                    ) from 3
                )
                from 1 for 6
            ),
            52, 1
        ),
        53, 1
    ),
    'hex'
)::uuid;
$$ LANGUAGE SQL VOLATILE;

-- Log successful initialization
\echo 'Successfully initialized extensions for misviaticos_tenant_1'

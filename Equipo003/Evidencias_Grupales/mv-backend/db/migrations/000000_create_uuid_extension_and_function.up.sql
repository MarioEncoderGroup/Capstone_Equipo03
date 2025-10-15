-- Crear la extensión si no existe
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Crear o reemplazar la función uuid_generate_v7
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

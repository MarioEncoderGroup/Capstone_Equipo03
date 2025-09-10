-- Generate seed data for the banks table

-- Get the country_id for Chile
DO $$
DECLARE
    chile_id UUID;
BEGIN
    -- Get the country_id for Chile
    SELECT id INTO chile_id FROM countries WHERE name = 'Chile';

    -- Insert banks for Chile
    INSERT INTO banks (name, country_id) VALUES
    ('Banco de Chile', chile_id),
    ('Banco Santander-Chile', chile_id),
    ('Banco de Crédito e Inversiones (BCI)', chile_id),
    ('BancoEstado', chile_id),
    ('Scotiabank Chile', chile_id),
    ('Itaú Corpbanca', chile_id),
    ('Banco Bice', chile_id),
    ('Banco Security', chile_id),
    ('Banco Falabella', chile_id),
    ('Banco Ripley', chile_id);
END $$;
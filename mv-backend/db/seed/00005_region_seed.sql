/* ===========
   2.  REGIONES DE CHILE
   =========== */
-- Seed idempotente: puede ejecutarse múltiples veces sin errores

INSERT INTO region (id, number, roman_number, name) VALUES
('AP', 15, 'XV' , 'Arica y Parinacota'),
('TA',  1, 'I'  , 'Tarapacá'),
('AN',  2, 'II' , 'Antofagasta'),
('AT',  3, 'III', 'Atacama'),
('CO',  4, 'IV' , 'Coquimbo'),
('VA',  5, 'V'  , 'Valparaíso'),
('RM', 13, 'XIII', 'Metropolitana de Santiago'),
('LI',  6, 'VI' , 'Libertador Gral. Bernardo O''Higgins'),
('ML',  7, 'VII', 'Maule'),
('NB', 16, 'XVI', 'Ñuble'),
('BI',  8, 'VIII', 'Biobío'),
('AR',  9, 'IX' , 'Araucanía'),
('LR', 14, 'XIV', 'Los Ríos'),
('LL', 10, 'X'  , 'Los Lagos'),
('AI', 11, 'XI' , 'Aisén del Gral. Carlos Ibáñez del Campo'),
('MA', 12, 'XII', 'Magallanes y de la Antártica Chilena')
ON CONFLICT (id) DO NOTHING;
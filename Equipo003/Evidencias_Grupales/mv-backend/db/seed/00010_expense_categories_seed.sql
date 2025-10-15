-- Seed de categorías de gastos empresariales
-- Categorías principales y subcategorías comunes

-- 1. TRANSPORTE
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('11111111-1111-1111-1111-111111111111', 'Transporte', 'Gastos de movilización y transporte', 'car', '#3B82F6', NULL, 50000, 500000, true, true);

-- Subcategorías de Transporte
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Taxi', 'Servicio de taxi', 'taxi', '#60A5FA', '11111111-1111-1111-1111-111111111111', 30000, 300000, true, true),
('Uber/Cabify', 'Aplicaciones de transporte', 'smartphone', '#60A5FA', '11111111-1111-1111-1111-111111111111', 30000, 300000, true, true),
('Combustible', 'Gasolina y diesel', 'fuel', '#60A5FA', '11111111-1111-1111-1111-111111111111', 60000, 400000, true, true),
('Peajes', 'Cobro de autopistas', 'road', '#60A5FA', '11111111-1111-1111-1111-111111111111', 10000, 50000, true, true),
('Estacionamiento', 'Parking y estacionamiento', 'parking', '#60A5FA', '11111111-1111-1111-1111-111111111111', 15000, 100000, true, true);

-- 2. ALIMENTACIÓN
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('22222222-2222-2222-2222-222222222222', 'Alimentación', 'Gastos de comida y bebida', 'utensils', '#10B981', NULL, 40000, 400000, true, true);

-- Subcategorías de Alimentación
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Desayuno', 'Desayuno de trabajo', 'coffee', '#34D399', '22222222-2222-2222-2222-222222222222', 10000, 100000, true, true),
('Almuerzo', 'Almuerzo de trabajo', 'lunch', '#34D399', '22222222-2222-2222-2222-222222222222', 20000, 250000, true, true),
('Cena', 'Cena de trabajo', 'dinner', '#34D399', '22222222-2222-2222-2222-222222222222', 25000, 200000, true, true);

-- 3. ALOJAMIENTO
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('33333333-3333-3333-3333-333333333333', 'Alojamiento', 'Gastos de hospedaje', 'hotel', '#F59E0B', NULL, 150000, 1000000, true, true);

-- Subcategorías de Alojamiento
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Hotel', 'Hospedaje en hotel', 'building', '#FBBF24', '33333333-3333-3333-3333-333333333333', 150000, 800000, true, true),
('Airbnb', 'Alojamiento Airbnb', 'home', '#FBBF24', '33333333-3333-3333-3333-333333333333', 100000, 600000, true, true);

-- 4. VIAJE
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('44444444-4444-4444-4444-444444444444', 'Viaje', 'Gastos de pasajes y viajes', 'plane', '#8B5CF6', NULL, 300000, 2000000, true, true);

-- Subcategorías de Viaje
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Vuelos', 'Pasajes aéreos', 'airplane', '#A78BFA', '44444444-4444-4444-4444-444444444444', 500000, 3000000, true, true),
('Buses', 'Pasajes de bus', 'bus', '#A78BFA', '44444444-4444-4444-4444-444444444444', 30000, 200000, true, true),
('Trenes', 'Pasajes de tren', 'train', '#A78BFA', '44444444-4444-4444-4444-444444444444', 50000, 300000, true, true);

-- 5. COMUNICACIONES
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('55555555-5555-5555-5555-555555555555', 'Comunicaciones', 'Gastos de telecomunicaciones', 'phone', '#EC4899', NULL, 30000, 150000, true, true);

-- Subcategorías de Comunicaciones
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Internet', 'Servicio de internet', 'wifi', '#F472B6', '55555555-5555-5555-5555-555555555555', 20000, 100000, true, true),
('Teléfono', 'Llamadas telefónicas', 'device', '#F472B6', '55555555-5555-5555-5555-555555555555', 15000, 80000, true, true);

-- 6. OFICINA
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('66666666-6666-6666-6666-666666666666', 'Oficina', 'Gastos de oficina y materiales', 'printer', '#6366F1', NULL, 50000, 300000, true, true);

-- Subcategorías de Oficina
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Materiales', 'Material de oficina', 'supplies', '#818CF8', '66666666-6666-6666-6666-666666666666', 30000, 200000, true, true),
('Impresiones', 'Impresiones y fotocopias', 'print', '#818CF8', '66666666-6666-6666-6666-666666666666', 20000, 100000, true, true);

-- 7. CAPACITACIÓN
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('77777777-7777-7777-7777-777777777777', 'Capacitación', 'Gastos de formación y capacitación', 'book', '#14B8A6', NULL, 100000, 500000, true, true);

-- Subcategorías de Capacitación
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Cursos', 'Cursos y talleres', 'course', '#2DD4BF', '77777777-7777-7777-7777-777777777777', 150000, 800000, true, true),
('Conferencias', 'Conferencias y seminarios', 'presentation', '#2DD4BF', '77777777-7777-7777-7777-777777777777', 200000, 1000000, true, true);

-- 8. CLIENTE
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('88888888-8888-8888-8888-888888888888', 'Cliente', 'Gastos relacionados con clientes', 'gift', '#EF4444', NULL, 80000, 500000, true, true);

-- Subcategorías de Cliente
INSERT INTO expense_categories (name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('Regalos', 'Regalos empresariales', 'present', '#F87171', '88888888-8888-8888-8888-888888888888', 50000, 300000, true, true),
('Atenciones', 'Atenciones a clientes', 'handshake', '#F87171', '88888888-8888-8888-8888-888888888888', 60000, 400000, true, true);

-- 9. OTROS
INSERT INTO expense_categories (id, name, description, icon, color, parent_id, daily_limit, monthly_limit, requires_receipt, is_active) VALUES
('99999999-9999-9999-9999-999999999999', 'Otros', 'Otros gastos no categorizados', 'briefcase', '#64748B', NULL, 100000, 500000, true, true);

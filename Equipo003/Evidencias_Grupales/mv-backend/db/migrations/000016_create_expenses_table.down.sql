-- Revertir creación de tabla expenses y tipos ENUM
DROP TABLE IF EXISTS expenses CASCADE;
DROP TYPE IF EXISTS expense_status CASCADE;
DROP TYPE IF EXISTS payment_method CASCADE;

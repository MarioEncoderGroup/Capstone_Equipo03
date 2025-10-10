-- Tabla de categorías de gastos
CREATE TABLE expense_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(20),
    parent_id UUID REFERENCES expense_categories(id) ON DELETE SET NULL,
    daily_limit DECIMAL(12,2),
    monthly_limit DECIMAL(12,2),
    requires_receipt BOOLEAN DEFAULT true,
    is_active BOOLEAN DEFAULT true,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT check_limits CHECK (
        (daily_limit IS NULL OR daily_limit > 0) AND
        (monthly_limit IS NULL OR monthly_limit > 0)
    )
);

-- Índices para optimizar búsquedas
CREATE INDEX idx_expense_categories_parent ON expense_categories(parent_id);
CREATE INDEX idx_expense_categories_active ON expense_categories(is_active) WHERE is_active = true;
CREATE INDEX idx_expense_categories_deleted ON expense_categories(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_expense_categories_name ON expense_categories(name);

-- Comentarios
COMMENT ON TABLE expense_categories IS 'Categorías de gastos empresariales';
COMMENT ON COLUMN expense_categories.parent_id IS 'Referencia a categoría padre para jerarquías (ej: Transporte > Taxi)';
COMMENT ON COLUMN expense_categories.daily_limit IS 'Límite diario permitido para esta categoría';
COMMENT ON COLUMN expense_categories.monthly_limit IS 'Límite mensual permitido para esta categoría';
COMMENT ON COLUMN expense_categories.requires_receipt IS 'Si requiere comprobante obligatorio';

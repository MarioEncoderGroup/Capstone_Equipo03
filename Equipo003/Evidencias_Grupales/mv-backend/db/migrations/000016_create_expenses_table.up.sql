-- Crear tipos ENUM para expenses
CREATE TYPE expense_status AS ENUM (
    'draft',
    'submitted',
    'approved',
    'rejected',
    'reimbursed'
);

CREATE TYPE payment_method AS ENUM (
    'cash',
    'card',
    'transfer'
);

-- Tabla principal de gastos
CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    policy_id UUID,
    category_id UUID NOT NULL REFERENCES expense_categories(id) ON DELETE RESTRICT,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    amount DECIMAL(12,2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(10) DEFAULT 'CLP',
    exchange_rate DECIMAL(10,4) DEFAULT 1.0,
    amount_clp DECIMAL(12,2) GENERATED ALWAYS AS (amount * exchange_rate) STORED,
    expense_date DATE NOT NULL,
    merchant_name VARCHAR(200),
    merchant_rut VARCHAR(20),
    receipt_number VARCHAR(100),
    payment_method payment_method NOT NULL,
    status expense_status DEFAULT 'draft',
    is_reimbursable BOOLEAN DEFAULT true,
    violation_reason TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT check_expense_date CHECK (expense_date <= CURRENT_DATE),
    CONSTRAINT check_exchange_rate CHECK (exchange_rate > 0)
);

-- Índices para optimizar búsquedas
CREATE INDEX idx_expenses_user ON expenses(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_policy ON expenses(policy_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_category ON expenses(category_id);
CREATE INDEX idx_expenses_status ON expenses(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_date ON expenses(expense_date);
CREATE INDEX idx_expenses_deleted ON expenses(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_created ON expenses(created);

-- Índice compuesto para queries comunes
CREATE INDEX idx_expenses_user_status ON expenses(user_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_user_date ON expenses(user_id, expense_date DESC) WHERE deleted_at IS NULL;

-- Comentarios
COMMENT ON TABLE expenses IS 'Tabla principal de gastos empresariales';
COMMENT ON COLUMN expenses.user_id IS 'Usuario que creó el gasto';
COMMENT ON COLUMN expenses.policy_id IS 'Política aplicada al gasto (si corresponde)';
COMMENT ON COLUMN expenses.amount_clp IS 'Monto convertido a CLP (calculado automáticamente)';
COMMENT ON COLUMN expenses.expense_date IS 'Fecha en que se realizó el gasto';
COMMENT ON COLUMN expenses.status IS 'Estado del gasto: draft, submitted, approved, rejected, reimbursed';
COMMENT ON COLUMN expenses.violation_reason IS 'Razón si el gasto viola alguna política';

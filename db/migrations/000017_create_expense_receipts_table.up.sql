-- Tabla de comprobantes de gastos
CREATE TABLE expense_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE,
    file_url VARCHAR(500) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_size BIGINT NOT NULL,
    ocr_data JSONB,
    ocr_confidence DECIMAL(5,2),
    is_primary BOOLEAN DEFAULT false,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT check_file_size CHECK (file_size > 0 AND file_size <= 10485760), -- Max 10MB
    CONSTRAINT check_file_type CHECK (
        file_type IN ('image/jpeg', 'image/png', 'image/jpg', 'application/pdf')
    ),
    CONSTRAINT check_confidence CHECK (
        ocr_confidence IS NULL OR (ocr_confidence >= 0 AND ocr_confidence <= 100)
    )
);

-- Índices
CREATE INDEX idx_receipts_expense ON expense_receipts(expense_id);
CREATE INDEX idx_receipts_primary ON expense_receipts(is_primary) WHERE is_primary = true;
CREATE INDEX idx_receipts_created ON expense_receipts(created);

-- Solo un comprobante principal por gasto
CREATE UNIQUE INDEX idx_receipts_unique_primary
    ON expense_receipts(expense_id)
    WHERE is_primary = true;

-- Comentarios
COMMENT ON TABLE expense_receipts IS 'Comprobantes y recibos de gastos';
COMMENT ON COLUMN expense_receipts.file_url IS 'URL del archivo en S3 o storage';
COMMENT ON COLUMN expense_receipts.ocr_data IS 'Datos extraídos por OCR en formato JSON';
COMMENT ON COLUMN expense_receipts.ocr_confidence IS 'Nivel de confianza del OCR (0-100)';
COMMENT ON COLUMN expense_receipts.is_primary IS 'Si es el comprobante principal (solo uno por gasto)';

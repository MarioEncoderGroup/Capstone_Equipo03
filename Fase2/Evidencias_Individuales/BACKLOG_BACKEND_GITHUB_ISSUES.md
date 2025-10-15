# 📋 BACKLOG - GITHUB ISSUES MV-BACKEND

**Proyecto:** MisViáticos Backend
**Repositorio:** [Tu repo de backend]
**Base:** Master Plan Viáticos
**Stack:** Go + Fiber + PostgreSQL + Redis

---

## 🏷️ LABELS A CREAR EN GITHUB

```
- phase-1-foundations (color: #0E8A16)
- phase-2-ocr (color: #1D76DB)
- phase-3-policies (color: #5319E7)
- phase-4-approvals (color: #D93F0B)
- phase-5-notifications (color: #FBCA04)
- phase-6-reports (color: #006B75)
- phase-7-optimization (color: #C5DEF5)
- phase-8-testing (color: #BFD4F2)
- priority-critical (color: #B60205)
- priority-high (color: #D93F0B)
- priority-medium (color: #FBCA04)
- priority-low (color: #0E8A16)
- type-feature (color: #84B6EB)
- type-bug (color: #D73A4A)
- type-refactor (color: #FEF2C0)
- type-docs (color: #C5DEF5)
- type-db (color: #D4C5F9)
- type-chore (color: #EDEDED)
```

---

## 🚀 FASE 1: FUNDAMENTOS DE VIÁTICOS (SPRINT 1-4)

### Issue #1: Crear migración de tabla expense_categories [COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-db`

**Descripción:**
Crear tabla para categorías de gastos.

**Archivo:** `db/migrations/000015_create_expense_categories_table.up.sql`

**Schema:**
```sql
CREATE TABLE expense_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(20),
    parent_id UUID REFERENCES expense_categories(id),
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

CREATE INDEX idx_expense_categories_parent ON expense_categories(parent_id);
CREATE INDEX idx_expense_categories_active ON expense_categories(is_active);
CREATE INDEX idx_expense_categories_deleted ON expense_categories(deleted_at);
```

**Tareas:**
- [ ] Crear archivo de migración .up.sql
- [ ] Crear archivo de migración .down.sql
- [ ] Ejecutar migración en dev
- [ ] Verificar constraints
- [ ] Verificar índices

**Criterios de Aceptación:**
- ✅ Tabla creada correctamente
- ✅ Migración reversible (down funciona)
- ✅ Índices creados
- ✅ Constraints funcionan

---

### Issue #2: Crear seed de categorías predefinidas[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-high`, `type-db`

**Descripción:**
Seed con categorías comunes de gastos empresariales.

**Archivo:** `db/seed/00010_expense_categories_seed.sql`

**Categorías:**
- 🚗 Transporte (Taxi, Uber, Combustible, Peajes, Estacionamiento)
- 🍽️ Alimentación (Desayuno, Almuerzo, Cena)
- 🏨 Alojamiento (Hotel, Airbnb)
- ✈️ Viaje (Vuelos, Buses, Trenes)
- 📱 Comunicaciones (Internet, Teléfono)
- 🖨️ Oficina (Materiales, Impresiones)
- 📚 Capacitación (Cursos, Conferencias)
- 🎁 Cliente (Regalos, Atenciones)
- 💼 Otros

**Tareas:**
- [ ] Definir categorías y subcategorías
- [ ] Asignar íconos (emojis)
- [ ] Asignar colores
- [ ] Definir límites razonables
- [ ] Crear seed SQL

**Criterios de Aceptación:**
- ✅ Al menos 20 categorías creadas
- ✅ Íconos y colores asignados
- ✅ Límites coherentes

---

### Issue #3: Crear migración de tabla expenses [COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-db`

**Descripción:**
Tabla principal de gastos.

**Archivo:** `db/migrations/000016_create_expenses_table.up.sql`

**Schema:**
```sql
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

CREATE TABLE expenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    policy_id UUID,
    category_id UUID NOT NULL REFERENCES expense_categories(id),
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

    CONSTRAINT check_expense_date CHECK (expense_date <= CURRENT_DATE)
);

-- Índices
CREATE INDEX idx_expenses_user ON expenses(user_id);
CREATE INDEX idx_expenses_policy ON expenses(policy_id);
CREATE INDEX idx_expenses_category ON expenses(category_id);
CREATE INDEX idx_expenses_status ON expenses(status);
CREATE INDEX idx_expenses_date ON expenses(expense_date);
CREATE INDEX idx_expenses_deleted ON expenses(deleted_at);
CREATE INDEX idx_expenses_created ON expenses(created);

-- Índice compuesto para queries comunes
CREATE INDEX idx_expenses_user_status ON expenses(user_id, status) WHERE deleted_at IS NULL;
```

**Criterios de Aceptación:**
- ✅ Tabla creada con todos los campos
- ✅ ENUMs definidos
- ✅ Columna calculada (amount_clp) funciona
- ✅ Índices optimizados
- ✅ Constraints funcionan

---

### Issue #4: Crear migración de tabla expense_receipts[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-db`

**Descripción:**
Tabla para comprobantes de gastos.

**Archivo:** `db/migrations/000017_create_expense_receipts_table.up.sql`

**Schema:**
```sql
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

CREATE INDEX idx_receipts_expense ON expense_receipts(expense_id);
CREATE INDEX idx_receipts_primary ON expense_receipts(is_primary) WHERE is_primary = true;

-- Solo un comprobante principal por gasto
CREATE UNIQUE INDEX idx_receipts_unique_primary
    ON expense_receipts(expense_id)
    WHERE is_primary = true;
```

**Criterios de Aceptación:**
- ✅ Tabla creada
- ✅ Cascade delete funciona
- ✅ Constraint de tamaño funciona
- ✅ Solo un comprobante principal

---

### Issue #5: Implementar dominio de Expense[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-feature`

**Descripción:**
Crear entidad y DTOs para expenses.

**Archivo:** `internal/core/expense/domain/expense.go`

**Estructuras:**
```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type ExpenseStatus string

const (
    ExpenseStatusDraft      ExpenseStatus = "draft"
    ExpenseStatusSubmitted  ExpenseStatus = "submitted"
    ExpenseStatusApproved   ExpenseStatus = "approved"
    ExpenseStatusRejected   ExpenseStatus = "rejected"
    ExpenseStatusReimbursed ExpenseStatus = "reimbursed"
)

type PaymentMethod string

const (
    PaymentMethodCash     PaymentMethod = "cash"
    PaymentMethodCard     PaymentMethod = "card"
    PaymentMethodTransfer PaymentMethod = "transfer"
)

type Expense struct {
    ID              uuid.UUID      `json:"id"`
    UserID          uuid.UUID      `json:"user_id"`
    PolicyID        *uuid.UUID     `json:"policy_id,omitempty"`
    CategoryID      uuid.UUID      `json:"category_id"`
    Title           string         `json:"title"`
    Description     *string        `json:"description,omitempty"`
    Amount          float64        `json:"amount"`
    Currency        string         `json:"currency"`
    ExchangeRate    float64        `json:"exchange_rate"`
    AmountCLP       float64        `json:"amount_clp"`
    ExpenseDate     time.Time      `json:"expense_date"`
    MerchantName    *string        `json:"merchant_name,omitempty"`
    MerchantRUT     *string        `json:"merchant_rut,omitempty"`
    ReceiptNumber   *string        `json:"receipt_number,omitempty"`
    PaymentMethod   PaymentMethod  `json:"payment_method"`
    Status          ExpenseStatus  `json:"status"`
    IsReimbursable  bool           `json:"is_reimbursable"`
    ViolationReason *string        `json:"violation_reason,omitempty"`
    Receipts        []ExpenseReceipt `json:"receipts,omitempty"`
    Category        *ExpenseCategory `json:"category,omitempty"`
    Created         time.Time      `json:"created_at"`
    Updated         time.Time      `json:"updated_at"`
    DeletedAt       *time.Time     `json:"deleted_at,omitempty"`
}

type ExpenseReceipt struct {
    ID             uuid.UUID       `json:"id"`
    ExpenseID      uuid.UUID       `json:"expense_id"`
    FileURL        string          `json:"file_url"`
    FileName       string          `json:"file_name"`
    FileType       string          `json:"file_type"`
    FileSize       int64           `json:"file_size"`
    OCRData        *map[string]any `json:"ocr_data,omitempty"`
    OCRConfidence  *float64        `json:"ocr_confidence,omitempty"`
    IsPrimary      bool            `json:"is_primary"`
    Created        time.Time       `json:"created_at"`
}

type ExpenseCategory struct {
    ID             uuid.UUID  `json:"id"`
    Name           string     `json:"name"`
    Description    *string    `json:"description,omitempty"`
    Icon           *string    `json:"icon,omitempty"`
    Color          *string    `json:"color,omitempty"`
    ParentID       *uuid.UUID `json:"parent_id,omitempty"`
    DailyLimit     *float64   `json:"daily_limit,omitempty"`
    MonthlyLimit   *float64   `json:"monthly_limit,omitempty"`
    RequiresReceipt bool      `json:"requires_receipt"`
    IsActive       bool       `json:"is_active"`
}
```

**Tareas:**
- [ ] Crear archivo domain/expense.go
- [ ] Definir constantes de status
- [ ] Definir constantes de payment method
- [ ] Definir struct Expense
- [ ] Definir struct ExpenseReceipt
- [ ] Definir struct ExpenseCategory
- [ ] Agregar tags JSON
- [ ] Documentación godoc

**Criterios de Aceptación:**
- ✅ Structs definidos
- ✅ Constantes creadas
- ✅ Tags JSON correctos
- ✅ Sin errores de compilación

---

### Issue #6: Crear DTOs para Expense[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-high`, `type-feature`

**Descripción:**
DTOs de request/response para expenses.

**Archivo:** `internal/core/expense/domain/dtos.go`

**DTOs:**
```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

// CreateExpenseDto - Request para crear gasto
type CreateExpenseDto struct {
    CategoryID     uuid.UUID     `json:"category_id" validate:"required,uuid"`
    Title          string        `json:"title" validate:"required,min=3,max=200"`
    Description    string        `json:"description" validate:"omitempty,max=500"`
    Amount         float64       `json:"amount" validate:"required,gt=0"`
    Currency       string        `json:"currency" validate:"required,oneof=CLP USD EUR"`
    ExpenseDate    time.Time     `json:"expense_date" validate:"required"`
    MerchantName   string        `json:"merchant_name" validate:"omitempty,max=200"`
    MerchantRUT    string        `json:"merchant_rut" validate:"omitempty,rut"`
    ReceiptNumber  string        `json:"receipt_number" validate:"omitempty,max=100"`
    PaymentMethod  PaymentMethod `json:"payment_method" validate:"required,oneof=cash card transfer"`
    IsReimbursable bool          `json:"is_reimbursable"`
}

// UpdateExpenseDto - Request para actualizar gasto
type UpdateExpenseDto struct {
    CategoryID     *uuid.UUID     `json:"category_id" validate:"omitempty,uuid"`
    Title          *string        `json:"title" validate:"omitempty,min=3,max=200"`
    Description    *string        `json:"description" validate:"omitempty,max=500"`
    Amount         *float64       `json:"amount" validate:"omitempty,gt=0"`
    Currency       *string        `json:"currency" validate:"omitempty,oneof=CLP USD EUR"`
    ExpenseDate    *time.Time     `json:"expense_date"`
    MerchantName   *string        `json:"merchant_name" validate:"omitempty,max=200"`
    MerchantRUT    *string        `json:"merchant_rut" validate:"omitempty,rut"`
    ReceiptNumber  *string        `json:"receipt_number" validate:"omitempty,max=100"`
    PaymentMethod  *PaymentMethod `json:"payment_method" validate:"omitempty,oneof=cash card transfer"`
    IsReimbursable *bool          `json:"is_reimbursable"`
}

// ExpenseFilters - Filtros para listar gastos
type ExpenseFilters struct {
    UserID     *uuid.UUID
    CategoryID *uuid.UUID
    Status     *ExpenseStatus
    DateFrom   *time.Time
    DateTo     *time.Time
    MinAmount  *float64
    MaxAmount  *float64
    Search     *string
    Limit      int
    Offset     int
}

// UploadReceiptDto - Request para subir comprobante
type UploadReceiptDto struct {
    ExpenseID  uuid.UUID `json:"expense_id" validate:"required,uuid"`
    FileName   string    `json:"file_name" validate:"required"`
    FileType   string    `json:"file_type" validate:"required,oneof=image/jpeg image/png image/jpg application/pdf"`
    FileSize   int64     `json:"file_size" validate:"required,gt=0,lte=10485760"` // Max 10MB
    IsPrimary  bool      `json:"is_primary"`
}

// ExpenseResponse - Response con datos completos
type ExpenseResponse struct {
    Expense
    Receipts []ExpenseReceipt `json:"receipts,omitempty"`
    Category *ExpenseCategory `json:"category,omitempty"`
}

// ExpensesResponse - Response para lista
type ExpensesResponse struct {
    Expenses []ExpenseResponse `json:"expenses"`
    Total    int               `json:"total"`
    Limit    int               `json:"limit"`
    Offset   int               `json:"offset"`
}
```

**Tareas:**
- [ ] Crear archivo dtos.go
- [ ] Definir CreateExpenseDto
- [ ] Definir UpdateExpenseDto
- [ ] Definir ExpenseFilters
- [ ] Definir UploadReceiptDto
- [ ] Definir responses
- [ ] Agregar validaciones

**Criterios de Aceptación:**
- ✅ DTOs definidos
- ✅ Validaciones completas
- ✅ Response types definidos

---

### Issue #7: Implementar ExpenseRepository[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-feature`

**Descripción:**
Repository para acceso a datos de expenses.

**Archivos:**
- `internal/core/expense/ports/repository.go` (interface)
- `internal/core/expense/adapters/postgresql.go` (implementación)

**Interface:**
```go
package ports

type ExpenseRepository interface {
    Create(ctx context.Context, expense *domain.Expense) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)
    GetAll(ctx context.Context, filters domain.ExpenseFilters) ([]domain.Expense, int, error)
    Update(ctx context.Context, expense *domain.Expense) error
    Delete(ctx context.Context, id uuid.UUID) error

    // Receipts
    AddReceipt(ctx context.Context, receipt *domain.ExpenseReceipt) error
    GetReceipts(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseReceipt, error)
    DeleteReceipt(ctx context.Context, receiptID uuid.UUID) error
    SetPrimaryReceipt(ctx context.Context, receiptID uuid.UUID) error
}
```

**Tareas:**
- [ ] Crear interface ExpenseRepository
- [ ] Implementar PostgresExpenseRepository
- [ ] Método Create con transaction
- [ ] Método GetByID con JOIN de receipts
- [ ] Método GetAll con filtros dinámicos
- [ ] Método Update
- [ ] Método Delete (soft delete)
- [ ] Métodos de receipts
- [ ] Tests unitarios con mocks

**Criterios de Aceptación:**
- ✅ Interface definida
- ✅ Implementación PostgreSQL completa
- ✅ Queries optimizadas
- ✅ Transactions donde corresponde
- ✅ Tests con 80%+ coverage

---

### Issue #8: Implementar CategoryRepository[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-high`, `type-feature`

**Descripción:**
Repository para categorías.

**Archivos:**
- `internal/core/expense/ports/category_repository.go`
- `internal/core/expense/adapters/postgresql_category.go`

**Interface:**
```go
type CategoryRepository interface {
    Create(ctx context.Context, category *domain.ExpenseCategory) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseCategory, error)
    GetAll(ctx context.Context, activeOnly bool) ([]domain.ExpenseCategory, error)
    GetByParent(ctx context.Context, parentID *uuid.UUID) ([]domain.ExpenseCategory, error)
    Update(ctx context.Context, category *domain.ExpenseCategory) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**Criterios de Aceptación:**
- ✅ CRUD completo
- ✅ Query para subcategorías
- ✅ Soft delete

---

### Issue #9: Implementar ExpenseService[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-feature`

**Descripción:**
Lógica de negocio para expenses.

**Archivos:**
- `internal/core/expense/ports/service.go` (interface)
- `internal/core/expense/services/expense_service.go` (implementación)

**Interface:**
```go
package ports

type ExpenseService interface {
    // CRUD
    Create(ctx context.Context, dto domain.CreateExpenseDto, userID uuid.UUID) (*domain.Expense, error)
    GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Expense, error)
    GetAll(ctx context.Context, filters domain.ExpenseFilters) ([]domain.Expense, int, error)
    Update(ctx context.Context, id uuid.UUID, dto domain.UpdateExpenseDto, userID uuid.UUID) (*domain.Expense, error)
    Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error

    // Receipts
    UploadReceipt(ctx context.Context, dto domain.UploadReceiptDto, fileData []byte) (*domain.ExpenseReceipt, error)
    DeleteReceipt(ctx context.Context, receiptID uuid.UUID, userID uuid.UUID) error

    // Business logic
    CanEdit(ctx context.Context, expenseID uuid.UUID, userID uuid.UUID) (bool, error)
    CanDelete(ctx context.Context, expenseID uuid.UUID, userID uuid.UUID) (bool, error)
    ChangeStatus(ctx context.Context, expenseID uuid.UUID, newStatus domain.ExpenseStatus) error
}
```

**Lógica de Negocio:**
- Solo el owner puede editar/eliminar
- Solo editable si status es "draft" o "rejected"
- Validar fecha de gasto (no futura)
- Sanitizar inputs
- Validar RUT si se proporciona
- Calcular amount_clp con exchange rate

**Tareas:**
- [ ] Crear interface ExpenseService
- [ ] Implementar service
- [ ] Método Create con validaciones
- [ ] Método GetByID con permisos
- [ ] Método GetAll con filtros
- [ ] Método Update con validaciones
- [ ] Método Delete con validaciones
- [ ] Métodos de receipts
- [ ] Lógica CanEdit/CanDelete
- [ ] Tests unitarios

**Criterios de Aceptación:**
- ✅ Service completo
- ✅ Validaciones de negocio
- ✅ Permisos verificados
- ✅ Tests con 80%+ coverage

---

### Issue #10: Crear endpoints REST para Expenses[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-critical`, `type-feature`

**Descripción:**
Controlador HTTP para expenses.

**Archivos:**
- `internal/controllers/expense_controller.go`
- `internal/routes/expense_routes.go`

**Endpoints:**
```go
// Expenses
GET    /api/v1/expenses              - Lista de gastos del usuario
GET    /api/v1/expenses/:id          - Detalle de gasto
POST   /api/v1/expenses              - Crear gasto
PUT    /api/v1/expenses/:id          - Actualizar gasto
DELETE /api/v1/expenses/:id          - Eliminar gasto (soft)

// Receipts
POST   /api/v1/expenses/:id/receipts - Upload comprobante
DELETE /api/v1/receipts/:id          - Eliminar comprobante

// Admin
GET    /api/v1/admin/expenses        - Todos los gastos (admin only)
```

**Middlewares:**
- AuthMiddleware (JWT)
- RequirePermission (según acción)
- DatabaseTenantMiddleware

**Tareas:**
- [ ] Crear ExpenseController
- [ ] Implementar GetAll
- [ ] Implementar GetByID
- [ ] Implementar Create
- [ ] Implementar Update
- [ ] Implementar Delete
- [ ] Implementar UploadReceipt
- [ ] Implementar DeleteReceipt
- [ ] Configurar rutas
- [ ] Agregar middlewares
- [ ] Swagger documentation

**Criterios de Aceptación:**
- ✅ Todos los endpoints funcionan
- ✅ Validaciones en el controller
- ✅ Responses consistentes
- ✅ Error handling robusto

---

### Issue #11: Configurar AWS S3 para comprobantes[COMPLETADO]
**Labels:** `phase-1-foundations`, `priority-high`, `type-feature`

**Descripción:**
Implementar cliente AWS S3 para almacenar comprobantes.

**Archivo:** `internal/libraries/aws/s3_client.go`

**Funciones:**
```go
package aws

type S3Client struct {
    client *s3.Client
    bucket string
}

func NewS3Client(region, accessKey, secretKey, bucket string) (*S3Client, error)

// UploadFile sube un archivo a S3 y retorna la URL
func (c *S3Client) UploadFile(ctx context.Context, key string, data []byte, contentType string) (string, error)

// DeleteFile elimina un archivo de S3
func (c *S3Client) DeleteFile(ctx context.Context, key string) error

// GetPresignedURL genera URL firmada para acceso temporal
func (c *S3Client) GetPresignedURL(ctx context.Context, key string, duration time.Duration) (string, error)
```

**Configuración (.env):**
```bash
AWS_REGION=us-west-2
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_S3_BUCKET_RECEIPTS=misviaticos-receipts
```

**Tareas:**
- [ ] Instalar SDK AWS Go v2
- [ ] Crear S3Client
- [ ] Método UploadFile
- [ ] Método DeleteFile
- [ ] Método GetPresignedURL
- [ ] Configurar CORS del bucket
- [ ] Tests de integración

**Criterios de Aceptación:**
- ✅ Upload funciona
- ✅ URLs generadas son accesibles
- ✅ Delete funciona
- ✅ Presigned URLs funcionan

---

## 🤖 FASE 2: OCR INTEGRATION (SPRINT 5-7)

### Issue #12: Configurar Google Cloud Vision API[COMPLETADO]
**Labels:** `phase-2-ocr`, `priority-critical`, `type-feature`

**Descripción:**
Integrar Google Vision API para OCR.

**Archivo:** `internal/libraries/ocr/google_vision.go`

**Configuración:**
1. Crear proyecto en Google Cloud
2. Habilitar Vision API
3. Crear Service Account
4. Descargar JSON de credenciales

**Variables de entorno:**
```bash
GOOGLE_CLOUD_PROJECT_ID=misviaticos-prod
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
```

**Tareas:**
- [ ] Instalar SDK Google Cloud Vision
- [ ] Configurar credenciales
- [ ] Crear GoogleVisionClient
- [ ] Método DetectText
- [ ] Método DetectDocument
- [ ] Manejo de errores
- [ ] Rate limiting
- [ ] Tests con imágenes de muestra

**Criterios de Aceptación:**
- ✅ API configurada
- ✅ Extrae texto correctamente
- ✅ Manejo de errores robusto

---

### Issue #13: Implementar ReceiptParser para recibos chilenos[COMPLETADO]
**Labels:** `phase-2-ocr`, `priority-critical`, `type-feature`

**Descripción:**
Parser que extrae datos estructurados de texto OCR.

**Archivo:** `internal/core/ocr/services/receipt_parser.go`

**Funciones:**
```go
package services

type OCRResult struct {
    Amount        float64            `json:"amount"`
    Date          *time.Time         `json:"date,omitempty"`
    MerchantRUT   string             `json:"merchant_rut,omitempty"`
    MerchantName  string             `json:"merchant_name,omitempty"`
    DocumentType  string             `json:"document_type"` // boleta, factura
    Confidence    float64            `json:"confidence"`
    RawText       string             `json:"raw_text"`
    ExtractedData map[string]float64 `json:"extracted_data"` // campo -> confidence
}

type ReceiptParser struct {}

// ParseChileanReceipt extrae datos de recibos chilenos
func (p *ReceiptParser) ParseChileanReceipt(text string) (*OCRResult, error)

// ExtractAmount busca montos en el texto
func (p *ReceiptParser) ExtractAmount(text string) (float64, float64, error) // monto, confidence, error

// ExtractDate busca fechas
func (p *ReceiptParser) ExtractDate(text string) (*time.Time, float64, error)

// ExtractRUT busca RUT chileno
func (p *ReceiptParser) ExtractRUT(text string) (string, float64, error)

// ValidateRUT valida dígito verificador
func (p *ReceiptParser) ValidateRUT(rut string) bool
```

**Regexes:**
- Monto: `\$\s*(\d{1,3}(?:\.\d{3})*|\d+)`
- RUT: `(\d{1,2}\.\d{3}\.\d{3}-[\dkK])`
- Fecha: varios formatos chilenos

**Tareas:**
- [ ] Crear ReceiptParser
- [ ] Regex para montos
- [ ] Regex para fechas
- [ ] Regex para RUT
- [ ] Validación de RUT
- [ ] Cálculo de confidence
- [ ] Tests con casos reales

**Criterios de Aceptación:**
- ✅ Extrae monto con 80%+ accuracy
- ✅ Extrae RUT con 70%+ accuracy
- ✅ Valida RUT correctamente
- ✅ Confidence scores precisos

---

### Issue #14: Implementar OCRService[COMPLETADO]
**Labels:** `phase-2-ocr`, `priority-critical`, `type-feature`

**Descripción:**
Servicio que coordina OCR y parsing.

**Archivo:** `internal/core/ocr/services/ocr_service.go`

**Interface:**
```go
package services

type OCRService interface {
    // ProcessReceipt procesa imagen y retorna datos extraídos
    ProcessReceipt(ctx context.Context, imageData []byte) (*OCRResult, error)

    // ProcessReceiptFromURL procesa desde URL (S3)
    ProcessReceiptFromURL(ctx context.Context, imageURL string) (*OCRResult, error)
}
```

**Implementación:**
```go
type ocrService struct {
    visionClient *GoogleVisionClient
    parser       *ReceiptParser
}

func (s *ocrService) ProcessReceipt(ctx context.Context, imageData []byte) (*OCRResult, error) {
    // 1. Llamar Google Vision API
    text, err := s.visionClient.DetectText(ctx, imageData)

    // 2. Parsear texto extraído
    result, err := s.parser.ParseChileanReceipt(text)

    // 3. Retornar resultado estructurado
    return result, nil
}
```

**Tareas:**
- [ ] Crear OCRService interface
- [ ] Implementar service
- [ ] Integrar Vision API
- [ ] Integrar Parser
- [ ] Preprocesamiento de imagen (opcional)
- [ ] Cache de resultados (Redis)
- [ ] Logging detallado
- [ ] Tests

**Criterios de Aceptación:**
- ✅ Service funcional
- ✅ Integración completa
- ✅ Cache funciona
- ✅ Tests con imágenes reales

---

### Issue #15: Crear endpoint de OCR[COMPLETADO]
**Labels:** `phase-2-ocr`, `priority-high`, `type-feature`

**Descripción:**
Endpoint para analizar comprobantes con OCR.

**Archivo:** `internal/controllers/ocr_controller.go`

**Endpoint:**
```
POST /api/v1/ocr/analyze
Content-Type: multipart/form-data

Body:
- image: file (jpg, png, pdf)

Response:
{
  "success": true,
  "data": {
    "amount": 15500,
    "date": "2025-10-05",
    "merchant_rut": "76.123.456-7",
    "merchant_name": "Supermercado ABC",
    "document_type": "boleta",
    "confidence": 0.92,
    "raw_text": "..."
  }
}
```

**Tareas:**
- [ ] Crear OCRController
- [ ] Endpoint POST /ocr/analyze
- [ ] Validación de archivo (tipo, tamaño)
- [ ] Llamar OCRService
- [ ] Response estructurado
- [ ] Rate limiting (5 req/min)
- [ ] Logging

**Criterios de Aceptación:**
- ✅ Endpoint funcional
- ✅ Validaciones robustas
- ✅ Response correcto
- ✅ Rate limiting funciona

---

## 📋 FASE 3: POLÍTICAS Y VALIDACIÓN (SPRINT 8-10)

### Issue #16: Crear migraciones de tablas de políticas[COMPLETADO]
**Labels:** `phase-3-policies`, `priority-critical`, `type-db`

**Descripción:**
Crear tablas para políticas de gastos.

**Archivos:**
- `000018_create_policies_table.up.sql`
- `000019_create_policy_rules_table.up.sql`
- `000020_create_policy_approvers_table.up.sql`
- `000021_create_policy_submitters_table.up.sql`

**Schemas en el Master Plan**

**Criterios de Aceptación:**
- ✅ 4 tablas creadas
- ✅ Relaciones correctas
- ✅ Índices optimizados

---

### Issue #17: Implementar PolicyService[COMPLETADO]
**Labels:** `phase-3-policies`, `priority-critical`, `type-feature`

**Descripción:**
Servicio para gestionar políticas.

**Interface:**
```go
type PolicyService interface {
    Create(ctx context.Context, dto CreatePolicyDto) (*Policy, error)
    GetAll(ctx context.Context, tenantID uuid.UUID) ([]Policy, error)
    GetByID(ctx context.Context, id uuid.UUID) (*Policy, error)
    Update(ctx context.Context, id uuid.UUID, dto UpdatePolicyDto) (*Policy, error)
    Delete(ctx context.Context, id uuid.UUID) error

    // Rules
    AddRule(ctx context.Context, dto CreatePolicyRuleDto) (*PolicyRule, error)
    UpdateRule(ctx context.Context, id uuid.UUID, dto UpdatePolicyRuleDto) (*PolicyRule, error)
    DeleteRule(ctx context.Context, id uuid.UUID) error
}
```

**Criterios de Aceptación:**
- ✅ CRUD completo
- ✅ Gestión de reglas
- ✅ Tests

---

### Issue #18: Implementar RuleEngine para validación[COMPLETADO]
**Labels:** `phase-3-policies`, `priority-critical`, `type-feature`

**Descripción:**
Motor de reglas para validar gastos contra políticas.

**Archivo:** `internal/core/policy/services/rule_engine.go`

**Interface:**
```go
type RuleEngine interface {
    // ValidateExpense valida un gasto contra una política
    ValidateExpense(ctx context.Context, expense *Expense, policy *Policy) ([]Violation, error)

    // CheckApprovalRequired determina si requiere aprobación
    CheckApprovalRequired(ctx context.Context, expense *Expense, policy *Policy) (bool, int, error) // required, level, error

    // GetApprovers retorna aprobadores necesarios según monto
    GetApprovers(ctx context.Context, expense *Expense, policy *Policy) ([]User, error)
}

type Violation struct {
    Field    string `json:"field"`
    Message  string `json:"message"`
    Severity string `json:"severity"` // error, warning
}
```

**Lógica:**
1. Evaluar condiciones de reglas
2. Ejecutar acciones según reglas
3. Calcular nivel de aprobación requerido
4. Detectar violaciones

**Criterios de Aceptación:**
- ✅ Evalúa reglas correctamente
- ✅ Detecta violaciones
- ✅ Determina aprobadores
- ✅ Tests exhaustivos

---

### Issue #19: Crear endpoint de validación de gastos[COMPLETADO]
**Labels:** `phase-3-policies`, `priority-high`, `type-feature`

**Descripción:**
Endpoint para validar gasto antes de crear.

**Endpoint:**
```
POST /api/v1/expenses/validate

Body:
{
  "category_id": "uuid",
  "amount": 75000,
  "expense_date": "2025-10-05"
}

Response:
{
  "success": true,
  "data": {
    "is_valid": false,
    "violations": [
      {
        "field": "amount",
        "message": "Excede límite diario de $50.000 para Transporte",
        "severity": "warning"
      }
    ],
    "requires_approval": true,
    "approval_level": 2,
    "approvers": [
      { "id": "uuid", "name": "Gerente Área" }
    ]
  }
}
```

**Criterios de Aceptación:**
- ✅ Validación en tiempo real
- ✅ Response detallado
- ✅ Performance < 200ms

---

## ✅ FASE 4: FLUJOS DE APROBACIÓN (SPRINT 11-14)

### Issue #20: Crear migraciones de tablas de aprobaciones[COMPLETADO]
**Labels:** `phase-4-approvals`, `priority-critical`, `type-db`

**Descripción:**
Tablas para reportes y aprobaciones.

**Archivos:**
- `000022_create_expense_reports_table.up.sql`
- `000023_create_expense_report_items_table.up.sql`
- `000024_create_approvals_table.up.sql`
- `000025_create_approval_history_table.up.sql`
- `000026_create_expense_comments_table.up.sql`

**Ver Master Plan para schemas completos**

**Criterios de Aceptación:**
- ✅ 5 tablas creadas
- ✅ Relaciones correctas
- ✅ Cascade deletes configurados

---

### Issue #21: Implementar ReportService[COMPLETADO]
**Labels:** `phase-4-approvals`, `priority-critical`, `type-feature`

**Descripción:**
Servicio para gestionar reportes de gastos.

**Interface:**
```go
type ReportService interface {
    Create(ctx context.Context, dto CreateReportDto, userID uuid.UUID) (*ExpenseReport, error)
    GetByID(ctx context.Context, id uuid.UUID) (*ExpenseReport, error)
    GetByUser(ctx context.Context, userID uuid.UUID) ([]ExpenseReport, error)
    Update(ctx context.Context, id uuid.UUID, dto UpdateReportDto) (*ExpenseReport, error)
    Delete(ctx context.Context, id uuid.UUID) error

    // Items
    AddExpenses(ctx context.Context, reportID uuid.UUID, expenseIDs []uuid.UUID) error
    RemoveExpense(ctx context.Context, reportID uuid.UUID, expenseID uuid.UUID) error

    // Workflow
    Submit(ctx context.Context, reportID uuid.UUID) (*ExpenseReport, error)
}
```

**Lógica:**
- Calcular total automáticamente
- Validar que gastos no estén en otro reporte
- Cambiar estado de gastos al agregar/quitar
- Crear aprobaciones al submit

**Criterios de Aceptación:**
- ✅ CRUD completo
- ✅ Gestión de items
- ✅ Submit crea aprobaciones
- ✅ Tests

---

### Issue #22: Implementar ApprovalService[COMPLETADO]
**Labels:** `phase-4-approvals`, `priority-critical`, `type-feature`

**Descripción:**
Servicio para gestionar aprobaciones.

**Interface:**
```go
type ApprovalService interface {
    // GetPendingApprovals retorna aprobaciones pendientes del usuario
    GetPendingApprovals(ctx context.Context, approverID uuid.UUID) ([]Approval, error)

    // Approve aprueba una solicitud
    Approve(ctx context.Context, approvalID uuid.UUID, approverID uuid.UUID, comments string) error

    // Reject rechaza una solicitud
    Reject(ctx context.Context, approvalID uuid.UUID, approverID uuid.UUID, reason string) error

    // GetHistory retorna historial de aprobaciones
    GetHistory(ctx context.Context, reportID uuid.UUID) ([]ApprovalHistory, error)
}
```

**Lógica:**
- Verificar que el usuario sea el aprobador asignado
- Actualizar estado de aprobación
- Crear registro en historial
- Si es aprobación multi-nivel, crear siguiente aprobación
- Si todos aprueban, cambiar estado de reporte a "approved"
- Enviar notificaciones

**Criterios de Aceptación:**
- ✅ Aprobar funciona
- ✅ Rechazar funciona
- ✅ Multi-nivel funciona
- ✅ Historial completo

---

### Issue #23: Implementar WorkflowEngine[COMPLETADO]
**Labels:** `phase-4-approvals`, `priority-critical`, `type-feature`

**Descripción:**
Motor de workflows para aprobaciones.

**Archivo:** `internal/core/approval/services/workflow_engine.go`

**Interface:**
```go
type WorkflowEngine interface {
    // CreateApprovals crea aprobaciones según política y monto
    CreateApprovals(ctx context.Context, report *ExpenseReport) ([]Approval, error)

    // ProcessApproval procesa una aprobación y determina siguiente paso
    ProcessApproval(ctx context.Context, approval *Approval) error

    // EscalateApproval escala al siguiente nivel
    EscalateApproval(ctx context.Context, approvalID uuid.UUID) error
}
```

**Lógica:**
1. Determinar aprobadores según policy y monto
2. Crear aprobaciones por nivel
3. Al aprobar nivel N, crear aprobación nivel N+1
4. Al rechazar, marcar reporte como rechazado
5. Al aprobar último nivel, marcar reporte como aprobado

**Criterios de Aceptación:**
- ✅ Crea aprobaciones correctas
- ✅ Escalamiento funciona
- ✅ Estado del reporte se actualiza

---

### Issue #24: Crear job de escalamiento automático[COMPLETADO]
**Labels:** `phase-4-approvals`, `priority-high`, `type-feature`

**Descripción:**
Cron job que escala aprobaciones pendientes > 24h.

**Archivo:** `internal/jobs/approval_escalation_job.go`

**Lógica:**
```go
func EscalateStaleApprovals() {
    // 1. Buscar aprobaciones pendientes > 24h
    // 2. Por cada aprobación:
    //    - Escalar al siguiente nivel
    //    - Enviar notificación al siguiente aprobador
    //    - Registrar en historial
}
```

**Schedule:** Ejecutar cada hora

**Criterios de Aceptación:**
- ✅ Job corre automáticamente
- ✅ Escala correctamente
- ✅ Notificaciones se envían

---

### Issue #25: Crear endpoints de reportes y aprobaciones
**Labels:** `phase-4-approvals`, `priority-critical`, `type-feature`

**Descripción:**
REST API para reportes y aprobaciones.

**Endpoints:**
```
// Reportes
POST   /api/v1/expense-reports              - Crear reporte
GET    /api/v1/expense-reports               - Mis reportes
GET    /api/v1/expense-reports/:id           - Detalle
PUT    /api/v1/expense-reports/:id           - Actualizar
DELETE /api/v1/expense-reports/:id           - Eliminar
POST   /api/v1/expense-reports/:id/submit    - Enviar a aprobación
POST   /api/v1/expense-reports/:id/expenses  - Agregar gastos
DELETE /api/v1/expense-reports/:id/expenses/:expenseId - Quitar gasto

// Aprobaciones
GET    /api/v1/approvals/pending             - Mis aprobaciones pendientes
GET    /api/v1/approvals/:id                 - Detalle
POST   /api/v1/approvals/:id/approve         - Aprobar
POST   /api/v1/approvals/:id/reject          - Rechazar
GET    /api/v1/approvals/reports/:id/history - Historial
```

**Criterios de Aceptación:**
- ✅ Todos los endpoints funcionan
- ✅ Permisos correctos
- ✅ Validaciones robustas

---

## 🔔 FASE 5: NOTIFICACIONES (SPRINT 15-16)

### Issue #26: Crear tabla de notificaciones[COMPLETADO]
**Labels:** `phase-5-notifications`, `priority-high`, `type-db`

**Archivo:** `000027_create_notifications_table.up.sql`

**Schema en Master Plan**

**Criterios de Aceptación:**
- ✅ Tabla creada
- ✅ Índices optimizados

---

### Issue #27: Implementar NotificationService[COMPLETADO]
**Labels:** `phase-5-notifications`, `priority-critical`, `type-feature`

**Interface:**
```go
type NotificationService interface {
    Create(ctx context.Context, notification *Notification) error
    GetByUser(ctx context.Context, userID uuid.UUID, unreadOnly bool) ([]Notification, error)
    MarkAsRead(ctx context.Context, notificationID uuid.UUID) error
    MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
    Delete(ctx context.Context, notificationID uuid.UUID) error

    // Helpers
    NotifyExpenseApproved(ctx context.Context, expense *Expense) error
    NotifyExpenseRejected(ctx context.Context, expense *Expense, reason string) error
    NotifyApprovalNeeded(ctx context.Context, approval *Approval) error
}
```

**Criterios de Aceptación:**
- ✅ CRUD completo
- ✅ Helpers de notificación
- ✅ Templates de mensajes

---

### Issue #28: Integrar RabbitMQ para cola de notificaciones[COMPLETADO]
**Labels:** `phase-5-notifications`, `priority-high`, `type-feature`

**Descripción:**
Queue para procesar notificaciones async.

**Tareas:**
- [ ] Configurar RabbitMQ
- [ ] Crear producer
- [ ] Crear consumer/worker
- [ ] Queue "notifications"
- [ ] Retry logic
- [ ] Dead letter queue

**Criterios de Aceptación:**
- ✅ Cola funcional
- ✅ Workers procesan mensajes
- ✅ Retry funciona

---

### Issue #29: Implementar WebSocket server[COMPLETADO]
**Labels:** `phase-5-notifications`, `priority-critical`, `type-feature`

**Descripción:**
WebSocket para notificaciones en tiempo real.

**Archivo:** `internal/websocket/server.go`

**Features:**
- Conexión por usuario (autenticada con JWT)
- Broadcast a usuarios específicos
- Heartbeat/ping-pong
- Reconnection handling

**Endpoint:** `ws://localhost:8080/ws`

**Criterios de Aceptación:**
- ✅ WebSocket funcional
- ✅ Autenticación funciona
- ✅ Broadcast funciona

---

### ✅ Issue #30: Configurar Firebase Cloud Messaging [COMPLETADO]
**Labels:** `phase-5-notifications`, `priority-medium`, `type-feature`

**Descripción:**
Push notifications para móviles.

**Tareas:**
- [x] Crear proyecto Firebase
- [x] Configurar FCM
- [x] Integrar SDK Firebase Admin
- [x] Guardar tokens de dispositivos
- [x] Enviar push notifications

**Criterios de Aceptación:**
- ✅ FCM configurado
- ✅ Push notifications funcionan

---

## 📊 FASE 6: REPORTES Y ANALYTICS (SPRINT 17-20)

### ✅ Issue #31: Configurar ClickHouse para analytics [COMPLETADO]
**Labels:** `phase-6-reports`, `priority-high`, `type-db`

**Descripción:**
Base de datos columnar para analytics.

**Tareas:**
- [x] Instalar ClickHouse
- [x] Crear database
- [x] Crear tabla de gastos agregados
- [x] Pipeline de ETL (PostgreSQL → ClickHouse)
- [x] Configurar refreshes

**Criterios de Aceptación:**
- ✅ ClickHouse funcionando
- ✅ Datos sincronizados
- ✅ Queries optimizadas

---

### ✅ Issue #32: Implementar AnalyticsService [COMPLETADO]
**Labels:** `phase-6-reports`, `priority-critical`, `type-feature`

**Interface:**
```go
type AnalyticsService interface {
    GetDashboardMetrics(ctx context.Context, userID uuid.UUID) (*DashboardMetrics, error)
    GetExpensesByCategory(ctx context.Context, filters AnalyticsFilters) ([]CategoryStats, error)
    GetExpensesByMonth(ctx context.Context, filters AnalyticsFilters) ([]MonthlyStats, error)
    GetTopSpenders(ctx context.Context, limit int) ([]UserStats, error)
    GetComplianceReport(ctx context.Context, policyID uuid.UUID) (*ComplianceReport, error)
}
```

**Criterios de Aceptación:**
- ✅ Métricas calculadas correctamente
- ✅ Performance < 500ms
- ✅ Cache con Redis

---

### ✅ Issue #33: Implementar ExportService (Excel/PDF) [COMPLETADO]
**Labels:** `phase-6-reports`, `priority-high`, `type-feature`

**Interface:**
```go
type ExportService interface {
    ExportToExcel(ctx context.Context, data ReportData) ([]byte, error)
    ExportToPDF(ctx context.Context, data ReportData) ([]byte, error)
}
```

**Librerías:**
- Excel: `excelize`
- PDF: `gofpdf`

**Criterios de Aceptación:**
- ✅ Excel generado correctamente
- ✅ PDF generado correctamente
- ✅ Formato profesional

---

### ✅ Issue #34: Crear endpoints de analytics [COMPLETADO]
**Labels:** `phase-6-reports`, `priority-high`, `type-feature`

**Endpoints:**
```
GET /api/v1/analytics/dashboard        - Métricas del dashboard
GET /api/v1/analytics/expenses/category - Por categoría
GET /api/v1/analytics/expenses/monthly  - Por mes
GET /api/v1/analytics/top-spenders      - Top gastadores
GET /api/v1/analytics/compliance/:policyId - Cumplimiento
GET /api/v1/reports/export/excel        - Exportar Excel
GET /api/v1/reports/export/pdf          - Exportar PDF
```

**Criterios de Aceptación:**
- ✅ Endpoints funcionan
- ✅ Performance óptima
- ✅ Export funciona

---

### Issue #35: Implementar ML para predicciones (Opcional)
**Labels:** `phase-6-reports`, `priority-low`, `type-feature`

**Descripción:**
Microservicio Python con TensorFlow para predicciones.

**Features:**
- Predicción de gastos futuros
- Detección de anomalías
- Recomendaciones

**Criterios de Aceptación:**
- ✅ API Python funcional
- ✅ Modelo entrenado
- ✅ Integración con Go

---

## 🔧 FASE 7: OPTIMIZACIONES (SPRINT 21-22)

### Issue #36: Optimizar queries con índices [COMPLETADO]
**Labels:** `phase-7-optimization`, `priority-high`, `type-db`

**Descripción:**
Análisis y optimización de queries lentas.

**Tareas:**
- [ ] Habilitar query logging
- [ ] Identificar queries > 100ms
- [ ] Crear índices adicionales
- [ ] Analizar EXPLAIN plans
- [ ] Refactorizar queries N+1

**Criterios de Aceptación:**
- ✅ Todas las queries < 100ms
- ✅ No hay N+1 queries

---

### Issue #37: Implementar cache con Redis [COMPLETADO]
**Labels:** `phase-7-optimization`, `priority-high`, `type-feature`

**Descripción:**
Cachear datos frecuentes.

**Datos a Cachear:**
- Políticas (TTL 1h)
- Categorías (TTL 24h)
- Permisos (TTL 1h)
- Analytics (TTL 5min)

**Criterios de Aceptación:**
- ✅ Cache funciona
- ✅ TTL correcto
- ✅ Invalidation funciona

---

### Issue #38: Configurar compression (Gzip) [COMPLETADO]
**Labels:** `phase-7-optimization`, `priority-medium`, `type-feature`

**Descripción:**
Comprimir responses HTTP.

**Criterios de Aceptación:**
- ✅ Responses comprimidos
- ✅ Reducción 60%+ en tamaño

---

## 🧪 FASE 8: TESTING Y DEPLOY (SPRINT 23-25)

### ✅ Issue #39: Escribir tests unitarios [COMPLETADO]
**Labels:** `phase-8-testing`, `priority-critical`, `type-chore`

**Descripción:**
Tests unitarios para services.

**Objetivo:** 80%+ coverage

**Criterios de Aceptación:**
- ✅ Tests pasan
- ✅ Coverage > 80%

---

### ✅ Issue #40: Escribir tests de integración [COMPLETADO]
**Labels:** `phase-8-testing`, `priority-high`, `type-chore`

**Descripción:**
Tests con base de datos real (dockertest).

**Criterios de Aceptación:**
- ✅ Tests de integración pasan
- ✅ Coverage > 70%

---

### Issue #41: Configurar GitHub Actions CI/CD
**Labels:** `phase-8-testing`, `priority-critical`, `type-chore`

**Workflow:**
1. Lint (golangci-lint)
2. Tests
3. Build
4. Deploy a staging

**Criterios de Aceptación:**
- ✅ Pipeline funciona
- ✅ Deploy automático

---

### Issue #42: Configurar Docker y Kubernetes
**Labels:** `phase-8-testing`, `priority-critical`, `type-chore`

**Tareas:**
- [ ] Dockerfile optimizado
- [ ] Docker Compose para local
- [ ] Kubernetes manifests
- [ ] Helm charts

**Criterios de Aceptación:**
- ✅ App corre en Kubernetes
- ✅ Escalamiento funciona

---

### Issue #43: Configurar monitoring (Prometheus + Grafana)
**Labels:** `phase-8-testing`, `priority-high`, `type-chore`

**Métricas:**
- Request rate
- Error rate
- Response time
- Database connections
- Cache hit rate

**Criterios de Aceptación:**
- ✅ Métricas exportadas
- ✅ Dashboards configurados
- ✅ Alertas configuradas

---

### Issue #44: Configurar logging (ELK Stack)
**Labels:** `phase-8-testing`, `priority-medium`, `type-chore`

**Stack:** Elasticsearch + Logstash + Kibana

**Criterios de Aceptación:**
- ✅ Logs centralizados
- ✅ Búsqueda funciona
- ✅ Dashboards configurados

---

### Issue #45: Security audit y penetration testing
**Labels:** `phase-8-testing`, `priority-critical`, `type-chore`

**Tareas:**
- [ ] OWASP Top 10 check
- [ ] SQL injection tests
- [ ] XSS tests
- [ ] CSRF tests
- [ ] Auth bypass tests

**Criterios de Aceptación:**
- ✅ No hay vulnerabilidades críticas
- ✅ Reporte de seguridad completo

---

## 📝 TAREAS GENERALES

### Issue #46: Actualizar documentación API (Swagger)
**Labels:** `type-docs`, `priority-medium`

**Tareas:**
- [ ] Swagger annotations en todos los endpoints
- [ ] Ejemplos de request/response
- [ ] Códigos de error documentados

**Criterios de Aceptación:**
- ✅ Swagger UI funcional
- ✅ Todos los endpoints documentados

---

### Issue #47: Crear README.md completo
**Labels:** `type-docs`, `priority-high`

**Contenido:**
- Descripción del proyecto
- Tech stack
- Instalación
- Configuración
- Uso
- Testing
- Deploy

**Criterios de Aceptación:**
- ✅ README completo y claro

---

## 📊 RESUMEN

**Total de Issues:** 47
**Fases:** 8
**Duración Estimada:** 25 semanas (~6 meses)

**Distribución:**
- Fase 1: 11 issues (4 semanas)
- Fase 2: 4 issues (3 semanas)
- Fase 3: 4 issues (3 semanas)
- Fase 4: 6 issues (4 semanas)
- Fase 5: 5 issues (2 semanas)
- Fase 6: 5 issues (4 semanas)
- Fase 7: 3 issues (2 semanas)
- Fase 8: 7 issues (3 semanas)
- General: 2 issues

---

## 🏁 CONVENCIONES DE BRANCHES

**Formato:** `{tipo}/issue-{numero}-{slug}`

**Tipos:**
- `feature/` - Nueva funcionalidad
- `bugfix/` - Corrección de bug
- `refactor/` - Refactorización
- `chore/` - Tareas de mantenimiento
- `docs/` - Documentación

**Ejemplos:**
```bash
git checkout -b feature/issue-1-expense-categories-migration
git checkout -b feature/issue-5-expense-domain
git checkout -b feature/issue-12-google-vision-integration
```

---

**Listo para crear issues en GitHub y comenzar implementación** 🚀

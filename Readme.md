# 🚀 MisViaticos Backend API

Backend para la aplicación MisViaticos - Sistema de gestión de gastos y viáticos empresariales construido con **Go**, **Fiber**, y **arquitectura hexagonal**.

## 🏗️ Arquitectura

Este proyecto implementa **Arquitectura Hexagonal (Clean Architecture)** con las siguientes características:

- ✅ **Multi-tenancy**: Manejo de múltiples empresas con bases de datos separadas
- ✅ **Separación por capas**: Dominio, servicios, adaptadores y controladores
- ✅ **Patrón Repository**: Para acceso a datos
- ✅ **Dependency Injection**: A través de interfaces (ports)
- ✅ **JWT Authentication**: Autenticación segura con refresh tokens
- ✅ **Rate Limiting**: Control de tasa de requests con Redis
- ✅ **Validaciones personalizadas**: Para RUT chileno, monedas, etc.

## 📁 Estructura del Proyecto

```
mv-backend/
├── cmd/api/                 # Punto de entrada de la aplicación
│   └── main.go
├── internal/                # Código privado de la aplicación
│   ├── core/               # Lógica de negocio (Hexágono Interior)
│   │   ├── auth/           # Autenticación y autorización
│   │   ├── user/           # Gestión de usuarios
│   │   ├── tenant/         # Gestión de empresas (multi-tenancy)
│   │   ├── expense/        # Gastos y viáticos
│   │   ├── receipt/        # Comprobantes y recibos
│   │   ├── category/       # Categorías de gastos
│   │   └── [module]/       # Cada módulo tiene:
│   │       ├── domain/     #   - Entidades y structs
│   │       ├── ports/      #   - Interfaces (Repository, Service)
│   │       ├── services/   #   - Lógica de negocio
│   │       └── adapters/   #   - Implementaciones (PostgreSQL, etc.)
│   ├── controllers/        # Controladores HTTP
│   ├── middleware/         # Middlewares (auth, validación, etc.)
│   ├── routes/             # Definición de rutas
│   ├── server/             # Configuración del servidor Fiber
│   ├── libraries/          # Clientes externos
│   │   ├── postgresql/     # Cliente PostgreSQL
│   │   ├── redis/          # Cliente Redis
│   │   └── aws/            # Servicios AWS S3
│   └── shared/             # Utilidades compartidas
│       ├── utils/          # Funciones utilitarias
│       └── validatorapi/   # Validaciones personalizadas
├── db/                     # Base de datos
│   ├── migrations/         # Migraciones BD control
│   ├── migrations-tenants/ # Migraciones BD tenants
│   └── seed/               # Datos iniciales
├── email_templates/        # Plantillas de email
├── scripts/               # Scripts de utilidades
└── [archivos config]      # Docker, env, Makefile, etc.
```

## 🛠️ Stack Tecnológico

### Core
- **Go 1.24.5** - Lenguaje de programación
- **Fiber v2** - Framework web (similar a Express.js)
- **PostgreSQL** - Base de datos principal
- **Redis** - Cache y rate limiting
- **JWT** - Autenticación y autorización

### Librerías Principales
- `pgx/v5` - Driver PostgreSQL de alto rendimiento
- `golang-migrate` - Migraciones de BD
- `validator/v10` - Validaciones de datos
- `aws-sdk-go-v2` - AWS S3 para almacenar comprobantes
- `resend-go` - Servicio de emails
- `uuid` - Generación de UUIDs

## 🚀 Comandos Disponibles

```bash
# Desarrollo
make dev-up          # Iniciar entorno de desarrollo (Docker)
make dev-down        # Detener entorno de desarrollo
make run-dev         # Ejecutar con hot reload (requiere air)

# Build
make build           # Compilar aplicación
make build-docker    # Construir imagen Docker

# Base de datos
make migrate-up      # Ejecutar migraciones
make migrate-down    # Revertir migración
make seed           # Poblar base de datos

# Testing
make test           # Ejecutar tests
make test-coverage  # Tests con coverage

# Calidad de código
make lint           # Ejecutar linter
make format         # Formatear código

# Dependencias
make deps           # Descargar dependencias
make deps-update    # Actualizar dependencias
```

## 🔧 Configuración

### Variables de Entorno

Copia `.env.example` a `.env` y configura:

```bash
# Aplicación
APP_NAME="MisViaticos API"
HOST=0.0.0.0
PORT=8080

# Base de datos
POSTGRESQL_CONTROL_DATABASE="misviaticos_control"
POSTGRESQL_DATABASE_TENANT="misviaticos_tenant"

# JWT
JWT_SECRET="tu_secret_super_seguro"

# Email
RESEND_API_KEY="tu_api_key"

# AWS (para comprobantes)
AWS_S3_BUCKET_RECEIPTS="misviaticos-receipts"
```

## 🏢 Sistema Multi-Tenant

### Arquitectura de Datos
- **BD Control**: Almacena usuarios, tenants, autenticación
- **BD por Tenant**: Cada empresa tiene su propia base de datos
- **Selección Dinámica**: El sistema conecta automáticamente a la BD del tenant activo

### Flujo de Trabajo
1. **Registro**: Usuario se registra en BD control
2. **Crear Empresa**: Usuario crea tenant (empresa)
3. **BD Específica**: Sistema crea BD para el tenant
4. **Selección**: Usuario selecciona tenant activo
5. **Operaciones**: Todas las operaciones se realizan en BD del tenant

## 📊 APIs Disponibles

### Rutas Públicas (`/api/v1`)
- `POST /auth/register` - Registro de usuario
- `POST /auth/login` - Inicio de sesión
- `POST /auth/forgot-password` - Recuperar contraseña
- `GET /info/currencies` - Monedas soportadas

### Rutas Privadas (requieren autenticación)

#### Gestión de Empresas
- `POST /tenant/create` - Crear empresa
- `PUT /tenant/update/:id` - Actualizar empresa
- `POST /tenant/select` - Seleccionar empresa activa

#### Gestión de Gastos
- `GET /expenses` - Listar gastos
- `POST /expenses` - Crear gasto
- `PUT /expenses/:id` - Actualizar gasto
- `POST /expenses/:id/submit` - Enviar para aprobación
- `POST /expenses/:id/approve` - Aprobar gasto

#### Comprobantes
- `POST /receipts/upload` - Subir comprobante
- `GET /receipts/:id` - Obtener comprobante

#### Reportes
- `GET /reports/expenses` - Reporte de gastos
- `GET /reports/export/excel` - Exportar a Excel
- `GET /reports/export/pdf` - Exportar a PDF

## 🔐 Autenticación

### JWT con Refresh Tokens
- **Access Token**: 24 horas de duración
- **Refresh Token**: 30 días de duración
- **Almacenamiento**: Tokens almacenados en Redis para revocación

### Estructura del JWT
```json
{
  "user_id": "uuid",
  "email": "usuario@empresa.com",
  "tenant_id": "uuid",
  "roles": ["employee", "admin"],
  "exp": 1640995200
}
```

## 🐳 Docker

### Desarrollo Local
```bash
# Iniciar todos los servicios
docker-compose up -d

# Solo base de datos y Redis
make dev-up
```

### Servicios Incluidos
- **PostgreSQL**: Base de datos principal
- **Redis**: Cache y rate limiting  
- **MinIO**: Almacenamiento S3-compatible para desarrollo

## 📝 Próximos Pasos

### Fase 1: Implementación Core ✅
- [x] Estructura base del proyecto
- [x] Servidor Fiber con middleware
- [x] Sistema de rutas públicas/privadas
- [x] Configuración Docker y variables de entorno
- [x] Librerías PostgreSQL y Redis

## 🤝 Contribución

Este proyecto sigue los patrones establecidos en `api-golang-2025`. Antes de contribuir:

1. **Entiende la arquitectura hexagonal**
2. **Respeta la estructura de directorios**
3. **Implementa interfaces en `ports/`**
4. **Sigue las convenciones de nomenclatura**
5. **Agrega validaciones personalizadas cuando sea necesario**

## 📄 Licencia

Proyecto propietario de MisViaticos © 2025
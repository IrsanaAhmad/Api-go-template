# 🚀 Go API Starter Kit

Template REST API berbasis **Go + Fiber** dengan arsitektur Clean Architecture, siap digunakan untuk production.

## Tech Stack

| Teknologi | Kegunaan |
|---|---|
| [Go](https://go.dev/) | Bahasa pemrograman |
| [Fiber v2](https://gofiber.io/) | HTTP framework |
| [PostgreSQL](https://www.postgresql.org/) / MSSQL | Database (multi-driver) |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Database migration |
| [JWT](https://github.com/golang-jwt/jwt) | Autentikasi |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Hashing password |
| [godotenv](https://github.com/joho/godotenv) | Environment variables |
| [Air](https://github.com/air-verse/air) | Hot reload |

---

## 📁 Struktur Proyek

```
.
├── cmd/
│   └── app/
│       └── main.go                 # Entry point aplikasi
│
├── internal/                        # Kode internal (tidak bisa diimport dari luar module)
│   ├── auth/
│   │   └── jwt.go                  # JWT generate, parse, hash token
│   │
│   ├── config/
│   │   ├── config.go               # Load .env → struct Config
│   │   └── config_schema.go        # Definisi struct Config + env tags
│   │
│   ├── database/
│   │   ├── db.go                   # Koneksi DB dinamis (PostgreSQL / MSSQL)
│   │   └── db_interface.go         # DBClient interface
│   │
│   ├── middleware/
│   │   ├── auth.go                 # JWT auth (baca dari cookie → fallback header)
│   │   ├── logger.go               # Request logger
│   │   ├── ratelimit.go            # Rate limiter per IP
│   │   ├── recovery.go             # (placeholder)
│   │   └── timeout.go              # (placeholder)
│   │
│   ├── modules/
│   │   └── auth/                   # Module autentikasi
│   │       ├── auth_route.go       # Route registration
│   │       ├── user_dto.go         # Request/Response DTO
│   │       ├── user_entity.go      # Entity / domain model
│   │       ├── user_handler.go     # HTTP handler (Fiber)
│   │       ├── user_repository.go  # Database queries
│   │       └── user_service.go     # Business logic
│   │
│   └── router/
│       ├── routes.go               # Root router (/api/v1)
│       └── auth_routes.go          # Bridge ke auth module
│
├── migrations/                      # SQL migration files
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_refresh_tokens_table.up.sql
│   └── 000002_create_refresh_tokens_table.down.sql
│
├── shared/                          # Shared utilities (bisa diimport siapapun)
│   ├── response/
│   │   └── response.go             # Standard API response (Success/Error)
│   └── utils/
│       ├── helper.go               # General helpers
│       ├── pagination.go           # Pagination utility
│       ├── sql_null.go             # SQL null type helpers
│       ├── sql_sanitizer.go        # SQL input sanitizer
│       └── validator.go            # Request validator
│
├── .env.example                     # Template environment variables
├── .air.toml                        # Air hot-reload config
├── go.mod
└── go.sum
```

---

## 🏗️ Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan alur dependency satu arah:

```
Handler → Service → Repository → Database
   ↑                                 ↑
  DTO                             Entity
```

| Layer | Tanggung Jawab | Contoh File |
|---|---|---|
| **Handler** | Parse request, validasi input, format response | `user_handler.go` |
| **Service** | Business logic, orchestration | `user_service.go` |
| **Repository** | Query database, mapping ke entity | `user_repository.go` |
| **Entity** | Domain model / struct data | `user_entity.go` |
| **DTO** | Request/Response contract | `user_dto.go` |

---

## ⚡ Quick Start

### 1. Clone & Install Dependencies

```bash
git clone https://github.com/IrsanaAhmad/Api-go-template.git
cd Api-go-template
go mod tidy
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env sesuai kebutuhan (DATABASE_URL, JWT_SECRET_KEY, dll)
```

### 3. Install Tools

```bash
# golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Air hot-reload
go install github.com/air-verse/air@latest
```

### 4. Jalankan Migration

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### 5. Run Server

```bash
# Development (hot-reload)
air

# Atau tanpa hot-reload
go run ./cmd/app
```

Server akan berjalan di `http://localhost:8080`.

---

## 🔐 Auth Endpoints

Semua endpoint auth berada di prefix `/api/v1/auth`.

| Method | Path | Auth? | Deskripsi |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | ❌ | Registrasi user baru |
| `POST` | `/api/v1/auth/login` | ❌ | Login, set token di cookie |
| `POST` | `/api/v1/auth/logout` | ❌ | Revoke refresh token, hapus cookie |
| `POST` | `/api/v1/auth/refresh` | ❌ | Rotasi refresh token via cookie |

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123","full_name":"John Doe"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}' \
  -c cookies.txt   # Simpan cookies
```

Response body berisi data user. Token disimpan di **HttpOnly cookies**:
- `access_token` — path `/`, expires sesuai `JWT_ACCESS_TOKEN_EXP_MINUTES`
- `refresh_token` — path `/api/v1/auth`, expires sesuai `JWT_REFRESH_TOKEN_EXP_DAYS`

### Logout

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout -b cookies.txt
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh -b cookies.txt -c cookies.txt
```

---

## 🔌 Menambahkan Module Baru

Contoh: membuat module `products`.

### 1. Buat Folder Module

```
internal/modules/products/
├── product_entity.go
├── product_dto.go
├── product_repository.go
├── product_service.go
├── product_handler.go
└── product_route.go
```

### 2. Entity (`product_entity.go`)

```go
package products

import "time"

type Product struct {
    ID        string
    Name      string
    Price     float64
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 3. DTO (`product_dto.go`)

```go
package products

type CreateProductRequest struct {
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type ProductDTO struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

### 4. Repository (`product_repository.go`)

```go
package products

import (
    "context"
    "github.com/IrsanaAhmad/go-starter-kit/internal/database"
)

type ProductRepository interface {
    Create(ctx context.Context, product *Product) (*Product, error)
    FindAll(ctx context.Context, limit, offset int) ([]Product, int, error)
}

type SQLProductRepository struct {
    db database.DBClient
}

func NewSQLProductRepository(db database.DBClient) *SQLProductRepository {
    return &SQLProductRepository{db: db}
}

// Implementasi method Create, FindAll, dll...
```

### 5. Service (`product_service.go`)

```go
package products

import "context"

type ProductService interface {
    Create(ctx context.Context, req *CreateProductRequest) (*ProductDTO, error)
}

type productService struct {
    repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
    return &productService{repo: repo}
}

// Implementasi business logic...
```

### 6. Handler (`product_handler.go`)

```go
package products

import (
    "github.com/IrsanaAhmad/go-starter-kit/shared/response"
    "github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
    service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}

// Implementasi handler methods...
```

### 7. Route (`product_route.go`)

```go
package products

import (
    "github.com/IrsanaAhmad/go-starter-kit/internal/database"
    "github.com/IrsanaAhmad/go-starter-kit/internal/middleware"
    "github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, db database.DBClient) {
    repo := NewSQLProductRepository(db)
    service := NewProductService(repo)
    handler := NewProductHandler(service)

    group := router.Group("/products")
    group.Get("/", handler.GetAll)
    group.Post("/", middleware.JWTAuth(), handler.Create)
}
```

### 8. Daftarkan ke Router

**`internal/router/product_routes.go`** (file baru):

```go
package router

import (
    "github.com/IrsanaAhmad/go-starter-kit/internal/database"
    "github.com/IrsanaAhmad/go-starter-kit/internal/modules/products"
    "github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(v1 fiber.Router, db database.DBClient) {
    products.RegisterRoutes(v1, db)
}
```

**`internal/router/routes.go`** (tambahkan 1 baris):

```go
func Register(app *fiber.App, db database.DBClient) {
    api := app.Group("/api")
    v1 := api.Group("/v1")

    RegisterAuthRoutes(v1, db)
    RegisterProductRoutes(v1, db)  // ← tambahkan ini
}
```

---

## 🗃️ Database Migration

### Membuat Migration Baru

```bash
migrate create -ext sql -dir migrations -seq nama_migration
```

Contoh:

```bash
migrate create -ext sql -dir migrations -seq create_products_table
```

Ini akan membuat 2 file:

```
migrations/000003_create_products_table.up.sql     # DDL CREATE
migrations/000003_create_products_table.down.sql   # DDL DROP (rollback)
```

### Menjalankan Migration

```bash
# Up (semua pending)
migrate -path migrations -database "$DATABASE_URL" up

# Up N step
migrate -path migrations -database "$DATABASE_URL" up 1

# Down N step (rollback)
migrate -path migrations -database "$DATABASE_URL" down 1

# Force versi (jika dirty state)
migrate -path migrations -database "$DATABASE_URL" force VERSION
```

### Konvensi Migration

- Gunakan `IF NOT EXISTS` / `IF EXISTS` untuk idempotency
- Selalu buat `*.down.sql` yang benar agar bisa rollback
- Gunakan `TIMESTAMPTZ` untuk kolom waktu (timezone-aware)
- Gunakan `UUID` sebagai primary key
- Tambahkan `created_at` dan `updated_at` di setiap tabel

---

## ⚙️ Konfigurasi

Semua konfigurasi dibaca dari `.env` (12-Factor App). Lihat `.env.example` untuk referensi.

### Database

```env
# Prioritas 1: Jika diisi, semua DB_* akan diabaikan
DATABASE_URL=postgresql://user:password@host/dbname?sslmode=require

# Prioritas 2: Fallback jika DATABASE_URL kosong
DB_CONNECTION=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydb
DB_USERNAME=postgres
DB_PASSWORD=secret
```

### JWT

```env
JWT_SECRET_KEY=your-secret-key-here
JWT_ACCESS_TOKEN_EXP_MINUTES=30
JWT_REFRESH_TOKEN_EXP_DAYS=7
```

---

## 🛡️ Middleware

| Middleware | Scope | Deskripsi |
|---|---|---|
| **Logger** | Global | Log setiap request (method, path, status, latency, IP) |
| **RateLimiter** | Global | 60 req/menit per IP (in-memory token bucket) |
| **JWTAuth** | Per-route | Validasi token dari cookie / Authorization header |

### Menggunakan JWTAuth di Route

```go
// Public route
group.Get("/products", handler.GetAll)

// Protected route
group.Post("/products", middleware.JWTAuth(), handler.Create)
```

Dalam handler, akses user info:

```go
func (h *Handler) Create(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    username := c.Locals("username").(string)
    // ...
}
```

---

## 📝 API Response Format

Semua response menggunakan format standar:

```json
// Success
{
    "success": true,
    "message": "operasi berhasil",
    "data": { ... }
}

// Error
{
    "success": false,
    "message": "deskripsi error",
    "data": null
}
```

---

## 📦 Shared Utilities

| File | Kegunaan |
|---|---|
| `response/response.go` | Standard `Success()` / `Error()` response |
| `utils/pagination.go` | Pagination helper (page, limit, links) |
| `utils/validator.go` | Request validation |
| `utils/sql_null.go` | SQL null type helpers |
| `utils/sql_sanitizer.go` | SQL input sanitizer |
| `utils/helper.go` | General helpers |

---

## 🚢 Deployment

### Build

```bash
go build -o server ./cmd/app
```

### Environment Variables di Production

Set `APP_ENV=production` untuk mengaktifkan:
- Cookie `Secure` flag (HTTPS only)
- Dan konfigurasi production lainnya

```bash
APP_ENV=production ./server
```

---

## License

MIT

# AGENTS.md — KasirAI

> Baca seluruh file ini sebelum menulis satu baris kode pun.
> Jangan berasumsi. Jangan mengarang library. Jangan skip seksi manapun.

---

## 0. ATURAN MUTLAK (WAJIB DIIKUTI, TIDAK BOLEH DILANGGAR)

1. **JANGAN pernah membuat file di luar struktur yang didefinisikan di Seksi 5.**
2. **JANGAN pernah menggunakan library/package yang tidak ada di Seksi 6.**
3. **JANGAN pernah membuat mock data atau dummy response untuk endpoint yang sudah memiliki handler nyata.**
4. **JANGAN pernah menggabungkan logic bisnis ke dalam handler HTTP.** Semua logic bisnis harus di `internal/service/`.
5. **JANGAN pernah menulis query SQL langsung di handler atau service.** Semua query ada di `internal/repository/`.
6. **JANGAN pernah hardcode credential, secret, atau URL.** Semua harus dari environment variable via `config/`.
7. **JANGAN pernah skip error handling.** Setiap `error` yang dikembalikan WAJIB ditangani.
8. **JANGAN pernah menggunakan `fmt.Println` untuk logging.** Gunakan `slog` dari standard library Go.
9. **JANGAN pernah membuat file `.env`.** Hanya boleh ada `.env.example`.
10. **JANGAN membuat migration baru jika tabel sudah ada di Seksi 7.** Gunakan migration yang sudah didefinisikan.
11. **Jika ragu antara dua pendekatan, pilih yang lebih sederhana dan tulis komentar alasannya.**
12. **Jangan pernah menghapus kode yang sudah ada tanpa instruksi eksplisit dari user.**

---

## 1. OVERVIEW PROYEK

**Nama Aplikasi:** KasirAI  
**Tagline:** Catat keuangan UMKM via WhatsApp, pantau lewat dashboard, siap KUR.  
**Target pengguna:** Pelaku UMKM Indonesia non-teknis (warung, pedagang, pengrajin).  
**Bahasa UI:** Bahasa Indonesia seluruhnya.

### Core Features (MVP — harus selesai semua)

| #   | Fitur                 | Deskripsi                                                                        |
| --- | --------------------- | -------------------------------------------------------------------------------- |
| 1   | WhatsApp Webhook      | Terima pesan teks dari UMKM via Fonnte, parse transaksi dengan AI                |
| 2   | AI Transaction Parser | DeepSeek V4-Flash parse natural language → structured transaction                |
| 3   | Financial Dashboard   | Vue 3 — tampilkan ringkasan pemasukan/pengeluaran/laba harian, mingguan, bulanan |
| 4   | KUR Readiness Score   | Hitung skor kesiapan kredit UMKM berdasarkan history transaksi                   |
| 5   | Laporan Bulanan       | Generate ringkasan teks laporan keuangan bulanan, kirim via WhatsApp             |
| 6   | Registrasi UMKM       | UMKM daftar via link web sederhana, dapat nomor WhatsApp bot                     |

### Out of Scope (JANGAN dibangun kecuali diminta)

- Voice note processing
- Multi-currency
- Integrasi marketplace
- Multi-user per UMKM
- Role-based access
- File upload

---

## 2. ARSITEKTUR SISTEM

```
┌─────────────────────────────────────────────────────────────┐
│                        FRONTEND                             │
│          Vue 3 + Vite  (port 5173 dev / Vercel prod)        │
│          Axios → REST API                                   │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP/JSON
┌───────────────────────────▼─────────────────────────────────┐
│                        BACKEND                              │
│          Go + Gin  (port 8080 dev / Render Free Tier prod)           │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────┐  │
│  │ handler/ │→ │ service/ │→ │repository/│→ │PostgreSQL│  │
│  └──────────┘  └────┬─────┘  └───────────┘  └──────────┘  │
│                     │                                       │
│              ┌──────▼──────┐  ┌──────────────┐             │
│              │  pkg/ai/    │  │  pkg/fonnte/ │             │
│              │  (DeepSeek) │  │  (WhatsApp)  │             │
│              └─────────────┘  └──────────────┘             │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────▼──────────────────┐
        │           EXTERNAL SERVICES          │
        │  DeepSeek API  |  Fonnte API         │
        └──────────────────────────────────────┘
```

### Request Flow WhatsApp

```
User WhatsApp → Fonnte → POST /webhook/whatsapp
  → handler.WebhookHandler
  → service.TransactionService.ParseAndSave(rawMessage)
    → pkg/ai.ParseTransaction(rawMessage) → DeepSeek API
    → repository.TransactionRepo.Create(parsedTx)
  → pkg/fonnte.SendMessage(reply)
```

### Request Flow Dashboard

```
Browser → GET /api/v1/dashboard/summary?umkm_id=X&period=monthly
  → middleware.AuthMiddleware (validate JWT)
  → handler.DashboardHandler
  → service.DashboardService.GetSummary(umkmID, period)
    → repository.TransactionRepo.SumByPeriod(...)
  → JSON response
```

---

## 3. KEPUTUSAN TEKNIS (TIDAK BOLEH DIUBAH TANPA DISKUSI)

### Backend

- **Framework:** Gin (github.com/gin-gonic/gin v1.9.1) — bukan Echo, bukan Fiber, bukan chi
- **Database driver:** pgx/v5 (github.com/jackc/pgx/v5) — bukan database/sql langsung, bukan gorm, bukan sqlx
- **Migration:** golang-migrate/migrate (github.com/golang-migrate/migrate/v4) — file SQL di `backend/migrations/`
- **JWT:** golang-jwt/jwt (github.com/golang-jwt/jwt/v5)
- **Env loading:** godotenv (github.com/joho/godotenv) — hanya untuk development
- **Logging:** `log/slog` (standard library, Go 1.21+) — bukan zerolog, bukan zap, bukan logrus
- **HTTP client:** `net/http` standard library — bukan resty, bukan heimdall
- **Validasi:** go-playground/validator (github.com/go-playground/validator/v10)
- **CORS:** github.com/gin-contrib/cors

### Frontend

- **Framework:** Vue 3 dengan Composition API (`<script setup>`) — BUKAN Options API
- **Build tool:** Vite (latest)
- **State management:** Pinia — bukan Vuex
- **Router:** Vue Router 4
- **HTTP client:** Axios
- **UI Components:** PrimeVue 4 — bukan Vuetify, bukan Element Plus, bukan Naive UI
- **Charts:** Chart.js via vue-chartjs — bukan ECharts, bukan ApexCharts
- **CSS:** PrimeFlex + Tailwind CSS (utility only, tidak pakai component Tailwind)
- **Formatter:** Prettier + ESLint

### Database

- **PostgreSQL 15** — bukan MySQL, bukan SQLite, bukan MongoDB
- Semua timestamp menggunakan `TIMESTAMPTZ` (timezone-aware)
- Semua ID menggunakan `UUID` (`gen_random_uuid()`)
- Soft delete menggunakan kolom `deleted_at TIMESTAMPTZ NULL`

### External Services

- **WhatsApp:** Fonnte API (https://fonnte.com) — bukan Twilio, bukan Meta Cloud API
  - Endpoint webhook: `POST /webhook/whatsapp`
  - Fonnte kirim `multipart/form-data` dengan field: `sender`, `message`, `device`
- **AI:** DeepSeek API — model `deepseek-v4-flash` — bukan openai, bukan gemini
  - Base URL: `https://api.deepseek.com`
  - Format: OpenAI-compatible (`/v1/chat/completions`)

---

## 4. ENVIRONMENT VARIABLES

File `.env.example` harus ada di root project. Jangan pernah buat file `.env`.

```env
# Server
APP_ENV=development
APP_PORT=8080
APP_BASE_URL=http://localhost:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=kasiraiai
DB_USER=postgres
DB_PASSWORD=
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5

# JWT
JWT_SECRET=ganti_dengan_secret_panjang_minimal_32_karakter
JWT_EXPIRY_HOURS=72

# DeepSeek AI
DEEPSEEK_API_KEY=
DEEPSEEK_BASE_URL=https://api.deepseek.com
DEEPSEEK_MODEL=deepseek-v4-flash
DEEPSEEK_MAX_TOKENS=1000
DEEPSEEK_TIMEOUT_SECONDS=30

# Fonnte WhatsApp
FONNTE_API_KEY=
FONNTE_BASE_URL=https://api.fonnte.com
FONNTE_DEVICE=

# Frontend (Vite env vars — prefix VITE_)
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

---

## 5. STRUKTUR DIREKTORI (WAJIB DIIKUTI PERSIS)

```
kasiraiai/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go              # Entry point. Hanya inisialisasi, tidak ada logic.
│   ├── config/
│   │   └── config.go                # Baca semua env var ke struct Config
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── auth_handler.go      # Register, Login
│   │   │   ├── webhook_handler.go   # Terima webhook dari Fonnte
│   │   │   ├── transaction_handler.go
│   │   │   ├── dashboard_handler.go
│   │   │   ├── umkm_handler.go
│   │   │   └── response.go          # Helper: SuccessResponse, ErrorResponse
│   │   ├── middleware/
│   │   │   ├── auth.go              # JWT validation middleware
│   │   │   ├── logger.go            # Request/response logging
│   │   │   └── cors.go              # CORS config
│   │   ├── model/
│   │   │   ├── umkm.go
│   │   │   ├── transaction.go
│   │   │   ├── user.go
│   │   │   └── kur_score.go
│   │   ├── repository/
│   │   │   ├── db.go                # Init pgx pool
│   │   │   ├── umkm_repo.go
│   │   │   ├── transaction_repo.go
│   │   │   └── user_repo.go
│   │   └── service/
│   │       ├── auth_service.go
│   │       ├── transaction_service.go  # Core: parse + save transaksi
│   │       ├── dashboard_service.go
│   │       ├── kur_service.go          # Hitung KUR readiness score
│   │       └── report_service.go       # Generate laporan bulanan
│   ├── pkg/
│   │   ├── ai/
│   │   │   ├── client.go            # DeepSeek HTTP client
│   │   │   └── parser.go            # ParseTransaction, prompt templates
│   │   └── fonnte/
│   │       └── client.go            # SendMessage, kirim WhatsApp
│   ├── migrations/
│   │   ├── 000001_init_schema.up.sql
│   │   ├── 000001_init_schema.down.sql
│   │   ├── 000002_seed_categories.up.sql
│   │   └── 000002_seed_categories.down.sql
│   ├── go.mod                       # module: kasiraiai/backend
│   ├── go.sum
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── main.js
│   │   ├── App.vue
│   │   ├── router/
│   │   │   └── index.js
│   │   ├── stores/
│   │   │   ├── auth.js
│   │   │   ├── transaction.js
│   │   │   └── dashboard.js
│   │   ├── api/
│   │   │   ├── axios.js             # Instance axios dengan interceptor
│   │   │   ├── auth.js
│   │   │   ├── transaction.js
│   │   │   └── dashboard.js
│   │   ├── views/
│   │   │   ├── LoginView.vue
│   │   │   ├── RegisterView.vue
│   │   │   ├── DashboardView.vue
│   │   │   ├── TransactionView.vue
│   │   │   ├── KurScoreView.vue
│   │   │   └── ReportView.vue
│   │   ├── components/
│   │   │   ├── layout/
│   │   │   │   ├── AppSidebar.vue
│   │   │   │   ├── AppTopbar.vue
│   │   │   │   └── AppLayout.vue
│   │   │   ├── dashboard/
│   │   │   │   ├── SummaryCard.vue
│   │   │   │   ├── IncomeExpenseChart.vue
│   │   │   │   └── RecentTransactions.vue
│   │   │   ├── transaction/
│   │   │   │   └── TransactionTable.vue
│   │   │   └── kur/
│   │   │       ├── KurScoreGauge.vue
│   │   │       └── KurRecommendations.vue
│   │   └── assets/
│   │       └── styles/
│   │           └── main.css
│   ├── index.html
│   ├── vite.config.js
│   ├── package.json
│   └── .env.example
├── docker-compose.yml               # PostgreSQL untuk development
├── .gitignore
└── AGENTS.md                        # File ini
```

---

## 6. DEPENDENCIES YANG DIIZINKAN

### Backend (go.mod)

```
module kasiraiai/backend

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gin-contrib/cors v1.5.0
    github.com/jackc/pgx/v5 v5.5.4
    github.com/golang-migrate/migrate/v4 v4.17.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/joho/godotenv v1.5.1
    github.com/go-playground/validator/v10 v10.19.0
    github.com/google/uuid v1.6.0
    golang.org/x/crypto v0.21.0
)
```

> **DILARANG menambah dependency tanpa konfirmasi eksplisit dari user.**

### Frontend (package.json dependencies)

```json
{
  "dependencies": {
    "vue": "^3.4.0",
    "vue-router": "^4.3.0",
    "pinia": "^2.1.0",
    "axios": "^1.6.0",
    "primevue": "^4.0.0",
    "primeicons": "^7.0.0",
    "primeflex": "^3.3.0",
    "chart.js": "^4.4.0",
    "vue-chartjs": "^5.3.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^5.0.0",
    "vite": "^5.0.0",
    "eslint": "^8.57.0",
    "prettier": "^3.2.0"
  }
}
```

---

## 7. DATABASE SCHEMA

Jalankan migration dalam urutan ini. Jangan buat tabel yang tidak ada di sini.

### Migration 000001: Init Schema

```sql
-- 000001_init_schema.up.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Tabel UMKM (pemilik usaha)
CREATE TABLE umkms (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    phone_number  VARCHAR(20)  NOT NULL UNIQUE, -- format: 628XXXXXXXXXX
    email         VARCHAR(255) UNIQUE,
    address       TEXT,
    business_type VARCHAR(100), -- contoh: 'kuliner', 'fashion', 'jasa', 'pertanian'
    password_hash VARCHAR(255) NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Tabel kategori transaksi
CREATE TABLE transaction_categories (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL, -- 'Penjualan', 'Pembelian Bahan', 'Gaji', dll
    type       VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    icon       VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabel transaksi
CREATE TABLE transactions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    umkm_id       UUID NOT NULL REFERENCES umkms(id),
    category_id   UUID REFERENCES transaction_categories(id),
    amount        BIGINT NOT NULL CHECK (amount > 0), -- dalam satuan Rupiah, BUKAN desimal
    type          VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    description   TEXT NOT NULL,           -- deskripsi asli dari AI
    raw_message   TEXT NOT NULL,           -- pesan WhatsApp asli dari UMKM
    transaction_date DATE NOT NULL,        -- tanggal transaksi (bisa berbeda dari created_at)
    source        VARCHAR(20) NOT NULL DEFAULT 'whatsapp' CHECK (source IN ('whatsapp', 'manual')),
    ai_confidence DECIMAL(3,2),            -- skor kepercayaan AI 0.00-1.00
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

-- Tabel KUR score history
CREATE TABLE kur_scores (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    umkm_id             UUID NOT NULL REFERENCES umkms(id),
    score               INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
    level               VARCHAR(20) NOT NULL CHECK (level IN ('rendah', 'sedang', 'baik', 'sangat_baik')),
    monthly_income_avg  BIGINT NOT NULL,   -- rata-rata pemasukan 3 bulan terakhir
    monthly_expense_avg BIGINT NOT NULL,
    profit_margin       DECIMAL(5,2),      -- persentase margin laba
    consistency_score   INTEGER,           -- konsistensi pencatatan (0-100)
    months_of_data      INTEGER NOT NULL,  -- berapa bulan data yang dianalisis
    recommendations     TEXT[],            -- array rekomendasi dalam Bahasa Indonesia
    calculated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_transactions_umkm_id ON transactions(umkm_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_deleted_at ON transactions(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_kur_scores_umkm_id ON kur_scores(umkm_id);
CREATE INDEX idx_umkms_phone ON umkms(phone_number);
```

```sql
-- 000001_init_schema.down.sql
DROP TABLE IF EXISTS kur_scores;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS transaction_categories;
DROP TABLE IF EXISTS umkms;
DROP EXTENSION IF EXISTS "pgcrypto";
```

### Migration 000002: Seed Categories

```sql
-- 000002_seed_categories.up.sql
INSERT INTO transaction_categories (id, name, type, icon) VALUES
    (gen_random_uuid(), 'Penjualan Produk', 'income', 'pi pi-shopping-bag'),
    (gen_random_uuid(), 'Penjualan Jasa', 'income', 'pi pi-briefcase'),
    (gen_random_uuid(), 'Pendapatan Lain', 'income', 'pi pi-plus-circle'),
    (gen_random_uuid(), 'Pembelian Bahan Baku', 'expense', 'pi pi-box'),
    (gen_random_uuid(), 'Gaji Karyawan', 'expense', 'pi pi-users'),
    (gen_random_uuid(), 'Sewa Tempat', 'expense', 'pi pi-home'),
    (gen_random_uuid(), 'Utilitas (Listrik/Air)', 'expense', 'pi pi-bolt'),
    (gen_random_uuid(), 'Transportasi', 'expense', 'pi pi-car'),
    (gen_random_uuid(), 'Pemasaran', 'expense', 'pi pi-megaphone'),
    (gen_random_uuid(), 'Peralatan', 'expense', 'pi pi-wrench'),
    (gen_random_uuid(), 'Pengeluaran Lain', 'expense', 'pi pi-minus-circle');
```

```sql
-- 000002_seed_categories.down.sql
DELETE FROM transaction_categories;
```

---

## 8. API CONTRACT

Base URL: `/api/v1`  
Format: JSON (`Content-Type: application/json`)  
Auth: Bearer JWT di header `Authorization: Bearer <token>`

### Response Format — WAJIB konsisten

```json
// Success
{
  "success": true,
  "message": "Berhasil",
  "data": { ... }
}

// Error
{
  "success": false,
  "message": "Pesan error dalam Bahasa Indonesia",
  "errors": ["detail error 1", "detail error 2"]  // opsional
}
```

### Endpoints

#### Auth

```
POST /api/v1/auth/register
Body: { "name", "business_name", "phone_number", "email"(opt), "password", "business_type"(opt) }
Response 201: { "success": true, "data": { "token": "...", "umkm": {...} } }

POST /api/v1/auth/login
Body: { "phone_number", "password" }
Response 200: { "success": true, "data": { "token": "...", "umkm": {...} } }
```

#### Webhook (NO AUTH — Fonnte tidak kirim JWT)

```
POST /webhook/whatsapp
Content-Type: multipart/form-data
Fields: sender (string), message (string), device (string)
Response 200: { "success": true, "message": "Pesan diterima" }

PENTING:
- Endpoint ini TIDAK di bawah /api/v1
- Endpoint ini TIDAK membutuhkan JWT
- Validasi dengan FONNTE_WEBHOOK_SECRET jika Fonnte support (cek doc Fonnte)
- Selalu return 200 ke Fonnte meskipun parsing gagal (kirim pesan error ke user via WhatsApp)
```

#### Transactions

```
GET  /api/v1/transactions
Query: page(int,default:1), limit(int,default:20), type(income|expense),
       start_date(YYYY-MM-DD), end_date(YYYY-MM-DD)
Response 200: { "success": true, "data": { "transactions": [...], "total": int, "page": int } }

POST /api/v1/transactions
Body: { "amount", "type", "description", "transaction_date", "category_id"(opt) }
Response 201: { "success": true, "data": { "transaction": {...} } }

DELETE /api/v1/transactions/:id
Response 200: { "success": true, "message": "Transaksi dihapus" }
// Soft delete: set deleted_at, JANGAN hapus dari database
```

#### Dashboard

```
GET /api/v1/dashboard/summary
Query: period (daily|weekly|monthly|yearly), date (YYYY-MM-DD, default: today)
Response 200:
{
  "success": true,
  "data": {
    "period": "monthly",
    "total_income": 5000000,
    "total_expense": 2000000,
    "net_profit": 3000000,
    "profit_margin": 60.0,
    "transaction_count": 45,
    "chart_data": [
      { "date": "2026-06-01", "income": 200000, "expense": 50000 }
    ]
  }
}

GET /api/v1/dashboard/categories
Query: period, type (income|expense)
Response 200: { "success": true, "data": { "categories": [ { "name", "total", "percentage" } ] } }
```

#### KUR Score

```
GET /api/v1/kur/score
Response 200:
{
  "success": true,
  "data": {
    "score": 72,
    "level": "baik",
    "monthly_income_avg": 4500000,
    "monthly_expense_avg": 2000000,
    "profit_margin": 55.5,
    "consistency_score": 80,
    "months_of_data": 3,
    "recommendations": [
      "Tingkatkan konsistensi pencatatan harian untuk meningkatkan skor",
      "Margin laba Anda sudah baik, pertahankan efisiensi pengeluaran"
    ],
    "calculated_at": "2026-06-08T10:00:00Z"
  }
}

POST /api/v1/kur/recalculate
Response 200: { "success": true, "data": { "score": {...} } }
```

#### Report

```
POST /api/v1/reports/monthly
Body: { "year": 2026, "month": 6 }
Response 200:
{
  "success": true,
  "data": {
    "report_text": "Laporan Keuangan Juni 2026...",
    "total_income": ...,
    "total_expense": ...,
    "net_profit": ...,
    "top_categories": [...]
  }
}
// Setelah generate, otomatis kirim via WhatsApp ke UMKM
```

---

## 9. LOGIKA AI — PARSING TRANSAKSI

### Prompt Template (WAJIB digunakan persis ini, simpan di `pkg/ai/parser.go`)

```
System prompt:
Kamu adalah asisten keuangan UMKM Indonesia. Tugasmu adalah mengekstrak informasi transaksi dari pesan WhatsApp pelaku UMKM.

Balas HANYA dengan JSON valid, tidak ada teks lain di luar JSON.

Format JSON yang harus dikembalikan:
{
  "is_transaction": boolean,
  "type": "income" | "expense" | null,
  "amount": integer (dalam Rupiah, tanpa desimal) | null,
  "description": "deskripsi singkat dalam Bahasa Indonesia" | null,
  "category_hint": "nama kategori yang paling cocok dari daftar" | null,
  "transaction_date": "YYYY-MM-DD" | null,
  "confidence": float (0.0 - 1.0),
  "reply_message": "pesan balasan singkat dan ramah dalam Bahasa Indonesia untuk dikirim ke pengguna"
}

Aturan parsing:
- Jika pesan bukan transaksi (salam, pertanyaan, dll): is_transaction=false, reply_message berisi respons yang membantu
- Tanggal: jika tidak disebutkan, gunakan hari ini
- Kata "jual/dapat/masuk/terima" → income
- Kata "beli/bayar/keluar/bayarin/belanja" → expense
- Format angka Indonesia: "5rb" = 5000, "2jt" = 2000000, "500k" = 500000
- Jika ambigu, confidence rendah (< 0.6) dan minta klarifikasi di reply_message

Daftar kategori: Penjualan Produk, Penjualan Jasa, Pendapatan Lain, Pembelian Bahan Baku, Gaji Karyawan, Sewa Tempat, Utilitas (Listrik/Air), Transportasi, Pemasaran, Peralatan, Pengeluaran Lain

User prompt:
Tanggal hari ini: {TODAY_DATE}
Pesan WhatsApp: "{RAW_MESSAGE}"
```

### Contoh Parsing

```
Input: "tadi jual nasi uduk 5 bungkus @15rb"
Output:
{
  "is_transaction": true,
  "type": "income",
  "amount": 75000,
  "description": "Penjualan nasi uduk 5 bungkus",
  "category_hint": "Penjualan Produk",
  "transaction_date": "2026-06-08",
  "confidence": 0.95,
  "reply_message": "✅ Tercatat! Pemasukan Rp75.000 dari penjualan nasi uduk. Total hari ini: Rp[total]"
}

Input: "bayar listrik 250rb"
Output:
{
  "is_transaction": true,
  "type": "expense",
  "amount": 250000,
  "description": "Pembayaran listrik",
  "category_hint": "Utilitas (Listrik/Air)",
  "transaction_date": "2026-06-08",
  "confidence": 0.98,
  "reply_message": "✅ Tercatat! Pengeluaran Rp250.000 untuk listrik."
}

Input: "halo, ini toko saya bisa dilihat nggak"
Output:
{
  "is_transaction": false,
  "type": null,
  "amount": null,
  "description": null,
  "category_hint": null,
  "transaction_date": null,
  "confidence": 1.0,
  "reply_message": "Halo! 👋 Saya KasirAI, asisten keuangan Anda. Kirim pesan seperti 'jual bakso 50rb' atau 'bayar bahan 200rb' untuk mencatat transaksi. Ketik *laporan* untuk melihat ringkasan keuangan Anda."
}
```

---

## 10. LOGIKA KUR READINESS SCORE

Hitung score di `internal/service/kur_service.go` berdasarkan formula ini.

### Formula

```
Total Score (0-100) =
  income_stability_score  (bobot 30%) +
  profit_margin_score     (bobot 25%) +
  consistency_score       (bobot 25%) +
  data_longevity_score    (bobot 20%)
```

### Detail Kalkulasi

```
income_stability_score (0-30):
  Ambil pemasukan 3 bulan terakhir
  Hitung koefisien variasi (CV = stddev / mean)
  CV < 0.2 → 30 poin (sangat stabil)
  CV 0.2-0.4 → 20 poin
  CV 0.4-0.6 → 10 poin
  CV > 0.6 → 0 poin (tidak stabil)

profit_margin_score (0-25):
  profit_margin = (total_income - total_expense) / total_income * 100
  margin >= 30% → 25 poin
  margin 20-30% → 18 poin
  margin 10-20% → 10 poin
  margin < 10% → 0 poin

consistency_score (0-25):
  Hitung berapa persen hari kerja dalam 30 hari terakhir ada transaksi
  >= 80% hari ada transaksi → 25 poin
  60-80% → 18 poin
  40-60% → 10 poin
  < 40% → 0 poin

data_longevity_score (0-20):
  Berapa bulan UMKM sudah mencatat di sistem
  >= 6 bulan → 20 poin
  3-5 bulan → 14 poin
  2 bulan → 8 poin
  1 bulan → 4 poin
  < 1 bulan → 0 poin
```

### Level KUR

```
Score 80-100 → "sangat_baik" → "Sangat siap mengajukan KUR"
Score 60-79  → "baik"        → "Siap mengajukan KUR dengan persiapan tambahan"
Score 40-59  → "sedang"      → "Perlu meningkatkan konsistensi pencatatan"
Score 0-39   → "rendah"      → "Fokus dulu pada pencatatan rutin minimal 3 bulan"
```

---

## 11. KONVENSI KODE

### Go

```go
// Penamaan package: lowercase, singkat, tanpa underscore
package repository

// Struct: PascalCase
type Transaction struct {
    ID          uuid.UUID  `json:"id" db:"id"`
    UmkmID      uuid.UUID  `json:"umkm_id" db:"umkm_id"`
    Amount      int64      `json:"amount" db:"amount"`
    // ...
}

// Interface untuk setiap repository dan service
type TransactionRepository interface {
    Create(ctx context.Context, tx *model.Transaction) error
    FindByUmkmID(ctx context.Context, umkmID uuid.UUID, filter TransactionFilter) ([]model.Transaction, int, error)
    SumByPeriod(ctx context.Context, umkmID uuid.UUID, start, end time.Time) (SummaryResult, error)
    SoftDelete(ctx context.Context, id uuid.UUID, umkmID uuid.UUID) error
}

// Error handling: selalu wrap error dengan context
if err != nil {
    return fmt.Errorf("transaction_repo.Create: %w", err)
}

// Context: selalu teruskan context dari handler ke repo
func (r *transactionRepo) Create(ctx context.Context, tx *model.Transaction) error {
    // ...
}

// Logging: gunakan slog, bukan fmt
slog.Info("transaksi berhasil dibuat", "id", tx.ID, "umkm_id", tx.UmkmID)
slog.Error("gagal membuat transaksi", "error", err)
```

### Vue

```vue
<!-- Composition API dengan <script setup> — WAJIB -->
<script setup>
import { ref, computed, onMounted } from "vue";
import { useDashboardStore } from "@/stores/dashboard";

const store = useDashboardStore();
const isLoading = ref(false);

onMounted(async () => {
  isLoading.value = true;
  try {
    await store.fetchSummary();
  } finally {
    isLoading.value = false;
  }
});
</script>

<!-- Semua teks UI dalam Bahasa Indonesia -->
<!-- Gunakan PrimeVue components, bukan buat dari scratch -->
<!-- Format angka selalu Rupiah: Rp1.500.000 (titik sebagai separator ribuan) -->
```

### Format Rupiah (JavaScript)

```javascript
// Gunakan fungsi ini di seluruh frontend, simpan di src/utils/format.js
export const formatRupiah = (amount) => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(amount);
};
// Output: "Rp1.500.000"
```

---

## 12. DOCKER COMPOSE (development only)

```yaml
# docker-compose.yml
version: "3.8"
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: kasiraiai
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 13. MAKEFILE (backend)

```makefile
.PHONY: run migrate-up migrate-down build test

run:
	go run ./cmd/server/main.go

migrate-up:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down 1

build:
	go build -o bin/kasiraiai ./cmd/server/main.go

test:
	go test ./... -v
```

---

## 14. URUTAN IMPLEMENTASI

Kerjakan dalam urutan ini. Jangan loncat ke step berikutnya sebelum step sekarang selesai dan bisa dijalankan.

```
Step 1: Setup database
  → docker-compose up
  → buat struktur direktori backend
  → buat go.mod, config/config.go, internal/repository/db.go
  → jalankan migration 000001 dan 000002
  → VERIFIKASI: bisa konek ke DB, tabel ada

Step 2: Auth backend
  → internal/model/umkm.go, user.go
  → internal/repository/umkm_repo.go (Create, FindByPhone)
  → internal/service/auth_service.go (Register, Login)
  → internal/handler/auth_handler.go
  → internal/middleware/auth.go
  → cmd/server/main.go (setup router)
  → VERIFIKASI: POST /auth/register dan /auth/login berjalan

Step 3: WhatsApp webhook + AI parser
  → pkg/ai/client.go, pkg/ai/parser.go
  → pkg/fonnte/client.go
  → internal/model/transaction.go
  → internal/repository/transaction_repo.go (Create)
  → internal/service/transaction_service.go (ParseAndSave)
  → internal/handler/webhook_handler.go
  → VERIFIKASI: kirim POST ke /webhook/whatsapp dengan pesan simulasi

Step 4: Transaction API
  → internal/repository/transaction_repo.go (FindByUmkmID, SumByPeriod, SoftDelete)
  → internal/service/dashboard_service.go
  → internal/handler/transaction_handler.go
  → internal/handler/dashboard_handler.go
  → VERIFIKASI: GET /api/v1/dashboard/summary return data benar

Step 5: KUR Score
  → internal/model/kur_score.go
  → internal/repository/transaction_repo.go (tambah query untuk KUR)
  → internal/service/kur_service.go
  → internal/handler/webhook_handler.go (tambah command "skor")
  → VERIFIKASI: POST /api/v1/kur/recalculate return skor

Step 6: Frontend setup
  → npm create vite frontend --template vue
  → install semua dependencies di Seksi 6
  → setup router, pinia stores, axios instance
  → VERIFIKASI: npm run dev jalan di port 5173

Step 7: Frontend views
  → LoginView, RegisterView
  → DashboardView + components
  → TransactionView
  → KurScoreView
  → ReportView
  → VERIFIKASI: semua halaman bisa dibuka dan data tampil

Step 8: Laporan bulanan
  → internal/service/report_service.go
  → internal/handler/webhook_handler.go (command "laporan")
  → POST /api/v1/reports/monthly
  → VERIFIKASI: kirim "laporan" via WhatsApp, terima ringkasan

Step 9: Deployment
  → Render Free Tier untuk backend (set env vars)
  → Vercel untuk frontend (set VITE_API_BASE_URL)
  → VERIFIKASI: semua endpoint accessible dari internet
```

---

## 15. COMMAND WHATSAPP YANG DIKENALI

Selain pesan transaksi bebas, sistem harus mengenali command ini:

| Pesan User                      | Aksi                                    |
| ------------------------------- | --------------------------------------- |
| `laporan` / `laporan bulan ini` | Kirim ringkasan keuangan bulan berjalan |
| `skor` / `kur saya`             | Kirim KUR readiness score               |
| `hari ini`                      | Ringkasan transaksi hari ini            |
| `minggu ini`                    | Ringkasan transaksi minggu ini          |
| `hapus terakhir`                | Soft delete transaksi terakhir          |
| `bantuan` / `help`              | Kirim daftar command                    |

Command detection dilakukan di `service.TransactionService.ParseAndSave()` sebelum memanggil AI. Jika pesan cocok dengan command, langsung handle tanpa panggil DeepSeek API (hemat biaya).

---

## 16. PESAN WHATSAPP STANDAR

Simpan semua template pesan di `pkg/fonnte/messages.go` sebagai konstanta.

```go
const (
    MsgWelcome = `Selamat datang di KasirAI! 👋

Saya adalah asisten keuangan UMKM Anda.

Cara pakai:
📥 Catat pemasukan: "jual nasi goreng 3 porsi 45rb"
📤 Catat pengeluaran: "beli bahan 150rb"
📊 Lihat ringkasan: ketik *laporan*
⭐ Cek skor KUR: ketik *skor*
❓ Bantuan: ketik *bantuan*`

    MsgHelp = `Perintah yang tersedia:
• *laporan* - Ringkasan keuangan bulan ini
• *hari ini* - Transaksi hari ini
• *minggu ini* - Transaksi minggu ini
• *skor* - Cek kesiapan KUR Anda
• *hapus terakhir* - Hapus transaksi terakhir

Untuk mencatat transaksi, cukup tulis seperti biasa:
"jual 5 kaos 250rb" atau "bayar listrik 200rb"`
)
```

---

## 17. KEAMANAN

- Password: bcrypt dengan cost 12 (`golang.org/x/crypto/bcrypt`)
- JWT: HS256, expiry 72 jam, simpan `umkm_id` dan `phone_number` di claims
- Webhook Fonnte: validasi header atau token dari Fonnte jika tersedia
- Input sanitasi: semua input dari user di-validate sebelum masuk ke query
- Query: SELALU gunakan parameterized query, TIDAK PERNAH string concatenation untuk query
- Rate limiting: TIDAK diimplementasi di MVP (tambahkan jika diminta)
- CORS: izinkan origin frontend saja, bukan `*` di production

---

## 18. ERROR MESSAGES (Bahasa Indonesia)

```go
// Simpan di internal/handler/response.go
const (
    ErrInvalidInput       = "Data yang dikirim tidak valid"
    ErrPhoneAlreadyExists = "Nomor WhatsApp sudah terdaftar"
    ErrInvalidCredential  = "Nomor WhatsApp atau password salah"
    ErrUnauthorized       = "Anda belum login atau sesi telah habis"
    ErrNotFound           = "Data tidak ditemukan"
    ErrInternalServer     = "Terjadi kesalahan sistem, coba lagi nanti"
    ErrAIParseFailure     = "Maaf, saya tidak bisa memahami transaksi tersebut. Coba tulis lebih jelas, contoh: 'jual nasi 30rb'"
)
```

---

_Versi AGENTS.md ini: 1.0.0 — KasirAI MVP_  
_Diperbarui: 2026-06-09_  
_Jangan modifikasi file ini kecuali ada perubahan arsitektur yang disepakati._

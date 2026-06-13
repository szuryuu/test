# KasirAI

> Catat keuangan UMKM via WhatsApp, pantau lewat dashboard, siap KUR.

AI-powered financial assistant for Indonesian MSMEs. Record transactions by texting WhatsApp in plain Indonesian — KasirAI parses your message, saves the record, and replies with a confirmation. A web dashboard shows income, expenses, profit trends, and a KUR (micro-credit) readiness score.

Built for the [IDCamp Developer Challenge #2](https://www.dicoding.com/challenges/639) — _Digitalization & Acceleration of MSMEs with Generative AI_.

---

## Features

- **WhatsApp Transaction Logging**  
  Send a WhatsApp message like _"jual nasi uduk 5 bungkus 75rb"_ — KasirAI parses it via DeepSeek, saves the structured transaction, and replies with a confirmation. No app install, no login.

- **AI-Powered Natural Language Parsing**  
  DeepSeek Chat extracts transaction type (income/expense), amount, date, and category from freeform Indonesian text. Handles slang amounts (_"5rb"_, _"2jt"_, _"500k"_) and ambiguous messages with low-confidence clarification prompts.

- **Financial Dashboard**  
  Vue 3 dashboard with income/expense charts (Chart.js), summary KPI cards, recent transactions table, and category breakdown. Filter by daily, weekly, monthly, or yearly periods.

- **KUR Readiness Score**  
  A composite 0–100 score computed from four weighted factors: income stability (30%), profit margin (25%), recording consistency (25%), and data history length (20%). Tells the MSME owner whether they're ready to apply for KUR micro-credit — directly from their bank's checklist.

- **Monthly Reports via WhatsApp**  
  Send the message _"laporan"_ to receive an instant monthly financial summary: total income, expenses, net profit, and top categories — all delivered back to the same WhatsApp chat.

- **Simple Web Registration**  
  MSME owners sign up once through a clean registration page. They get a phone-number-linked account and can immediately start texting the WhatsApp bot.

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      FRONTEND                           │
│        Vue 3 + Vite   (port 5173 dev / Vercel prod)    │
│        Axios → REST API                                 │
└─────────────────────────┬───────────────────────────────┘
                          │ HTTP/JSON
┌─────────────────────────▼───────────────────────────────┐
│                      BACKEND                            │
│        Go + Gin   (port 8080 dev / Render prod)         │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌───────┐ │
│  │ handler/ │→ │ service/ │→ │repository/│→ │  PG   │ │
│  └──────────┘  └────┬─────┘  └───────────┘  └───────┘ │
│                     │                                   │
│              ┌──────▼──────┐  ┌──────────────┐         │
│              │  pkg/ai/    │  │  pkg/fonnte/ │         │
│              │  (DeepSeek) │  │  (WhatsApp)  │         │
│              └─────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────┘
                          │
       ┌──────────────────▼──────────────────┐
       │         EXTERNAL SERVICES           │
       │  DeepSeek API  |  Fonnte API        │
       └─────────────────────────────────────┘
```

### System Components

| Layer               | Technology                            | Role                                                                 |
| :------------------ | :------------------------------------ | :------------------------------------------------------------------- |
| **Frontend**        | Vue 3, Vite, PrimeVue, Chart.js       | SPA dashboard with KPI cards, charts, transaction table, KUR gauge   |
| **Backend**         | Go 1.26, Gin, pgx/v5                  | REST API, AI orchestration, business logic, WhatsApp webhook ingress |
| **Database**        | PostgreSQL 17                         | Stores UMKM profiles, transactions (amount in Rupiah), KUR history   |
| **AI Parsing**      | DeepSeek Chat API                     | Natural-language → structured transaction JSON from WhatsApp messages |
| **WhatsApp Gateway**| Fonnte API                            | Receives webhook, sends reply messages back to MSME owners           |

---

## How It Works

### WhatsApp Transaction Flow

1. **MSME sends a message**  
   User texts the bot: _"jual nasi uduk 5 bungkus @15rb"_.

2. **Fonnte forwards the message**  
   Fonnte POSTs `multipart/form-data` (sender, message, device) to `POST /webhook/whatsapp`.

3. **Command detection (no-AI fast path)**  
   If the message matches a known command (_"laporan"_, _"skor"_, _"bantuan"_), it is handled locally without calling DeepSeek — saving cost and latency.

4. **AI parsing**  
   For unrecognized messages, the DeepSeek prompt template extracts `type`, `amount`, `description`, `category_hint`, `transaction_date`, and `confidence`. The response is validated and unmarshalled.

5. **Transaction saved**  
   A new row is inserted into `transactions` with `source = 'whatsapp'` and the AI confidence score.

6. **Reply sent**  
   Fonnte sends a confirmation WhatsApp message back to the user, including running totals when available.

### Dashboard Request Flow

1. **Browser loads the SPA**  
   User logs in with phone number + password, receives a JWT (72h expiry).

2. **API call**  
   `GET /api/v1/dashboard/summary?period=monthly` with `Authorization: Bearer <token>`.

3. **Auth middleware** validates the JWT and injects `umkm_id` into the Gin context.

4. **Dashboard service** queries the repository for aggregated sums grouped by period.

5. **JSON response** returns `total_income`, `total_expense`, `net_profit`, `profit_margin`, `transaction_count`, and `chart_data` array for the chart component.

---

## Tech Stack

| Layer        | Technology                                                              |
| :----------- | :---------------------------------------------------------------------- |
| **Language** | Go 1.26 (backend), JavaScript (frontend)                                |
| **HTTP**     | Gin Web Framework                                                       |
| **Database** | PostgreSQL 17 + pgx/v5 (connection pooling, parameterized queries)      |
| **Migrations**| golang-migrate/migrate                                                 |
| **Auth**     | JWT (HS256, golang-jwt/jwt/v5), bcrypt password hashing (cost 12)      |
| **AI**       | DeepSeek Chat API (OpenAI-compatible `/v1/chat/completions`)            |
| **WhatsApp** | Fonnte API (webhook ingress, message sending)                           |
| **Frontend** | Vue 3 (Composition API), Vite 6, Pinia, Vue Router 4, Axios            |
| **UI Kit**   | PrimeVue 4 + PrimeFlex + Tailwind CSS 4                                 |
| **Charts**   | Chart.js via vue-chartjs                                                |
| **DevOps**   | Docker Compose (dev), Dockerfiles (backend/frontend), Render + Vercel   |

---

## Getting Started

### Prerequisites

- Go 1.26+
- Node.js 22+
- PostgreSQL 17 (_or Docker for the included `docker-compose.yml`_)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) (`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`)
- DeepSeek API key ([platform.deepseek.com](https://platform.deepseek.com))
- Fonnte API key and device ID ([fonnte.com](https://fonnte.com))

### 1. Clone & Environment

```bash
git clone <repo-url> kasiraiai
cd kasiraiai

# Copy and edit environment files
cp .env.example .env
cp frontend/.env.example frontend/.env
```

Fill in the required secrets in `.env`:

```env
JWT_SECRET=your_secret_at_least_32_chars
DEEPSEEK_API_KEY=sk-...
FONNTE_API_KEY=your_fonnte_key
FONNTE_DEVICE=your_fonnte_device_id
DB_PASSWORD=your_db_password
```

### 2. Database

```bash
# Option A: Docker Compose (recommended)
docker compose up -d postgres

# Option B: Local PostgreSQL — create the database manually
createdb kasiraiai
```

Run migrations:

```bash
cd backend
make migrate-up
```

### 3. Backend

```bash
cd backend
go mod download
make run
```

Server starts on `http://localhost:8080`.  
Health check: `curl http://localhost:8080/health`.

### 4. Frontend

```bash
cd frontend
npm install
npm run dev
```

Dev server starts on `http://localhost:5173`. API calls are proxied to the backend via Vite's dev proxy.

### 5. Register & Test

1. Open `http://localhost:5173/register` and create an account.
2. Simulate a WhatsApp message (no Fonnte key yet? POST directly to webhook):
   ```bash
   curl -X POST http://localhost:8080/webhook/whatsapp \
     -F "sender=6281234567890" \
     -F "message=jual nasi goreng 5 porsi 75rb" \
     -F "device=test"
   ```

### All-in-One (Docker)

```bash
docker compose up -d
```

Starts PostgreSQL, backend (port 8080), and frontend (port 80) with health-checked startup ordering.

---

## Project Structure

```
kasiraiai/
├── backend/
│   ├── cmd/server/main.go              # Entry point — wire dependencies, start Gin
│   ├── config/config.go                # Env var loading into typed Config struct
│   ├── internal/
│   │   ├── handler/                    # HTTP handlers (auth, webhook, dashboard, KUR, reports, transactions)
│   │   ├── middleware/                 # Auth (JWT), Logger (slog), CORS
│   │   ├── model/                      # Domain structs (umkm, transaction, kur_score)
│   │   ├── repository/                 # pgx queries, DB init, migration runner
│   │   └── service/                    # Business logic (auth, transaction parsing, dashboard, KUR, reports)
│   ├── pkg/
│   │   ├── ai/                         # DeepSeek client + parser (prompt template, JSON cleanup)
│   │   └── fonnte/                     # Fonnte client (SendMessage) + message templates
│   ├── migrations/                     # SQL migration files (000001_init_schema, 000002_seed_categories)
│   ├── Dockerfile
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── api/                        # Axios instance + API modules (auth, dashboard, transactions, KUR)
│   │   ├── components/                 # Reusable Vue components (layout, dashboard, transaction, KUR)
│   │   ├── router/index.js             # Vue Router (auth guard, lazy-loaded views)
│   │   ├── stores/                     # Pinia stores (auth, dashboard, transaction)
│   │   ├── utils/format.js             # Rupiah formatting (Intl.NumberFormat id-ID)
│   │   └── views/                      # Page-level components (Login, Register, Dashboard, Transactions, KUR, Reports)
│   ├── index.html
│   ├── vite.config.js                  # Vite config (Vue plugin, Tailwind CSS plugin, dev proxy)
│   ├── nginx.conf                      # Production nginx config for SPA routing
│   └── Dockerfile
├── docker-compose.yml                  # PostgreSQL 17 + backend + frontend (dev setup)
├── .env.example                        # Environment variable template
└── AGENTS.md                           # Project specification & coding standards (authoritative)
```

---

## API Overview

Base URL: `/api/v1`  
Auth: `Authorization: Bearer <JWT>` (except `/auth/*` and `/webhook/*`)

### Public Endpoints

| Method | Path                  | Description                                   |
| :----- | :-------------------- | :-------------------------------------------- |
| POST   | `/api/v1/auth/register` | Register a new MSME account                  |
| POST   | `/api/v1/auth/login`    | Login with phone number + password, get JWT  |
| POST   | `/webhook/whatsapp`     | Fonnte webhook — receives WhatsApp messages  |
| GET    | `/health`               | Health check (returns `{"status":"ok"}`)     |

### Protected Endpoints

| Method | Path                            | Description                                    |
| :----- | :------------------------------ | :--------------------------------------------- |
| GET    | `/api/v1/dashboard/summary`     | Income/expense/profit summary (filters: period, date) |
| GET    | `/api/v1/dashboard/categories`  | Spending/income breakdown by category           |
| GET    | `/api/v1/transactions`          | List transactions (pagination, date range, type filter) |
| POST   | `/api/v1/transactions`          | Manually create a transaction                   |
| DELETE | `/api/v1/transactions/:id`      | Soft-delete a transaction                       |
| GET    | `/api/v1/kur/score`             | Get latest KUR readiness score + recommendations |
| POST   | `/api/v1/kur/recalculate`       | Force-recalculate KUR score from transaction history |
| POST   | `/api/v1/reports/monthly`       | Generate monthly report, send via WhatsApp      |
| GET    | `/api/v1/umkm/profile`          | Get authenticated user's profile                |

---

## WhatsApp Commands

| Message         | Action                                          |
| :-------------- | :---------------------------------------------- |
| `bantuan`       | Show command list and usage help                |
| `laporan`       | Send monthly financial summary                  |
| `hari ini`      | Today's transaction summary                     |
| `minggu ini`    | This week's transaction summary                 |
| `skor`          | Show KUR readiness score and recommendations    |
| `hapus terakhir`| Soft-delete the last recorded transaction       |

Any other message is treated as a transaction attempt and parsed by the AI.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

# ISSUES — KasirAI

## How to Use This File

Each issue has an ID, severity, and exact location. Fix in priority order.
Mark done by changing `[ ]` to `[x]` and adding the fix summary.

---

## 🔍 QA Verification — 2026-06-10

### ISSUE-008: Transaction Filter Disappears on Empty State

**Status:** ❌ NOT CONFIRMED — no bug exists (User Verification)

**Investigation:** Inspected `TransactionTable.vue` template structure. The filter controls
(type select, start_date, end_date, Reset button) are rendered unconditionally at the top
of the component (lines 65–97). No `v-if`, `v-show`, or conditional rendering wraps the
filter bar — it is always visible regardless of `store.transactions.length` or `store.isLoading`.

The `<div v-else>` block containing the table includes a proper empty-state row that displays
"Coba ubah filter pencarian" when filters are active and results are zero, and "Tambahkan
transaksi baru untuk memulai" when no filters are set. Both cases render correctly inside
the table container alongside the always-visible filter bar.

**Verdict:** Filter visibility is decoupled from transaction data. No fix needed.

---

### ISSUE-009: Navbar Rubber-Banding / Overscroll

**Status:** ✅ CONFIRMED

**Files:** `frontend/src/components/layout/AppLayout.vue`, `frontend/src/assets/styles/main.css`

**Root Cause:** The `<main>` scrollable container in `AppLayout.vue` uses
`overflow-y-auto` but has no `overscroll-behavior: none` (or Tailwind's `overscroll-none`).
Neither `body` nor `html` in `main.css` declare overscroll-behavior.

The topbar (`AppTopbar.vue`) uses `sticky top-0` — when the browser rubber-bands during
overscroll at the top or bottom of the page, the sticky navbar moves with the container
because there is no overscroll containment.

**Grep result:** Zero matches for `overscroll` across the entire `src/` directory.

**Fix:** Add `overscroll-none` to the `<main>` element in `AppLayout.vue`, or add a global
`html, body { overscroll-behavior: none; }` rule in `main.css` to prevent rubber-banding
on all scrollable surfaces.

- [x] Fixed — added `overscroll-none` class to `<main>` in `AppLayout.vue` line 49.

---

## 🔴 FATAL — Security & Stability Audit (2026-06-10)

### ISSUE-010: `JSON.parse` Crash on Corrupted localStorage

**Status:** ✅ CONFIRMED — FATAL

**File:** `frontend/src/stores/auth.js`, line 7

**Root Cause:** Store initialization calls `JSON.parse(localStorage.getItem("umkm") || "null")`
without try/catch. The `|| "null"` fallback only activates when the key does NOT exist
(`localStorage.getItem` returns `null`). If the key exists but contains invalid JSON
(e.g., corrupted by a browser extension, manual edit, or disk error), `JSON.parse` throws
an unhandled `SyntaxError` synchronously during Pinia store setup. This crashes the entire
application before any route or component renders — white screen, no recovery.

**Why it's fatal:** Pinia stores initialize lazily on first `useAuthStore()` call, which
happens early in the app lifecycle (router guard, layout components). A single corrupted
`localStorage.umkm` entry renders the app permanently dead until the user manually clears
localStorage.

**Fix:** Wrap in try/catch with graceful fallback:

```js
let parsed = null;
try {
  parsed = JSON.parse(localStorage.getItem("umkm") || "null");
} catch {
  localStorage.removeItem("umkm"); // clear corrupt entry
  parsed = null;
}
const umkm = ref(parsed);
```

- [x] Fixed — wrapped `JSON.parse` in try/catch at `stores/auth.js`: catches `SyntaxError`, removes corrupt `localStorage.umkm`, falls back to `null`.

### ISSUE-011: No Indonesian Phone Number Format Validation

**Status:** ✅ CONFIRMED — VALIDATION MISSING

**Scope:** Full-stack gap — backend, frontend, and database all accept arbitrary
phone number strings.

**Files & evidence:**

- `backend/internal/service/auth_service.go`, line 31:

  ```go
  PhoneNumber  string `json:"phone_number" binding:"required"`
  ```

  Only `required` tag — no `len`, `min`, `max`, `startswith`, or custom validator.
  The AGENTS.md schema mandates format `628XXXXXXXXXX` (11–13 digits starting
  with 628) but this is never enforced in code.

- `backend/internal/handler/auth_handler.go`, line 47:
  Login phone field also only `binding:"required"` — no format enforcement.

- `backend/internal/model/umkm.go`, line 13:
  No struct tag for phone format validation.

- `backend/migrations/000001_init_schema.up.sql`, line 10:

  ```sql
  phone_number  VARCHAR(20)  NOT NULL UNIQUE, -- format: 628XXXXXXXXXX
  ```

  Only a `VARCHAR(20)` length cap and `UNIQUE` constraint. No `CHECK` constraint
  enforcing the `628` prefix or digit-only pattern. A value like `"hello"` passes
  the DB validation.

- `frontend/src/views/LoginView.vue`, line 80 and
  `frontend/src/views/RegisterView.vue`, line 115:
  Phone input uses `type="text"` with no `pattern`, `inputmode="numeric"`,
  `maxlength`, or JS-side format validation. The placeholder `"6281234567890"`
  is purely cosmetic.

- **Grep result:** Zero matches for `regexp`, `MustCompile`, `MatchString`,
  `pattern`, or `validator` across the entire backend. Gin's built-in
  `go-playground/validator` is available via `binding` tags but never used
  for phone format validation.

**Root Cause:** The Indonesian phone format `628XXXXXXXXXX` is documented in
AGENTS.md and the SQL migration comment but never enforced in any software layer.
The project has a structural assumption that Fonnte will always provide valid
`628...` sender IDs, but:

1. Users can self-register with arbitrary strings (e.g., `"081234"`, `"test"`).
2. The webhook handler (`webhook_handler.go`, line 32) reads `sender` from
   Fonnte's form data and passes it directly to `FindByPhone` — if Fonnte ever
   sends a differently-formatted number, the lookup silently returns `nil`.
3. Corrupt phone numbers in the `umkms` table prevent WhatsApp webhook delivery
   because `FindByPhone` does exact string matching.

**Impact:**

- Invalid phone numbers break the WhatsApp webhook loop — registered users with
  malformed numbers never receive replies.
- No defense against format drift (e.g., `08xx` vs `628xx`).
- DB uniqueness is the only guard, but uniqueness on arbitrary strings doesn't
  guarantee format correctness.

- [x] Fixed — backend: added `binding:"required,startswith=628,numeric,min=11,max=13"` to RegisterRequest and login PhoneNumber; frontend: added `type="tel" inputmode="numeric" maxlength="13" pattern="628[0-9]{8,10}"` to phone inputs in LoginView and RegisterView.: Expense Category Amounts Rendered as Positive in Report

**Status:** ✅ CONFIRMED — DESIGN BUG

**Observed Output:**

```
Rincian Kategori:
Tanpa Kategori: Rp443 (100%)
```

An expense of Rp443 is displayed without a negative sign, visually
indistinguishable from an income entry (aside from the 📤 emoji prefix).

**Files:**

- `backend/internal/repository/transaction_repo.go`, lines 219–234 —
  `CategorySummary` SQL query:

  ```sql
  WITH category_totals AS (
      SELECT
          COALESCE(tc.name, 'Tanpa Kategori') AS name,
          SUM(t.amount) AS total      -- ← ALWAYS POSITIVE
      FROM transactions t
      LEFT JOIN transaction_categories tc ON t.category_id = tc.id
      WHERE t.umkm_id = $1
        AND t.deleted_at IS NULL
        AND t.type = $2              -- discriminator used for filtering
        AND t.transaction_date >= $3::date
        AND t.transaction_date <  $4::date
      GROUP BY tc.name
  )
  ```

  The `t.type = $2` filter separates income from expense rows (the query is
  called once for `"income"` and once for `"expense"`), but `SUM(t.amount)`
  always yields a positive integer because the `amount` column has a schema
  constraint `CHECK (amount > 0)`. Expense values are not negated in the query.

- `backend/internal/service/report_service.go`, lines 79–84 — expense category
  formatting:

  ```go
  for _, c := range expCats {
      topCategories = append(topCategories,
          fmt.Sprintf("📤 %s: Rp%s (%.0f%%)",
              c.Name, formatRupiahReport(c.Total), c.Percentage))
  }
  ```

  `c.Total` is always positive (inherited from the SQL query).
  `formatRupiahReport` (line 157) does handle negative `int64` values — its
  first branch prepends `"-"` — but it never receives a negative value for
  expense categories.

- `backend/internal/service/report_service.go`, line 157 — `formatRupiahReport`:
  ```go
  func formatRupiahReport(amount int64) string {
      if amount < 0 {
          return "-" + formatRupiahReport(-amount)  // dead code for categories
      }
      ...
  }
  ```

**Root Cause:** The database schema stores all amounts as positive integers
(distinguishing income vs. expense via the `type` column, not via sign). The
`CategorySummary` query preserves this positivity for both income and expense
categories. The report service's category loop prints the raw positive value
without context-aware negation. The 📤 emoji is the only structural distinction,
and it may be stripped or rendered invisibly on some WhatsApp clients.

**Design Note:** The top-level summary section (`Pemasukan` / `Pengeluaran` /
`Laba Bersih`) correctly shows the expense context because `NetProfit` is
computed as `TotalIncome - TotalExpense` in `SumByPeriod`. The category
breakdown has no equivalent adjustment.

- [x] Fixed — prefixed `-` before `Rp` in expense category format string at `report_service.go:72`, so "📤 Tanpa Kategori: -Rp443 (100%)" now clearly distinguishes expenses from income.

---

## 🐛 UI/UX & Validation — Register Page (2026-06-11)

### ISSUE-013: Triple Border on Focused Inputs

**Status:** ✅ CONFIRMED — CSS CONFLICT

**Files:**

- `frontend/src/views/RegisterView.vue`, every input element (lines 89–137)
- `frontend/src/assets/styles/main.css`, lines 98–100

**Root Cause:** Three concentric visual border layers stack on input focus:

| Layer | Source | Visual |
|---|---|---|
| 1 | CSS class `border border-[var(--color-border)]` | 1px border (gray → green via inline style) |
| 2 | Inline `@focus` → `boxShadow = '0 0 0 3px rgba(16,185,129,0.1)'` | 3px solid-looking ring (no blur = reads as border) |
| 3 | Global `:focus-visible { outline: 2px solid var(--color-brand-500); outline-offset: 2px }` in `main.css:98-100` | 2px outline, 2px outside element |

The global `:focus-visible` outline was added as a catch-all accessibility rule
but collides with inputs that already have their own focus ring via `boxShadow`.
The component intends to use box-shadow as its focus indicator (matching the
login page pattern), but the global outline adds a second competing ring.

**Also affects:** `LoginView.vue` inputs (same inline `@focus` handlers +
global `:focus-visible`). The select dropdown at RegisterView line 148 has the
same pattern but without `boxShadow` — only border-color change + outline.

- [x] Fixed — added `focus-visible:outline-none` to all 7 form inputs/selects across `RegisterView.vue` (5) and `LoginView.vue` (2), overriding the global `:focus-visible` outline rule that clashed with per-element box-shadow focus rings.: No Inline Per-Field Validation on Register Form

**Status:** ✅ CONFIRMED — MISSING FEATURE

**File:** `frontend/src/views/RegisterView.vue`, `handleRegister()` lines 30–48

**Root Cause:** Client-side validation runs only on submit — no real-time
per-field feedback as user types.

Current flow:
```js
async function handleRegister() {
  // Only fires on button click
  if (!form.value.name || ...) {
    errorMsg.value = "Semua field wajib diisi kecuali jenis usaha";
    return;
  }
  if (form.value.password.length < 6) {
    errorMsg.value = "Password minimal 6 karakter";
    return;
  }
  // ... API call
}
```

What's missing:
- Phone field has `pattern="628[0-9]{8,10}"` but no JS-side regex check or
  inline error message when user types invalid format.
- No per-field `error` state/tracking — only one global `errorMsg` ref.
- No `@input` / `@blur` validators; no visual feedback (red border, error text
  below field) while user fills the form.
- HTML5 constraint validation never fires because the form uses `@click` on a
  `<button>` not a native `<form @submit>`, and `type="tel"` doesn't trigger
  `pattern` validation in all browsers.

**Impact:** User only discovers validation errors after clicking "Daftar".
If backend rejects with "Data yang dikirim tidak valid", user has no clue
which field is wrong.

- [x] Fixed — added `errors` reactive state, per-field `validate*()` functions triggered on `@blur`, red border + inline error text below each input, and `getInputClass`/`getInputStyle`/`onFocus`/`onBlur` helpers that respect error state. Phone validates against `^628[0-9]{8,10}$` regex.: Backend Returns Generic Error — Discards Per-Field Validation Details

**Status:** ✅ CONFIRMED — ERROR SUPPRESSION

**Files:**

- `backend/internal/handler/auth_handler.go`, lines 22–25
- `backend/internal/handler/response.go`, line 38
- `backend/internal/service/auth_service.go`, lines 29–34 (RegisterRequest binding tags)

**Root Cause:** Gin's `ShouldBindJSON` returns detailed validation errors
(per-field, with tag info like `"PhoneNumber" failed on 'startswith'`),
but the handler discards them:

```go
// auth_handler.go:22-25
var req service.RegisterRequest
if err := c.ShouldBindJSON(&req); err != nil {
    BadRequest(c, ErrInvalidInput)  // ← err thrown away, only generic message
    return
}
```

`ErrInvalidInput` = `"Data yang dikirim tidak valid"` — no field names,
no reason (e.g., "Nomor WhatsApp harus diawali 628").

The `RegisterRequest` struct already has precise binding rules:
```go
PhoneNumber string `json:"phone_number" binding:"required,startswith=628,numeric,min=11,max=13"`
Password    string `json:"password" binding:"required,min=6"`
```

But the API response is always:
```json
{ "success": false, "message": "Data yang dikirim tidak valid" }
```

The `APIResponse` struct has an `Errors []string` field (response.go:14) that
could carry per-field messages, but it's never populated from Gin's validator
output.

**Impact:** Frontend receives zero field-specific guidance. Combined with
ISSUE-014 (no inline validation), user gets an opaque error with no hint what
to fix.

- [x] Fixed — added `validationErrors()` helper in `response.go` that maps Gin validator tags to Indonesian messages per-field (e.g., `"Nomor WhatsApp harus diawali 628"`), with `fieldLabels` and `tagMessages` lookup tables. Both `Register` and `Login` handlers now return field-specific errors via `ErrorResponse(c, ..., fieldErrors...)`.

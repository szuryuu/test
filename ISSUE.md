# ISSUES — KasirAI

## How to Use This File

Each issue has an ID, severity, and exact location. Fix in priority order.
Mark done by changing `[ ]` to `[x]` and adding the fix summary.

---

## 🔴 CRITICAL — Fix Before Testing

### ISSUE-001: `rows.Err()` Not Checked After Row Iteration

**Files:** `backend/internal/repository/transaction_repo.go`
**Methods:** `FindByUmkmID`, `ChartData`, `CategorySummary`, `MonthlyIncomeValues`
**Risk:** Silent partial data return on network error mid-iteration. No error returned to caller.

**Pattern to find:**

```go
for rows.Next() { ... }
return results, nil  // ← rows.Err() missing here
```

**Fix pattern (add after every rows loop):**

```go
if err := rows.Err(); err != nil {
    return nil, 0, fmt.Errorf("transactionRepo.MethodName rows: %w", err)
}
```

Note: return signature differs per method — adjust accordingly.

- `FindByUmkmID` returns `([]model.Transaction, int, error)`
- `ChartData` returns `([]ChartDataPoint, error)`
- `CategorySummary` returns `([]CategorySummary, error)`
- `MonthlyIncomeValues` returns `([]int64, error)`

- [x] Fixed — already had rows.Err() in all 4 methods, no change needed.

---

### ISSUE-002: Double Report Generation in `report_handler.go`

**File:** `backend/internal/handler/report_handler.go`
**Risk:** Two separate DB queries, two category aggregations. Report sent to WhatsApp is recalculated independently from report returned to client — they may differ.

**Root cause:**

```go
report, err := h.svc.GenerateMonthly(...)       // generates once
h.svc.GenerateAndSend(...)                        // calls GenerateMonthly AGAIN internally
```

**Fix:** Add a `SendReport` method to `ReportService` that accepts pre-generated text:

```go
// internal/service/report_service.go — add to interface and struct
SendReport(ctx context.Context, phoneNumber, reportText string) error

func (s *reportService) SendReport(ctx context.Context, phoneNumber, reportText string) error {
    return s.fonnte.SendMessage(phoneNumber, reportText)
}
```

```go
// internal/handler/report_handler.go — use already-generated report
report, err := h.svc.GenerateMonthly(ctx, umkmID, req.Year, req.Month)
if err != nil { ... }

if umkm != nil {
    if sendErr := h.svc.SendReport(ctx, umkm.PhoneNumber, report.ReportText); sendErr != nil {
        slog.Error("gagal kirim laporan via WhatsApp", "error", sendErr)
    }
}
SuccessResponse(c, http.StatusOK, "Laporan berhasil dibuat dan dikirim via WhatsApp", report)
```

- [x] Fixed — added `SendReport` to `report_service.go`, replaced `GenerateAndSend` with `SendReport` in `report_handler.go`.

---

### ISSUE-003: Webhook Processes AI Synchronously — Fonnte Will Timeout and Retry

**File:** `backend/internal/handler/webhook_handler.go`
**Risk:** AI call takes 3-10 seconds. Fonnte default timeout is ~10 seconds. On slow AI response, Fonnte retries → duplicate transactions.

**Current flow (broken):**

```
Fonnte → POST /webhook → [waits 3-10s for AI] → return 200
```

**Required flow:**

```
Fonnte → POST /webhook → return 200 immediately → [AI processes in background]
```

**Fix:** Move ParseAndSave + reply to a goroutine. Return 200 to Fonnte immediately.

```go
func (h *WebhookHandler) Handle(c *gin.Context) {
    sender := c.PostForm("sender")
    message := c.PostForm("message")

    if sender == "" || message == "" {
        c.String(http.StatusOK, "OK")
        return
    }

    slog.Info("webhook diterima", "sender", sender, "message_len", len(message))

    umkm, err := h.umkmRepo.FindByPhone(c.Request.Context(), sender)
    if err != nil {
        slog.Error("gagal mencari UMKM", "sender", sender, "error", err)
        c.String(http.StatusOK, "OK")
        return
    }

    if umkm == nil {
        go func() {
            ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
            defer cancel()
            if sendErr := h.fonnte.SendMessage(sender, fonnte.MsgWelcome); sendErr != nil {
                slog.Error("gagal kirim welcome", "sender", sender, "error", sendErr)
            }
        }()
        c.String(http.StatusOK, "OK")
        return
    }

    // Capture values for goroutine — do NOT pass gin context
    umkmID := umkm.ID
    go func() {
        ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
        defer cancel()

        parsed, err := h.svc.ParseAndSave(ctx, umkmID, message, sender)
        if err != nil {
            slog.Error("gagal ParseAndSave", "umkm_id", umkmID, "error", err)
            _ = h.fonnte.SendMessage(sender, fonnte.MsgAIParseFailure)
            return
        }
        if parsed.ReplyMessage != "" {
            if sendErr := h.fonnte.SendMessage(sender, parsed.ReplyMessage); sendErr != nil {
                slog.Error("gagal kirim balasan", "sender", sender, "error", sendErr)
            }
        }
    }()

    // Return 200 to Fonnte immediately — before goroutine finishes
    c.String(http.StatusOK, "OK")
}
```

IMPORTANT: Never pass `c.Request.Context()` into the goroutine — it is cancelled after `c.String` returns.

- [x] Fixed — moved ParseAndSave + reply to goroutine with 45s timeout context, return 200 immediately. Also made welcome/error messages async.

---

## 🟡 MEDIUM — Fix Before Production Deploy

### ISSUE-004: Fonnte Client Uses `http.DefaultClient` (No Timeout)

**File:** `backend/pkg/fonnte/client.go`
**Risk:** `http.DefaultClient` has no timeout. If Fonnte API hangs, the goroutine leaks indefinitely.

**Fix:** Add `httpClient` field to `Client` struct with 15s timeout.

```go
type Client struct {
    cfg        *config.Config
    baseURL    string
    httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
    return &Client{
        cfg:     cfg,
        baseURL: cfg.FonnteBaseURL,
        httpClient: &http.Client{
            Timeout: 15 * time.Second,
        },
    }
}

// In SendMessage, replace:
resp, err := http.DefaultClient.Do(req)
// With:
resp, err := c.httpClient.Do(req)
```

- [x] Fixed — added `httpClient` field with 15s timeout, replaced `http.DefaultClient.Do` with `c.httpClient.Do`.

---

### ISSUE-005: Error String Comparison in `auth_handler.go`

**File:** `backend/internal/handler/auth_handler.go`
**Risk:** If the error message string changes in `auth_service.go`, the 409 branch silently becomes 500. Brittle.

**Fix:** Use a sentinel error variable.

```go
// internal/service/auth_service.go — add at package level
var ErrPhoneAlreadyExists = errors.New("nomor WhatsApp sudah terdaftar")

// In Register(), replace:
return nil, "", fmt.Errorf("nomor WhatsApp sudah terdaftar")
// With:
return nil, "", ErrPhoneAlreadyExists
```

```go
// internal/handler/auth_handler.go — replace string comparison
if errors.Is(err, service.ErrPhoneAlreadyExists) {
    ErrorResponse(c, http.StatusConflict, ErrPhoneAlreadyExists)
    return
}
```

- [x] Fixed — added `ErrPhoneAlreadyExists` sentinel to `auth_service.go`, used `errors.Is` in `auth_handler.go`.

---

## 🟢 LOW — Nice to Have Before Submission

### ISSUE-006: docker-compose Frontend Uses `service_started` Instead of `service_healthy`

**File:** `docker-compose.yml`
**Risk:** Frontend container starts proxying to backend before backend is ready to accept connections.

**Fix:**

```yaml
frontend:
  depends_on:
    backend:
      condition: service_healthy # was: service_started
```

Backend Dockerfile already has a `HEALTHCHECK` on `/health`. This change makes it effective.

- [x] Fixed — changed `service_started` to `service_healthy` in docker-compose.yml.

---

### ISSUE-007: KUR Score Cache Has No Staleness Indicator

**File:** `backend/internal/service/kur_service.go`
**Risk:** `GetScore` returns cached score indefinitely. After many new transactions, user sees stale score without knowing.

**Fix:** Add `is_stale` boolean to `KurScoreResult`.

```go
// internal/service/kur_service.go — add to struct
type KurScoreResult struct {
    // ... existing fields ...
    IsStale bool `json:"is_stale"`
}

// In GetScore(), after fetching cached:
isStale := time.Since(cached.CalculatedAt) > 24*time.Hour
result := kurScoreToResult(cached)
result.IsStale = isStale
return result, nil
```

```vue
<!-- frontend/src/views/KurScoreView.vue — add banner if stale -->
<div v-if="score?.is_stale" style="...warning banner...">
  Data mungkin sudah tidak akurat. Klik "Hitung Ulang" untuk memperbarui.
</div>
```

- [x] Fixed — added `IsStale` to `KurScoreResult` struct, set in `GetScore()` when cached >24h, added warning banner to `KurScoreView.vue`.

---

## Completion Checklist

```
[x] ISSUE-001 rows.Err() — 4 methods in transaction_repo.go (already fixed)
[x] ISSUE-002 double report generation — report_handler.go + report_service.go
[x] ISSUE-003 async webhook — webhook_handler.go
[x] ISSUE-004 fonnte timeout — pkg/fonnte/client.go
[x] ISSUE-005 sentinel error — auth_service.go + auth_handler.go
[x] ISSUE-006 docker-compose healthy — docker-compose.yml
[x] ISSUE-007 KUR stale flag — kur_service.go + KurScoreView.vue
```

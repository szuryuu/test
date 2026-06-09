# SKILL: PostgreSQL + pgx/v5 Best Practice — KasirAI

> Baca ini sebelum menulis query SQL atau kode database apapun.

## Driver: pgx/v5

- Package: `github.com/jackc/pgx/v5`
- Pool: `pgxpool.Pool` (bukan `sql.DB`)
- JANGAN pernah pakai `database/sql` langsung
- JANGAN pernah pakai GORM atau SQLX

---

## INISIALISASI POOL

```go
// repository/db.go
func NewDB(cfg *config.Config) *pgxpool.Pool {
    dsn := fmt.Sprintf(
        "host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
    )

    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        log.Fatalf("db config parse failed: %v", err)
    }

    config.MaxConns = int32(cfg.DBMaxConnections)       // 25
    config.MinConns = int32(cfg.DBMaxIdleConnections)   // 5
    config.MaxConnLifetime = 1 * time.Hour
    config.MaxConnIdleTime = 30 * time.Minute

    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("db connection failed: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := pool.Ping(ctx); err != nil {
        log.Fatalf("db ping failed: %v", err)
    }

    slog.Info("database connected", "host", cfg.DBHost, "db", cfg.DBName)
    return pool
}
```

---

## QUERY PATTERNS

### Single Row

```go
var tx model.Transaction
err := r.db.QueryRow(ctx, `
    SELECT id, umkm_id, amount, type, description, transaction_date, created_at
    FROM transactions
    WHERE id = $1 AND umkm_id = $2 AND deleted_at IS NULL
`, id, umkmID).Scan(
    &tx.ID, &tx.UmkmID, &tx.Amount, &tx.Type,
    &tx.Description, &tx.TransactionDate, &tx.CreatedAt,
)

if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrNotFound
    }
    return nil, fmt.Errorf("transactionRepo.FindByID: %w", err)
}
return &tx, nil
```

### Multiple Rows

```go
rows, err := r.db.Query(ctx, `
    SELECT id, umkm_id, amount, type, description, transaction_date, created_at
    FROM transactions
    WHERE umkm_id = $1 AND deleted_at IS NULL
    ORDER BY transaction_date DESC
    LIMIT $2 OFFSET $3
`, umkmID, limit, offset)

if err != nil {
    return nil, fmt.Errorf("transactionRepo.List: %w", err)
}
defer rows.Close() // WAJIB

var results []model.Transaction
for rows.Next() {
    var tx model.Transaction
    if err := rows.Scan(&tx.ID, &tx.UmkmID, &tx.Amount, &tx.Type, &tx.Description, &tx.TransactionDate, &tx.CreatedAt); err != nil {
        return nil, fmt.Errorf("transactionRepo.List scan: %w", err)
    }
    results = append(results, tx)
}

// WAJIB cek error setelah iterasi
if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("transactionRepo.List rows.Err: %w", err)
}

return results, nil
```

### Insert

```go
_, err := r.db.Exec(ctx, `
    INSERT INTO transactions (id, umkm_id, category_id, amount, type, description, raw_message, transaction_date, source, ai_confidence)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`, tx.ID, tx.UmkmID, tx.CategoryID, tx.Amount, tx.Type, tx.Description, tx.RawMessage, tx.TransactionDate, tx.Source, tx.AIConfidence)

if err != nil {
    // Cek duplicate key
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) && pgErr.Code == "23505" {
        return ErrDuplicate
    }
    return fmt.Errorf("transactionRepo.Create: %w", err)
}
```

### Soft Delete

```go
result, err := r.db.Exec(ctx, `
    UPDATE transactions
    SET deleted_at = NOW(), updated_at = NOW()
    WHERE id = $1 AND umkm_id = $2 AND deleted_at IS NULL
`, id, umkmID)

if err != nil {
    return fmt.Errorf("transactionRepo.SoftDelete: %w", err)
}

// Cek apakah row benar-benar di-update
if result.RowsAffected() == 0 {
    return ErrNotFound
}
return nil
```

### Aggregate Query (untuk Dashboard)

```go
type SummaryResult struct {
    TotalIncome   int64   `db:"total_income"`
    TotalExpense  int64   `db:"total_expense"`
    TxCount       int     `db:"tx_count"`
}

var s SummaryResult
err := r.db.QueryRow(ctx, `
    SELECT
        COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0)  AS total_income,
        COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense,
        COUNT(*) AS tx_count
    FROM transactions
    WHERE umkm_id = $1
      AND transaction_date BETWEEN $2 AND $3
      AND deleted_at IS NULL
`, umkmID, startDate, endDate).Scan(&s.TotalIncome, &s.TotalExpense, &s.TxCount)
```

---

## MIGRATION — GOLANG-MIGRATE

```go
// repository/db.go — tambahkan fungsi ini
func RunMigrations(cfg *config.Config) {
    dbURL := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
    )

    m, err := migrate.New("file://migrations", dbURL)
    if err != nil {
        log.Fatalf("migration init failed: %v", err)
    }
    defer m.Close()

    if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        log.Fatalf("migration up failed: %v", err)
    }

    slog.Info("migrations applied successfully")
}
```

---

## QUERY UNTUK KUR SCORE

```sql
-- Rata-rata pemasukan per bulan (3 bulan terakhir)
SELECT
    DATE_TRUNC('month', transaction_date) AS month,
    SUM(amount) AS monthly_total
FROM transactions
WHERE umkm_id = $1
  AND type = 'income'
  AND transaction_date >= NOW() - INTERVAL '3 months'
  AND deleted_at IS NULL
GROUP BY DATE_TRUNC('month', transaction_date)
ORDER BY month;

-- Konsistensi: berapa hari dalam 30 hari terakhir ada transaksi
SELECT COUNT(DISTINCT transaction_date) AS active_days
FROM transactions
WHERE umkm_id = $1
  AND transaction_date >= NOW() - INTERVAL '30 days'
  AND deleted_at IS NULL;

-- Berapa bulan data yang ada
SELECT
    DATE_PART('year', AGE(MAX(transaction_date), MIN(transaction_date))) * 12 +
    DATE_PART('month', AGE(MAX(transaction_date), MIN(transaction_date))) AS months_diff
FROM transactions
WHERE umkm_id = $1 AND deleted_at IS NULL;
```

---

## ATURAN SQL

```sql
-- 1. SELALU filter deleted_at IS NULL untuk soft delete
WHERE deleted_at IS NULL  -- WAJIB di semua query transaksi

-- 2. SELALU gunakan COALESCE untuk SUM agar tidak return NULL
COALESCE(SUM(amount), 0)  -- bukan SUM(amount) langsung

-- 3. SELALU gunakan parameterized query ($1, $2, ...)
-- JANGAN string concatenation

-- 4. Amount selalu dalam Rupiah sebagai INTEGER (BIGINT)
-- JANGAN DECIMAL/FLOAT untuk uang (floating point error)
-- "Rp1.500.500" = 1500500 (bigint)

-- 5. Timestamp selalu TIMESTAMPTZ (dengan timezone)
-- JANGAN TIMESTAMP WITHOUT TIME ZONE

-- 6. Date-only field (tanggal transaksi) pakai DATE, bukan TIMESTAMP
transaction_date DATE  -- bukan TIMESTAMPTZ

-- 7. Index pada kolom yang sering di-filter
CREATE INDEX idx_transactions_umkm_date
ON transactions(umkm_id, transaction_date)
WHERE deleted_at IS NULL;
```

---

## COMMON ERRORS DAN PENANGANAN

```go
import (
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

// No rows found
if errors.Is(err, pgx.ErrNoRows) {
    return nil, ErrNotFound
}

// Duplicate key (unique constraint violation)
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) {
    switch pgErr.Code {
    case "23505": // unique_violation
        return ErrAlreadyExists
    case "23503": // foreign_key_violation
        return ErrInvalidReference
    case "23514": // check_violation
        return ErrInvalidData
    }
}
```

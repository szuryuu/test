package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kasiraiai/backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionFilter struct {
	Type      string
	StartDate string
	EndDate   string
	Page      int
	Limit     int
}

type SummaryResult struct {
	TotalIncome  int64
	TotalExpense int64
	NetProfit    int64
	Count        int
}

type ChartDataPoint struct {
	Date   string `json:"date"`
	Income int64  `json:"income"`
	Expense int64  `json:"expense"`
}

type CategorySummary struct {
	Name       string  `json:"name"`
	Total      int64   `json:"total"`
	Percentage float64 `json:"percentage"`
}

type TransactionRepo interface {
	Create(ctx context.Context, tx *model.Transaction) error
	FindByUmkmID(ctx context.Context, umkmID uuid.UUID, filter TransactionFilter) ([]model.Transaction, int, error)
	FindLastByUmkmID(ctx context.Context, umkmID uuid.UUID) (*model.Transaction, error)
	SoftDelete(ctx context.Context, id uuid.UUID, umkmID uuid.UUID) error
	SumByPeriod(ctx context.Context, umkmID uuid.UUID, start, end time.Time) (SummaryResult, error)
	ChartData(ctx context.Context, umkmID uuid.UUID, start, end time.Time) ([]ChartDataPoint, error)
	CategorySummary(ctx context.Context, umkmID uuid.UUID, start, end time.Time, txType string) ([]CategorySummary, error)
	MonthlyAverages(ctx context.Context, umkmID uuid.UUID, months int) (incomeAvg, expenseAvg int64, err error)
	MonthlyIncomeValues(ctx context.Context, umkmID uuid.UUID, months int) ([]int64, error)
	ConsistencyPct(ctx context.Context, umkmID uuid.UUID, days int) (float64, error)
	DataMonths(ctx context.Context, umkmID uuid.UUID) (int, error)
}

type transactionRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionRepo(pool *pgxpool.Pool) TransactionRepo {
	return &transactionRepo{pool: pool}
}

func (r *transactionRepo) Create(ctx context.Context, tx *model.Transaction) error {
	query := `
		INSERT INTO transactions (umkm_id, category_id, amount, type, description, raw_message,
		                          transaction_date, source, ai_confidence)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		tx.UmkmID, tx.CategoryID, tx.Amount, tx.Type,
		tx.Description, tx.RawMessage, tx.TransactionDate,
		tx.Source, tx.AIConfidence,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		return fmt.Errorf("transaction_repo.Create: %w", err)
	}
	return nil
}

func (r *transactionRepo) FindByUmkmID(ctx context.Context, umkmID uuid.UUID, filter TransactionFilter) ([]model.Transaction, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	var conditions []string
	args := []interface{}{umkmID}
	argIdx := 2

	conditions = append(conditions, fmt.Sprintf("umkm_id = $1"))

	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, filter.Type)
		argIdx++
	}
	if filter.StartDate != "" {
		conditions = append(conditions, fmt.Sprintf("transaction_date >= $%d::date", argIdx))
		args = append(args, filter.StartDate)
		argIdx++
	}
	if filter.EndDate != "" {
		conditions = append(conditions, fmt.Sprintf("transaction_date <= $%d::date", argIdx))
		args = append(args, filter.EndDate)
		argIdx++
	}

	whereClause := "deleted_at IS NULL AND " + strings.Join(conditions, " AND ")

	// Count total
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions WHERE %s", whereClause)
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("transaction_repo.FindByUmkmID: count: %w", err)
	}

	// Fetch page
	offset := (filter.Page - 1) * filter.Limit
	selectQuery := fmt.Sprintf(`
		SELECT id, umkm_id, category_id, amount, type, description, raw_message,
		       transaction_date::text, source, ai_confidence, created_at, updated_at, deleted_at
		FROM transactions
		WHERE %s
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, filter.Limit, offset)

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("transaction_repo.FindByUmkmID: query: %w", err)
	}
	defer rows.Close()

	var txs []model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(
			&tx.ID, &tx.UmkmID, &tx.CategoryID, &tx.Amount, &tx.Type,
			&tx.Description, &tx.RawMessage, &tx.TransactionDate,
			&tx.Source, &tx.AIConfidence, &tx.CreatedAt, &tx.UpdatedAt, &tx.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("transaction_repo.FindByUmkmID: scan: %w", err)
		}
		txs = append(txs, tx)
	}

	return txs, total, nil
}

func (r *transactionRepo) FindLastByUmkmID(ctx context.Context, umkmID uuid.UUID) (*model.Transaction, error) {
	query := `
		SELECT id, umkm_id, category_id, amount, type, description, raw_message,
		       transaction_date::text, source, ai_confidence, created_at, updated_at, deleted_at
		FROM transactions
		WHERE umkm_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1`

	tx := &model.Transaction{}
	err := r.pool.QueryRow(ctx, query, umkmID).Scan(
		&tx.ID, &tx.UmkmID, &tx.CategoryID, &tx.Amount, &tx.Type,
		&tx.Description, &tx.RawMessage, &tx.TransactionDate,
		&tx.Source, &tx.AIConfidence, &tx.CreatedAt, &tx.UpdatedAt, &tx.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("transaction_repo.FindLastByUmkmID: %w", err)
	}
	return tx, nil
}

func (r *transactionRepo) SoftDelete(ctx context.Context, id uuid.UUID, umkmID uuid.UUID) error {
	query := `UPDATE transactions SET deleted_at = NOW() WHERE id = $1 AND umkm_id = $2 AND deleted_at IS NULL`
	_, err := r.pool.Exec(ctx, query, id, umkmID)
	if err != nil {
		return fmt.Errorf("transaction_repo.SoftDelete: %w", err)
	}
	return nil
}

func (r *transactionRepo) SumByPeriod(ctx context.Context, umkmID uuid.UUID, start, end time.Time) (SummaryResult, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expense,
			COUNT(*)
		FROM transactions
		WHERE umkm_id = $1
		  AND deleted_at IS NULL
		  AND transaction_date >= $2::date
		  AND transaction_date <  $3::date`

	var result SummaryResult
	err := r.pool.QueryRow(ctx, query, umkmID, start.Format("2006-01-02"), end.Format("2006-01-02")).Scan(
		&result.TotalIncome, &result.TotalExpense, &result.Count,
	)
	if err != nil && err != pgx.ErrNoRows {
		return result, fmt.Errorf("transaction_repo.SumByPeriod: %w", err)
	}

	result.NetProfit = result.TotalIncome - result.TotalExpense
	return result, nil
}

func (r *transactionRepo) ChartData(ctx context.Context, umkmID uuid.UUID, start, end time.Time) ([]ChartDataPoint, error) {
	query := `
		SELECT
			transaction_date::text AS date,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE umkm_id = $1
		  AND deleted_at IS NULL
		  AND transaction_date >= $2::date
		  AND transaction_date <  $3::date
		GROUP BY transaction_date
		ORDER BY transaction_date ASC`

	rows, err := r.pool.Query(ctx, query, umkmID, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("transaction_repo.ChartData: %w", err)
	}
	defer rows.Close()

	var data []ChartDataPoint
	for rows.Next() {
		var d ChartDataPoint
		if err := rows.Scan(&d.Date, &d.Income, &d.Expense); err != nil {
			return nil, fmt.Errorf("transaction_repo.ChartData: scan: %w", err)
		}
		data = append(data, d)
	}

	return data, nil
}

func (r *transactionRepo) CategorySummary(ctx context.Context, umkmID uuid.UUID, start, end time.Time, txType string) ([]CategorySummary, error) {
	query := `
		WITH category_totals AS (
			SELECT
				COALESCE(tc.name, 'Tanpa Kategori') AS name,
				SUM(t.amount) AS total
			FROM transactions t
			LEFT JOIN transaction_categories tc ON t.category_id = tc.id
			WHERE t.umkm_id = $1
			  AND t.deleted_at IS NULL
			  AND t.type = $2
			  AND t.transaction_date >= $3::date
			  AND t.transaction_date <  $4::date
			GROUP BY tc.name
		)
		SELECT name, total,
			CASE WHEN SUM(total) OVER() > 0
				THEN ROUND(total * 100.0 / SUM(total) OVER(), 1)
				ELSE 0
			END AS percentage
		FROM category_totals
		ORDER BY total DESC`

	rows, err := r.pool.Query(ctx, query, umkmID, txType, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("transaction_repo.CategorySummary: %w", err)
	}
	defer rows.Close()

	var cats []CategorySummary
	for rows.Next() {
		var c CategorySummary
		if err := rows.Scan(&c.Name, &c.Total, &c.Percentage); err != nil {
			return nil, fmt.Errorf("transaction_repo.CategorySummary: scan: %w", err)
		}
		cats = append(cats, c)
	}

	return cats, nil
}

func (r *transactionRepo) MonthlyAverages(ctx context.Context, umkmID uuid.UUID, months int) (int64, int64, error) {
	query := `
		SELECT
			COALESCE(AVG(monthly_income), 0)::bigint,
			COALESCE(AVG(monthly_expense), 0)::bigint
		FROM (
			SELECT
				DATE_TRUNC('month', transaction_date) AS month,
				COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS monthly_income,
				COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS monthly_expense
			FROM transactions
			WHERE umkm_id = $1
			  AND deleted_at IS NULL
			  AND transaction_date >= DATE_TRUNC('month', NOW()) - ($2::int || ' months')::interval
			GROUP BY DATE_TRUNC('month', transaction_date)
		) monthly_data`

	var incomeAvg, expenseAvg int64
	err := r.pool.QueryRow(ctx, query, umkmID, months).Scan(&incomeAvg, &expenseAvg)
	if err != nil {
		return 0, 0, fmt.Errorf("transaction_repo.MonthlyAverages: %w", err)
	}
	return incomeAvg, expenseAvg, nil
}

func (r *transactionRepo) MonthlyIncomeValues(ctx context.Context, umkmID uuid.UUID, months int) ([]int64, error) {
	query := `
		SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) AS monthly_income
		FROM transactions
		WHERE umkm_id = $1
		  AND deleted_at IS NULL
		  AND transaction_date >= DATE_TRUNC('month', NOW()) - ($2::int || ' months')::interval
		GROUP BY DATE_TRUNC('month', transaction_date)
		ORDER BY DATE_TRUNC('month', transaction_date) ASC`

	rows, err := r.pool.Query(ctx, query, umkmID, months)
	if err != nil {
		return nil, fmt.Errorf("transaction_repo.MonthlyIncomeValues: %w", err)
	}
	defer rows.Close()

	var values []int64
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("transaction_repo.MonthlyIncomeValues: scan: %w", err)
		}
		values = append(values, v)
	}
	return values, nil
}

func (r *transactionRepo) ConsistencyPct(ctx context.Context, umkmID uuid.UUID, days int) (float64, error) {
	query := `
		SELECT
			CASE WHEN $2 <= 0 THEN 0.0
				ELSE ROUND(
					COUNT(DISTINCT transaction_date)::numeric / $2::numeric * 100.0, 1
				)
			END
		FROM transactions
		WHERE umkm_id = $1
		  AND deleted_at IS NULL
		  AND transaction_date >= CURRENT_DATE - ($2::int || ' days')::interval`

	var pct float64
	err := r.pool.QueryRow(ctx, query, umkmID, days).Scan(&pct)
	if err != nil {
		return 0, fmt.Errorf("transaction_repo.ConsistencyPct: %w", err)
	}
	return pct, nil
}

func (r *transactionRepo) DataMonths(ctx context.Context, umkmID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(DISTINCT DATE_TRUNC('month', transaction_date))::int
		FROM transactions
		WHERE umkm_id = $1 AND deleted_at IS NULL`

	var count int
	err := r.pool.QueryRow(ctx, query, umkmID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("transaction_repo.DataMonths: %w", err)
	}
	return count, nil
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type KurScore struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	UmkmID             uuid.UUID  `json:"umkm_id" db:"umkm_id"`
	Score              int        `json:"score" db:"score"`
	Level              string     `json:"level" db:"level"` // rendah | sedang | baik | sangat_baik
	MonthlyIncomeAvg   int64      `json:"monthly_income_avg" db:"monthly_income_avg"`
	MonthlyExpenseAvg  int64      `json:"monthly_expense_avg" db:"monthly_expense_avg"`
	ProfitMargin       *float64   `json:"profit_margin,omitempty" db:"profit_margin"`
	ConsistencyScore   *int       `json:"consistency_score,omitempty" db:"consistency_score"`
	MonthsOfData       int        `json:"months_of_data" db:"months_of_data"`
	Recommendations    []string   `json:"recommendations" db:"recommendations"`
	CalculatedAt       time.Time  `json:"calculated_at" db:"calculated_at"`
}

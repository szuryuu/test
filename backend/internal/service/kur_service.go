package service

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"kasiraiai/backend/internal/model"
	"kasiraiai/backend/internal/repository"

	"github.com/google/uuid"
)

type KurScoreResult struct {
	Score             int      `json:"score"`
	Level             string   `json:"level"`
	MonthlyIncomeAvg  int64    `json:"monthly_income_avg"`
	MonthlyExpenseAvg int64    `json:"monthly_expense_avg"`
	ProfitMargin      float64  `json:"profit_margin"`
	ConsistencyScore  int      `json:"consistency_score"`
	MonthsOfData      int      `json:"months_of_data"`
	Recommendations   []string `json:"recommendations"`
	CalculatedAt      string   `json:"calculated_at"`
}

type KurService interface {
	GetScore(ctx context.Context, umkmID uuid.UUID) (*KurScoreResult, error)
	Recalculate(ctx context.Context, umkmID uuid.UUID) (*KurScoreResult, error)
}

type kurService struct {
	repo repository.TransactionRepo
}

func NewKurService(repo repository.TransactionRepo) KurService {
	return &kurService{repo: repo}
}

func (s *kurService) GetScore(ctx context.Context, umkmID uuid.UUID) (*KurScoreResult, error) {
	// Coba ambil skor terbaru dari DB dulu
	cached, err := s.repo.GetLatestKurScore(ctx, umkmID)
	if err != nil {
		slog.Warn("gagal ambil skor KUR dari cache, hitung ulang", "umkm_id", umkmID, "error", err)
		return s.calculate(ctx, umkmID)
	}
	if cached != nil {
		return kurScoreToResult(cached), nil
	}
	// Belum ada skor — hitung baru
	return s.calculate(ctx, umkmID)
}

func (s *kurService) Recalculate(ctx context.Context, umkmID uuid.UUID) (*KurScoreResult, error) {
	result, err := s.calculate(ctx, umkmID)
	if err != nil {
		return nil, err
	}

	// Persist ke database
	margin := result.ProfitMargin
	consistency := result.ConsistencyScore
	ks := &model.KurScore{
		UmkmID:            umkmID,
		Score:             result.Score,
		Level:             result.Level,
		MonthlyIncomeAvg:  result.MonthlyIncomeAvg,
		MonthlyExpenseAvg: result.MonthlyExpenseAvg,
		ProfitMargin:      &margin,
		ConsistencyScore:  &consistency,
		MonthsOfData:      result.MonthsOfData,
		Recommendations:   result.Recommendations,
	}
	if saveErr := s.repo.SaveKurScore(ctx, ks); saveErr != nil {
		slog.Error("gagal simpan skor KUR ke DB", "umkm_id", umkmID, "error", saveErr)
		// Tetap kembalikan hasil walau gagal simpan
	}

	return result, nil
}

func (s *kurService) calculate(ctx context.Context, umkmID uuid.UUID) (*KurScoreResult, error) {
	incomeAvg, expenseAvg, err := s.repo.MonthlyAverages(ctx, umkmID, 3)
	if err != nil {
		return nil, fmt.Errorf("kur_service.calculate: averages: %w", err)
	}

	consistencyPct, err := s.repo.ConsistencyPct(ctx, umkmID, 30)
	if err != nil {
		return nil, fmt.Errorf("kur_service.calculate: consistency: %w", err)
	}

	months, err := s.repo.DataMonths(ctx, umkmID)
	if err != nil {
		return nil, fmt.Errorf("kur_service.calculate: data_months: %w", err)
	}

	// Default margin
	margin := 0.0
	if incomeAvg > 0 {
		margin = float64(incomeAvg-expenseAvg) / float64(incomeAvg) * 100
	}

	// Income stability — hitung koefisien variasi dari pemasukan 3 bulan terakhir
	incomeStability := 20 // default sedang
	monthlyIncomes, err := s.repo.MonthlyIncomeValues(ctx, umkmID, 3)
	if err != nil {
		slog.Warn("gagal ambil data stabilitas pendapatan", "umkm_id", umkmID, "error", err)
	} else if len(monthlyIncomes) >= 2 {
		mean := float64(incomeAvg)
		if mean > 0 {
			// Hitung varians
			var varianceSum float64
			for _, v := range monthlyIncomes {
				diff := float64(v) - mean
				varianceSum += diff * diff
			}
			variance := varianceSum / float64(len(monthlyIncomes))
			stddev := math.Sqrt(variance)
			cv := stddev / mean

			switch {
			case cv < 0.2:
				incomeStability = 30
			case cv < 0.4:
				incomeStability = 20
			case cv < 0.6:
				incomeStability = 10
			default:
				incomeStability = 0
			}
			slog.Info("stabilitas pendapatan dihitung", "umkm_id", umkmID, "cv", cv, "months", len(monthlyIncomes), "score", incomeStability)
		}
	}

	profitMarginScore := calcMarginScore(margin)
	consistencyScore := calcConsistencyScore(int(consistencyPct))
	longevityScore := calcLongevityScore(months)

	totalScore := incomeStability + profitMarginScore + consistencyScore + longevityScore
	if totalScore > 100 {
		totalScore = 100
	}

	level := scoreLevel(totalScore)
	recommendations := generateRecommendations(totalScore, margin, int(consistencyPct), months)

	slog.Info("KUR score dihitung", "umkm_id", umkmID, "score", totalScore, "level", level)

	return &KurScoreResult{
		Score:             totalScore,
		Level:             level,
		MonthlyIncomeAvg:  incomeAvg,
		MonthlyExpenseAvg: expenseAvg,
		ProfitMargin:      math.Round(margin*100) / 100,
		ConsistencyScore:  consistencyScore,
		MonthsOfData:      months,
		Recommendations:   recommendations,
		CalculatedAt:      time.Now().Format(time.RFC3339),
	}, nil
}

func calcMarginScore(margin float64) int {
	switch {
	case margin >= 30:
		return 25
	case margin >= 20:
		return 18
	case margin >= 10:
		return 10
	default:
		return 0
	}
}

func calcConsistencyScore(pct int) int {
	switch {
	case pct >= 80:
		return 25
	case pct >= 60:
		return 18
	case pct >= 40:
		return 10
	default:
		return 0
	}
}

func calcLongevityScore(months int) int {
	switch {
	case months >= 6:
		return 20
	case months >= 3:
		return 14
	case months >= 2:
		return 8
	case months >= 1:
		return 4
	default:
		return 0
	}
}

func scoreLevel(score int) string {
	switch {
	case score >= 80:
		return "sangat_baik"
	case score >= 60:
		return "baik"
	case score >= 40:
		return "sedang"
	default:
		return "rendah"
	}
}

func generateRecommendations(score int, margin float64, consistency, months int) []string {
	var recs []string
	if consistency < 80 {
		recs = append(recs, "Tingkatkan konsistensi pencatatan harian untuk meningkatkan skor")
	}
	if margin < 20 {
		recs = append(recs, "Kurangi pengeluaran atau tingkatkan pemasukan untuk memperbaiki margin laba")
	} else {
		recs = append(recs, "Margin laba Anda sudah baik, pertahankan efisiensi pengeluaran")
	}
	if months < 3 {
		recs = append(recs, fmt.Sprintf("Lanjutkan pencatatan rutin minimal %d bulan lagi untuk meningkatkan skor", 3-months))
	}
	return recs
}

// kurScoreToResult mengkonversi model.KurScore (DB) ke KurScoreResult (API response).
func kurScoreToResult(ks *model.KurScore) *KurScoreResult {
	margin := 0.0
	if ks.ProfitMargin != nil {
		margin = *ks.ProfitMargin
	}
	consistency := 0
	if ks.ConsistencyScore != nil {
		consistency = *ks.ConsistencyScore
	}
	return &KurScoreResult{
		Score:             ks.Score,
		Level:             ks.Level,
		MonthlyIncomeAvg:  ks.MonthlyIncomeAvg,
		MonthlyExpenseAvg: ks.MonthlyExpenseAvg,
		ProfitMargin:      margin,
		ConsistencyScore:  consistency,
		MonthsOfData:      ks.MonthsOfData,
		Recommendations:   ks.Recommendations,
		CalculatedAt:      ks.CalculatedAt.Format(time.RFC3339),
	}
}

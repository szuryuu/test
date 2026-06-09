package service

import (
	"context"
	"fmt"
	"time"

	"kasiraiai/backend/internal/repository"

	"github.com/google/uuid"
)

type DashboardSummary struct {
	Period           string                     `json:"period"`
	TotalIncome      int64                      `json:"total_income"`
	TotalExpense     int64                      `json:"total_expense"`
	NetProfit        int64                      `json:"net_profit"`
	ProfitMargin     float64                    `json:"profit_margin"`
	TransactionCount int                        `json:"transaction_count"`
	ChartData        []repository.ChartDataPoint `json:"chart_data"`
}

type CategoryBreakdown struct {
	Categories []repository.CategorySummary `json:"categories"`
}

type DashboardService interface {
	GetSummary(ctx context.Context, umkmID uuid.UUID, period, date string) (*DashboardSummary, error)
	GetCategoryBreakdown(ctx context.Context, umkmID uuid.UUID, period, txType string) (*CategoryBreakdown, error)
}

type dashboardService struct {
	repo repository.TransactionRepo
}

func NewDashboardService(repo repository.TransactionRepo) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetSummary(ctx context.Context, umkmID uuid.UUID, period, date string) (*DashboardSummary, error) {
	start, end := periodRange(period, date)
	result, err := s.repo.SumByPeriod(ctx, umkmID, start, end)
	if err != nil {
		return nil, fmt.Errorf("dashboard_service.GetSummary: %w", err)
	}

	margin := 0.0
	if result.TotalIncome > 0 {
		margin = float64(result.NetProfit) / float64(result.TotalIncome) * 100
	}

	chartData, err := s.repo.ChartData(ctx, umkmID, start, end)
	if err != nil {
		return nil, fmt.Errorf("dashboard_service.GetSummary: chart: %w", err)
	}

	return &DashboardSummary{
		Period:           period,
		TotalIncome:      result.TotalIncome,
		TotalExpense:     result.TotalExpense,
		NetProfit:        result.NetProfit,
		ProfitMargin:     margin,
		TransactionCount: result.Count,
		ChartData:        chartData,
	}, nil
}

func (s *dashboardService) GetCategoryBreakdown(ctx context.Context, umkmID uuid.UUID, period, txType string) (*CategoryBreakdown, error) {
	start, end := periodRange(period, "")
	cats, err := s.repo.CategorySummary(ctx, umkmID, start, end, txType)
	if err != nil {
		return nil, fmt.Errorf("dashboard_service.GetCategoryBreakdown: %w", err)
	}
	return &CategoryBreakdown{Categories: cats}, nil
}

func periodRange(period, dateStr string) (start, end time.Time) {
	now := time.Now()
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			now = parsed
		}
	}

	year, month, day := now.Date()

	switch period {
	case "daily":
		start = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		end = start.Add(24 * time.Hour)
	case "weekly":
		weekday := now.Weekday()
		if weekday == 0 {
			weekday = 7
		}
		start = time.Date(year, month, day-int(weekday-1), 0, 0, 0, 0, now.Location())
		end = start.Add(7 * 24 * time.Hour)
	case "monthly":
		start = time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0)
	case "yearly":
		start = time.Date(year, 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0)
	default:
		start = time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0)
	}

	return start, end
}

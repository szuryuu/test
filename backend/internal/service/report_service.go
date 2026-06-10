package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"kasiraiai/backend/internal/repository"
	"kasiraiai/backend/pkg/fonnte"

	"github.com/google/uuid"
)

type MonthlyReport struct {
	ReportText    string   `json:"report_text"`
	TotalIncome   int64    `json:"total_income"`
	TotalExpense  int64    `json:"total_expense"`
	NetProfit     int64    `json:"net_profit"`
	TopCategories []string `json:"top_categories"`
}

type ReportService interface {
	GenerateMonthly(ctx context.Context, umkmID uuid.UUID, year, month int) (*MonthlyReport, error)
	GenerateAndSend(ctx context.Context, umkmID uuid.UUID, phoneNumber string, year, month int) (string, error)
	SendReport(ctx context.Context, phoneNumber, reportText string) error
}

type reportService struct {
	repo   repository.TransactionRepo
	fonnte *fonnte.Client
}

func NewReportService(repo repository.TransactionRepo, fonnteClient *fonnte.Client) ReportService {
	return &reportService{repo: repo, fonnte: fonnteClient}
}

var monthNames = []string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func (s *reportService) GenerateMonthly(ctx context.Context, umkmID uuid.UUID, year, month int) (*MonthlyReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)

	// Ringkasan periode
	summary, err := s.repo.SumByPeriod(ctx, umkmID, start, end)
	if err != nil {
		return nil, fmt.Errorf("report_service.GenerateMonthly: sum: %w", err)
	}

	// Kategori pengeluaran
	expCats, err := s.repo.CategorySummary(ctx, umkmID, start, end, "expense")
	if err != nil {
		return nil, fmt.Errorf("report_service.GenerateMonthly: expense cats: %w", err)
	}

	// Kategori pemasukan
	incCats, err := s.repo.CategorySummary(ctx, umkmID, start, end, "income")
	if err != nil {
		return nil, fmt.Errorf("report_service.GenerateMonthly: income cats: %w", err)
	}

	// Susun top categories
	var topCategories []string
	for _, c := range incCats {
		topCategories = append(topCategories, fmt.Sprintf("📥 %s: Rp%s (%.0f%%)", c.Name, formatRupiahReport(c.Total), c.Percentage))
	}
	for _, c := range expCats {
		topCategories = append(topCategories, fmt.Sprintf("📤 %s: Rp%s (%.0f%%)", c.Name, formatRupiahReport(c.Total), c.Percentage))
	}

	// Hitung margin
	margin := 0.0
	if summary.TotalIncome > 0 {
		margin = float64(summary.NetProfit) / float64(summary.TotalIncome) * 100
	}

	// Status kesehatan keuangan
	healthStatus := "sehat"
	if margin < 10 {
		healthStatus = "perlu perhatian"
	} else if margin < 20 {
		healthStatus = "cukup baik"
	}

	// Bangun teks laporan
	reportText := fmt.Sprintf(
		`📊 *LAPORAN KEUANGAN %s %d*

💰 *Pemasukan:* Rp%s
📤 *Pengeluaran:* Rp%s
📈 *Laba Bersih:* Rp%s
📐 *Margin Laba:* %.1f%%

📝 *Total Transaksi:* %d transaksi
💚 *Status:* Keuangan %s

📋 *Rincian Kategori:*
%s

---
📱 *KasirAI* — catat keuangan UMKM via WhatsApp`,
		strings.ToUpper(monthNames[month-1]),
		year,
		formatRupiahReport(summary.TotalIncome),
		formatRupiahReport(summary.TotalExpense),
		formatRupiahReport(summary.NetProfit),
		margin,
		summary.Count,
		healthStatus,
		strings.Join(topCategories, "\n"),
	)

	report := &MonthlyReport{
		ReportText:    reportText,
		TotalIncome:   summary.TotalIncome,
		TotalExpense:  summary.TotalExpense,
		NetProfit:     summary.NetProfit,
		TopCategories: topCategories,
	}

	slog.Info("laporan bulanan dibuat",
		"umkm_id", umkmID,
		"period", fmt.Sprintf("%d-%02d", year, month),
		"income", summary.TotalIncome,
		"expense", summary.TotalExpense,
		"profit", summary.NetProfit,
	)

	return report, nil
}

// SendReport mengirim teks laporan yang sudah di-generate via WhatsApp.
// Tidak melakukan query ulang — gunakan report yang sudah ada.
func (s *reportService) SendReport(ctx context.Context, phoneNumber, reportText string) error {
	return s.fonnte.SendMessage(phoneNumber, reportText)
}

// GenerateAndSend membuat laporan dan mengirimnya via WhatsApp.
// Dipanggil dari WhatsApp command handler.
func (s *reportService) GenerateAndSend(ctx context.Context, umkmID uuid.UUID, phoneNumber string, year, month int) (string, error) {
	report, err := s.GenerateMonthly(ctx, umkmID, year, month)
	if err != nil {
		return "❌ Gagal membuat laporan bulanan. Coba lagi nanti.", err
	}

	if err := s.fonnte.SendMessage(phoneNumber, report.ReportText); err != nil {
		slog.Error("gagal kirim laporan via WhatsApp", "phone", phoneNumber, "error", err)
		return "❌ Laporan berhasil dibuat tapi gagal dikirim via WhatsApp.", nil
	}

	slog.Info("laporan terkirim via WhatsApp", "umkm_id", umkmID, "phone", phoneNumber)
	return "✅ Laporan bulanan berhasil dikirim ke WhatsApp Anda!", nil
}

// formatRupiahReport memformat int64 ke string dengan separator titik.
func formatRupiahReport(amount int64) string {
	if amount < 0 {
		return "-" + formatRupiahReport(-amount)
	}
	if amount < 1000 {
		return fmt.Sprintf("%d", amount)
	}
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	var parts []string
	for i := n; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, ".")
}

package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"kasiraiai/backend/internal/model"
	"kasiraiai/backend/internal/repository"
	"kasiraiai/backend/pkg/ai"
	"kasiraiai/backend/pkg/fonnte"

	"github.com/google/uuid"
)

type TransactionService interface {
	ParseAndSave(ctx context.Context, umkmID uuid.UUID, rawMessage string, senderPhone string) (*ai.ParsedTransaction, error)
	Create(ctx context.Context, tx *model.Transaction) error
	FindByUmkmID(ctx context.Context, umkmID uuid.UUID, filter repository.TransactionFilter) ([]model.Transaction, int, error)
	SoftDelete(ctx context.Context, id uuid.UUID, umkmID uuid.UUID) error
	DeleteLast(ctx context.Context, umkmID uuid.UUID) error
}

type transactionService struct {
	repo         repository.TransactionRepo
	parser       *ai.Parser
	fonnte       *fonnte.Client
	dashboardSvc DashboardService
	reportSvc    ReportService
	kurSvc       KurService
}

func NewTransactionService(
	repo repository.TransactionRepo,
	parser *ai.Parser,
	fonnteClient *fonnte.Client,
	dashboardSvc DashboardService,
	reportSvc ReportService,
	kurSvc KurService,
) TransactionService {
	return &transactionService{
		repo:         repo,
		parser:       parser,
		fonnte:       fonnteClient,
		dashboardSvc: dashboardSvc,
		reportSvc:    reportSvc,
		kurSvc:       kurSvc,
	}
}

// ParseAndSave menerima pesan WhatsApp mentah, memproses dengan AI, dan menyimpan ke database.
// Jika pesan adalah command (bukan transaksi), command ditangani tanpa panggil AI.
func (s *transactionService) ParseAndSave(ctx context.Context, umkmID uuid.UUID, rawMessage string, senderPhone string) (*ai.ParsedTransaction, error) {
	// Deteksi command — hemat panggilan AI
	if cmd := detectCommand(rawMessage); cmd != "" {
		slog.Info("command terdeteksi", "cmd", cmd, "umkm_id", umkmID)
		return &ai.ParsedTransaction{
			IsTransaction: false,
			ReplyMessage:  s.handleCommand(ctx, umkmID, cmd, senderPhone),
			Confidence:    1.0,
		}, nil
	}

	parsed, err := s.parser.Parse(rawMessage)
	if err != nil {
		slog.Error("gagal parsing transaksi", "error", err, "umkm_id", umkmID)
		return &ai.ParsedTransaction{
			IsTransaction: false,
			ReplyMessage:  fonnte.MsgAIParseFailure,
		}, nil
	}

	if !parsed.IsTransaction {
		return parsed, nil
	}

	tx := &model.Transaction{
		UmkmID:          umkmID,
		Amount:          parsed.Amount,
		Type:            parsed.Type,
		Description:     parsed.Description,
		RawMessage:      rawMessage,
		TransactionDate: parsed.TransactionDate,
		Source:          "whatsapp",
		AIConfidence:    &parsed.Confidence,
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("transaction_service.ParseAndSave: create: %w", err)
	}

	slog.Info("transaksi tersimpan",
		"id", tx.ID,
		"umkm_id", umkmID,
		"type", tx.Type,
		"amount", tx.Amount,
		"confidence", parsed.Confidence,
	)

	return parsed, nil
}

func (s *transactionService) FindByUmkmID(ctx context.Context, umkmID uuid.UUID, filter repository.TransactionFilter) ([]model.Transaction, int, error) {
	return s.repo.FindByUmkmID(ctx, umkmID, filter)
}

func (s *transactionService) Create(ctx context.Context, tx *model.Transaction) error {
	if err := s.repo.Create(ctx, tx); err != nil {
		return fmt.Errorf("transaction_service.Create: %w", err)
	}
	slog.Info("transaksi manual tersimpan", "id", tx.ID, "umkm_id", tx.UmkmID, "type", tx.Type, "amount", tx.Amount)
	return nil
}

func (s *transactionService) SoftDelete(ctx context.Context, id uuid.UUID, umkmID uuid.UUID) error {
	if err := s.repo.SoftDelete(ctx, id, umkmID); err != nil {
		return fmt.Errorf("transaction_service.SoftDelete: %w", err)
	}
	slog.Info("transaksi dihapus", "id", id, "umkm_id", umkmID)
	return nil
}

func (s *transactionService) DeleteLast(ctx context.Context, umkmID uuid.UUID) error {
	tx, err := s.repo.FindLastByUmkmID(ctx, umkmID)
	if err != nil {
		return fmt.Errorf("transaction_service.DeleteLast: %w", err)
	}
	return s.SoftDelete(ctx, tx.ID, umkmID)
}

func (s *transactionService) handleCommand(ctx context.Context, umkmID uuid.UUID, cmd string, senderPhone string) string {
	switch cmd {
	case "bantuan":
		return fonnte.MsgHelp

	case "hapus_terakhir":
		if err := s.DeleteLast(ctx, umkmID); err != nil {
			slog.Error("gagal hapus transaksi terakhir", "umkm_id", umkmID, "error", err)
			return "❌ Gagal menghapus transaksi terakhir. Coba lagi nanti."
		}
		return "✅ Transaksi terakhir berhasil dihapus."

	case "hari_ini":
		summary, err := s.dashboardSvc.GetSummary(ctx, umkmID, "daily", "")
		if err != nil {
			slog.Error("gagal ambil ringkasan harian", "umkm_id", umkmID, "error", err)
			return "❌ Gagal mengambil ringkasan hari ini."
		}
		return fmt.Sprintf(
			"📊 *Ringkasan Hari Ini*\n\n💰 Pemasukan: Rp%s\n📤 Pengeluaran: Rp%s\n📈 Laba Bersih: Rp%s\n📝 Total Transaksi: %d",
			formatRupiah(summary.TotalIncome),
			formatRupiah(summary.TotalExpense),
			formatRupiah(summary.NetProfit),
			summary.TransactionCount,
		)

	case "minggu_ini":
		summary, err := s.dashboardSvc.GetSummary(ctx, umkmID, "weekly", "")
		if err != nil {
			slog.Error("gagal ambil ringkasan mingguan", "umkm_id", umkmID, "error", err)
			return "❌ Gagal mengambil ringkasan minggu ini."
		}
		return fmt.Sprintf(
			"📊 *Ringkasan Minggu Ini*\n\n💰 Pemasukan: Rp%s\n📤 Pengeluaran: Rp%s\n📈 Laba Bersih: Rp%s\n📝 Total Transaksi: %d",
			formatRupiah(summary.TotalIncome),
			formatRupiah(summary.TotalExpense),
			formatRupiah(summary.NetProfit),
			summary.TransactionCount,
		)

	case "skor":
		result, err := s.kurSvc.Recalculate(ctx, umkmID)
		if err != nil {
			slog.Error("gagal hitung KUR via WhatsApp", "umkm_id", umkmID, "error", err)
			return "❌ Gagal menghitung skor KUR. Coba lagi nanti."
		}

		levelDesc := map[string]string{
			"sangat_baik": "Sangat siap mengajukan KUR! 🎉",
			"baik":        "Siap KUR dengan persiapan tambahan 👍",
			"sedang":      "Perlu tingkatkan konsistensi pencatatan 📝",
			"rendah":      "Fokus pencatatan rutin minimal 3 bulan 💪",
		}

		recs := ""
		for i, r := range result.Recommendations {
			recs += fmt.Sprintf("\n%d. %s", i+1, r)
		}

		return fmt.Sprintf(
			"⭐ *Skor KUR Anda: %d/100*\n\n"+
				"Status: %s\n\n"+
				"💰 Rata-rata pemasukan: Rp%s/bulan\n"+
				"📊 Margin laba: %.1f%%\n\n"+
				"📋 *Rekomendasi:*%s",
			result.Score,
			levelDesc[result.Level],
			formatRupiah(result.MonthlyIncomeAvg),
			result.ProfitMargin,
			recs,
		)

	case "laporan":
		now := time.Now()
		msg, err := s.reportSvc.GenerateAndSend(ctx, umkmID, senderPhone, now.Year(), int(now.Month()))
		if err != nil {
			slog.Error("gagal generate laporan via command", "umkm_id", umkmID, "error", err)
			return "❌ Gagal membuat laporan bulanan."
		}
		return msg

	default:
		return fonnte.MsgHelp
	}
}

// formatRupiah memformat int64 ke string Rupiah dengan separator ribuan.
func formatRupiah(amount int64) string {
	if amount < 0 {
		return "-" + formatRupiah(-amount)
	}
	if amount < 1000 {
		return fmt.Sprintf("%d", amount)
	}
	// Format manual: reverse, insert dots every 3 digits, reverse back
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

// detectCommand mengembalikan command keyword jika pesan cocok, string kosong jika bukan.
// Command detection dilakukan sebelum memanggil AI untuk menghemat biaya.
func detectCommand(msg string) string {
	msg = strings.TrimSpace(strings.ToLower(msg))

	// Normalisasi: hapus tanda baca di awal/akhir
	msg = strings.Trim(msg, "!.?*#")

	switch msg {
	case "bantuan", "help", "halo", "hai", "hi", "hello":
		return "bantuan"
	case "laporan", "laporan bulan ini", "laporan bulanan":
		return "laporan"
	case "skor", "kur saya", "skor kur", "cek skor", "skorku":
		return "skor"
	case "hari ini", "ringkasan hari ini", "today":
		return "hari_ini"
	case "minggu ini", "ringkasan minggu ini":
		return "minggu_ini"
	case "hapus terakhir", "hapus", "delete terakhir", "batalkan", "undo":
		return "hapus_terakhir"
	default:
		return ""
	}
}

package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"kasiraiai/backend/internal/repository"
	"kasiraiai/backend/internal/service"
	"kasiraiai/backend/pkg/fonnte"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	svc     service.TransactionService
	umkmRepo repository.UmkmRepo
	fonnte   *fonnte.Client
}

func NewWebhookHandler(
	svc service.TransactionService,
	umkmRepo repository.UmkmRepo,
	fonnteClient *fonnte.Client,
) *WebhookHandler {
	return &WebhookHandler{svc: svc, umkmRepo: umkmRepo, fonnte: fonnteClient}
}

// Handle menerima webhook dari Fonnte (multipart/form-data).
// Selalu return 200 ke Fonnte meskipun parsing gagal —
// pesan error dikirim ke user via WhatsApp, bukan via HTTP response.
func (h *WebhookHandler) Handle(c *gin.Context) {
	sender := c.PostForm("sender")
	message := c.PostForm("message")

	if sender == "" || message == "" {
		c.String(http.StatusOK, "OK")
		return
	}

	slog.Info("webhook diterima", "sender", sender, "message_len", len(message))

	// Cari UMKM berdasarkan nomor pengirim
	umkm, err := h.umkmRepo.FindByPhone(c.Request.Context(), sender)
	if err != nil {
		slog.Error("gagal mencari UMKM", "sender", sender, "error", err)
		// Kirim pesan error async — jangan tahan response Fonnte
		go func() {
			_ = h.fonnte.SendMessage(sender, "Mohon maaf, sistem sedang sibuk. Silakan coba lagi nanti.")
		}()
		c.String(http.StatusOK, "OK")
		return
	}

	// Pengirim belum terdaftar — kirim pesan selamat datang async
	if umkm == nil {
		slog.Info("pengirim belum terdaftar", "sender", sender)
		go func() {
			if sendErr := h.fonnte.SendMessage(sender, fonnte.MsgWelcome); sendErr != nil {
				slog.Error("gagal kirim welcome", "sender", sender, "error", sendErr)
			}
		}()
		c.String(http.StatusOK, "OK")
		return
	}

	// Parse dan simpan transaksi di background — jangan tahan response Fonnte.
	// AI call bisa 3-10 detik; Fonnte timeout ~10 detik dan akan retry.
	// JANGAN gunakan c.Request.Context() — dibatalkan setelah c.String return.
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

	// Return 200 ke Fonnte segera — sebelum goroutine selesai
	c.String(http.StatusOK, "OK")
}

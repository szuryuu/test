package handler

import (
	"log/slog"
	"net/http"

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
		c.String(http.StatusOK, "OK")
		return
	}

	// Pengirim belum terdaftar — kirim pesan selamat datang
	if umkm == nil {
		slog.Info("pengirim belum terdaftar", "sender", sender)
		if sendErr := h.fonnte.SendMessage(sender, fonnte.MsgWelcome); sendErr != nil {
			slog.Error("gagal kirim welcome", "sender", sender, "error", sendErr)
		}
		c.String(http.StatusOK, "OK")
		return
	}

	// Parse dan simpan transaksi
	parsed, err := h.svc.ParseAndSave(c.Request.Context(), umkm.ID, message, sender)
	if err != nil {
		slog.Error("gagal ParseAndSave", "umkm_id", umkm.ID, "error", err)
		_ = h.fonnte.SendMessage(sender, fonnte.MsgAIParseFailure)
		c.String(http.StatusOK, "OK")
		return
	}

	// Kirim balasan ke user via WhatsApp
	if parsed.ReplyMessage != "" {
		if sendErr := h.fonnte.SendMessage(sender, parsed.ReplyMessage); sendErr != nil {
			slog.Error("gagal kirim balasan", "sender", sender, "error", sendErr)
		}
	}

	c.String(http.StatusOK, "OK")
}

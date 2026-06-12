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
	svc      service.TransactionService
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

// Handle menerima webhook dari Fonnte.
// Mampu memproses Content-Type application/json maupun application/x-www-form-urlencoded.
// Selalu return 200 ke Fonnte meskipun parsing gagal.
func (h *WebhookHandler) Handle(c *gin.Context) {
	var payload struct {
		Sender  string `json:"sender" form:"sender"`
		Message string `json:"message" form:"message"`
		Device  string `json:"device" form:"device"`
	}

	// Coba bind otomatis berdasarkan Content-Type (JSON atau Form) dari Fonnte
	if err := c.ShouldBind(&payload); err != nil {
		// Fallback darurat jika header Fonnte tidak standar atau berantakan
		payload.Sender = c.PostForm("sender")
		payload.Message = c.PostForm("message")
		payload.Device = c.PostForm("device")
	}

	sender := payload.Sender
	message := payload.Message

	// Jika masih kosong, tolak secara diam-diam agar Fonnte tidak retry terus
	if sender == "" || message == "" {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "payload kosong atau tidak valid"})
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

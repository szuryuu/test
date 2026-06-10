package handler

import (
	"log/slog"
	"net/http"

	"kasiraiai/backend/internal/repository"
	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	svc      service.ReportService
	umkmRepo repository.UmkmRepo
}

func NewReportHandler(svc service.ReportService, umkmRepo repository.UmkmRepo) *ReportHandler {
	return &ReportHandler{svc: svc, umkmRepo: umkmRepo}
}

func (h *ReportHandler) GenerateMonthly(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	var req struct {
		Year  int `json:"year" binding:"required,gte=2024"`
		Month int `json:"month" binding:"required,gte=1,lte=12"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, ErrInvalidInput)
		return
	}

	report, err := h.svc.GenerateMonthly(c.Request.Context(), umkmID, req.Year, req.Month)
	if err != nil {
		slog.Error("gagal generate laporan bulanan", "umkm_id", umkmID, "error", err)
		InternalServerError(c, ErrInternalServer)
		return
	}

	// Kirim via WhatsApp sesuai API contract
	umkm, err := h.umkmRepo.FindByID(c.Request.Context(), umkmID)
	if err != nil {
		slog.Error("gagal cari UMKM untuk kirim laporan", "umkm_id", umkmID, "error", err)
		// Tetap return laporan meskipun WhatsApp gagal
		SuccessResponse(c, http.StatusOK, "Laporan berhasil dibuat (gagal kirim via WhatsApp)", report)
		return
	}
	if umkm != nil {
		if sendErr := h.svc.SendReport(c.Request.Context(), umkm.PhoneNumber, report.ReportText); sendErr != nil {
			slog.Error("gagal kirim laporan via WhatsApp", "umkm_id", umkmID, "error", sendErr)
		}
	}

	SuccessResponse(c, http.StatusOK, "Laporan berhasil dibuat dan dikirim via WhatsApp", report)
}

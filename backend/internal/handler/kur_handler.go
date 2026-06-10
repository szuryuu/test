package handler

import (
	"log/slog"
	"net/http"

	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type KurHandler struct {
	svc service.KurService
}

func NewKurHandler(svc service.KurService) *KurHandler {
	return &KurHandler{svc: svc}
}

func (h *KurHandler) GetScore(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	result, err := h.svc.GetScore(c.Request.Context(), umkmID)
	if err != nil {
		slog.Error("gagal ambil skor KUR", "umkm_id", umkmID, "error", err)
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Skor KUR berhasil diambil", result)
}

func (h *KurHandler) Recalculate(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	result, err := h.svc.Recalculate(c.Request.Context(), umkmID)
	if err != nil {
		slog.Error("gagal hitung ulang skor KUR", "umkm_id", umkmID, "error", err)
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Skor KUR berhasil dihitung ulang", result)
}

package handler

import (
	"log/slog"
	"net/http"

	"kasiraiai/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type UmkmHandler struct {
	repo repository.UmkmRepo
}

func NewUmkmHandler(repo repository.UmkmRepo) *UmkmHandler {
	return &UmkmHandler{repo: repo}
}

func (h *UmkmHandler) GetProfile(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	umkm, err := h.repo.FindByID(c.Request.Context(), umkmID)
	if err != nil {
		slog.Error("gagal ambil profil UMKM", "umkm_id", umkmID, "error", err)
		InternalServerError(c, ErrInternalServer)
		return
	}
	if umkm == nil {
		NotFound(c, ErrNotFound)
		return
	}

	SuccessResponse(c, http.StatusOK, "Profil berhasil diambil", umkm)
}

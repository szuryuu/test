package handler

import (
	"net/http"

	"kasiraiai/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc service.DashboardService
}

func NewDashboardHandler(svc service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Summary(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	period := c.DefaultQuery("period", "monthly")
	date := c.DefaultQuery("date", "")

	summary, err := h.svc.GetSummary(c.Request.Context(), umkmID, period, date)
	if err != nil {
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Ringkasan berhasil diambil", summary)
}

func (h *DashboardHandler) Categories(c *gin.Context) {
	umkmID, err := getUmkmID(c)
	if err != nil {
		Unauthorized(c, ErrUnauthorized)
		return
	}

	period := c.DefaultQuery("period", "monthly")
	txType := c.DefaultQuery("type", "expense")

	cats, err := h.svc.GetCategoryBreakdown(c.Request.Context(), umkmID, period, txType)
	if err != nil {
		InternalServerError(c, ErrInternalServer)
		return
	}

	SuccessResponse(c, http.StatusOK, "Kategori berhasil diambil", cats)
}
